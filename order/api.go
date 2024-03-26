package order

// The Orders API is exposed as the JSON equivalents will be used to start Orders via the local API.

// Item represents an item being ordered.
// All fields are required.
type Item struct {
	SKU      string
	Quantity int32
}

// OrderInput is the input for an Order workflow.
// All fields are required.
type OrderInput struct {
	ID    string
	Items []Item
}

// OrderResult is the result of an Order workflow.
type OrderResult struct{}
