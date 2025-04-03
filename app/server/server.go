package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/temporalio/reference-app-orders-go/app/billing"
	"github.com/temporalio/reference-app-orders-go/app/config"
	"github.com/temporalio/reference-app-orders-go/app/db"
	"github.com/temporalio/reference-app-orders-go/app/fraud"
	"github.com/temporalio/reference-app-orders-go/app/order"
	"github.com/temporalio/reference-app-orders-go/app/shipment"
	"github.com/uber-go/tally/v4"
	"github.com/uber-go/tally/v4/prometheus"
	"go.temporal.io/sdk/client"
	sdktally "go.temporal.io/sdk/contrib/tally"
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
	logger := slog.Default()

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
		Logger:    log.NewStructuredLogger(logger),
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

	endpoint := os.Getenv("TEMPORAL_METRICS_ENDPOINT")
	if endpoint != "" {
		scope, err := newPrometheusScope(prometheus.Configuration{
			ListenAddress: endpoint,
			TimerType:     "histogram",
		}, logger)
		if err != nil {
			return clientOpts, fmt.Errorf("failed to create metrics scope: %w", err)
		}
		clientOpts.MetricsHandler = sdktally.NewMetricsHandler(scope)
	}

	return clientOpts, nil
}

func newPrometheusScope(c prometheus.Configuration, logger *slog.Logger) (tally.Scope, error) {
	reporter, err := c.NewReporter(
		prometheus.ConfigurationOptions{
			Registry: prom.NewRegistry(),
			OnError: func(err error) {
				logger.Error("error in prometheus reporter", "error", err)
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error creating prometheus reporter: %w", err)
	}
	scopeOpts := tally.ScopeOptions{
		CachedReporter:  reporter,
		Separator:       prometheus.DefaultSeparator,
		SanitizeOptions: &sdktally.PrometheusSanitizeOptions,
	}
	scope, _ := tally.NewRootScope(scopeOpts, time.Second)
	scope = sdktally.NewPrometheusNamingScope(scope)

	return scope, nil
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

// RunAPIServers runs API servers for the requested services.
func RunAPIServers(ctx context.Context, config config.AppConfig, client client.Client, services []string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	db := db.CreateDB(config)

	if slices.Contains(services, "orders") || slices.Contains(services, "shipment") {
		err := db.Connect(context.TODO())
		if err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}

		if err := db.Setup(); err != nil {
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
