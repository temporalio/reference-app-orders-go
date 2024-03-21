package ordersapi

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
