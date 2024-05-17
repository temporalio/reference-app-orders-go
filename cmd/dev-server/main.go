package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/temporalio/orders-reference-app-go/app/server"
	"github.com/temporalio/orders-reference-app-go/app/shipment"
	"go.temporal.io/sdk/client"
)

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

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
