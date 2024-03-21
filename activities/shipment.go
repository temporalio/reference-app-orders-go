package activities

import (
	"context"

	"github.com/temporalio/orders-reference-app-go/pkg/ordersapi"
)

type RegisterShipmentInput struct {
	OrderID ordersapi.OrderID
	Items   []ordersapi.Item
}

type RegisterShipmentResult struct{}

func (a *Activities) RegisterShipment(ctx context.Context, input RegisterShipmentInput) (RegisterShipmentResult, error) {
	return RegisterShipmentResult{}, nil
}
