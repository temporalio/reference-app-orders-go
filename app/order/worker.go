package order

import (
	"context"

	"github.com/temporalio/orders-reference-app-go/app/internal/temporalutil"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

// Config is the configuration for the Order system.
type Config struct {
	BillingURL string
	OrderURL   string
}

// RunWorker runs a Workflow and Activity worker for the Order system.
func RunWorker(ctx context.Context, client client.Client, config Config) error {
	w := worker.New(client, TaskQueue, worker.Options{})

	w.RegisterWorkflow(Order)
	w.RegisterActivity(&Activities{BillingURL: config.BillingURL, OrderURL: config.OrderURL})

	return w.Run(temporalutil.WorkerInterruptFromContext(ctx))
}
