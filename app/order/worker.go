package order

import (
	"context"

	"github.com/temporalio/reference-app-orders-go/app/config"
	"github.com/temporalio/reference-app-orders-go/app/temporalutil"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

// RunWorker runs a Workflow and Activity worker for the Order system.
func RunWorker(ctx context.Context, config config.AppConfig, client client.Client) error {
	w := worker.New(client, TaskQueue, worker.Options{
		MaxConcurrentWorkflowTaskPollers: 8,
		MaxConcurrentActivityTaskPollers: 8,
	})

	w.RegisterWorkflow(Order)
	w.RegisterActivity(&Activities{BillingURL: config.BillingURL, OrderURL: config.OrderURL})

	return w.Run(temporalutil.WorkerInterruptFromContext(ctx))
}
