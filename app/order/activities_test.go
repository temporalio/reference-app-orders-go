package order_test

import (
	"log/slog"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/temporalio/reference-app-orders-go/app/config"
	"github.com/temporalio/reference-app-orders-go/app/db"
	"github.com/temporalio/reference-app-orders-go/app/order"
	"go.temporal.io/sdk/testsuite"
)

func TestInsertOrderIsIdempotent(t *testing.T) {
	t.Chdir(t.TempDir())
	ctx := t.Context()
	store := db.CreateDB(config.AppConfig{})
	require.NoError(t, store.Connect(ctx))
	t.Cleanup(func() { require.NoError(t, store.Close()) })
	require.NoError(t, store.Setup())
	api := httptest.NewServer(order.Router(nil, store, slog.Default()))
	t.Cleanup(api.Close)
	activities := &order.Activities{OrderURL: api.URL}
	input := &order.OrderStatusInsert{
		ID:         "order-1",
		CustomerID: "customer-1",
		ReceivedAt: time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC),
		Status:     order.OrderStatusPending,
	}

	require.NoError(t, activities.InsertOrder(ctx, input))
	require.NoError(t, activities.UpdateOrderStatus(ctx, &order.OrderStatusUpdate{
		ID: input.ID, Status: order.OrderStatusProcessing,
	}))
	require.NoError(t, activities.InsertOrder(ctx, input))

	var orders []db.OrderStatus
	require.NoError(t, store.GetOrders(ctx, &orders))
	require.Len(t, orders, 1)
	require.Equal(t, input.ID, orders[0].ID)
	require.True(t, input.ReceivedAt.Equal(orders[0].ReceivedAt))
	require.Equal(t, order.OrderStatusProcessing, orders[0].Status)
}

func TestFulfillOrderZeroItems(t *testing.T) {
	testSuite := testsuite.WorkflowTestSuite{}

	var a *order.Activities

	env := testSuite.NewTestActivityEnvironment()
	env.RegisterActivity(a.ReserveItems)

	input := order.ReserveItemsInput{
		Items: []*order.Item{},
	}

	future, err := env.ExecuteActivity(a.ReserveItems, &input)
	require.NoError(t, err)

	var result order.ReserveItemsResult
	require.NoError(t, future.Get(&result))

	expected := order.ReserveItemsResult{}

	require.Equal(t, expected, result)
}

func TestFulfillOrderOneItem(t *testing.T) {
	testSuite := testsuite.WorkflowTestSuite{}

	var a *order.Activities

	env := testSuite.NewTestActivityEnvironment()
	env.RegisterActivity(a.ReserveItems)

	input := order.ReserveItemsInput{
		OrderID: "test",
		Items: []*order.Item{
			{SKU: "Hiking Boots", Quantity: 2},
		},
	}

	future, err := env.ExecuteActivity(a.ReserveItems, &input)
	require.NoError(t, err)

	var result order.ReserveItemsResult
	require.NoError(t, future.Get(&result))

	expected := order.ReserveItemsResult{
		Reservations: []*order.Reservation{
			{
				Available: true,
				Location:  "Warehouse A",
				Items: []*order.Item{
					{SKU: "Hiking Boots", Quantity: 2},
				},
			},
		},
	}

	require.Equal(t, expected, result)
}

func TestFulfillOrderTwoItems(t *testing.T) {
	testSuite := testsuite.WorkflowTestSuite{}

	var a *order.Activities

	env := testSuite.NewTestActivityEnvironment()
	env.RegisterActivity(a.ReserveItems)

	input := order.ReserveItemsInput{
		OrderID: "test",
		Items: []*order.Item{
			{SKU: "Hiking Boots", Quantity: 2},
			{SKU: "Tennis Shoes", Quantity: 1},
		},
	}

	future, err := env.ExecuteActivity(a.ReserveItems, &input)
	require.NoError(t, err)

	var result order.ReserveItemsResult
	require.NoError(t, future.Get(&result))

	expected := order.ReserveItemsResult{
		Reservations: []*order.Reservation{
			{
				Available: true,
				Location:  "Warehouse A",
				Items: []*order.Item{
					{SKU: "Hiking Boots", Quantity: 2},
				},
			},
			{
				Available: true,
				Location:  "Warehouse B",
				Items: []*order.Item{
					{SKU: "Tennis Shoes", Quantity: 1},
				},
			},
		},
	}

	require.Equal(t, expected, result)
}
