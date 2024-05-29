package shipment

import (
	"context"

	"github.com/temporalio/orders-reference-app-go/app/internal/temporalutil"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

// Config is the configuration for the Order system.
type Config struct {
	ShipmentURL string
}

// RunWorker runs a Workflow and Activity worker for the Shipment system.
func RunWorker(ctx context.Context, client client.Client, config Config) error {
	w := worker.New(client, TaskQueue, worker.Options{})

	w.RegisterWorkflow(Shipment)
	w.RegisterActivity(&Activities{ShipmentURL: config.ShipmentURL})

	return w.Run(temporalutil.WorkerInterruptFromContext(ctx))
}
