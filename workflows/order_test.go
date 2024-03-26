package workflows_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/temporalio/orders-reference-app-go/activities"
	"github.com/temporalio/orders-reference-app-go/pkg/ordersapi"
	"github.com/temporalio/orders-reference-app-go/workflows"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"
)

func TestOrderWorkflow(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	var a *activities.Activities

	env.RegisterWorkflow(workflows.Order)
	env.RegisterActivity(a.FulfillOrder)

	orderInput := ordersapi.OrderInput{
		Items: []ordersapi.Item{
			{SKU: "test1", Quantity: 1},
			{SKU: "test2", Quantity: 3},
		},
	}

	env.OnWorkflow(workflows.Shipment, mock.Anything, mock.Anything).Return(func(ctx workflow.Context, input workflows.ShipmentInput) (workflows.ShipmentResult, error) {
		return workflows.ShipmentResult{}, nil
	})

	env.ExecuteWorkflow(
		workflows.Order,
		orderInput,
	)

	env.AssertWorkflowNumberOfCalls(t, "Shipment", 2)

	var result ordersapi.OrderResult
	err := env.GetWorkflowResult(&result)
	assert.NoError(t, err)
}
