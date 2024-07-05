package fraud

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
)

// FraudLimitInput is the input for the SetLimit API.
type FraudLimitInput struct {
	Limit int32 `json:"limit"`
}

// FraudSettingsResult is the result for the GetSettings API.
type FraudSettingsResult struct {
	Limit           int32 `json:"limit"`
	MaintenanceMode bool  `json:"maintenanceMode"`
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
	maintenanceMode     bool
	tallyLock           sync.Mutex
	customerChargeTally map[string]int32
	logger              *slog.Logger
}

// Router implements the http.Handler interface for the Billing API
func Router(logger *slog.Logger) http.Handler {
	r := http.NewServeMux()
	h := handlers{customerChargeTally: make(map[string]int32), logger: logger}

	r.HandleFunc("GET /settings", h.handleGetSettings)
	r.HandleFunc("POST /limit", h.handleSetLimit)
	r.HandleFunc("POST /maintenance", h.handleSetMaintenanceMode)
	r.HandleFunc("POST /reset", h.handleReset)
	r.HandleFunc("POST /check", h.handleRunCheck)

	return r
}

func (h *handlers) handleGetSettings(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(FraudSettingsResult{
		Limit:           h.limit,
		MaintenanceMode: h.maintenanceMode,
	})
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
	h.maintenanceMode = false
}

func (h *handlers) handleSetMaintenanceMode(w http.ResponseWriter, _ *http.Request) {
	h.maintenanceMode = true

	w.WriteHeader(http.StatusOK)
}

func (h *handlers) handleRunCheck(w http.ResponseWriter, r *http.Request) {
	var input FraudCheckInput

	if h.maintenanceMode {
		http.Error(w, "Fraud service is in maintenance mode", http.StatusServiceUnavailable)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		h.logger.Error("Failed to decode charge input", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.tallyLock.Lock()
	declined := h.limit > 0 && input.Charge+h.customerChargeTally[input.CustomerID] > h.limit
	if !declined {
		h.customerChargeTally[input.CustomerID] += input.Charge
	}
	h.tallyLock.Unlock()
	result := FraudCheckResult{Declined: declined}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		h.logger.Error("Failed to encode charge result", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
