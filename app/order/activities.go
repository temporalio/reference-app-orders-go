package order

import "context"

type Activities struct{}

var a Activities

// FulfillOrderInput is the input to the FulfillOrder activity.
type FulfillOrderInput struct {
	Items []Item
}

// Fulfillment holds a set of items that will be delivered in one shipment (due to location and stock level).
type Fulfillment struct {
	// Location is the address for courier pickup (the warehouse).
	Location string
	// Items is the set of the items that will be part of this shipment.
	Items []Item
}

// FulfillOrderResult is the result from the FulfillOrder activity.
type FulfillOrderResult struct {
	// A set of Fulfillments.
	Fulfillments []Fulfillment
}

func (a *Activities) FulfillOrder(ctx context.Context, input *FulfillOrderInput) (*FulfillOrderResult, error) {
	if len(input.Items) < 1 {
		return nil, nil
	}

	var fulfillments []Fulfillment

	// Hard coded. Open discussion where this stub boundary should live.

	// First item from one warehouse
	fulfillments = append(
		fulfillments,
		Fulfillment{
			Location: "Warehouse A",
			Items:    input.Items[0:1],
		},
	)

	if len(input.Items) > 1 {
		// Second fulfillment with all other items
		fulfillments = append(
			fulfillments,
			Fulfillment{
				Location: "Warehouse B",
				Items:    input.Items[1:len(input.Items)],
			},
		)
	}

	return &FulfillOrderResult{
		Fulfillments: fulfillments,
	}, nil
}
