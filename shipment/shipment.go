package shipment

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

// Item represents an item being ordered.
// All fields are required.
type Item struct {
	SKU      string
	Quantity int32
}

// ShipmentInput is the input for a Shipment workflow.
// All fields are required.
type ShipmentInput struct {
	OrderID string
	Items   []Item
}

// ShipmentUpdateSignalName is the name for a signal to update a shipment's status.
const ShipmentUpdateSignalName = "ShipmentUpdate"

// ShipmentStatus holds a shipment's status.
type ShipmentStatus int

const (
	// ShipmentStatusRegistered represents a shipment which has been registered but not dispatched.
	ShipmentStatusRegistered ShipmentStatus = iota
	// ShipmentStatusDispatched represents a shipment which has been dispatched but not delivered.
	ShipmentStatusDispatched
	// ShipmentStatusDelivered represents a shipment which has been delivered.
	ShipmentStatusDelivered
)

// ShipmentUpdateSignal is used by couriers to update a shipment's status.
type ShipmentUpdateSignal struct {
	Status ShipmentStatus
}

// ShipmentResult is the result of a Shipment workflow.
type ShipmentResult struct {
	CourierReference string
}

type shipmentImpl struct {
	status ShipmentStatus
}

// Shipment implements the Shipment workflow.
func Shipment(ctx workflow.Context, input ShipmentInput) (ShipmentResult, error) {
	return new(shipmentImpl).run(ctx, input)
}

func (s *shipmentImpl) run(ctx workflow.Context, input ShipmentInput) (ShipmentResult, error) {
	workflow.Go(ctx, s.statusUpdater)

	var result ShipmentResult

	ctx = workflow.WithActivityOptions(ctx,
		workflow.ActivityOptions{
			StartToCloseTimeout: 5 * time.Second,
		},
	)

	err := workflow.ExecuteActivity(ctx,
		a.RegisterShipment,
		RegisterShipmentInput{
			OrderID: input.OrderID,
			Items:   input.Items,
		},
	).Get(ctx, nil)
	if err != nil {
		return result, err
	}

	err = workflow.ExecuteActivity(ctx,
		a.ShipmentCreatedNotification,
		ShipmentCreatedNotificationInput{
			OrderID: input.OrderID,
		},
	).Get(ctx, nil)
	if err != nil {
		return result, err
	}

	s.waitForStatus(ctx, ShipmentStatusDispatched)

	err = workflow.ExecuteActivity(ctx,
		a.ShipmentDispatchedNotification,
		ShipmentDispatchedNotificationInput{
			OrderID: input.OrderID,
		},
	).Get(ctx, nil)
	if err != nil {
		return result, err
	}

	s.waitForStatus(ctx, ShipmentStatusDelivered)

	err = workflow.ExecuteActivity(ctx,
		a.ShipmentDeliveredNotification,
		ShipmentDeliveredNotificationInput{
			OrderID: input.OrderID,
		},
	).Get(ctx, nil)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (s *shipmentImpl) statusUpdater(ctx workflow.Context) {
	var signal ShipmentUpdateSignal

	ch := workflow.GetSignalChannel(ctx, ShipmentUpdateSignalName)
	for {
		ch.Receive(ctx, &signal)
		s.status = signal.Status
	}
}

func (s *shipmentImpl) waitForStatus(ctx workflow.Context, status ShipmentStatus) {
	workflow.Await(ctx, func() bool {
		return s.status == status
	})
}
