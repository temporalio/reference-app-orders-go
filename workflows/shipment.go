package workflows

import (
	"time"

	"github.com/temporalio/orders-reference-app-go/activities"
	"github.com/temporalio/orders-reference-app-go/internal/shipmentapi"
	"go.temporal.io/sdk/workflow"
)

func Shipment(ctx workflow.Context, input shipmentapi.ShipmentInput) (shipmentapi.ShipmentResult, error) {
	var result shipmentapi.ShipmentResult
	var status shipmentapi.ShipmentStatus

	workflow.Go(ctx, handleStatusUpdates(ctx, &status))

	aCtx := workflow.WithActivityOptions(ctx,
		workflow.ActivityOptions{
			StartToCloseTimeout: 5 * time.Second,
		},
	)

	err := workflow.ExecuteActivity(aCtx,
		a.RegisterShipment,
		activities.RegisterShipmentInput{
			OrderID: input.OrderID,
			Items:   input.Items,
		},
	).Get(ctx, nil)
	if err != nil {
		return result, err
	}

	err = workflow.ExecuteActivity(aCtx,
		a.ShipmentCreatedNotification,
		activities.ShipmentCreatedNotificationInput{
			OrderID: input.OrderID,
		},
	).Get(ctx, nil)
	if err != nil {
		return result, err
	}

	workflow.Await(ctx, func() bool {
		return status == shipmentapi.ShipmentStatusDispatched
	})

	err = workflow.ExecuteActivity(aCtx,
		a.ShipmentDispatchedNotification,
		activities.ShipmentDispatchedNotificationInput{
			OrderID: input.OrderID,
		},
	).Get(ctx, nil)
	if err != nil {
		return result, err
	}

	workflow.Await(ctx, func() bool {
		return status == shipmentapi.ShipmentStatusDelivered
	})

	err = workflow.ExecuteActivity(aCtx,
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

func handleStatusUpdates(ctx workflow.Context, status *shipmentapi.ShipmentStatus) func(workflow.Context) {
	ch := workflow.GetSignalChannel(ctx, shipmentapi.ShipmentUpdateSignalName)

	return func(ctx workflow.Context) {
		var signal shipmentapi.ShipmentUpdateSignal

		for {
			ch.Receive(ctx, &signal)
			*status = signal.Status
		}
	}
}
