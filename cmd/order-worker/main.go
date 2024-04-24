package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/temporalio/orders-reference-app-go/app/order"
)

var rootCmd = &cobra.Command{
	Use:   "order-worker",
	Short: "Worker for Order system",
	RunE: func(*cobra.Command, []string) error {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)

		go func() {
			<-sigCh
			log.Printf("Interrupt signal received, shutting down...")
			cancel()
		}()

		return order.RunWorker(ctx)
	},
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
