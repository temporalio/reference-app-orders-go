package activities

import (
	"context"

	"github.com/temporalio/orders-reference-app-go/pkg/ordersapi"
)

// RegisterShipmentInput is the input for the RegisterShipment operation.
// All fields are required.
type RegisterShipmentInput struct {
	OrderID ordersapi.OrderID
	Items   []ordersapi.Item
}

// RegisterShipmentResult is the result for the RegisterShipment operation.
// CourierReference is recorded where available, to allow tracking enquiries.
type RegisterShipmentResult struct {
	CourierReference string
}

// RegisterShipment registers a shipment with a courier.
func (a *Activities) RegisterShipment(ctx context.Context, input RegisterShipmentInput) (RegisterShipmentResult, error) {
	return RegisterShipmentResult{}, nil
}
