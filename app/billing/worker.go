package billing

import (
	"context"

	"github.com/temporalio/orders-reference-app-go/app/config"
	"github.com/temporalio/orders-reference-app-go/app/util"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

// RunWorker runs a Workflow and Activity worker for the Billing system.
func RunWorker(ctx context.Context, config config.AppConfig, client client.Client) error {
	w := worker.New(client, TaskQueue, worker.Options{})

	w.RegisterWorkflow(Charge)
	w.RegisterActivity(&Activities{FraudCheckURL: config.FraudURL})

	return w.Run(util.WorkerInterruptFromContext(ctx))
}
