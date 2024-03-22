package workflows

import (
	"github.com/temporalio/orders-reference-app-go/internal/shipmentapi"
	"github.com/temporalio/orders-reference-app-go/pkg/ordersapi"
	"go.temporal.io/sdk/workflow"
)

type order struct {
	ID          ordersapi.OrderID
	fulfilments []fulfilment
	shipments   []workflow.Future
}

type fulfilment struct {
	Items []ordersapi.Item
}

func Order(ctx workflow.Context, input ordersapi.OrderInput) (ordersapi.OrderResult, error) {
	return new(order).run(ctx, input)
}

func (o *order) run(ctx workflow.Context, order ordersapi.OrderInput) (ordersapi.OrderResult, error) {
	var result ordersapi.OrderResult

	o.ID = order.ID
	o.fulfil(order.Items)
	o.createShipments(ctx)

	return result, o.waitForDeliveries(ctx)
}

func (o *order) fulfil(items []ordersapi.Item) {
	// Hard coded. Open discussion where this stub boundary should live.

	// First item from one warehouse
	o.fulfilments = append(
		o.fulfilments,
		fulfilment{
			Items: items[0:1],
		},
	)

	if len(items) <= 1 {
		return
	}

	// Second fulfilment with all other items
	o.fulfilments = append(
		o.fulfilments,
		fulfilment{
			Items: items[1 : len(items)-1],
		},
	)
}

func (o *order) createShipments(ctx workflow.Context) {
	for i, f := range o.fulfilments {
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

func (o *order) waitForDeliveries(ctx workflow.Context) error {
	for _, d := range o.shipments {
		err := d.Get(ctx, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
