package billing

import (
	"context"

	"github.com/temporalio/orders-reference-app-go/app/internal/temporalutil"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

// RunWorker runs a Workflow and Activity worker for the Billing system.
func RunWorker(ctx context.Context, client client.Client) error {
	w := worker.New(client, TaskQueue, worker.Options{})

	w.RegisterWorkflow(Charge)
	w.RegisterActivity(GenerateInvoice)
	w.RegisterActivity(ChargeCustomer)

	return w.Run(temporalutil.WorkerInterruptFromContext(ctx))
}
