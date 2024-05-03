package shipment

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// Item represents an item being shipped.
type Item struct {
	SKU      string `json:"sku"`
	Quantity int32  `json:"quantity"`
}

// ShipmentInput is the input for a Shipment workflow.
type ShipmentInput struct {
	RequestorWID string

	ID    string
	Items []Item
}

// ShipmentStatusAttr is a Custom Search Attribute that indicates current status of a shipment
var ShipmentStatusAttr = temporal.NewSearchAttributeKeyKeyword("ShipmentStatus")

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

// ShipmentStatusUpdatedSignal is used to notify the requestor of an update to a shipment's status.
type ShipmentStatusUpdatedSignal struct {
	ShipmentID string    `json:"shipmentID"`
	Status     string    `json:"status"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// ShipmentResult is the result of a Shipment workflow.
type ShipmentResult struct {
	CourierReference string
}

type shipmentImpl struct {
	requestorWID string

	id        string
	status    string
	updatedAt time.Time
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
	s.requestorWID = input.RequestorWID
	s.id = input.ID
	s.status = ShipmentStatusPending

	err := workflow.UpsertTypedSearchAttributes(ctx, ShipmentStatusAttr.ValueSet(ShipmentStatusPending))
	if err != nil {
		return err
	}

	return workflow.SetQueryHandler(ctx, StatusQuery, func() (*ShipmentStatus, error) {
		return &ShipmentStatus{
			ID:        s.id,
			Status:    s.status,
			UpdatedAt: s.updatedAt,
			Items:     input.Items,
		}, nil
	})
}

func (s *shipmentImpl) run(ctx workflow.Context, input *ShipmentInput) (*ShipmentResult, error) {
	ctx = workflow.WithActivityOptions(ctx,
		workflow.ActivityOptions{
			StartToCloseTimeout: 5 * time.Second,
		},
	)

	var result BookShipmentResult

	err := workflow.ExecuteActivity(ctx,
		a.BookShipment,
		BookShipmentInput{
			Reference: s.id,
			Items:     input.Items,
		},
	).Get(ctx, &result)
	if err != nil {
		return nil, err
	}

	s.updateStatus(ctx, ShipmentStatusBooked)

	err = s.handleCarrierUpdates(ctx)

	return &ShipmentResult{
		CourierReference: result.CourierReference,
	}, err
}

func (s *shipmentImpl) handleCarrierUpdates(ctx workflow.Context) error {
	ch := workflow.GetSignalChannel(ctx, ShipmentCarrierUpdateSignalName)

	var signal ShipmentCarrierUpdateSignal

	for s.status != ShipmentStatusDelivered {
		ch.Receive(ctx, &signal)
		s.updateStatus(ctx, signal.Status)
	}

	return nil
}

func (s *shipmentImpl) updateStatus(ctx workflow.Context, status string) error {
	s.status = status
	s.updatedAt = workflow.Now(ctx)

	if err := s.notifyRequestorOfStatus(ctx); err != nil {
		return fmt.Errorf("failed to notify requestor of status: %w", err)
	}

	return workflow.UpsertTypedSearchAttributes(ctx, ShipmentStatusAttr.ValueSet(status))
}

func (s *shipmentImpl) notifyRequestorOfStatus(ctx workflow.Context) error {
	return workflow.SignalExternalWorkflow(ctx,
		s.requestorWID, "",
		ShipmentStatusUpdatedSignalName,
		ShipmentStatusUpdatedSignal{
			ShipmentID: s.id,
			Status:     s.status,
			UpdatedAt:  s.updatedAt,
		},
	).Get(ctx, nil)
}
