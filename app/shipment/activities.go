package shipment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Activities implements the shipment package's Activities.
// Any state shared by the worker among the activities is stored here.
type Activities struct {
	ShipmentURL string
}

var a Activities

// BookShipmentInput is the input for the BookShipment operation.
// All fields are required.
type BookShipmentInput struct {
	Reference string
	Items     []Item
}

// BookShipmentResult is the result for the BookShipment operation.
// CourierReference is recorded where available, to allow tracking enquiries.
type BookShipmentResult struct {
	CourierReference string
}

// BookShipment engages a courier who can deliver the shipment to the customer
func (a *Activities) BookShipment(_ context.Context, input *BookShipmentInput) (*BookShipmentResult, error) {
	return &BookShipmentResult{
		CourierReference: input.Reference + ":1234",
	}, nil
}

// UpdateShipmentStatus stores the Order status to the database.
func (a *Activities) UpdateShipmentStatus(ctx context.Context, status *ShipmentStatusUpdate) error {
	jsonInput, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("unable to encode status: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.ShipmentURL+"/shipments/"+status.ID, bytes.NewReader(jsonInput))
	if err != nil {
		return fmt.Errorf("unable to build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("%s: %s", http.StatusText(res.StatusCode), body)
	}

	return nil
}
