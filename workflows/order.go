package workflows

import (
	"go.temporal.io/sdk/workflow"

	"github.com/temporalio/orders-reference-app-go/pkg/ordersapi"

	"github.com/temporalio/orders-reference-app-go/internal/shipmentapi"
)

func Order(ctx workflow.Context, order ordersapi.OrderInput) (ordersapi.OrderResult, error) {
	shipment := createShipment(ctx, order.ID, order.Items)

	err := waitForDelivery(ctx, shipment)
	if err != nil {
		return ordersapi.OrderResult{}, err
	}

	return ordersapi.OrderResult{}, nil
}

func createShipment(ctx workflow.Context, orderID string, items []ordersapi.Item) workflow.Future {
	cctx := workflow.WithChildOptions(ctx,
		workflow.ChildWorkflowOptions{
			WorkflowID: shipmentapi.ShipmentWorkflowID(orderID),
		},
	)

	return workflow.ExecuteChildWorkflow(cctx,
		Shipment,
		shipmentapi.ShipmentInput{
			OrderID: orderID,
			Items:   items,
		},
	)
}

func waitForDelivery(ctx workflow.Context, shipment workflow.Future) error {
	return shipment.Get(ctx, nil)
}
