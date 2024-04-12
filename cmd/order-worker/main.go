package main

import (
	"github.com/spf13/cobra"
	"github.com/temporalio/orders-reference-app-go/app/order"
)

var rootCmd = &cobra.Command{
	Use:   "order-worker",
	Short: "Worker for Order system",
	RunE: func(cmd *cobra.Command, args []string) error {
		return order.Worker()
	},
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
