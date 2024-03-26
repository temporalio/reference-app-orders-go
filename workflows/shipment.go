package workflows

import (
	"time"

	"github.com/temporalio/orders-reference-app-go/activities"
	"github.com/temporalio/orders-reference-app-go/internal/shipmentapi"
	"go.temporal.io/sdk/workflow"
)

type shipmentImpl struct {
	status shipmentapi.ShipmentStatus
}

// Shipment implements the Shipment workflow.
func Shipment(ctx workflow.Context, input shipmentapi.ShipmentInput) (shipmentapi.ShipmentResult, error) {
	return new(shipmentImpl).run(ctx, input)
}

func (s *shipmentImpl) run(ctx workflow.Context, input shipmentapi.ShipmentInput) (shipmentapi.ShipmentResult, error) {
	workflow.Go(ctx, s.statusUpdater)

	var result shipmentapi.ShipmentResult

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

	s.waitForStatus(ctx, shipmentapi.ShipmentStatusDispatched)

	err = workflow.ExecuteActivity(ctx,
		a.ShipmentDispatchedNotification,
		activities.ShipmentDispatchedNotificationInput{
			OrderID: input.OrderID,
		},
	).Get(ctx, nil)
	if err != nil {
		return result, err
	}

	s.waitForStatus(ctx, shipmentapi.ShipmentStatusDelivered)

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

func (s *shipmentImpl) statusUpdater(ctx workflow.Context) {
	var signal shipmentapi.ShipmentUpdateSignal

	ch := workflow.GetSignalChannel(ctx, shipmentapi.ShipmentUpdateSignalName)
	for {
		ch.Receive(ctx, &signal)
		s.status = signal.Status
	}
}

func (s *shipmentImpl) waitForStatus(ctx workflow.Context, status shipmentapi.ShipmentStatus) {
	workflow.Await(ctx, func() bool {
		return s.status == status
	})
}
