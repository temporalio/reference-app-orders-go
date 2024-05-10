package fraudcheck

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
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
	Approved bool `json:"approved"`
}

type handlers struct {
	limit               int32
	tallyLock           sync.Mutex
	customerChargeTally map[string]int32
}

// RunServer runs a FraudCheck API HTTP server on the given port.
func RunServer(ctx context.Context, port int) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		Handler: Router(),
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
func Router() *mux.Router {
	r := mux.NewRouter()
	h := handlers{customerChargeTally: make(map[string]int32)}

	r.HandleFunc("/limit", h.handleSetLimit).Methods("POST")
	r.HandleFunc("/reset", h.handleReset).Methods("POST")
	r.HandleFunc("/check", h.handleRunCheck).Methods("POST")

	return r
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
	approved := h.limit == 0 || h.customerChargeTally[input.CustomerID] < h.limit
	h.tallyLock.Unlock()
	result := FraudCheckResult{Approved: approved}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Printf("Failed to encode charge result: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
