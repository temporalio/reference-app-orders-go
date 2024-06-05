package order_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/temporalio/reference-app-orders-go/app/order"
	"go.temporal.io/sdk/testsuite"
)

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
