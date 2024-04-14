package order_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/temporalio/orders-reference-app-go/app/order"
	"github.com/temporalio/orders-reference-app-go/app/shipment"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"
)

func TestOrderWorkflow(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	var a *order.Activities

	env.RegisterWorkflow(order.Order)
	env.RegisterActivity(a.FulfillOrder)

	orderInput := order.OrderInput{
		Items: []order.Item{
			{SKU: "test1", Quantity: 1},
			{SKU: "test2", Quantity: 3},
		},
	}

	env.OnWorkflow(shipment.Shipment, mock.Anything, mock.Anything).Return(func(ctx workflow.Context, input shipment.ShipmentInput) (shipment.ShipmentResult, error) {
		return shipment.ShipmentResult{}, nil
	})

	env.ExecuteWorkflow(
		order.Order,
		orderInput,
	)

	env.AssertWorkflowNumberOfCalls(t, "Shipment", 2)

	var result order.OrderResult
	err := env.GetWorkflowResult(&result)
	assert.NoError(t, err)
}
