package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/temporalio/orders-reference-app-go/app/billing"
	"github.com/temporalio/orders-reference-app-go/app/order"
	"github.com/temporalio/orders-reference-app-go/app/shipment"
)

var rootCmd = &cobra.Command{
	Use:   "dev-server",
	Short: "Workers and API Servers for Order/Shipment/Billing system",
	RunE: func(*cobra.Command, []string) error {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)

		doneCh := make(chan interface{})
		errCh := make(chan error, 1)

		go func() { errCh <- billing.RunWorker(doneCh) }()
		go func() { errCh <- shipment.RunWorker(doneCh) }()
		go func() { errCh <- order.RunWorker(doneCh) }()

		go func() { errCh <- billing.RunServer(ctx, 8082) }()
		go func() { errCh <- order.RunServer(ctx, 8083) }()
		go func() { errCh <- shipment.RunServer(ctx, 8081) }()

		select {
		case <-sigCh:
			doneCh <- true
			cancel()
		case err := <-errCh:
			doneCh <- true
			cancel()
			log.Fatalf("Error: %v", err)
		}

		return nil
	},
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
