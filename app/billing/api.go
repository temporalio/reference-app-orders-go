package billing

// Item represents an item being ordered.
// All fields are required.
type Item struct {
	SKU      string
	Quantity int32
}

// GenerateInvoiceInput is the input for the GenerateInvoice workflow.
// All fields are required.
type GenerateInvoiceInput struct {
	CustomerID     string
	OrderReference string
	Items          []Item
}

// GenerateInvoiceResult is the result for the GenerateInvoice workflow.
type GenerateInvoiceResult struct {
	InvoiceReference string
	SubTotal         int32
	Shipping         int32
	Tax              int32
}

// ChargeCustomerInput is the input for the ChargeCustomer workflow.
// All fields are required.
type ChargeCustomerInput struct {
	CustomerID string
	Reference  string
	Charge     int32
}

// ChargeCustomerResult is the result for the GenerateInvoice workflow.
type ChargeCustomerResult struct {
	Success  bool
	AuthCode string
}
