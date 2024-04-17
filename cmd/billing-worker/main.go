package main

import (
	"github.com/spf13/cobra"
	"github.com/temporalio/orders-reference-app-go/app/billing"
	"go.temporal.io/sdk/worker"
)

var rootCmd = &cobra.Command{
	Use:   "billing-worker",
	Short: "Worker for Billing system",
	RunE: func(cmd *cobra.Command, args []string) error {
		return billing.RunWorker(worker.InterruptCh())
	},
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
