package shipment

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

	w := worker.New(c, TaskQueue, worker.Options{})

	w.RegisterWorkflow(Shipment)
	w.RegisterActivity(&Activities{})

	return w.Run(worker.InterruptCh())
}
