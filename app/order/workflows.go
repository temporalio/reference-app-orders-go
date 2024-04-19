package order

import (
	"fmt"
	"time"

	"github.com/temporalio/orders-reference-app-go/app/billing"
	"github.com/temporalio/orders-reference-app-go/app/shipment"
	"go.temporal.io/sdk/workflow"
)

type orderImpl struct {
	ID         string
	CustomerID string
}

// Order Workflow process an order from a customer.
func Order(ctx workflow.Context, input *OrderInput) (*OrderResult, error) {
	if input.ID == "" {
		return nil, fmt.Errorf("ID is required")
	}

	if input.CustomerID == "" {
		return nil, fmt.Errorf("CustomerID is required")
	}

	if len(input.Items) == 0 {
		return nil, fmt.Errorf("order must contain items")
	}

	return new(orderImpl).run(ctx, input)
}

func (o *orderImpl) run(ctx workflow.Context, order *OrderInput) (*OrderResult, error) {
	var result OrderResult

	o.ID = order.ID
	o.CustomerID = order.CustomerID

	fulfillments, err := o.fulfill(ctx, order.Items)
	if err != nil {
		return nil, err
	}

	completed := 0
	for i, fulfillment := range fulfillments {
		workflow.Go(ctx, func(ctx workflow.Context) {
			err := o.processFulfillment(ctx, fulfillment, i)
			if err != nil {
				workflow.GetLogger(ctx).Error("fulfillment error", "order", order.ID, "fulfillment", i, "error", err)
			}
			completed++
		})
	}

	workflow.Await(ctx, func() bool { return completed == len(fulfillments) })

	return &result, nil
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

func (o *orderImpl) processFulfillment(ctx workflow.Context, fulfillment Fulfillment, fulfillmentID int) error {
	ref := ShipmentWorkflowID(o.ID, fulfillmentID)

	var billingItems []billing.Item
	for _, i := range fulfillment.Items {
		billingItems = append(billingItems, billing.Item{SKU: i.SKU, Quantity: i.Quantity})
	}

	var charge ChargeResult

	ctx = workflow.WithActivityOptions(ctx,
		workflow.ActivityOptions{
			StartToCloseTimeout: 30 * time.Second,
		},
	)

	cwf := workflow.ExecuteActivity(ctx,
		a.Charge,
		&ChargeInput{
			CustomerID: o.CustomerID,
			Reference:  ref,
			Items:      billingItems,
		},
	)
	err := cwf.Get(ctx, &charge)
	if err != nil {
		// TODO: Payment specific errors for business logic
		return err
	}

	ctx = workflow.WithChildOptions(ctx,
		workflow.ChildWorkflowOptions{
			TaskQueue:  shipment.TaskQueue,
			WorkflowID: ref,
		},
	)

	var shippingItems []shipment.Item
	for _, i := range fulfillment.Items {
		shippingItems = append(shippingItems, shipment.Item{SKU: i.SKU, Quantity: i.Quantity})
	}

	shipment := workflow.ExecuteChildWorkflow(ctx,
		shipment.Shipment,
		shipment.ShipmentInput{
			OrderID: o.ID,
			Items:   shippingItems,
		},
	)
	if err := shipment.Get(ctx, nil); err != nil {
		// TODO: On shipment failure, prompt user if they'd like to re-ship or cancel
		return err
	}

	return nil
}

// ShipmentWorkflowID creates a shipment workflow ID from an order ID.
func ShipmentWorkflowID(orderID string, fulfillmentID int) string {
	return fmt.Sprintf("shipment:%s:%d", orderID, fulfillmentID)
}
