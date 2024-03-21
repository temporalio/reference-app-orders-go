package workflows

import (
	"go.temporal.io/sdk/workflow"

	"github.com/temporalio/orders-reference-app-go/activities"
	"github.com/temporalio/orders-reference-app-go/internal/shipmentapi"
)

func Shipment(ctx workflow.Context, input shipmentapi.ShipmentInput) (shipmentapi.ShipmentResult, error) {
	var result shipmentapi.ShipmentResult

	err := workflow.ExecuteActivity(ctx,
		a.CreateShipment,
		activities.CreateShipmentInput{
			OrderID: input.OrderID,
			Items:   input.Items,
		},
	).Get(ctx, nil)
	if err != nil {
		return result, err
	}

	err = workflow.ExecuteActivity(ctx,
		a.ShipmentCreatedNotification,
		activities.ShipmentCreatedNotificationInput{
			OrderID: input.OrderID,
		},
	).Get(ctx, nil)
	if err != nil {
		return result, err
	}

	workflow.GetSignalChannel(ctx,
		shipmentapi.ShipmentDispatchedSignalName,
	).Receive(ctx, nil)

	err = workflow.ExecuteActivity(ctx,
		a.ShipmentDispatchedNotification,
		activities.ShipmentDispatchedNotificationInput{
			OrderID: input.OrderID,
		},
	).Get(ctx, nil)
	if err != nil {
		return result, err
	}

	workflow.GetSignalChannel(ctx,
		shipmentapi.ShipmentDeliveredSignalName,
	).Receive(ctx, nil)

	err = workflow.ExecuteActivity(ctx,
		a.ShipmentDeliveredNotification,
		activities.ShipmentDeliveredNotificationInput{
			OrderID: input.OrderID,
		},
	).Get(ctx, nil)
	if err != nil {
		return result, err
	}

	return result, nil
}
