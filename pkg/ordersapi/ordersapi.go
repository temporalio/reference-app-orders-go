package ordersapi

type Item struct {
	SKU      string
	Quantity int32
}

type OrderInput struct {
	ID    string
	Items []Item
}

type OrderResult struct{}
