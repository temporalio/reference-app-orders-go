package workflows

import (
	"time"

	"github.com/temporalio/orders-reference-app-go/activities"
	"github.com/temporalio/orders-reference-app-go/internal/shipmentapi"
	"go.temporal.io/sdk/workflow"
)

type shipmentStatus struct {
	status shipmentapi.ShipmentStatus
}

func Shipment(ctx workflow.Context, input shipmentapi.ShipmentInput) (shipmentapi.ShipmentResult, error) {
	var result shipmentapi.ShipmentResult
	var status shipmentStatus

	workflow.Go(ctx, status.handleUpdates(ctx))

	ctx = workflow.WithActivityOptions(ctx,
		workflow.ActivityOptions{
			StartToCloseTimeout: 5 * time.Second,
		},
	)

	err := workflow.ExecuteActivity(ctx,
		a.RegisterShipment,
		activities.RegisterShipmentInput{
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

	status.waitUntil(ctx, shipmentapi.ShipmentStatusDispatched)

	err = workflow.ExecuteActivity(ctx,
		a.ShipmentDispatchedNotification,
		activities.ShipmentDispatchedNotificationInput{
			OrderID: input.OrderID,
		},
	).Get(ctx, nil)
	if err != nil {
		return result, err
	}

	status.waitUntil(ctx, shipmentapi.ShipmentStatusDelivered)

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

func (s *shipmentStatus) handleUpdates(ctx workflow.Context) func(workflow.Context) {
	ch := workflow.GetSignalChannel(ctx, shipmentapi.ShipmentUpdateSignalName)

	return func(ctx workflow.Context) {
		var signal shipmentapi.ShipmentUpdateSignal

		for {
			ch.Receive(ctx, &signal)
			s.status = signal.Status
		}
	}
}

func (s *shipmentStatus) waitUntil(ctx workflow.Context, status shipmentapi.ShipmentStatus) {
	workflow.Await(ctx, func() bool {
		return s.status == status
	})
}
