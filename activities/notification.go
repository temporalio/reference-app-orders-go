package activities

import (
	"context"

	"github.com/temporalio/orders-reference-app-go/pkg/ordersapi"
)

type ShipmentCreatedNotificationInput struct {
	OrderID ordersapi.OrderID
}

func (a *Activities) ShipmentCreatedNotification(ctx context.Context, input ShipmentCreatedNotificationInput) error {
	return nil
}

type ShipmentDispatchedNotificationInput struct {
	OrderID ordersapi.OrderID
}

func (a *Activities) ShipmentDispatchedNotification(ctx context.Context, input ShipmentDispatchedNotificationInput) error {
	return nil
}

type ShipmentDeliveredNotificationInput struct {
	OrderID ordersapi.OrderID
}

func (a *Activities) ShipmentDeliveredNotification(ctx context.Context, input ShipmentDeliveredNotificationInput) error {
	return nil
}
