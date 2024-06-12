package order_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/temporalio/reference-app-orders-go/app/order"
	"github.com/temporalio/reference-app-orders-go/app/shipment"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"
)

func TestOrderWorkflow(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	var a *order.Activities

	env.RegisterActivity(a.ReserveItems)
	env.OnActivity(a.Charge, mock.Anything, mock.Anything).Return(func(ctx context.Context, input *order.ChargeInput) (*order.ChargeResult, error) {
		return &order.ChargeResult{Success: true}, nil
	})
	env.OnActivity(a.UpdateOrderStatus, mock.Anything, mock.Anything).Return(func(ctx context.Context, input *order.OrderStatusUpdate) error {
		return nil
	})
	env.OnWorkflow(shipment.Shipment, mock.Anything, mock.Anything).Return(func(ctx workflow.Context, input *shipment.ShipmentInput) (*shipment.ShipmentResult, error) {
		return &shipment.ShipmentResult{CourierReference: "test"}, nil
	}).Times(2)

	orderInput := order.OrderInput{
		ID:         "1234",
		CustomerID: "1234",
		Items: []*order.Item{
			{SKU: "test1", Quantity: 1},
			{SKU: "test2", Quantity: 3},
		},
	}

	env.ExecuteWorkflow(
		order.Order,
		&orderInput,
	)

	var result order.OrderResult
	err := env.GetWorkflowResult(&result)
	assert.NoError(t, err)

	env.AssertWorkflowNumberOfCalls(t, "Shipment", 2)
}

func TestOrderShipmentStatus(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	var a *order.Activities

	env.RegisterActivity(a.ReserveItems)
	env.OnActivity(a.Charge, mock.Anything, mock.Anything).Return(func(ctx context.Context, input *order.ChargeInput) (*order.ChargeResult, error) {
		return &order.ChargeResult{Success: true}, nil
	})
	env.OnActivity(a.UpdateOrderStatus, mock.Anything, mock.Anything).Return(func(ctx context.Context, input *order.OrderStatusUpdate) error {
		return nil
	})
	env.OnWorkflow(shipment.Shipment, mock.Anything, mock.Anything).Return(func(ctx workflow.Context, input *shipment.ShipmentInput) (*shipment.ShipmentResult, error) {
		env.SignalWorkflow(
			shipment.ShipmentStatusUpdatedSignalName,
			shipment.ShipmentStatusUpdatedSignal{
				ShipmentID: input.ID,
				Status:     shipment.ShipmentStatusDelivered,
				UpdatedAt:  env.Now(),
			},
		)

		return &shipment.ShipmentResult{CourierReference: "test"}, nil
	}).Times(2)

	orderInput := order.OrderInput{
		ID:         "1234",
		CustomerID: "1234",
		Items: []*order.Item{
			{SKU: "test1", Quantity: 1},
		},
	}

	env.ExecuteWorkflow(
		order.Order,
		&orderInput,
	)

	var result order.OrderResult
	err := env.GetWorkflowResult(&result)
	assert.NoError(t, err)

	var status order.OrderStatus
	v, err := env.QueryWorkflow(order.StatusQuery, nil)
	assert.NoError(t, err)

	err = v.Get(&status)
	assert.NoError(t, err)

	f := status.Fulfillments[0]
	assert.Equal(t, shipment.ShipmentStatusDelivered, f.Shipment.Status)
}

