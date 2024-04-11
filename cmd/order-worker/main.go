package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/temporalio/orders-reference-app-go/order"
	"github.com/temporalio/orders-reference-app-go/shipment"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var rootCmd = &cobra.Command{
	Use:   "order-worker",
	Short: "Worker for Order system",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.Dial(client.Options{})
		if err != nil {
			log.Fatalf("client error: %v", err)
		}
		defer c.Close()

		w := worker.New(c, "order", worker.Options{})

		w.RegisterWorkflow(order.Order)
		w.RegisterActivity(&order.Activities{})
		w.RegisterWorkflow(shipment.Shipment)
		w.RegisterActivity(&shipment.Activities{SMTPStub: true})

		err = w.Run(worker.InterruptCh())
		if err != nil {
			log.Fatalf("worker exited: %v", err)
		}
	},
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
