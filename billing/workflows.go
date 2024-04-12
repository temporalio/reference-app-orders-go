package billing

import (
	"fmt"
	"math/rand"

	"go.temporal.io/sdk/workflow"
)

// GenerateInvoice workflow creates an invoice for a fulfillment.
func GenerateInvoice(ctx workflow.Context, input GenerateInvoiceInput) (GenerateInvoiceResult, error) {
	var result GenerateInvoiceResult

	if input.CustomerID == "" {
		return GenerateInvoiceResult{}, fmt.Errorf("CustomerID is required")
	}
	if input.OrderReference == "" {
		return GenerateInvoiceResult{}, fmt.Errorf("OrderReference is required")
	}
	if len(input.Items) == 0 {
		return GenerateInvoiceResult{}, fmt.Errorf("invoice must have items")
	}

	result.InvoiceReference = input.OrderReference

	for _, item := range input.Items {
		cost, tax := calculateCosts(item)
		result.SubTotal += cost
		result.Tax += tax
		result.Shipping += calculateShippingCost(item)
	}

	return result, nil
}

// calculateCosts calculates the cost and tax for an item.
func calculateCosts(item Item) (int32, int32) {
	// This is just a simulation, so make up a cost
	// Normally this would looked up on the SKU
	costPerUnit := rand.Int31n(10000)
	// Return 0 tax for now
	return costPerUnit * int32(item.Quantity), 0
}

// calculateShippingCost calculates the shipping cost for an item.
func calculateShippingCost(item Item) int32 {
	// This is just a simulation, so make up a cost
	// Normally this would looked up on the SKU
	costPerUnit := rand.Int31n(500)
	return costPerUnit * int32(item.Quantity)
}

// ChargeCustomer workflow charges a customer for a fulfillment.
func ChargeCustomer(ctx workflow.Context, input ChargeCustomerInput) (ChargeCustomerResult, error) {
	var result ChargeCustomerResult

	// Return success for now
	result.Success = true
	result.AuthCode = "1234"

	return result, nil
}