func TestOrderAmendWithUnavailableItems(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	var a *order.Activities

	env.RegisterActivity(a.ReserveItems)
	env.OnActivity(a.Charge, mock.Anything, mock.Anything).Return(func(ctx context.Context, input *order.ChargeInput) (*order.ChargeResult, error) {
		return &order.ChargeResult{Success: true}, nil
	})
	env.OnActivity(a.UpdateOrderStatus, mock.Anything, mock.Anything).Return(func(ctx context.Context, input *order.OrderStatusUpdate) error {
		return nil
	})
	env.OnWorkflow(shipment.Shipment, mock.Anything, mock.Anything).Return(func(ctx workflow.Context, input *shipment.ShipmentInput) (*shipment.ShipmentResult, error) {
		return &shipment.ShipmentResult{CourierReference: "test"}, nil
	})

	orderInput := order.OrderInput{
		ID:         "1234",
		CustomerID: "1234",
		Items: []*order.Item{
			{SKU: "Adidas", Quantity: 1},
			{SKU: "test2", Quantity: 3},
		},
	}

	env.RegisterDelayedCallback(func() {
		var status order.OrderStatus
		v, err := env.QueryWorkflow(order.StatusQuery, nil)
		assert.NoError(t, err)

		err = v.Get(&status)
		assert.Equal(t, order.OrderStatus{
			ID:         "1234",
			CustomerID: "1234",
			Status:     order.OrderStatusCustomerActionRequired,
			Fulfillments: []*order.Fulfillment{
				{
					ID:     "1234:1",
					Status: order.FulfillmentStatusUnavailable,
					Items: []*order.Item{
						{SKU: "Adidas", Quantity: 1},
					},
				},
				{
					ID:       "1234:2",
					Status:   order.FulfillmentStatusPending,
					Location: "Warehouse A",
					Items: []*order.Item{
						{SKU: "test2", Quantity: 3},
					},
				},
			},
		}, status)
	}, time.Second*1)

	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow(
			order.CustomerActionSignalName,
			order.CustomerActionSignal{
				Action: order.CustomerActionAmend,
			},
		)

	}, time.Second*2)

	env.ExecuteWorkflow(
		order.Order,
		&orderInput,
	)

	var result order.OrderResult
	err := env.GetWorkflowResult(&result)
	assert.NoError(t, err)

	var status order.OrderStatus
	v, err := env.QueryWorkflow(order.StatusQuery, nil)
	assert.NoError(t, err)

	err = v.Get(&status)
	assert.Len(t, status.Fulfillments, 2)

	f := status.Fulfillments[0]
	assert.Equal(t, order.FulfillmentStatusCancelled, f.Status)

	f = status.Fulfillments[1]
	assert.Equal(t, order.PaymentStatusSuccess, f.Payment.Status)
	assert.Equal(t, f.ID, f.Shipment.ID)

	env.AssertWorkflowNumberOfCalls(t, "Shipment", 1)
}

func TestOrderCancelWithUnavailableItems(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	var a *order.Activities

	env.RegisterActivity(a.ReserveItems)
	env.OnActivity(a.UpdateOrderStatus, mock.Anything, mock.Anything).Return(func(ctx context.Context, input *order.OrderStatusUpdate) error {
		return nil
	})

	orderInput := order.OrderInput{
		ID:         "1234",
		CustomerID: "1234",
		Items: []*order.Item{
			{SKU: "Adidas", Quantity: 1},
			{SKU: "test2", Quantity: 3},
		},
	}

	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow(
			order.CustomerActionSignalName,
			order.CustomerActionSignal{
				Action: order.CustomerActionCancel,
			},
		)
	}, 1)

	env.ExecuteWorkflow(
		order.Order,
		&orderInput,
	)

	var result order.OrderResult
	err := env.GetWorkflowResult(&result)
	assert.NoError(t, err)
}

func TestOrderCancelAfterTimeout(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	var a *order.Activities

	env.RegisterActivity(a.ReserveItems)
	env.OnActivity(a.UpdateOrderStatus, mock.Anything, mock.Anything).Return(func(ctx context.Context, input *order.OrderStatusUpdate) error {
		return nil
	})

	orderInput := order.OrderInput{
		ID:         "1234",
		CustomerID: "1234",
		Items: []*order.Item{
			{SKU: "Adidas", Quantity: 1},
			{SKU: "test2", Quantity: 3},
		},
	}

	env.ExecuteWorkflow(
		order.Order,
		&orderInput,
	)

	var result order.OrderResult
	err := env.GetWorkflowResult(&result)
	assert.NoError(t, err)

	assert.Equal(t, order.OrderStatusTimedOut, result.Status)
}
