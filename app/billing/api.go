package billing

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/temporalio/reference-app-orders-go/app/config"
	"github.com/temporalio/reference-app-orders-go/app/util"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
)

// TaskQueue is the default task queue for the Billing system.
const TaskQueue = "billing"

// Item represents an item being ordered.
type Item struct {
	SKU      string `json:"sku"`
	Quantity int32  `json:"quantity"`
}

// ChargeInput is the input for the Charge workflow.
type ChargeInput struct {
	CustomerID     string `json:"customerId"`
	Reference      string `json:"orderReference"`
	Items          []Item `json:"items"`
	IdempotencyKey string `json:"idempotencyKey,omitempty"`
}

// ChargeResult is the result for the Charge workflow.
type ChargeResult struct {
	InvoiceReference string `json:"invoiceReference"`
	SubTotal         int32  `json:"subTotal"`
	Shipping         int32  `json:"shipping"`
	Tax              int32  `json:"tax"`
	Total            int32  `json:"total"`

	Success  bool   `json:"success"`
	AuthCode string `json:"authCode"`
}

// GenerateInvoiceInput is the input for the GenerateInvoice activity.
type GenerateInvoiceInput struct {
	CustomerID string `json:"customerId"`
	Reference  string `json:"orderReference"`
	Items      []Item `json:"items"`
}

// GenerateInvoiceResult is the result for the GenerateInvoice activity.
type GenerateInvoiceResult struct {
	InvoiceReference string `json:"invoiceReference"`
	SubTotal         int32  `json:"subTotal"`
	Shipping         int32  `json:"shipping"`
	Tax              int32  `json:"tax"`
	Total            int32  `json:"total"`
}

// ChargeCustomerInput is the input for the ChargeCustomer activity.
type ChargeCustomerInput struct {
	CustomerID string `json:"customerId"`
	Reference  string `json:"reference"`
	Charge     int32  `json:"charge"`
}

// ChargeCustomerResult is the result for the GenerateInvoice activity.
type ChargeCustomerResult struct {
	Success  bool   `json:"success"`
	AuthCode string `json:"authCode"`
}

type handlers struct {
	temporal client.Client
	logger   *slog.Logger
}

// RunServer runs a Billing API HTTP server on the given port.
func RunServer(ctx context.Context, config config.AppConfig, client client.Client) error {
	logger := slog.Default().With("service", "billing")

	hostPort := fmt.Sprintf("%s:%d", config.BindOnIP, config.BillingPort)
	srv := &http.Server{
		Addr:    hostPort,
		Handler: util.LoggingMiddleware(logger, Router(client, logger)),
	}

	logger.Info("Listening", "endpoint", "http://"+hostPort)

	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()

	select {
	case <-ctx.Done():
		srv.Close()
	case err := <-errCh:
		return err
	}

	return nil
}

// Router implements the http.Handler interface for the Billing API
func Router(c client.Client, logger *slog.Logger) http.Handler {
	r := http.NewServeMux()
	h := handlers{temporal: c, logger: logger}

	r.HandleFunc("POST /charge", h.handleCharge)

	return r
}

// ChargeWorkflowID returns the workflow ID for a Charge workflow.
func ChargeWorkflowID(input ChargeInput) string {
	// If an idempotency key is provided, use it as the workflow ID.
	// This ensures that the same charge is not processed multiple times.
	key := input.IdempotencyKey
	if key == "" {
		// If no idempotency key is provided, generate a random one.
		// This will not offer any idempotency guarantees.
		key = uuid.NewString()
	}

	return fmt.Sprintf("Charge:%s", key)
}

func (h *handlers) handleCharge(w http.ResponseWriter, r *http.Request) {
	var input ChargeInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		h.logger.Error("Failed to decode charge input", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Start the Charge workflow.
	// If the workflow is already running, or is finished but still within retention period, this will return the existing workflow.
	// If an idempotency key was provided, this provides idempotency guarantees for the Charge operation.
	wf, err := h.temporal.ExecuteWorkflow(context.Background(),
		client.StartWorkflowOptions{
			TaskQueue:             TaskQueue,
			ID:                    ChargeWorkflowID(input),
			WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE,
		},
		Charge,
		&input,
	)
	if err != nil {
		h.logger.Error("Failed to start charge workflow", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result ChargeResult
	err = wf.Get(r.Context(), &result)
	if err != nil {
		h.logger.Error("Failed to get charge result", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		h.logger.Error("Failed to encode charge result", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
