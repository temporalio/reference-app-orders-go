package shipmentapi

import (
	"fmt"

	"github.com/temporalio/orders-reference-app-go/pkg/ordersapi"
)

func ShipmentWorkflowID(orderID ordersapi.OrderID) string {
	return fmt.Sprintf("shipment:%s", orderID)
}

type ShipmentInput struct {
	OrderID ordersapi.OrderID
	Items   []ordersapi.Item
}

const ShipmentDispatchedSignalName = "ShipmentDispatched"

type ShipmentDispatchedSignal struct{}

const ShipmentDeliveredSignalName = "ShipmentDelivered"

type ShipmentDeliveredSignal struct{}

type ShipmentResult struct{}
