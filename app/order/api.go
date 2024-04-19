package order

// TaskQueue is the default task queue for the Order system.
const TaskQueue = "orders"

// Item represents an item being ordered.
// All fields are required.
type Item struct {
	SKU      string
	Quantity int32
}

// OrderInput is the input for an Order workflow.
type OrderInput struct {
	ID         string
	CustomerID string
	Items      []Item
}

// OrderResult is the result of an Order workflow.
type OrderResult struct{}
