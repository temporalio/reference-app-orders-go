package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/temporalio/orders-reference-app-go/app/billing"
)

var port int

var rootCmd = &cobra.Command{
	Use:   "billing-api-server",
	Short: "API Server for Billing",
	Run: func(_ *cobra.Command, _ []string) {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)

		errCh := make(chan error, 1)
		go func() { errCh <- billing.RunServer(ctx, port) }()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)

		select {
		case <-sigCh:
			cancel()
		case err := <-errCh:
			log.Fatalf("Server error: %v", err)
		}
	},
}

func main() {
	rootCmd.Flags().IntVar(&port, "port", 8082, "Port to listen on")

	cobra.CheckErr(rootCmd.Execute())
}
