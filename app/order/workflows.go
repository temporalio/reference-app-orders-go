package order

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/temporalio/reference-app-orders-go/app/billing"
	"github.com/temporalio/reference-app-orders-go/app/shipment"
	"go.temporal.io/sdk/workflow"
)

type orderImpl struct {
	id           string
	customerID   string
	status       string
	fulfillments []*Fulfillment
}

// Aggressively low for demo purposes.
const customerActionTimeout = 30 * time.Second

// Order Workflow process an order from a customer.
func Order(ctx workflow.Context, input *OrderInput) (*OrderResult, error) {
	wf := new(orderImpl)

	if err := wf.setup(ctx, input); err != nil {
		return nil, err
	}

	return wf.run(ctx, input)
}

func (wf *orderImpl) setup(ctx workflow.Context, input *OrderInput) error {
	if input.ID == "" {
		return fmt.Errorf("ID is required")
	}

	if input.CustomerID == "" {
		return fmt.Errorf("CustomerID is required")
	}

	if len(input.Items) == 0 {
		return fmt.Errorf("order must contain items")
	}

	wf.id = input.ID
	wf.customerID = input.CustomerID

	return workflow.SetQueryHandler(ctx, StatusQuery, func() (*OrderStatus, error) {
		return &OrderStatus{
			ID:           wf.id,
			Status:       wf.status,
			CustomerID:   wf.customerID,
			Fulfillments: wf.fulfillments,
		}, nil
	})
}

func (wf *orderImpl) run(ctx workflow.Context, order *OrderInput) (*OrderResult, error) {
	err := wf.buildFulfillments(ctx, order.Items)
	if err != nil {
		return nil, err
	}

	if wf.customerActionRequired() {
		err = wf.updateStatus(ctx, OrderStatusCustomerActionRequired)
		if err != nil {
			return nil, err
		}

		action, err := wf.waitForCustomer(ctx)
		if err != nil {
			return nil, err
		}

		switch action {
		case CustomerActionCancel:
			err := wf.updateStatus(ctx, OrderStatusCancelled)
			return &OrderResult{Status: wf.status}, err
		case CustomerActionTimedOut:
			err := wf.updateStatus(ctx, OrderStatusTimedOut)
			return &OrderResult{Status: wf.status}, err
		case CustomerActionAmend:
			wf.cancelUnavailableFulfillments()
		default:
			return nil, fmt.Errorf("unhandled customer action %q", action)
		}
	}

	if err := wf.updateStatus(ctx, OrderStatusProcessing); err != nil {
		return nil, err
	}

	workflow.Go(ctx, wf.handleShipmentStatusUpdates)

	completed := 0
	for _, f := range wf.fulfillments {
		f := f
		workflow.Go(ctx, func(ctx workflow.Context) {
			f.process(ctx)
			completed++
		})
	}

	workflow.Await(ctx, func() bool { return completed == len(wf.fulfillments) })

	if err := wf.updateStatus(ctx, OrderStatusCompleted); err != nil {
		return nil, err
	}

	return &OrderResult{Status: wf.status}, nil
}

func (wf *orderImpl) updateStatus(ctx workflow.Context, status string) error {
	wf.status = status

	update := &OrderStatusUpdate{
		ID:     wf.id,
		Status: wf.status,
	}

	ctx = workflow.WithLocalActivityOptions(ctx, workflow.LocalActivityOptions{
		ScheduleToCloseTimeout: 5 * time.Second,
	})
	return workflow.ExecuteLocalActivity(ctx, a.UpdateOrderStatus, update).Get(ctx, nil)
}

func (wf *orderImpl) buildFulfillments(ctx workflow.Context, items []*Item) error {
	ctx = workflow.WithActivityOptions(ctx,
		workflow.ActivityOptions{
			StartToCloseTimeout: 30 * time.Second,
		},
	)

	var result ReserveItemsResult

	err := workflow.ExecuteActivity(ctx,
		a.ReserveItems,
		ReserveItemsInput{
			OrderID: wf.id,
			Items:   items,
		},
	).Get(ctx, &result)
	if err != nil {
		return err
	}

	for i, r := range result.Reservations {
		f := &Fulfillment{
			orderID:    wf.id,
			customerID: wf.customerID,

			ID:       fmt.Sprintf("%s:%d", wf.id, i+1),
			Items:    r.Items,
			Location: r.Location,
			Status:   FulfillmentStatusPending,
		}
		if !r.Available {
			f.Status = FulfillmentStatusUnavailable
		}
		wf.fulfillments = append(wf.fulfillments, f)
	}

	return nil
}

func (wf *orderImpl) customerActionRequired() bool {
	for _, f := range wf.fulfillments {
		if f.Status == FulfillmentStatusUnavailable {
			return true
		}
	}

	return false
}

