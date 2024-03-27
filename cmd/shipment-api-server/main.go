package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/temporalio/orders-reference-app-go/shipment"

	"go.temporal.io/sdk/client"
)

var port int

var rootCmd = &cobra.Command{
	Use:   "shipment-api-server",
	Short: "API Server for Shipments",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.Dial(client.Options{})
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		defer c.Close()

		srv := &http.Server{
			Addr:    fmt.Sprintf("0.0.0.0:%d", port),
			Handler: shipment.Router(c),
		}

		fmt.Printf("Listening on http://0.0.0.0:%d\n", port)

		errCh := make(chan error, 1)
		go func() { errCh <- srv.ListenAndServe() }()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)

		select {
		case <-sigCh:
			srv.Close()
		case err = <-errCh:
			log.Fatalf("error: %v", err)
		}
	},
}

func main() {
	rootCmd.Flags().IntVar(&port, "port", 8081, "Port to listen on")

	cobra.CheckErr(rootCmd.Execute())
}
