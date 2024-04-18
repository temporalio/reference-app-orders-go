package order

// Item represents an item being ordered.
// All fields are required.
type Item struct {
	SKU      string
	Quantity int32
}

// OrderInput is the input for an Order workflow.
// All fields are required.
type OrderInput struct {
	ID         string
	CustomerID string
	Items      []Item
}

// OrderResult is the result of an Order workflow.
type OrderResult struct{}
