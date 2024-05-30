package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/temporalio/orders-reference-app-go/app/billing"
	"github.com/temporalio/orders-reference-app-go/app/config"
	"github.com/temporalio/orders-reference-app-go/app/fraudcheck"
	"github.com/temporalio/orders-reference-app-go/app/internal/temporalutil"
	"github.com/temporalio/orders-reference-app-go/app/order"
	"github.com/temporalio/orders-reference-app-go/app/shipment"
	"go.temporal.io/sdk/client"
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
	if temporalutil.IsTemporalCloud(hostPort) && namespaceName == "" {
		return client.Options{}, fmt.Errorf("Namespace name unspecified; required for Temporal Cloud")
	}

	if namespaceName == "" {
		namespaceName = "default"
		fmt.Printf("Namespace name unspecified; using value '%s'\n", namespaceName)
	}

	clientOpts := client.Options{
		HostPort:  hostPort,
		Namespace: namespaceName,
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

// RunServer runs all the workers and API servers for the Order/Shipment/Fraud/Billing system.
func RunServer(ctx context.Context, config config.AppConfig, client client.Client) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return billing.RunServer(ctx, config, client)
	})
	g.Go(func() error {
		return order.RunServer(ctx, config, client)
	})
	g.Go(func() error {
		return shipment.RunServer(ctx, config, client)
	})
	g.Go(func() error {
		return fraudcheck.RunServer(ctx, config)
	})

	g.Go(func() error {
		return billing.RunWorker(ctx, config, client)
	})
	g.Go(func() error {
		return shipment.RunWorker(ctx, config, client)
	})
	g.Go(func() error {
		return order.RunWorker(ctx, config, client)
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
