package shipmentapi

import (
	"fmt"

	"github.com/temporalio/orders-reference-app-go/pkg/ordersapi"
)

// The Shipment interfaces are only used internally.

// Signals will be sent by couriers via a local API service,
// so we don't need to expose these.

// ShipmentWorkflowID creates a shipment workflow ID from an order ID.
func ShipmentWorkflowID(orderID ordersapi.OrderID) string {
	return fmt.Sprintf("shipment:%s", orderID)
}

// ShipmentInput is the input for a Shipment workflow.
// All fields are required.
type ShipmentInput struct {
	OrderID ordersapi.OrderID
	Items   []ordersapi.Item
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
