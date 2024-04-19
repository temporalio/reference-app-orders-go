package main

import (
	"github.com/spf13/cobra"
	"github.com/temporalio/orders-reference-app-go/app/order"
	"go.temporal.io/sdk/worker"
)

var rootCmd = &cobra.Command{
	Use:   "order-worker",
	Short: "Worker for Order system",
	RunE: func(*cobra.Command, []string) error {
		return order.RunWorker(worker.InterruptCh())
	},
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
