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
	"golang.org/x/sync/errgroup"
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

		g, ctx := errgroup.WithContext(ctx)

		g.Go(func() error {
			return billing.RunWorker(ctx)
		})
		g.Go(func() error {
			return shipment.RunWorker(ctx)
		})
		g.Go(func() error {
			return order.RunWorker(ctx)
		})
		g.Go(func() error {
			return billing.RunServer(ctx, 8082)
		})
		g.Go(func() error {
			return order.RunServer(ctx, 8083)
		})
		g.Go(func() error {
			return shipment.RunServer(ctx, 8081)
		})

		go func() {
			<-sigCh
			log.Printf("Interrupt signal received, shutting down...")
			cancel()
		}()

		if err := g.Wait(); err != nil {
			return err
		}

		return nil
	},
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
