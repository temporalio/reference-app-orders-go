package activities

import (
	"context"

	"github.com/temporalio/orders-reference-app-go/internal/shipmentapi"
	"github.com/temporalio/orders-reference-app-go/pkg/ordersapi"
)

type ShipmentCreatedNotificationInput struct {
	OrderID ordersapi.OrderID
}
type ShipmentCreatedNotificationResult struct{}

func (a *Activities) ShipmentCreatedNotification(ctx context.Context, input ShipmentCreatedNotificationInput) (ShipmentCreatedNotificationResult, error) {
	// Auto-dispatch for now
	err := a.client.SignalWorkflow(ctx,
		shipmentapi.ShipmentWorkflowID(input.OrderID), "",
		shipmentapi.ShipmentDispatchedSignalName,
		shipmentapi.ShipmentDispatchedSignal{},
	)

	// TODO: Hit Notification API

	return ShipmentCreatedNotificationResult{}, err
}

type ShipmentDispatchedNotificationInput struct {
	OrderID ordersapi.OrderID
}
type ShipmentDispatchedNotificationResult struct{}

func (a *Activities) ShipmentDispatchedNotification(ctx context.Context, input ShipmentDispatchedNotificationInput) (ShipmentDispatchedNotificationResult, error) {
	// Auto-deliver for now
	err := a.client.SignalWorkflow(ctx,
		shipmentapi.ShipmentWorkflowID(input.OrderID), "",
		shipmentapi.ShipmentDeliveredSignalName,
		shipmentapi.ShipmentDeliveredSignal{},
	)

	// TODO: Hit Notification API

	return ShipmentDispatchedNotificationResult{}, err
}

type ShipmentDeliveredNotificationInput struct {
	OrderID ordersapi.OrderID
}
type ShipmentDeliveredNotificationResult struct{}

func (a *Activities) ShipmentDeliveredNotification(ctx context.Context, input ShipmentDeliveredNotificationInput) (ShipmentDeliveredNotificationResult, error) {
	// TODO: Hit Notification API

	return ShipmentDeliveredNotificationResult{}, nil
}