func (wf *orderImpl) cancelUnavailableFulfillments() {
	for _, f := range wf.fulfillments {
		if f.Status == FulfillmentStatusUnavailable {
			f.Status = FulfillmentStatusCancelled
		}
	}
}

func (wf *orderImpl) waitForCustomer(ctx workflow.Context) (string, error) {
	var signal CustomerActionSignal

	s := workflow.NewSelector(ctx)

	timerCtx, cancelTimer := workflow.WithCancel(ctx)
	t := workflow.NewTimer(timerCtx, customerActionTimeout)

	var err error

	s.AddFuture(t, func(f workflow.Future) {
		if err = f.Get(timerCtx, nil); err != nil {
			return
		}

		signal.Action = CustomerActionTimedOut
	})

	ch := workflow.GetSignalChannel(ctx, CustomerActionSignalName)
	s.AddReceive(ch, func(c workflow.ReceiveChannel, _ bool) {
		c.Receive(ctx, &signal)
		cancelTimer()
	})
	s.Select(ctx)

	if err != nil {
		return "", err
	}

	switch signal.Action {
	case CustomerActionAmend:
	case CustomerActionCancel:
	case CustomerActionTimedOut:
	default:
		return "", fmt.Errorf("invalid customer action %q", signal.Action)
	}

	return signal.Action, nil
}

func (wf *orderImpl) handleShipmentStatusUpdates(ctx workflow.Context) {
	ch := workflow.GetSignalChannel(ctx, shipment.ShipmentStatusUpdatedSignalName)

	for {
		var signal shipment.ShipmentStatusUpdatedSignal
		_ = ch.Receive(ctx, &signal)
		for _, f := range wf.fulfillments {
			if f.ID == signal.ShipmentID {
				f.Shipment.Status = signal.Status
				f.Shipment.UpdatedAt = signal.UpdatedAt
				break
			}
		}
	}
}

func (f *Fulfillment) process(ctx workflow.Context) error {
	if f.Status == FulfillmentStatusCancelled {
		return nil
	}

	f.Status = FulfillmentStatusProcessing

	err := f.processPayment(ctx)
	if err != nil || f.Payment.Status != PaymentStatusSuccess {
		f.Status = FulfillmentStatusFailed
		return err
	}

	if err := f.processShipment(ctx); err != nil {
		f.Status = FulfillmentStatusFailed
		return err
	}

	f.Status = FulfillmentStatusCompleted

	return nil
}

func (f *Fulfillment) processPayment(ctx workflow.Context) error {
	var billingItems []billing.Item
	for _, i := range f.Items {
		billingItems = append(billingItems, billing.Item{SKU: i.SKU, Quantity: i.Quantity})
	}

	var charge ChargeResult

	ctx = workflow.WithActivityOptions(ctx,
		workflow.ActivityOptions{
			StartToCloseTimeout: 30 * time.Second,
		},
	)

	f.Payment = &PaymentStatus{Status: PaymentStatusPending}

	var chargeKey string
	v := workflow.SideEffect(ctx, func(_ workflow.Context) any {
		return uuid.NewString()
	})
	if err := v.Get(&chargeKey); err != nil {
		f.Payment.Status = PaymentStatusFailed
		return err
	}

	c := workflow.ExecuteActivity(ctx,
		a.Charge,
		&ChargeInput{
			CustomerID:     f.customerID,
			Reference:      f.ID,
			Items:          billingItems,
			IdempotencyKey: chargeKey,
		},
	)
	if err := c.Get(ctx, &charge); err != nil {
		f.Payment.Status = PaymentStatusFailed
		return err
	}

	p := f.Payment

	p.SubTotal = charge.SubTotal
	p.Tax = charge.Tax
	p.Shipping = charge.Shipping
	p.Total = charge.Total
	if charge.Success {
		p.Status = PaymentStatusSuccess
	} else {
		p.Status = PaymentStatusFailed
	}

	return nil
}

func (f *Fulfillment) processShipment(ctx workflow.Context) error {
	ctx = workflow.WithChildOptions(ctx,
		workflow.ChildWorkflowOptions{
			TaskQueue:  shipment.TaskQueue,
			WorkflowID: shipment.ShipmentWorkflowID(f.ID),
		},
	)

	var shippingItems []shipment.Item
	for _, i := range f.Items {
		shippingItems = append(shippingItems, shipment.Item{SKU: i.SKU, Quantity: i.Quantity})
	}

	f.Shipment = &ShipmentStatus{
		ID:        f.ID,
		Status:    shipment.ShipmentStatusPending,
		UpdatedAt: workflow.Now(ctx),
	}

	err := workflow.ExecuteChildWorkflow(ctx,
		shipment.Shipment,
		shipment.ShipmentInput{
			RequestorWID: workflow.GetInfo(ctx).WorkflowExecution.ID,

			ID:    f.ID,
			Items: shippingItems,
		},
	).Get(ctx, nil)

	return err
}
