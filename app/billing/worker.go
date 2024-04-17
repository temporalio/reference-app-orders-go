package billing

import (
	"log"

	"github.com/temporalio/orders-reference-app-go/app/internal/temporalutil"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func RunWorker(intCh <-chan interface{}) error {
	clientOptions, err := temporalutil.CreateClientOptionsFromEnv()
	if err != nil {
		log.Fatalf("failed to create client options: %v", err)
	}

	c, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalf("client error: %v", err)
	}
	defer c.Close()

	w := worker.New(c, TASK_QUEUE, worker.Options{})

	w.RegisterWorkflow(GenerateInvoice)
	w.RegisterWorkflow(ChargeCustomer)

	return w.Run(intCh)
}
