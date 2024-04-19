package order

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/temporalio/orders-reference-app-go/app/billing"
)

// Activities implements the order package's Activities.
// Any state shared by the worker among the activities is stored here.
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

// FulfillOrder creates fulfillments to satisfy an order.
// In a real system this would involve an inventory database of some kind.
// For our purposes we just split orders arbitrarily.
func (a *Activities) FulfillOrder(_ context.Context, input *FulfillOrderInput) (*FulfillOrderResult, error) {
	if len(input.Items) < 1 {
		return &FulfillOrderResult{}, nil
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

// ChargeInput is the input to the Charge activity.
type ChargeInput = billing.ChargeInput

// ChargeResult is the result of the Charge activity.
type ChargeResult = billing.ChargeResult

// Charge charges a customer for a fulfillment via the Billing API
func (a *Activities) Charge(ctx context.Context, input *ChargeInput) (*ChargeResult, error) {
	jsonInput, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("unable to encode input: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:8082/charge", bytes.NewReader(jsonInput))
	if err != nil {
		return nil, fmt.Errorf("unable to build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("%s: %s", http.StatusText(res.StatusCode), body)
	}

	var result ChargeResult

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
