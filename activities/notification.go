package activities

import (
	"context"

	"github.com/temporalio/orders-reference-app-go/pkg/ordersapi"
)

// ShipmentCreatedNotificationInput is the input for a ShipmentCreated notification.
type ShipmentCreatedNotificationInput struct {
	OrderID ordersapi.OrderID
}

// ShipmentCreatedNotification sends a ShipmentCreated notification to a user.
func (a *Activities) ShipmentCreatedNotification(ctx context.Context, input ShipmentCreatedNotificationInput) error {
	return nil
}

// ShipmentDispatchedNotificationInput is the input for a ShipmentDispatched notification.
type ShipmentDispatchedNotificationInput struct {
	OrderID ordersapi.OrderID
}

// ShipmentDispatchedNotification sends a ShipmentDispatched notification to a user.
func (a *Activities) ShipmentDispatchedNotification(ctx context.Context, input ShipmentDispatchedNotificationInput) error {
	return nil
}

// ShipmentDeliveredNotificationInput is the input for a ShipmentDelivered notification.
type ShipmentDeliveredNotificationInput struct {
	OrderID ordersapi.OrderID
}

// ShipmentDeliveredNotification sends a ShipmentDelivered notification to a user.
func (a *Activities) ShipmentDeliveredNotification(ctx context.Context, input ShipmentDeliveredNotificationInput) error {
	return nil
}
