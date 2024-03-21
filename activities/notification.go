package activities

import (
	"context"

	"github.com/temporalio/orders-reference-app-go/pkg/ordersapi"
)

type ShipmentCreatedNotificationInput struct {
	OrderID ordersapi.OrderID
}
type ShipmentCreatedNotificationResult struct{}

func (a *Activities) ShipmentCreatedNotification(ctx context.Context, input ShipmentCreatedNotificationInput) (ShipmentCreatedNotificationResult, error) {
	return ShipmentCreatedNotificationResult{}, nil
}

type ShipmentDispatchedNotificationInput struct {
	OrderID ordersapi.OrderID
}
type ShipmentDispatchedNotificationResult struct{}

func (a *Activities) ShipmentDispatchedNotification(ctx context.Context, input ShipmentDispatchedNotificationInput) (ShipmentDispatchedNotificationResult, error) {
	return ShipmentDispatchedNotificationResult{}, nil
}

type ShipmentDeliveredNotificationInput struct {
	OrderID ordersapi.OrderID
}
type ShipmentDeliveredNotificationResult struct{}

func (a *Activities) ShipmentDeliveredNotification(ctx context.Context, input ShipmentDeliveredNotificationInput) (ShipmentDeliveredNotificationResult, error) {
	return ShipmentDeliveredNotificationResult{}, nil
}
