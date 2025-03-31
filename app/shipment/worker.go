package shipment

import (
	"context"

	"github.com/temporalio/reference-app-orders-go/app/config"
	"github.com/temporalio/reference-app-orders-go/app/temporalutil"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

// Config is the configuration for the Order system.
type Config struct {
	ShipmentURL string
}

// RunWorker runs a Workflow and Activity worker for the Shipment system.
func RunWorker(ctx context.Context, config config.AppConfig, client client.Client) error {
	w := worker.New(client, TaskQueue, worker.Options{
		MaxConcurrentWorkflowTaskPollers: 8,
		MaxConcurrentActivityTaskPollers: 8,
	})

	w.RegisterWorkflow(Shipment)
	w.RegisterActivity(&Activities{ShipmentURL: config.ShipmentURL})

	return w.Run(temporalutil.WorkerInterruptFromContext(ctx))
}
