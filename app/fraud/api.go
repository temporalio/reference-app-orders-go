package fraud

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/temporalio/reference-app-orders-go/app/config"
	"github.com/temporalio/reference-app-orders-go/app/util"
)

// FraudLimitInput is the input for the SetLimit API.
type FraudLimitInput struct {
	Limit int32 `json:"limit"`
}

// FraudLimitResult is the result for the GetLimit API.
type FraudLimitResult struct {
	Limit int32 `json:"limit"`
}

// FraudCheckInput is the input for the check endpoint.
type FraudCheckInput struct {
	CustomerID string `json:"customerId"`
	Charge     int32  `json:"charge"`
}

// FraudCheckResult is the result for the check endpoint.
type FraudCheckResult struct {
	Declined bool `json:"declined"`
}

type handlers struct {
	limit               int32
	tallyLock           sync.Mutex
	customerChargeTally map[string]int32
	logger              *slog.Logger
}

// RunServer runs a FraudCheck API HTTP server on the given port.
func RunServer(ctx context.Context, config config.AppConfig) error {
	logger := slog.Default().With("service", "fraud")

	hostPort := fmt.Sprintf("%s:%d", config.BindOnIP, config.FraudPort)
	srv := &http.Server{
		Addr:    hostPort,
		Handler: util.LoggingMiddleware(logger, Router(logger)),
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
func Router(logger *slog.Logger) http.Handler {
	r := http.NewServeMux()
	h := handlers{customerChargeTally: make(map[string]int32), logger: logger}

	r.HandleFunc("GET /limit", h.handleGetLimit)
	r.HandleFunc("POST /limit", h.handleSetLimit)
	r.HandleFunc("POST /reset", h.handleReset)
	r.HandleFunc("POST /check", h.handleRunCheck)

	return r
}

func (h *handlers) handleGetLimit(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(FraudLimitResult{Limit: h.limit})
	if err != nil {
		h.logger.Error("Failed to encode limit result", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *handlers) handleSetLimit(w http.ResponseWriter, r *http.Request) {
	var input FraudLimitInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		h.logger.Error("Failed to decode limit input", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.limit = input.Limit
}

func (h *handlers) handleReset(http.ResponseWriter, *http.Request) {
	h.tallyLock.Lock()
	h.customerChargeTally = make(map[string]int32)
	h.tallyLock.Unlock()

	h.limit = 0
}

func (h *handlers) handleRunCheck(w http.ResponseWriter, r *http.Request) {
	var input FraudCheckInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		h.logger.Error("Failed to decode charge input", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.tallyLock.Lock()
	h.customerChargeTally[input.CustomerID] += input.Charge
	declined := h.limit > 0 && h.customerChargeTally[input.CustomerID] > h.limit
	h.tallyLock.Unlock()
	result := FraudCheckResult{Declined: declined}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		h.logger.Error("Failed to encode charge result", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
