package shipment

import (
	"time"

	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"
)

const TASK_QUEUE = "shipments"

// Name of the Custom Search Attribute (int) indicating the shipment status
const STATUS_CSA_NAME = "ShipmentStatus"

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
	// Represents a shipment acknowledged by a courier, but not yet picked up
	ShipmentStatusBooked ShipmentStatus = iota
	// Represents a shipment picked up by a courier, but not yet delivered to the customer
	ShipmentStatusDispatched
	// Represents a shipment that has been delivered to the customer
	ShipmentStatusDelivered
)

// ShipmentUpdateSignal is used by a courier to update a shipment's status.
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

var logger log.Logger

// Shipment implements the Shipment workflow.
func Shipment(ctx workflow.Context, input ShipmentInput) (ShipmentResult, error) {
	logger = workflow.GetLogger(ctx)
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
		a.BookShipment,
		BookShipmentInput{
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

	// Set the initial status in the custom search attribute
	s.upsertStatusAttrbute(ctx, s.status)

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

		// update the status in the custom search attribute
		s.upsertStatusAttrbute(ctx, s.status)
	}
}

func (s *shipmentImpl) waitForStatus(ctx workflow.Context, status ShipmentStatus) {
	workflow.Await(ctx, func() bool {
		return s.status == status
	})
}

func (s *shipmentImpl) upsertStatusAttrbute(ctx workflow.Context, status ShipmentStatus) {
	attributes := map[string]interface{}{
		STATUS_CSA_NAME: s.status,
	}

	err := workflow.UpsertSearchAttributes(ctx, attributes)
	if err != nil {
		logger.Error("error upserting shipment status attribute", "Error", err)
	}
}
