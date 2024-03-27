package order

import (
	"fmt"
	"time"

	"github.com/temporalio/orders-reference-app-go/shipment"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"
)

type orderImpl struct {
	ID string
}

func Order(ctx workflow.Context, input OrderInput) (OrderResult, error) {
	return new(orderImpl).run(ctx, input)
}

func (o *orderImpl) run(ctx workflow.Context, order OrderInput) (OrderResult, error) {
	var result OrderResult

	o.ID = order.ID

	fulfillments, err := o.fulfill(ctx, order.Items)
	if err != nil {
		return result, err
	}

	o.processShipments(ctx, fulfillments)

	return result, nil
}

func (o *orderImpl) fulfill(ctx workflow.Context, items []Item) ([]Fulfillment, error) {
	ctx = workflow.WithActivityOptions(ctx,
		workflow.ActivityOptions{
			StartToCloseTimeout: 30 * time.Second,
		},
	)

	var result FulfillOrderResult

	err := workflow.ExecuteActivity(ctx,
		a.FulfillOrder,
		FulfillOrderInput{
			Items: items,
		},
	).Get(ctx, &result)
	if err != nil {
		return []Fulfillment{}, err
	}

	return result.Fulfillments, nil
}

func (o *orderImpl) processShipments(ctx workflow.Context, fulfillments []Fulfillment) {
	s := workflow.NewSelector(ctx)

	for i, f := range fulfillments {
		ctx = workflow.WithChildOptions(ctx,
			workflow.ChildWorkflowOptions{
				WorkflowID: ShipmentWorkflowID(o.ID, i),
			},
		)

		var items []shipment.Item
		for _, i := range f.Items {
			items = append(items, shipment.Item{SKU: i.SKU, Quantity: i.Quantity})
		}

		shipment := workflow.ExecuteChildWorkflow(ctx,
			shipment.Shipment,
			shipment.ShipmentInput{
				OrderID: o.ID,
				Items:   items,
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
func ShipmentWorkflowID(orderID string, fulfillmentID int) string {
	return fmt.Sprintf("shipment:%s:%d", orderID, fulfillmentID)
}
