package billing

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
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

const statsInterval = 30

// ChargeStatsResult holds the stats for the Charge system.
type ChargeStatsResult struct {
	WorkerCount  int64   `json:"workerCount"`
	CompleteRate float64 `json:"completeRate"`
	Backlog      int64   `json:"backlog"`
}

type handlers struct {
	temporal client.Client
	logger   *slog.Logger
}

// Router implements the http.Handler interface for the Billing API
func Router(c client.Client, logger *slog.Logger) http.Handler {
	r := http.NewServeMux()
	h := handlers{temporal: c, logger: logger}

	r.HandleFunc("GET /charge/stats", h.handleGetStats)
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

func (h *handlers) handleGetStats(w http.ResponseWriter, _ *http.Request) {
	resp, err := h.temporal.DescribeTaskQueueEnhanced(context.Background(), client.DescribeTaskQueueEnhancedOptions{
		TaskQueue:     TaskQueue,
		ReportStats:   true,
		ReportPollers: true,
	})
	if err != nil {
		h.logger.Error("Failed to get task queue stats", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	recentPollerWindow := time.Now().Add(-1 * time.Minute)

	var backlog int64
	var workerCount int
	for _, versionInfo := range resp.VersionsInfo {
		for _, typeInfo := range versionInfo.TypesInfo {
			for _, pollerInfo := range typeInfo.Pollers {
				if pollerInfo.LastAccessTime.After(recentPollerWindow) {
					workerCount++
				}
			}
			if typeInfo.Stats != nil {
				backlog += typeInfo.Stats.ApproximateBacklogCount
			}
		}
	}

	closedSince := time.Now().Add(-statsInterval * time.Second).Format(time.RFC3339)
	countResp, err := h.temporal.CountWorkflow(context.Background(), &workflowservice.CountWorkflowExecutionsRequest{
		Query: fmt.Sprintf("WorkflowType='Charge' AND ExecutionStatus='Completed' AND CloseTime > %q", closedSince),
	})
	if err != nil {
		h.logger.Error("Failed to count workflows", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var completeRate float64
	if countResp.GetCount() > 0 {
		completeRate = float64(countResp.GetCount()) / float64(statsInterval)
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(ChargeStatsResult{
		WorkerCount:  int64(workerCount),
		CompleteRate: completeRate,
		Backlog:      backlog,
	})
	if err != nil {
		h.logger.Error("Failed to encode backlog", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
