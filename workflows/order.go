package workflows

import (
	"time"

	"github.com/temporalio/orders-reference-app-go/activities"
	"github.com/temporalio/orders-reference-app-go/internal/shipmentapi"
	"github.com/temporalio/orders-reference-app-go/pkg/ordersapi"
	"go.temporal.io/sdk/workflow"
)

type orderImpl struct {
	ID           ordersapi.OrderID
	fulfillments []activities.Fulfillment
	shipments    []workflow.Future
}

func Order(ctx workflow.Context, input ordersapi.OrderInput) (ordersapi.OrderResult, error) {
	return new(orderImpl).run(ctx, input)
}

func (o *orderImpl) run(ctx workflow.Context, order ordersapi.OrderInput) (ordersapi.OrderResult, error) {
	var result ordersapi.OrderResult

	o.ID = order.ID
	err := o.fulfill(ctx, order.Items)
	if err != nil {
		return result, err
	}
	o.createShipments(ctx)

	return result, o.waitForDeliveries(ctx)
}

func (o *orderImpl) fulfill(ctx workflow.Context, items []ordersapi.Item) error {
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
		return err
	}

	o.fulfillments = result.Fulfillments

	return nil
}

func (o *orderImpl) createShipments(ctx workflow.Context) {
	for i, f := range o.fulfillments {
		ctx = workflow.WithChildOptions(ctx,
			workflow.ChildWorkflowOptions{
				WorkflowID: shipmentapi.ShipmentWorkflowID(o.ID, i),
			},
		)

		o.shipments = append(o.shipments, workflow.ExecuteChildWorkflow(ctx,
			Shipment,
			shipmentapi.ShipmentInput{
				OrderID: o.ID,
				Items:   f.Items,
			},
		))
	}
}

func (o *orderImpl) waitForDeliveries(ctx workflow.Context) error {
	for _, d := range o.shipments {
		err := d.Get(ctx, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
