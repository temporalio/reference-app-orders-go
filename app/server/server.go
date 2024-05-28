package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/temporalio/orders-reference-app-go/app/billing"
	"github.com/temporalio/orders-reference-app-go/app/encryption"
	"github.com/temporalio/orders-reference-app-go/app/fraudcheck"
	"github.com/temporalio/orders-reference-app-go/app/internal/temporalutil"
	"github.com/temporalio/orders-reference-app-go/app/order"
	"github.com/temporalio/orders-reference-app-go/app/shipment"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
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
//	CLIENT_ENCRYPTION_KEY_ID: If set, Client encrypts payloads using key referenced by ID
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

	// The CLIENT_ENCRYPTION_KEY_ID environment variable contains a value
	// that can be used to look up an encryption key (e.g., from a key
	// management system). If this environment variable is set, then inputs
	// to Workflows and Activities, as well as the outputs returned by the
	// Workflows and Activities, will be encrypted with that key before
	// being transmitted by the Client in this application. This Client
	// will likewise decrypt them upon receipt. The Temporal CLI and Web
	// UI will be unable to view the original (unencrypted) data unless
	// you run a Codec server and configure them to use it.
	if encKeyID := os.Getenv("CLIENT_ENCRYPTION_KEY_ID"); encKeyID != "" {
		clientOpts.DataConverter = encryption.NewEncryptionDataConverter(
			converter.GetDefaultDataConverter(),
			encryption.DataConverterOptions{EncryptionKeyID: encKeyID},
		)
	}

	return clientOpts, nil
}

// RunServer runs all the workers and API servers for the Order/Shipment/Fraud/Billing system.
func RunServer(ctx context.Context, client client.Client) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return billing.RunServer(ctx, 8082, client)
	})
	g.Go(func() error {
		return order.RunServer(ctx, 8083, client)
	})
	g.Go(func() error {
		return shipment.RunServer(ctx, 8081, client)
	})
	g.Go(func() error {
		return fraudcheck.RunServer(ctx, 8084)
	})

	g.Go(func() error {
		return billing.RunWorker(ctx, client, billing.Config{FraudCheckURL: "http://localhost:8084"})
	})
	g.Go(func() error {
		return shipment.RunWorker(ctx, client)
	})
	g.Go(func() error {
		return order.RunWorker(ctx, client, order.Config{BillingURL: "http://localhost:8082"})
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
