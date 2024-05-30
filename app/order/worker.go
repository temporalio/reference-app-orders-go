package order

import (
	"context"

	"github.com/temporalio/orders-reference-app-go/app/config"
	"github.com/temporalio/orders-reference-app-go/app/internal/temporalutil"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

// RunWorker runs a Workflow and Activity worker for the Order system.
func RunWorker(ctx context.Context, config config.AppConfig, client client.Client) error {
	w := worker.New(client, TaskQueue, worker.Options{})

	w.RegisterWorkflow(Order)
	w.RegisterActivity(&Activities{BillingURL: config.BillingURL})

	return w.Run(temporalutil.WorkerInterruptFromContext(ctx))
}
