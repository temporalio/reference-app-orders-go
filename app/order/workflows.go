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

	err = o.processBilling(ctx, fulfillments)
	if err != nil {
		return result, err
	}

	err = o.processShipments(ctx, fulfillments)
	if err != nil {
		return result, err
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

func (o *orderImpl) processBilling(ctx workflow.Context, fulfillments []Fulfillment) error {
	for i, f := range fulfillments {
		var items []billing.Item
		for _, i := range f.Items {
			items = append(items, billing.Item{SKU: i.SKU, Quantity: i.Quantity})
		}

		var invoice billing.GenerateInvoiceResult

		cwf := workflow.ExecuteChildWorkflow(ctx,
			billing.GenerateInvoice,
			billing.GenerateInvoiceInput{
				CustomerID:     o.CustomerID,
				OrderReference: ShipmentWorkflowID(o.ID, i),
				Items:          items,
			},
		)
		err := cwf.Get(ctx, &invoice)
		if err != nil {
			// TODO: Explore invoicing failure handling.
			return err
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
		err = cwf.Get(ctx, &charge)
		if err != nil {
			// TODO: Explore payment failure handling.
			return err
		}
	}

	return nil
}

func (o *orderImpl) processShipments(ctx workflow.Context, fulfillments []Fulfillment) error {
	s := workflow.NewSelector(ctx)
	var err error

	for i, f := range fulfillments {
		ctx = workflow.WithChildOptions(ctx,
			workflow.ChildWorkflowOptions{
				TaskQueue:  shipment.TASK_QUEUE,
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
			err = f.Get(ctx, nil)
		})
	}

	// Handle each shipment success/failure as they happen.
	for range fulfillments {
		s.Select(ctx)
		if err != nil {
			// TODO: Explore shipping failure cases/handling.
			return err
		}
	}

	return nil
}

// ShipmentWorkflowID creates a shipment workflow ID from an order ID.
func ShipmentWorkflowID(orderID string, fulfillmentID int) string {
	return fmt.Sprintf("shipment:%s:%d", orderID, fulfillmentID)
}
