package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/temporalio/orders-reference-app-go/app/encryption"
	"github.com/temporalio/orders-reference-app-go/app/server"
	"github.com/temporalio/orders-reference-app-go/app/shipment"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
)

var encryptionKeyID string

var rootCmd = &cobra.Command{
	Use:   "dev-server",
	Short: "Workers and API Servers for Order/Shipment/Billing system",
	RunE: func(*cobra.Command, []string) error {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)

		clientOptions, err := server.CreateClientOptionsFromEnv()
		if err != nil {
			return fmt.Errorf("failed to create client options: %w", err)
		}

		if encryptionKeyID != "" {
			log.Printf("Enabling encrypting Data Converter using key ID '%s'", encryptionKeyID)
			ddc := converter.GetDefaultDataConverter()
			clientOptions.DataConverter = encryption.NewEncryptionDataConverter(ddc, encryptionKeyID)
		}

		client, err := client.Dial(clientOptions)
		if err != nil {
			return fmt.Errorf("client error: %w", err)
		}
		defer client.Close()

		err = shipment.EnsureValidTemporalEnv(ctx, client, clientOptions)
		if err != nil {
			return fmt.Errorf("environment is not valid for shipment system: %w", err)
		}

		go func() {
			<-sigCh
			log.Printf("Interrupt signal received, shutting down...")
			cancel()
		}()

		if err := server.RunServer(ctx, client); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	// The encryption key ID is a string that can be used to look up an encryption
	// key (e.g., from a key management system). If this option is specified, then
	// inputs to Workflows and Activities, as well as the outputs returned by the
	// Workflows and Activities, will be encrypted with that key before being sent
	// by the Client in this application. This Client will likewise decrypt them
	// upon receipt. The Temporal CLI and Web UI will be unable to view the original
	// (unencrypted) data unless you run a Codec server and configure them to use it.
	rootCmd.PersistentFlags().StringVarP(&encryptionKeyID, "encryption-key-id", "e", "",
		"ID of key used to encrypt payload data (optional)")
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
