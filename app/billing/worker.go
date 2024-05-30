package billing

import (
	"context"

	"github.com/temporalio/orders-reference-app-go/app/temporalutil"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

// Config for the Billing system.
type Config struct {
	FraudCheckURL string
}

// RunWorker runs a Workflow and Activity worker for the Billing system.
func RunWorker(ctx context.Context, client client.Client, config Config) error {
	w := worker.New(client, TaskQueue, worker.Options{})

	w.RegisterWorkflow(Charge)
	w.RegisterActivity(&Activities{FraudCheckURL: config.FraudCheckURL})

	return w.Run(temporalutil.WorkerInterruptFromContext(ctx))
}
