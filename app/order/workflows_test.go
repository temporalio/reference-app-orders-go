package order_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/temporalio/reference-app-orders-go/app/order"
	"github.com/temporalio/reference-app-orders-go/app/shipment"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"
)

func TestOrderReceivedAt(t *testing.T) {
	startTime := time.Date(2026, time.January, 2, 12, 0, 0, 0, time.UTC)
	receivedAt := startTime.Add(-time.Hour)
	for _, test := range []struct {
		name       string
		receivedAt time.Time
		expected   time.Time
		legacy     bool
	}{
		{name: "supplied receipt time", receivedAt: receivedAt, expected: receivedAt},
		{name: "input without receipt time", expected: startTime},
		{name: "legacy workflow history", expected: startTime, legacy: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			s := testsuite.WorkflowTestSuite{}
			env := s.NewTestWorkflowEnvironment()
			env.SetStartTime(startTime)
			if test.legacy {
				env.OnGetVersion("order-received-at", workflow.DefaultVersion, 1).Return(workflow.DefaultVersion)
			}
			var a *order.Activities
			var insertedAt time.Time
			env.OnActivity(a.InsertOrder, mock.Anything, mock.Anything).
				Return(func(_ context.Context, input *order.OrderStatusInsert) error {
					insertedAt = input.ReceivedAt
					return nil
				}).Once()
			env.OnActivity(a.ReserveItems, mock.Anything, mock.Anything).
				Return(&order.ReserveItemsResult{}, nil).Once()
			env.OnActivity(a.UpdateOrderStatus, mock.Anything, mock.Anything).Return(nil)

			env.ExecuteWorkflow(order.Order, &order.OrderInput{
				ID: "order-1", CustomerID: "customer-1",
				Items: []*order.Item{{SKU: "test", Quantity: 1}}, ReceivedAt: test.receivedAt,
			})
			require.NoError(t, env.GetWorkflowError())
			result, err := env.QueryWorkflow(order.StatusQuery)
			require.NoError(t, err)
			var status order.OrderStatus
			require.NoError(t, result.Get(&status))
			require.WithinDuration(t, test.expected, insertedAt, 0)
			require.WithinDuration(t, insertedAt, status.ReceivedAt, 0)
			env.AssertExpectations(t)
		})
	}
}

func TestOrderWorkflow(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	var a *order.Activities

	env.RegisterActivity(a.ReserveItems)
	env.OnActivity(a.InsertOrder, mock.Anything, mock.Anything).Return(func(ctx context.Context, input *order.OrderStatusInsert) error {
		return nil
	})
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
	env.OnActivity(a.InsertOrder, mock.Anything, mock.Anything).Return(func(_ctx context.Context, _input *order.OrderStatusInsert) error {
		return nil
	})
	env.OnActivity(a.Charge, mock.Anything, mock.Anything).Return(func(_ctx context.Context, _input *order.ChargeInput) (*order.ChargeResult, error) {
		return &order.ChargeResult{Success: true}, nil
	})
	env.OnActivity(a.UpdateOrderStatus, mock.Anything, mock.Anything).Return(func(_ctx context.Context, _input *order.OrderStatusUpdate) error {
		return nil
	})
	env.OnWorkflow(shipment.Shipment, mock.Anything, mock.Anything).Return(func(_ctx workflow.Context, input *shipment.ShipmentInput) (*shipment.ShipmentResult, error) {
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
	env.OnActivity(a.InsertOrder, mock.Anything, mock.Anything).Return(func(_ctx context.Context, _input *order.OrderStatusInsert) error {
		return nil
	})
	env.OnActivity(a.Charge, mock.Anything, mock.Anything).Return(func(_ctx context.Context, _input *order.ChargeInput) (*order.ChargeResult, error) {
		return &order.ChargeResult{Success: true}, nil
	})
	env.OnActivity(a.UpdateOrderStatus, mock.Anything, mock.Anything).Return(func(_ctx context.Context, _input *order.OrderStatusUpdate) error {
		return nil
	})
	env.OnWorkflow(shipment.Shipment, mock.Anything, mock.Anything).Return(func(_ctx workflow.Context, _input *shipment.ShipmentInput) (*shipment.ShipmentResult, error) {
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
		assert.NoError(t, err)
		assert.Equal(t, "1234", status.ID)
		assert.Equal(t, "1234", status.CustomerID)
		assert.Equal(t, order.OrderStatusCustomerActionRequired, status.Status)
		assert.NotZero(t, status.ReceivedAt)
		assert.Len(t, status.Fulfillments, 2)

		assert.Equal(t, "1234:1", status.Fulfillments[0].ID)
		assert.Equal(t, order.FulfillmentStatusUnavailable, status.Fulfillments[0].Status)
		assert.Equal(t, []*order.Item{{SKU: "Adidas", Quantity: 1}}, status.Fulfillments[0].Items)

		assert.Equal(t, "1234:2", status.Fulfillments[1].ID)
		assert.Equal(t, order.FulfillmentStatusPending, status.Fulfillments[1].Status)
		assert.Equal(t, "Warehouse A", status.Fulfillments[1].Location)
		assert.Equal(t, []*order.Item{{SKU: "test2", Quantity: 3}}, status.Fulfillments[1].Items)
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
	env.OnActivity(a.InsertOrder, mock.Anything, mock.Anything).Return(func(_ctx context.Context, _input *order.OrderStatusInsert) error {
		return nil
	})
	env.OnActivity(a.UpdateOrderStatus, mock.Anything, mock.Anything).Return(func(_ctx context.Context, _input *order.OrderStatusUpdate) error {
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
	env.OnActivity(a.InsertOrder, mock.Anything, mock.Anything).Return(func(ctx context.Context, input *order.OrderStatusInsert) error {
		return nil
	})
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

func TestOrderIgnoresInvalidCustomerAction(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	var a *order.Activities

	env.RegisterActivity(a.ReserveItems)
	env.OnActivity(a.InsertOrder, mock.Anything, mock.Anything).Return(func(_ctx context.Context, _input *order.OrderStatusInsert) error {
		return nil
	})
	env.OnActivity(a.UpdateOrderStatus, mock.Anything, mock.Anything).Return(func(_ctx context.Context, _input *order.OrderStatusUpdate) error {
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

	// A customer action outside the known set must be ignored: the order
	// must keep waiting for a valid action instead of failing.
	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow(
			order.CustomerActionSignalName,
			order.CustomerActionSignal{
				Action: "refund",
			},
		)
	}, time.Second*1)

	env.RegisterDelayedCallback(func() {
		v, err := env.QueryWorkflow(order.StatusQuery)
		assert.NoError(t, err)
		var status order.OrderStatus
		assert.NoError(t, v.Get(&status))
		assert.Equal(t, order.OrderStatusCustomerActionRequired, status.Status)
	}, time.Second*2)

	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow(
			order.CustomerActionSignalName,
			order.CustomerActionSignal{
				Action: order.CustomerActionCancel,
			},
		)
	}, time.Second*3)

	env.ExecuteWorkflow(
		order.Order,
		&orderInput,
	)

	var result order.OrderResult
	err := env.GetWorkflowResult(&result)
	assert.NoError(t, err)

	assert.Equal(t, order.OrderStatusCancelled, result.Status)
}
