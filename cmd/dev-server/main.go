package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/spf13/cobra"
	"github.com/temporalio/orders-reference-app-go/app/billing"
	"github.com/temporalio/orders-reference-app-go/app/order"
	"github.com/temporalio/orders-reference-app-go/app/shipment"
)

var rootCmd = &cobra.Command{
	Use:   "dev-server",
	Short: "Workers and API Servers for Order/Shipment/Billing system",
	Run: func(*cobra.Command, []string) {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)

		doneCh := make(chan interface{}, 3)
		errCh := make(chan error, 6)

		var wg sync.WaitGroup

		wg.Add(6)
		go func() {
			defer wg.Done()
			errCh <- billing.RunWorker(doneCh)
		}()
		go func() {
			defer wg.Done()
			errCh <- shipment.RunWorker(doneCh)
		}()
		go func() {
			defer wg.Done()
			errCh <- order.RunWorker(doneCh)
		}()

		go func() {
			defer wg.Done()
			errCh <- billing.RunServer(ctx, 8082)
		}()
		go func() {
			defer wg.Done()
			errCh <- order.RunServer(ctx, 8083)
		}()
		go func() {
			defer wg.Done()
			errCh <- shipment.RunServer(ctx, 8081)
		}()

		select {
		case <-sigCh:
			doneCh <- true
			doneCh <- true
			doneCh <- true
			cancel()
		case err := <-errCh:
			log.Printf("Error: %v\n", err)
			doneCh <- true
			doneCh <- true
			doneCh <- true
			cancel()
		}

		wg.Wait()
	},
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
