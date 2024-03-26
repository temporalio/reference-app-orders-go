package workflows

import (
	"fmt"
	"time"

	"github.com/temporalio/orders-reference-app-go/activities"
	"github.com/temporalio/orders-reference-app-go/pkg/ordersapi"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"
)

type orderImpl struct {
	ID ordersapi.OrderID
}

func Order(ctx workflow.Context, input ordersapi.OrderInput) (ordersapi.OrderResult, error) {
	return new(orderImpl).run(ctx, input)
}

func (o *orderImpl) run(ctx workflow.Context, order ordersapi.OrderInput) (ordersapi.OrderResult, error) {
	var result ordersapi.OrderResult

	o.ID = order.ID

	fulfillments, err := o.fulfill(ctx, order.Items)
	if err != nil {
		return result, err
	}

	o.processShipments(ctx, fulfillments)

	return result, nil
}

func (o *orderImpl) fulfill(ctx workflow.Context, items []ordersapi.Item) ([]activities.Fulfillment, error) {
	ctx = workflow.WithActivityOptions(ctx,
		workflow.ActivityOptions{
			StartToCloseTimeout: 30 * time.Second,
		},
	)

	var result activities.FulfillOrderResult

	err := workflow.ExecuteActivity(ctx,
		a.FulfillOrder,
		activities.FulfillOrderInput{
			Items: items,
		},
	).Get(ctx, &result)
	if err != nil {
		return []activities.Fulfillment{}, err
	}

	return result.Fulfillments, nil
}

func (o *orderImpl) processShipments(ctx workflow.Context, fulfillments []activities.Fulfillment) {
	s := workflow.NewSelector(ctx)

	for i, f := range fulfillments {
		ctx = workflow.WithChildOptions(ctx,
			workflow.ChildWorkflowOptions{
				WorkflowID: ShipmentWorkflowID(o.ID, i),
			},
		)

		shipment := workflow.ExecuteChildWorkflow(ctx,
			Shipment,
			ShipmentInput{
				OrderID: o.ID,
				Items:   f.Items,
			},
		)
		s.AddFuture(shipment, func(f workflow.Future) {
			err := f.Get(ctx, nil)
			if err != nil {
				// TODO: Explore shipping failure cases/handling.
				log.With(workflow.GetLogger(ctx), "order", o.ID).Error("Shipment Error", "error", err)
			}
		})
	}

	// Handle each shipment success/failure as they happen.
	for range fulfillments {
		s.Select(ctx)
	}
}

// ShipmentWorkflowID creates a shipment workflow ID from an order ID.
func ShipmentWorkflowID(orderID ordersapi.OrderID, fulfillmentID int) string {
	return fmt.Sprintf("shipment:%s:%d", orderID, fulfillmentID)
}
