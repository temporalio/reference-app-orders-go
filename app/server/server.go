package server

import (
	"context"
	"crypto/tls"
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"slices"

	"github.com/jmoiron/sqlx"
	"github.com/temporalio/reference-app-orders-go/app/billing"
	"github.com/temporalio/reference-app-orders-go/app/config"
	"github.com/temporalio/reference-app-orders-go/app/fraudcheck"
	"github.com/temporalio/reference-app-orders-go/app/order"
	"github.com/temporalio/reference-app-orders-go/app/shipment"
	"github.com/temporalio/reference-app-orders-go/app/temporalutil"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/opentelemetry"
	tlog "go.temporal.io/sdk/log"
	"golang.org/x/sync/errgroup"
)

const meterName = "github.com/temporalio/reference-app-orders-go"

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
	if temporalutil.IsTemporalCloud(hostPort) && namespaceName == "" {
		return client.Options{}, fmt.Errorf("Namespace name unspecified; required for Temporal Cloud")
	}

	if namespaceName == "" {
		namespaceName = "default"
		slog.Warn("Namespace name unspecified; using default", "namespace", namespaceName)
	}

	clientOpts := client.Options{
		HostPort:  hostPort,
		Namespace: namespaceName,
		Logger:    tlog.NewStructuredLogger(slog.Default()),
	}

	if os.Getenv("PROMETHEUS_ENDPOINT") != "" {
		exporter, err := prometheus.New()
		if err != nil {
			return client.Options{}, err
		}
		provider := metric.NewMeterProvider(metric.WithReader(exporter))
		meter := provider.Meter(meterName)

		clientOpts.MetricsHandler = opentelemetry.NewMetricsHandler(
			opentelemetry.MetricsHandlerOptions{Meter: meter},
		)
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

//go:embed schema.sql
var schema string

// SetupDB creates the necessary tables in the database.
func SetupDB(db *sqlx.DB) error {
	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create the database schema: %w", err)
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
				slog.Info("Starting worker", "service", "billing")
				return billing.RunWorker(ctx, config, client)
			})
		case "order":
			g.Go(func() error {
				slog.Info("Starting worker", "service", "order")
				return order.RunWorker(ctx, config, client)
			})
		case "shipment":
			g.Go(func() error {
				slog.Info("Starting worker", "service", "shipment")
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

// RunAPIServers runs API servers for the requested services.
func RunAPIServers(ctx context.Context, config config.AppConfig, client client.Client, services []string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	var db *sqlx.DB
	var err error

	if slices.Contains(services, "orders") || slices.Contains(services, "shipment") {
		db, err = sqlx.Connect("sqlite", "./api-store.db")
		db.SetMaxOpenConns(1) // SQLite does not support concurrent writes
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}

		if err := SetupDB(db); err != nil {
			return err
		}
	}

	for _, service := range services {
		switch service {
		case "billing":
			g.Go(func() error {
				slog.Info("Starting API server", "service", "billing", "port", config.BillingPort)
				return billing.RunServer(ctx, config, client)
			})
		case "order":
			g.Go(func() error {
				slog.Info("Starting API server", "service", "order", "port", config.OrderPort)
				return order.RunServer(ctx, config, client, db)
			})
		case "shipment":
			g.Go(func() error {
				slog.Info("Starting API server", "service", "shipment", "port", config.ShipmentPort)
				return shipment.RunServer(ctx, config, client, db)
			})
		case "fraudcheck":
			g.Go(func() error {
				slog.Info("Starting API server", "service", "fraudcheck", "port", config.FraudCheckPort)
				return fraudcheck.RunServer(ctx, config)
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
