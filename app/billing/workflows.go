package billing

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

// Charge Workflow invoices and processes payment for a fulfillment.
func Charge(ctx workflow.Context, input *ChargeInput) (*ChargeResult, error) {
	logger := workflow.GetLogger(ctx)
	ctx = workflow.WithActivityOptions(ctx,
		workflow.ActivityOptions{
			ScheduleToCloseTimeout: 30 * time.Second,
		},
	)

	var invoice GenerateInvoiceResult

	cwf := workflow.ExecuteActivity(ctx,
		a.GenerateInvoice,
		GenerateInvoiceInput{
			CustomerID: input.CustomerID,
			Reference:  input.Reference,
			Items:      input.Items,
		},
	)
	err := cwf.Get(ctx, &invoice)
	if err != nil {
		return nil, err
	}

	var charge ChargeCustomerResult

	cwf = workflow.ExecuteActivity(ctx,
		a.ChargeCustomer,
		ChargeCustomerInput{
			CustomerID: input.CustomerID,
			Reference:  invoice.InvoiceReference,
			Charge:     invoice.Total,
		},
	)
	if err := cwf.Get(ctx, &charge); err != nil {
		logger.Warn("Charge failed", "customer_id", input.CustomerID, "error", err)
		charge.Success = false
	}

	return &ChargeResult{
		InvoiceReference: invoice.InvoiceReference,
		SubTotal:         invoice.SubTotal,
		Tax:              invoice.Tax,
		Shipping:         invoice.Shipping,
		Total:            invoice.Total,

		Success:  charge.Success,
		AuthCode: charge.AuthCode,
	}, nil
}
