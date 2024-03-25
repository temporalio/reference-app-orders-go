package ordersapi

// The Orders API is exposed as the JSON equivalents will be used to start Orders via the local API.

// TODO: Do we want to do pass-through like this or separate and translate API definitions to Workflow types?
// TODO: Protobufs?

type OrderID string

type Item struct {
	SKU      string
	Quantity int32
}

type OrderInput struct {
	ID    OrderID
	Items []Item
}

type OrderResult struct{}
