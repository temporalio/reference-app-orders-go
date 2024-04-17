package order

import (
	"fmt"
	"time"

	"github.com/temporalio/orders-reference-app-go/app/billing"
	"github.com/temporalio/orders-reference-app-go/app/shipment"
	"go.temporal.io/sdk/workflow"
)

const TASK_QUEUE = "orders"

type orderImpl struct {
	ID         string
	CustomerID string
}

func Order(ctx workflow.Context, input OrderInput) (OrderResult, error) {
	if input.ID == "" {
		return OrderResult{}, fmt.Errorf("ID is required")
	}

	if input.CustomerID == "" {
		return OrderResult{}, fmt.Errorf("CustomerID is required")
	}

	if len(input.Items) == 0 {
		return OrderResult{}, fmt.Errorf("Order must contain items")
	}

	return new(orderImpl).run(ctx, input)
}

func (o *orderImpl) run(ctx workflow.Context, order OrderInput) (OrderResult, error) {
	var result OrderResult

	o.ID = order.ID
	o.CustomerID = order.CustomerID

	fulfillments, err := o.fulfill(ctx, order.Items)
	if err != nil {
		return result, err
	}

	s := workflow.NewSelector(ctx)
	for i, fulfillment := range fulfillments {
		s.AddFuture(
			o.processFulfillment(ctx, fulfillment, i),
			func(f workflow.Future) {
				if err := f.Get(ctx, nil); err != nil {
					// TODO: Figure out business logic for error handling
					workflow.GetLogger(ctx).Error("fulfillment error", "error", err)
				}
			},
		)
	}

	for range fulfillments {
		s.Select(ctx)
	}

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

func (o *orderImpl) processFulfillment(ctx workflow.Context, fulfillment Fulfillment, fulfillmentID int) workflow.Future {
	f, s := workflow.NewFuture(ctx)

	workflow.Go(ctx, func(ctx workflow.Context) {
		ctx = workflow.WithChildOptions(ctx,
			workflow.ChildWorkflowOptions{
				TaskQueue: billing.TASK_QUEUE,
			},
		)

		var billingItems []billing.Item
		for _, i := range fulfillment.Items {
			billingItems = append(billingItems, billing.Item{SKU: i.SKU, Quantity: i.Quantity})
		}

		var invoice billing.GenerateInvoiceResult

		cwf := workflow.ExecuteChildWorkflow(ctx,
			billing.GenerateInvoice,
			billing.GenerateInvoiceInput{
				CustomerID:     o.CustomerID,
				OrderReference: ShipmentWorkflowID(o.ID, fulfillmentID),
				Items:          billingItems,
			},
		)
		err := cwf.Get(ctx, &invoice)
		if err != nil {
			s.SetError(err)
			return
		}

		var charge billing.ChargeCustomerResult

		cwf = workflow.ExecuteChildWorkflow(ctx,
			billing.ChargeCustomer,
			billing.ChargeCustomerInput{
				CustomerID: o.CustomerID,
				Reference:  invoice.InvoiceReference,
				Charge:     invoice.SubTotal + invoice.Tax + invoice.Shipping,
			},
		)
		if err := cwf.Get(ctx, &charge); err != nil {
			s.SetError(err)
			return
		}

		ctx = workflow.WithChildOptions(ctx,
			workflow.ChildWorkflowOptions{
				TaskQueue:  shipment.TASK_QUEUE,
				WorkflowID: ShipmentWorkflowID(o.ID, fulfillmentID),
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
			s.SetError(err)
			return
		}

		// TODO: Anything useful to set here?
		s.SetValue(nil)
	})

	return f
}

// ShipmentWorkflowID creates a shipment workflow ID from an order ID.
func ShipmentWorkflowID(orderID string, fulfillmentID int) string {
	return fmt.Sprintf("shipment:%s:%d", orderID, fulfillmentID)
}
