package main

import (
	"github.com/spf13/cobra"
	"github.com/temporalio/orders-reference-app-go/app/shipment"
)

var port int

var rootCmd = &cobra.Command{
	Use:   "shipment-api-server",
	Short: "API Server for Shipments",
	RunE: func(cmd *cobra.Command, args []string) error {
		return shipment.Server(port)
	},
}

func main() {
	rootCmd.Flags().IntVar(&port, "port", 8081, "Port to listen on")

	cobra.CheckErr(rootCmd.Execute())
}
