package workflows_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/temporalio/orders-reference-app-go/activities"
	"github.com/temporalio/orders-reference-app-go/internal/shipmentapi"
	"github.com/temporalio/orders-reference-app-go/pkg/ordersapi"
	"github.com/temporalio/orders-reference-app-go/workflows"

	"go.temporal.io/sdk/testsuite"
)

func TestShipmentWorkflow(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	var a *activities.Activities

	shipmentInput := shipmentapi.ShipmentInput{
		OrderID: "test",
		Items: []ordersapi.Item{
			{SKU: "test1", Quantity: 1},
			{SKU: "test2", Quantity: 3},
		},
	}

	env.RegisterActivity(a.RegisterShipment)

	env.OnActivity(a.ShipmentCreatedNotification, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, input activities.ShipmentCreatedNotificationInput) (activities.ShipmentCreatedNotificationResult, error) {
			env.SignalWorkflow(
				shipmentapi.ShipmentUpdateSignalName,
				shipmentapi.ShipmentUpdateSignal{
					Status: shipmentapi.ShipmentStatusDispatched,
				},
			)

			return activities.ShipmentCreatedNotificationResult{}, nil
		},
	)

	env.OnActivity(a.ShipmentDispatchedNotification, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, input activities.ShipmentDispatchedNotificationInput) (activities.ShipmentDispatchedNotificationResult, error) {
			env.SignalWorkflow(
				shipmentapi.ShipmentUpdateSignalName,
				shipmentapi.ShipmentUpdateSignal{
					Status: shipmentapi.ShipmentStatusDelivered,
				},
			)

			return activities.ShipmentDispatchedNotificationResult{}, nil
		},
	)

	env.RegisterActivity(a.ShipmentDeliveredNotification)

	env.ExecuteWorkflow(
		workflows.Shipment,
		shipmentInput,
	)

	var result shipmentapi.ShipmentResult
	err := env.GetWorkflowResult(&result)
	assert.NoError(t, err)
}
