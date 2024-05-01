package shipment_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/temporalio/orders-reference-app-go/app/shipment"
	"go.temporal.io/sdk/testsuite"
)

func TestShipmentWorkflow(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	a := &shipment.Activities{}

	shipmentInput := shipment.ShipmentInput{
		OrderID:         "test",
		OrderWorkflowID: "mywfid",
		Items: []shipment.Item{
			{SKU: "test1", Quantity: 1},
			{SKU: "test2", Quantity: 3},
		},
	}

	env.RegisterActivity(a.BookShipment)

	env.OnActivity(a.ShipmentBookedNotification, mock.Anything, mock.Anything).Return(
		func(_ context.Context, input *shipment.ShipmentBookedNotificationInput) error {
			env.SignalWorkflow(
				shipment.ShipmentCarrierUpdateSignalName,
				shipment.ShipmentCarrierUpdateSignal{
					Status: shipment.ShipmentStatusDispatched,
				},
			)

			return nil
		},
	)

	env.OnActivity(a.ShipmentDispatchedNotification, mock.Anything, mock.Anything).Return(
		func(_ context.Context, input *shipment.ShipmentDispatchedNotificationInput) error {
			env.SignalWorkflow(
				shipment.ShipmentCarrierUpdateSignalName,
				shipment.ShipmentCarrierUpdateSignal{
					Status: shipment.ShipmentStatusDelivered,
				},
			)

			return nil
		},
	)

	env.RegisterActivity(a.ShipmentDeliveredNotification)

	env.OnSignalExternalWorkflow(mock.Anything, mock.Anything, mock.Anything, shipment.ShipmentStatusUpdatedSignalName, mock.Anything).Return(nil)

	env.ExecuteWorkflow(
		shipment.Shipment,
		&shipmentInput,
	)

	var result shipment.ShipmentResult
	err := env.GetWorkflowResult(&result)
	assert.NoError(t, err)
}
