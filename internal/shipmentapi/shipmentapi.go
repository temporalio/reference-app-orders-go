package shipmentapi

import (
	"fmt"

	"github.com/temporalio/orders-reference-app-go/pkg/ordersapi"
)

// The Shipment interfaces are only used internally.

// Signals will be sent by the courier via a local API service,
// so we don't need to expose these.

func ShipmentWorkflowID(orderID ordersapi.OrderID) string {
	return fmt.Sprintf("shipment:%s", orderID)
}

type ShipmentInput struct {
	OrderID ordersapi.OrderID
	Items   []ordersapi.Item
}

const ShipmentUpdateSignalName = "ShipmentUpdate"

type ShipmentStatus int

const (
	ShipmentStatusRegistered ShipmentStatus = iota
	ShipmentStatusDispatched
	ShipmentStatusDelivered
)

type ShipmentUpdateSignal struct {
	Status ShipmentStatus
}

type ShipmentResult struct{}
