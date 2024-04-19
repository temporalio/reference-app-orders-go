package order_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/temporalio/orders-reference-app-go/app/billing"
	"github.com/temporalio/orders-reference-app-go/app/order"
	"github.com/temporalio/orders-reference-app-go/app/shipment"
	"go.temporal.io/sdk/testsuite"
)

func TestOrderWorkflow(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()
	var a *order.Activities

	env.RegisterActivity(billing.GenerateInvoice)
	env.RegisterActivity(billing.ChargeCustomer)
	env.OnWorkflow(shipment.Shipment, mock.Anything, mock.Anything).Return(func(ctx context.Context, input *shipment.ShipmentInput) (*shipment.ShipmentResult, error) {
		return nil, nil
	})

	env.RegisterActivity(a.FulfillOrder)
	env.OnActivity(a.Charge, mock.Anything, mock.Anything).Return(func(ctx context.Context, input *billing.ChargeInput) (*order.ChargeResult, error) {
		return &order.ChargeResult{
			InvoiceReference: input.Reference,
			SubTotal:         1,
			Tax:              0,
			Shipping:         1,
			Success:          true,
			AuthCode:         "1234",
		}, nil
	})

	orderInput := order.OrderInput{
		ID:         "1234",
		CustomerID: "1234",
		Items: []order.Item{
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

	env.AssertActivityNumberOfCalls(t, "Charge", 2)
	env.AssertWorkflowNumberOfCalls(t, "Shipment", 2)
}
