package order

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/temporalio/orders-reference-app-go/app/billing"
)

// Activities implements the order package's Activities.
// Any state shared by the worker among the activities is stored here.
type Activities struct{}

var a Activities

// ReserveItemsInput is the input to the ReserveItems activity.
type ReserveItemsInput struct {
	OrderID string
	Items   []*Item
}

// Reservation is a reservation of items for an order.
type Reservation struct {
	Available bool
	Location  string
	Items     []*Item
}

// ReserveItemsResult is the result from the ReserveItems activity.
type ReserveItemsResult struct {
	Reservations []*Reservation
}

// ReserveItems reserves items to satisfy an order. It returns a list of reservations for the items.
// Any unavailable items will be returned in a Reservation with Available set to false.
// In a real system this would involve an inventory database of some kind.
// For our purposes we just split orders arbitrarily.
func (a *Activities) ReserveItems(_ context.Context, input *ReserveItemsInput) (*ReserveItemsResult, error) {
	if len(input.Items) < 1 {
		return &ReserveItemsResult{}, nil
	}

	var reservations []*Reservation
	var unavailableItems []*Item
	var availableItems []*Item

	for _, item := range input.Items {
		if strings.Contains(item.SKU, "Adidas") {
			unavailableItems = append(unavailableItems, item)
		} else {
			availableItems = append(availableItems, item)
		}
	}

	if len(unavailableItems) > 0 {
		reservations = append(
			reservations,
			&Reservation{
				Available: false,
				Items:     unavailableItems,
			},
		)
	}

	// First item from one warehouse
	reservations = append(
		reservations,
		&Reservation{
			Available: true,
			Location:  "Warehouse A",
			Items:     availableItems[0:1],
		},
	)

	if len(availableItems) > 1 {
		// Second fulfillment with all other items
		reservations = append(
			reservations,
			&Reservation{
				Available: true,
				Location:  "Warehouse B",
				Items:     availableItems[1:],
			},
		)
	}

	return &ReserveItemsResult{
		Reservations: reservations,
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://127.0.0.1:8082/charge", bytes.NewReader(jsonInput))
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
