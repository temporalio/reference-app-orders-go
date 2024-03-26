package workflows_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/temporalio/orders-reference-app-go/activities"
	"github.com/temporalio/orders-reference-app-go/pkg/ordersapi"
	"github.com/temporalio/orders-reference-app-go/workflows"
	"go.temporal.io/sdk/testsuite"
)

func TestShipmentWorkflow(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	a := &activities.Activities{
		SMTPStub: true,
	}

	shipmentInput := workflows.ShipmentInput{
		OrderID: "test",
		Items: []ordersapi.Item{
			{SKU: "test1", Quantity: 1},
			{SKU: "test2", Quantity: 3},
		},
	}

	env.RegisterActivity(a.RegisterShipment)

	env.OnActivity(a.ShipmentCreatedNotification, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, input activities.ShipmentCreatedNotificationInput) error {
			env.SignalWorkflow(
				workflows.ShipmentUpdateSignalName,
				workflows.ShipmentUpdateSignal{
					Status: workflows.ShipmentStatusDispatched,
				},
			)

			return nil
		},
	)

	env.OnActivity(a.ShipmentDispatchedNotification, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, input activities.ShipmentDispatchedNotificationInput) error {
			env.SignalWorkflow(
				workflows.ShipmentUpdateSignalName,
				workflows.ShipmentUpdateSignal{
					Status: workflows.ShipmentStatusDelivered,
				},
			)

			return nil
		},
	)

	env.RegisterActivity(a.ShipmentDeliveredNotification)

	env.ExecuteWorkflow(
		workflows.Shipment,
		shipmentInput,
	)

	var result workflows.ShipmentResult
	err := env.GetWorkflowResult(&result)
	assert.NoError(t, err)
}
