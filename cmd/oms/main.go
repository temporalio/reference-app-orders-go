package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"

	"github.com/spf13/cobra"
	"github.com/temporalio/reference-app-orders-go/app/config"
	"github.com/temporalio/reference-app-orders-go/app/server"
	"github.com/temporalio/reference-app-orders-go/app/util"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
	_ "modernc.org/sqlite"
)

const (
	defaultCodecPort = 8089
)

var (
	codecPort       int
	codecCorsURL    string
	encryptionKeyID string
	workers         []string
	apis            []string
)

var rootCmd = &cobra.Command{
	Use:   "oms",
	Short: "Order Management System",
	PersistentPreRun: func(*cobra.Command, []string) {
		level := slog.LevelInfo
		if os.Getenv("DEBUG") != "" {
			level = slog.LevelDebug
		}
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		}))
		slog.SetDefault(logger)
	},
}

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Workers for Order/Shipment/Billing system",
	RunE: func(*cobra.Command, []string) error {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)

		config, err := config.AppConfigFromEnv()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		clientOptions, err := server.CreateClientOptionsFromEnv()
		if err != nil {
			return fmt.Errorf("failed to create client options: %w", err)
		}

		if encryptionKeyID != "" {
			log.Printf("Enabling encrypting Data Converter using key ID '%s'", encryptionKeyID)
			ddc := converter.GetDefaultDataConverter()
			clientOptions.DataConverter = util.NewEncryptionDataConverter(ddc, encryptionKeyID)
		}

		client, err := client.Dial(clientOptions)
		if err != nil {
			return fmt.Errorf("client error: %w", err)
		}
		defer client.Close()

		go func() {
			<-sigCh
			log.Printf("Interrupt signal received, shutting down...")
			cancel()
		}()

		if err := server.RunWorkers(ctx, config, client, workers); err != nil {
			return err
		}

		return nil
	},
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "API Servers for Order/Shipment/Billing system",
	RunE: func(*cobra.Command, []string) error {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)

		config, err := config.AppConfigFromEnv()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		clientOptions, err := server.CreateClientOptionsFromEnv()
		if err != nil {
			return fmt.Errorf("failed to create client options: %w", err)
		}

		if encryptionKeyID != "" {
			log.Printf("Enabling encrypting Data Converter using key ID '%s'", encryptionKeyID)
			ddc := converter.GetDefaultDataConverter()
			clientOptions.DataConverter = util.NewEncryptionDataConverter(ddc, encryptionKeyID)
		}

		client, err := client.Dial(clientOptions)
		if err != nil {
			return fmt.Errorf("client error: %w", err)
		}
		defer client.Close()

		go func() {
			<-sigCh
			log.Printf("Interrupt signal received, shutting down...")
			cancel()
		}()

		if err := server.RunAPIServers(ctx, config, client, apis); err != nil {
			return err
		}

		return nil
	},
}

var codecCmd = &cobra.Command{
	Use:   "codec-server",
	Short: "Codec Server decrypts payloads for display by Temporal CLI and Web UI",
	RunE: func(*cobra.Command, []string) error {
		if codecCorsURL != "" {
			log.Printf("Codec Server will allow requests from Temporal Web UI at: %s\n", codecCorsURL)

			if strings.HasSuffix(codecCorsURL, "/") {
				// In my experience, a slash character at the end of the URL will
				// result in a "Codec server could not connect" in the Web UI and
				// the cause will not be obvious. I don't want to strip it off, in
				// case there really is a valid reason to have one, but warning the
				// user could help them to more quickly spot the problem otherwise.
				log.Println("Warning: Temporal Web UI base URL ends with '/'")
			}
		}

		log.Printf("Starting Codec Server on port %d\n", codecPort)
		err := util.RunCodecServer(codecPort, codecCorsURL)

		return err
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
	workerCmd.PersistentFlags().StringVarP(&encryptionKeyID, "encryption-key-id", "k", "",
		"ID of key used to encrypt payload data (optional)")
	apiCmd.PersistentFlags().StringVarP(&encryptionKeyID, "encryption-key-id", "k", "",
		"ID of key used to encrypt payload data (optional)")

	workerCmd.PersistentFlags().StringSliceVarP(&workers, "services", "s", []string{"order", "shipment", "billing"}, "Workers to run")
	apiCmd.PersistentFlags().StringSliceVarP(&apis, "services", "s", []string{"order", "shipment", "billing", "fraud"}, "API Servers to run")

	codecCmd.PersistentFlags().IntVarP(&codecPort, "port", "p", defaultCodecPort,
		"Port number on which the Codec Server will listen for requests")
	codecCmd.PersistentFlags().StringVarP(&codecCorsURL, "url", "u", "",
		"Temporal Web UI base URL (allow CORS for that origin)")

	rootCmd.AddCommand(workerCmd)
	rootCmd.AddCommand(apiCmd)
	rootCmd.AddCommand(codecCmd)
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
