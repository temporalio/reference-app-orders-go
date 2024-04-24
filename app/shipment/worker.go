package shipment

import (
	"context"
	"fmt"

	"github.com/temporalio/orders-reference-app-go/app/internal/temporalutil"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

// RunWorker runs a Workflow and Activity worker for the Shipment system.
func RunWorker(ctx context.Context) error {
	clientOptions, err := temporalutil.CreateClientOptionsFromEnv()
	if err != nil {
		return fmt.Errorf("failed to create client options: %w", err)
	}

	c, err := client.Dial(clientOptions)
	if err != nil {
		return fmt.Errorf("client error: %w", err)
	}
	defer c.Close()

	errCh := make(chan error, 1)

	w := worker.New(c, TaskQueue, worker.Options{
		OnFatalError: func(err error) {
			errCh <- err
		},
	})

	w.RegisterWorkflow(Shipment)
	w.RegisterActivity(&Activities{})

	err = w.Start()
	if err != nil {
		return err
	}

	select {
	case err = <-errCh:
		return err
	case <-ctx.Done():
		w.Stop()
	}

	return nil
}
