package order_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/temporalio/orders-reference-app-go/order"
	"go.temporal.io/sdk/testsuite"
)

func TestFulfillOrderOneItem(t *testing.T) {
	testSuite := testsuite.WorkflowTestSuite{}

	var a *order.Activities

	env := testSuite.NewTestActivityEnvironment()
	env.RegisterActivity(a.FulfillOrder)

	input := order.FulfillOrderInput{
		Items: []order.Item{
			{SKU: "Hiking Boots", Quantity: 2},
		},
	}

	future, err := env.ExecuteActivity(a.FulfillOrder, input)
	require.NoError(t, err)

	var result order.FulfillOrderResult
	require.NoError(t, future.Get(&result))

	fulfillments := result.Fulfillments
	require.Equal(t, 1, len(fulfillments))

	fulfillment := fulfillments[0]
	require.Equal(t, "Warehouse A", fulfillment.Location)

	require.Equal(t, 1, len(fulfillment.Items))

	item := fulfillment.Items[0]
	require.Equal(t, "Hiking Boots", item.SKU)
	require.Equal(t, int32(2), item.Quantity)
}

func TestFulfillOrderTwoItems(t *testing.T) {
	testSuite := testsuite.WorkflowTestSuite{}

	var a *order.Activities

	env := testSuite.NewTestActivityEnvironment()
	env.RegisterActivity(a.FulfillOrder)

	input := order.FulfillOrderInput{
		Items: []order.Item{
			{SKU: "Hiking Boots", Quantity: 2},
			{SKU: "Tennis Shoes", Quantity: 1},
		},
	}

	future, err := env.ExecuteActivity(a.FulfillOrder, input)
	require.NoError(t, err)

	var result order.FulfillOrderResult
	require.NoError(t, future.Get(&result))

	fulfillments := result.Fulfillments
	require.Equal(t, 2, len(fulfillments))

	fulfillmentA := fulfillments[0]
	require.Equal(t, "Warehouse A", fulfillmentA.Location)

	require.Equal(t, 1, len(fulfillmentA.Items))

	itemOne := fulfillmentA.Items[0]
	require.Equal(t, "Hiking Boots", itemOne.SKU)
	require.Equal(t, int32(2), itemOne.Quantity)

	fulfillmentB := fulfillments[1]
	require.Equal(t, "Warehouse B", fulfillmentB.Location)

	require.Equal(t, 1, len(fulfillmentB.Items))

	itemTwo := fulfillmentB.Items[0]
	require.Equal(t, "Tennis Shoes", itemTwo.SKU)
	require.Equal(t, int32(1), itemTwo.Quantity)
}
