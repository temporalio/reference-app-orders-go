package server

import (
	"context"
	"crypto/tls"
	_ "embed"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/temporalio/reference-app-orders-go/app/billing"
	"github.com/temporalio/reference-app-orders-go/app/config"
	"github.com/temporalio/reference-app-orders-go/app/fraud"
	"github.com/temporalio/reference-app-orders-go/app/order"
	"github.com/temporalio/reference-app-orders-go/app/shipment"
	mongodb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
	"golang.org/x/sync/errgroup"
)

// CreateClientOptionsFromEnv creates a client.Options instance, configures
// it based on environment variables, and returns that instance. It
// supports the following environment variables:
//
//	TEMPORAL_ADDRESS: Host and port (formatted as host:port) of the Temporal Frontend Service
//	TEMPORAL_NAMESPACE: Namespace to be used by the Client
//	TEMPORAL_TLS_CERT: Path to the x509 certificate
//	TEMPORAL_TLS_KEY: Path to the private certificate key
//
// If these environment variables are not set, the client.Options
// instance returned will be based on the SDK's default configuration.
func CreateClientOptionsFromEnv() (client.Options, error) {
	hostPort := os.Getenv("TEMPORAL_ADDRESS")
	namespaceName := os.Getenv("TEMPORAL_NAMESPACE")

	// Must explicitly set the Namepace for non-cloud use.
	if strings.Contains(hostPort, ".tmprl.cloud:") && namespaceName == "" {
		return client.Options{}, fmt.Errorf("Namespace name unspecified; required for Temporal Cloud")
	}

	if namespaceName == "" {
		namespaceName = "default"
		fmt.Printf("Namespace name unspecified; using value '%s'\n", namespaceName)
	}

	clientOpts := client.Options{
		HostPort:  hostPort,
		Namespace: namespaceName,
		Logger:    log.NewStructuredLogger(slog.Default()),
	}

	if certPath := os.Getenv("TEMPORAL_TLS_CERT"); certPath != "" {
		cert, err := tls.LoadX509KeyPair(certPath, os.Getenv("TEMPORAL_TLS_KEY"))
		if err != nil {
			return clientOpts, fmt.Errorf("failed loading key pair: %w", err)
		}

		clientOpts.ConnectionOptions.TLS = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}

	return clientOpts, nil
}

// SetupDB creates indexes in the database.
func SetupDB(db *mongodb.Database) error {
	orders := db.Collection(order.OrdersCollection)
	_, err := orders.Indexes().CreateOne(context.TODO(), mongodb.IndexModel{
		Keys: map[string]interface{}{"received_at": 1},
	})
	if err != nil {
		return fmt.Errorf("failed to create database index: %w", err)
	}

	shipments := db.Collection(shipment.ShipmentCollection)
	_, err = shipments.Indexes().CreateOne(context.TODO(), mongodb.IndexModel{
		Keys: map[string]interface{}{"booked_at": 1},
	})
	if err != nil {
		return fmt.Errorf("failed to create database index: %w", err)
	}

	return nil
}

// RunWorkers runs workers for the requested services.
func RunWorkers(ctx context.Context, config config.AppConfig, client client.Client, services []string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	for _, service := range services {
		switch service {
		case "billing":
			g.Go(func() error {
				return billing.RunWorker(ctx, config, client)
			})
		case "order":
			g.Go(func() error {
				return order.RunWorker(ctx, config, client)
			})
		case "shipment":
			g.Go(func() error {
				return shipment.RunWorker(ctx, config, client)
			})
		default:
			return fmt.Errorf("unknown service: %s", service)
		}
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

type instrumentedResponseWriter struct {
	http.ResponseWriter
	status int
}

func (r *instrumentedResponseWriter) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		iw := instrumentedResponseWriter{w, http.StatusOK}

		next.ServeHTTP(&iw, r)

		level := slog.LevelDebug
		if iw.status >= 500 {
			level = slog.LevelError
		}

		logger.Log(
			context.Background(), level,
			fmt.Sprintf("%d %s %s", iw.status, r.Method, r.URL.Path),
			"method", r.Method, "status", iw.status, "path", r.URL.Path,
		)
	})
}

// RunAPIServer runs a API HTTP server for the given service.
func runAPIServer(ctx context.Context, hostPort string, router http.Handler, logger *slog.Logger) error {
	srv := &http.Server{
		Addr:    hostPort,
		Handler: loggingMiddleware(logger, router),
	}

	logger.Info("Listening", "endpoint", "http://"+hostPort)

	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()

	select {
	case <-ctx.Done():
		srv.Close()
	case err := <-errCh:
		return err
	}

	return nil
}

func serverHostPort(config config.AppConfig, port int32) string {
	return fmt.Sprintf("%s:%d", config.BindOnIP, port)
}

// RunAPIServers runs API servers for the requested services.
func RunAPIServers(ctx context.Context, config config.AppConfig, client client.Client, services []string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	var db *mongodb.Database

	if slices.Contains(services, "orders") || slices.Contains(services, "shipment") {
		c, err := mongodb.Connect(context.TODO(), options.Client().ApplyURI(config.MongoURL))
		db = c.Database("orders")
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}

		if err := SetupDB(db); err != nil {
			return err
		}
	}

	for _, service := range services {
		logger := slog.Default().With("service", service)
		port, err := config.ServiceHostPort(service)
		if err != nil {
			return err
		}

		switch service {
		case "billing":
			g.Go(func() error {
				return runAPIServer(ctx, port, billing.Router(client, logger), logger)
			})
		case "fraud":
			g.Go(func() error {
				return runAPIServer(ctx, port, fraud.Router(logger), logger)
			})
		case "order":
			g.Go(func() error {
				return runAPIServer(ctx, port, order.Router(client, db, logger), logger)
			})
		case "shipment":
			g.Go(func() error {
				return runAPIServer(ctx, port, shipment.Router(client, db, logger), logger)
			})
		default:
			return fmt.Errorf("unknown service: %s", service)
		}
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
