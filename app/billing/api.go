package billing

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
}

// RunServer runs a Billing API HTTP server on the given port.
func RunServer(ctx context.Context, port int, client client.Client) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		Handler: Router(client),
	}

	fmt.Printf("Listening on http://127.0.0.1:%d\n", port)

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
func Router(c client.Client) *mux.Router {
	r := mux.NewRouter()
	h := handlers{temporal: c}

	r.HandleFunc("/charge", h.handleCharge)

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
		log.Printf("Failed to decode charge input: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Start the Charge workflow.
	// If the workflow is already running, this will return the existing workflow.
	// If an idempotency key was provided, this provides idempotency guarantees for the Charge operation.
	wf, err := h.temporal.ExecuteWorkflow(context.Background(),
		client.StartWorkflowOptions{
			TaskQueue: TaskQueue,
			ID:        ChargeWorkflowID(input),
		},
		Charge,
		&input,
	)
	if err != nil {
		log.Printf("Failed to start charge workflow: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result ChargeResult
	err = wf.Get(r.Context(), &result)
	if err != nil {
		log.Printf("Failed to get charge result: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Printf("Failed to encode charge result: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
