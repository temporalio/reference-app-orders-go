package shipment

import (
	"context"
)

// Activities implements the shipment package's Activities.
// Any state shared by the worker among the activities is stored here.
type Activities struct{}

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
