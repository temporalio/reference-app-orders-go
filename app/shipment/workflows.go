package shipment

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// Item represents an item being ordered.
type Item struct {
	SKU      string
	Quantity int32
}

// ShipmentInput is the input for a Shipment workflow.
type ShipmentInput struct {
	OrderID string
	Items   []Item
}

// Custom Search Attribute indicating current status of the shipment
var shipmentStatusAttr = temporal.NewSearchAttributeKeyKeyword("ShipmentStatus")

// ShipmentCarrierUpdateSignalName is the name for a signal to update a shipment's status from the carrier.
const ShipmentCarrierUpdateSignalName = "ShipmentCarrierUpdate"

// ShipmentStatusUpdatedSignalName is the name for a signal to notify of an update to a shipment's status.
const ShipmentStatusUpdatedSignalName = "ShipmentStatusUpdated"

const (
	// ShipmentStatusPending represents a shipment that has not yet been booked with a carrier
	ShipmentStatusPending = "pending"
	// ShipmentStatusBooked represents a shipment acknowledged by a carrier, but not yet picked up
	ShipmentStatusBooked = "booked"
	// ShipmentStatusDispatched represents a shipment picked up by a carrier, but not yet delivered to the customer
	ShipmentStatusDispatched = "dispatched"
	// ShipmentStatusDelivered represents a shipment that has been delivered to the customer
	ShipmentStatusDelivered = "delivered"
)

// ShipmentCarrierUpdateSignal is used by a carrier to update a shipment's status.
type ShipmentCarrierUpdateSignal struct {
	Status string `json:"status"`
}

// ShipmentStatusUpdatedSignal is used to notify an order of an update to a shipment's status.
type ShipmentStatusUpdatedSignal struct {
	ShipmentID string `json:"shipmentID"`
	Status     string `json:"status"`
}

// ShipmentResult is the result of a Shipment workflow.
type ShipmentResult struct {
	CourierReference string
}

type shipmentImpl struct {
	id      string
	orderID string
	status  string
}

// Shipment implements the Shipment workflow.
func Shipment(ctx workflow.Context, input *ShipmentInput) (*ShipmentResult, error) {
	wf := new(shipmentImpl)

	if err := wf.setup(ctx, input); err != nil {
		return nil, err
	}

	return wf.run(ctx, input)
}

func (s *shipmentImpl) setup(ctx workflow.Context, input *ShipmentInput) error {
	s.id = workflow.GetInfo(ctx).WorkflowExecution.ID
	s.orderID = input.OrderID

	s.status = ShipmentStatusPending
	err := workflow.UpsertTypedSearchAttributes(ctx, shipmentStatusAttr.ValueSet(ShipmentStatusPending))
	if err != nil {
		return err
	}

	return nil
}

func (s *shipmentImpl) run(ctx workflow.Context, input *ShipmentInput) (*ShipmentResult, error) {
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
		return nil, err
	}

	s.updateStatus(ctx, ShipmentStatusBooked)

	err = s.handleCarrierUpdates(ctx)

	return &ShipmentResult{}, err
}

func (s *shipmentImpl) handleCarrierUpdates(ctx workflow.Context) error {
	var signal ShipmentCarrierUpdateSignal

	ch := workflow.GetSignalChannel(ctx, ShipmentCarrierUpdateSignalName)
	for s.status != ShipmentStatusDelivered {
		ch.Receive(ctx, &signal)
		s.updateStatus(ctx, signal.Status)
	}

	return nil
}

func (s *shipmentImpl) updateStatus(ctx workflow.Context, status string) error {
	s.status = status
	if err := s.notifyOrderOfStatus(ctx); err != nil {
		return fmt.Errorf("failed to notify order of status: %w", err)
	}
	if err := s.notifyCustomerOfStatus(ctx); err != nil {
		workflow.GetLogger(ctx).Error("failed to notify order of status", "error", err)
	}

	err := workflow.UpsertTypedSearchAttributes(ctx, shipmentStatusAttr.ValueSet(status))
	if err != nil {
		return err
	}

	return nil
}

func (s *shipmentImpl) notifyOrderOfStatus(ctx workflow.Context) error {
	return workflow.SignalExternalWorkflow(ctx,
		s.orderID, "",
		ShipmentStatusUpdatedSignalName,
		ShipmentStatusUpdatedSignal{
			ShipmentID: s.id,
			Status:     s.status,
		},
	).Get(ctx, nil)
}

func (s *shipmentImpl) notifyCustomerOfStatus(ctx workflow.Context) error {
	ctx = workflow.WithActivityOptions(ctx,
		workflow.ActivityOptions{
			StartToCloseTimeout: 5 * time.Second,
		},
	)

	switch s.status {
	case ShipmentStatusBooked:
		return workflow.ExecuteActivity(ctx,
			a.ShipmentBookedNotification,
			ShipmentBookedNotificationInput{
				OrderID: s.orderID,
			},
		).Get(ctx, nil)
	case ShipmentStatusDispatched:
		return workflow.ExecuteActivity(ctx,
			a.ShipmentDispatchedNotification,
			ShipmentDispatchedNotificationInput{
				OrderID: s.orderID,
			},
		).Get(ctx, nil)
	case ShipmentStatusDelivered:
		return workflow.ExecuteActivity(ctx,
			a.ShipmentDeliveredNotification,
			ShipmentDeliveredNotificationInput{
				OrderID: s.orderID,
			},
		).Get(ctx, nil)
	}

	return nil
}

func (s *shipmentImpl) waitForStatus(ctx workflow.Context, status string) {
	workflow.Await(ctx, func() bool {
		return s.status == status
	})
}
