package shipment_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/temporalio/reference-app-orders-go/app/shipment"
	"go.temporal.io/sdk/testsuite"
)

func TestShipmentWorkflow(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	a := &shipment.Activities{}

	shipmentInput := shipment.ShipmentInput{
		RequestorWID: "parentwid",
		ID:           "test",
		Items: []shipment.Item{
			{SKU: "test1", Quantity: 1},
			{SKU: "test2", Quantity: 3},
		},
	}

	env.RegisterActivity(a.BookShipment)

	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow(
			shipment.ShipmentCarrierUpdateSignalName,
			shipment.ShipmentCarrierUpdateSignal{
				Status: shipment.ShipmentStatusDispatched,
			},
		)
	}, time.Second*1)

	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow(
			shipment.ShipmentCarrierUpdateSignalName,
			shipment.ShipmentCarrierUpdateSignal{
				Status: shipment.ShipmentStatusDelivered,
			},
		)
	}, time.Second*2)

	env.OnSignalExternalWorkflow(mock.Anything,
		"parentwid", "",
		shipment.ShipmentStatusUpdatedSignalName,
		mock.MatchedBy(func(arg shipment.ShipmentStatusUpdatedSignal) bool {
			return arg.Status == shipment.ShipmentStatusBooked
		}),
	).Return(nil).Once()

	env.OnSignalExternalWorkflow(mock.Anything,
		"parentwid", "",
		shipment.ShipmentStatusUpdatedSignalName,
		mock.MatchedBy(func(arg shipment.ShipmentStatusUpdatedSignal) bool {
			return arg.Status == shipment.ShipmentStatusDispatched
		}),
	).Return(nil).Once()

	env.OnSignalExternalWorkflow(mock.Anything,
		"parentwid", "",
		shipment.ShipmentStatusUpdatedSignalName,
		mock.MatchedBy(func(arg shipment.ShipmentStatusUpdatedSignal) bool {
			return arg.Status == shipment.ShipmentStatusDelivered
		}),
	).Return(nil).Once()

	env.ExecuteWorkflow(
		shipment.Shipment,
		&shipmentInput,
	)

	var result shipment.ShipmentResult
	err := env.GetWorkflowResult(&result)
	assert.NoError(t, err)
}

func TestShipmentWorkflowIgnoresUnknownCarrierStatus(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	a := &shipment.Activities{}

	shipmentInput := shipment.ShipmentInput{
		RequestorWID: "parentwid",
		ID:           "test",
		Items: []shipment.Item{
			{SKU: "test1", Quantity: 1},
		},
	}

	env.RegisterActivity(a.BookShipment)

	// A carrier update with a status outside the known set must be ignored:
	// it must neither become the shipment's status nor be forwarded to the
	// requestor.
	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow(
			shipment.ShipmentCarrierUpdateSignalName,
			shipment.ShipmentCarrierUpdateSignal{
				Status: "lost-in-transit",
			},
		)
	}, time.Second*1)

	env.RegisterDelayedCallback(func() {
		v, err := env.QueryWorkflow(shipment.StatusQuery)
		assert.NoError(t, err)
		var status shipment.ShipmentStatus
		assert.NoError(t, v.Get(&status))
		assert.Equal(t, shipment.ShipmentStatusBooked, status.Status)
	}, time.Second*2)

	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow(
			shipment.ShipmentCarrierUpdateSignalName,
			shipment.ShipmentCarrierUpdateSignal{
				Status: shipment.ShipmentStatusDelivered,
			},
		)
	}, time.Second*3)

	env.OnSignalExternalWorkflow(mock.Anything,
		"parentwid", "",
		shipment.ShipmentStatusUpdatedSignalName,
		mock.MatchedBy(func(arg shipment.ShipmentStatusUpdatedSignal) bool {
			return arg.Status == shipment.ShipmentStatusBooked
		}),
	).Return(nil).Once()

	env.OnSignalExternalWorkflow(mock.Anything,
		"parentwid", "",
		shipment.ShipmentStatusUpdatedSignalName,
		mock.MatchedBy(func(arg shipment.ShipmentStatusUpdatedSignal) bool {
			return arg.Status == shipment.ShipmentStatusDelivered
		}),
	).Return(nil).Once()

	env.ExecuteWorkflow(
		shipment.Shipment,
		&shipmentInput,
	)

	var result shipment.ShipmentResult
	err := env.GetWorkflowResult(&result)
	assert.NoError(t, err)
	env.AssertExpectations(t)
}
