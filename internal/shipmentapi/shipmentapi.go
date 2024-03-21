package shipmentapi

import (
	"fmt"

	"github.com/temporalio/orders-reference-app-go/pkg/ordersapi"
)

func ShipmentWorkflowID(orderID string) string {
	return fmt.Sprintf("shipment:%s", orderID)
}

type ShipmentInput struct {
	OrderID string
	Items   []ordersapi.Item
}
