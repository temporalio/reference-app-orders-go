package fraudcheck

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/temporalio/orders-reference-app-go/app/config"
)

// FraudLimitInput is the input for the SetLimit API.
type FraudLimitInput struct {
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
}

// RunServer runs a FraudCheck API HTTP server on the given port.
func RunServer(ctx context.Context, config config.AppConfig) error {
	hostPort := fmt.Sprintf("%s:%d", config.BindOnIP, config.FraudPort)
	srv := &http.Server{
		Addr:    hostPort,
		Handler: Router(),
	}

	fmt.Printf("Listening on http://%s\n", hostPort)

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
func Router() *mux.Router {
	r := mux.NewRouter()
	h := handlers{customerChargeTally: make(map[string]int32)}

	r.HandleFunc("/limit", h.handleGetLimit).Methods("GET")
	r.HandleFunc("/limit", h.handleSetLimit).Methods("POST")
	r.HandleFunc("/reset", h.handleReset).Methods("POST")
	r.HandleFunc("/check", h.handleRunCheck).Methods("POST")

	return r
}

func (h *handlers) handleGetLimit(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(FraudLimitInput{Limit: h.limit})
	if err != nil {
		log.Printf("Failed to encode limit: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *handlers) handleSetLimit(w http.ResponseWriter, r *http.Request) {
	var input FraudLimitInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Printf("Failed to decode limit input: %v", err)
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
		log.Printf("Failed to decode charge input: %v", err)
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
		log.Printf("Failed to encode charge result: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
