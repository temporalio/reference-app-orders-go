package activities

import (
	"context"

	"github.com/temporalio/orders-reference-app-go/pkg/ordersapi"
)

type CreateShipmentInput struct {
	OrderID ordersapi.OrderID
	Items   []ordersapi.Item
}

type CreateShipmentResult struct{}

func (a *Activities) CreateShipment(ctx context.Context, input CreateShipmentInput) (CreateShipmentResult, error) {
	// TODO: Hit shipment API

	return CreateShipmentResult{}, nil
}
