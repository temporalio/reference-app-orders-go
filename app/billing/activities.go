package billing

import (
	"context"
	"fmt"
	"math/rand"

	"go.temporal.io/sdk/activity"
)

// GenerateInvoice activity creates an invoice for a fulfillment.
func GenerateInvoice(_ context.Context, input *GenerateInvoiceInput) (*GenerateInvoiceResult, error) {
	var result GenerateInvoiceResult

	if input.CustomerID == "" {
		return nil, fmt.Errorf("CustomerID is required")
	}
	if input.Reference == "" {
		return nil, fmt.Errorf("OrderReference is required")
	}
	if len(input.Items) == 0 {
		return nil, fmt.Errorf("invoice must have items")
	}

	result.InvoiceReference = input.Reference

	for _, item := range input.Items {
		cost, tax := calculateCosts(item)
		result.SubTotal += cost
		result.Tax += tax
		result.Shipping += calculateShippingCost(item)
		result.Total += result.SubTotal + result.Tax + result.Shipping
	}

	return &result, nil
}

// calculateCosts calculates the cost and tax for an item.
func calculateCosts(item Item) (cost int32, tax int32) {
	// This is just a simulation, so make up a cost
	// Normally this would be looked up on the SKU
	costPerUnit := rand.Int31n(10000)
	// Return tax at 20%
	return costPerUnit * int32(item.Quantity), costPerUnit * int32(item.Quantity) / 5
}

// calculateShippingCost calculates the shipping cost for an item.
func calculateShippingCost(item Item) int32 {
	// This is just a simulation, so make up a cost
	// Normally this would looked up on the SKU
	costPerUnit := rand.Int31n(500)
	return costPerUnit * int32(item.Quantity)
}

// ChargeCustomer activity charges a customer for a fulfillment.
func ChargeCustomer(ctx context.Context, input *ChargeCustomerInput) (*ChargeCustomerResult, error) {
	var result ChargeCustomerResult

	activity.GetLogger(ctx).Info(
		"Charge",
		"Customer", input.CustomerID,
		"Amount", input.Charge,
		"Reference", input.Reference,
	)

	// Return success for now
	result.Success = true
	result.AuthCode = "1234"

	return &result, nil
}
