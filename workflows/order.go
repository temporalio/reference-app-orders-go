package workflows

import (
	"github.com/temporalio/orders-reference-app-go/internal/shipmentapi"
	"github.com/temporalio/orders-reference-app-go/pkg/ordersapi"
	"go.temporal.io/sdk/workflow"
)

func Order(ctx workflow.Context, order ordersapi.OrderInput) (ordersapi.OrderResult, error) {
	var result ordersapi.OrderResult

	ctx = workflow.WithChildOptions(ctx,
		workflow.ChildWorkflowOptions{
			WorkflowID: shipmentapi.ShipmentWorkflowID(order.ID),
		},
	)

	err := workflow.ExecuteChildWorkflow(ctx,
		Shipment,
		shipmentapi.ShipmentInput{
			OrderID: order.ID,
			Items:   order.Items,
		},
	).Get(ctx, nil)
	if err != nil {
		return result, err
	}

	return result, nil
}
