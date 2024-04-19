package shipment_test

import (
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
		OrderID: "test",
		Items: []shipment.Item{
			{SKU: "test1", Quantity: 1},
			{SKU: "test2", Quantity: 3},
		},
	}

	env.RegisterActivity(a.BookShipment)

	env.OnActivity(a.ShipmentBookedNotification, mock.Anything, mock.Anything).Return(
		func(input *shipment.ShipmentBookedNotificationInput) error {
			env.SignalWorkflow(
				shipment.ShipmentUpdateSignalName,
				shipment.ShipmentUpdateSignal{
					Status: shipment.ShipmentStatusDispatched,
				},
			)

			return nil
		},
	)

	env.OnActivity(a.ShipmentDispatchedNotification, mock.Anything, mock.Anything).Return(
		func(input *shipment.ShipmentDispatchedNotificationInput) error {
			env.SignalWorkflow(
				shipment.ShipmentUpdateSignalName,
				shipment.ShipmentUpdateSignal{
					Status: shipment.ShipmentStatusDelivered,
				},
			)

			return nil
		},
	)

	env.RegisterActivity(a.ShipmentDeliveredNotification)

	env.ExecuteWorkflow(
		shipment.Shipment,
		&shipmentInput,
	)

	var result shipment.ShipmentResult
	err := env.GetWorkflowResult(&result)
	assert.NoError(t, err)
}
