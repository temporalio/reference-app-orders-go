package shipment

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/temporalio/reference-app-orders-go/app/db"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/sdk/client"
)

// TaskQueue is the default task queue for the Shipment system.
const TaskQueue = "shipments"

// StatusQuery is the name of the query to use to fetch a Shipment's status.
const StatusQuery = "status"

// ShipmentWorkflowID returns the workflow ID for a Shipment.
func ShipmentWorkflowID(id string) string {
	return "Shipment:" + id
}

// ShipmentIDFromWorkflowID returns the ID for a Shipment from a WorkflowID.
func ShipmentIDFromWorkflowID(id string) string {
	return strings.TrimPrefix(id, "Shipment:")
}

type handlers struct {
	temporal client.Client
	db       db.DB
	logger   *slog.Logger
}

// ShipmentStatus holds the status of a Shipment.
type ShipmentStatus struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updatedAt"`
	Items     []Item    `json:"items"`
}

// ShipmentStatusUpdate is used to update the status of a Shipment.
type ShipmentStatusUpdate struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// ListShipmentEntry is an entry in the Shipment list.
type ListShipmentEntry struct {
	ID     string `json:"id" db:"id" bson:"id"`
	Status string `json:"status" db:"id" bson:"status"`
}

// Router implements the http.Handler interface for the Shipment API
func Router(client client.Client, db db.DB, logger *slog.Logger) http.Handler {
	r := http.NewServeMux()

	h := handlers{temporal: client, db: db, logger: logger}

	r.HandleFunc("GET /shipments", h.handleListShipments)
	r.HandleFunc("GET /shipments/{id}", h.handleGetShipment)
	r.HandleFunc("POST /shipments/{id}", h.handleUpdateShipmentStatus)
	r.HandleFunc("POST /shipments/{id}/status", h.handleUpdateShipmentCarrierStatus)

	return r
}

func (h *handlers) handleListShipments(w http.ResponseWriter, _ *http.Request) {
	shipments := []ListShipmentEntry{}

	err := h.db.GetShipments(context.Background(), &shipments)
	if err != nil {
		h.logger.Error("Failed to list shipments: %v", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(shipments); err != nil {
		h.logger.Error("Failed to encode orders: %v", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *handlers) handleGetShipment(w http.ResponseWriter, r *http.Request) {
	var status ShipmentStatus

	q, err := h.temporal.QueryWorkflow(r.Context(),
		ShipmentWorkflowID(r.PathValue("id")), "",
		StatusQuery,
	)
	if err != nil {
		if _, ok := err.(*serviceerror.NotFound); ok {
			h.logger.Error("Failed to query shipment workflow: %v", "error", err)
			http.Error(w, "Shipment not found", http.StatusNotFound)
		} else {
			h.logger.Error("Failed to query shipment workflow: %v", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := q.Get(&status); err != nil {
		h.logger.Error("Failed to get query result: %v", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(status); err != nil {
		h.logger.Error("Failed to encode shipment status: %v", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *handlers) handleUpdateShipmentStatus(w http.ResponseWriter, r *http.Request) {
	var status ShipmentStatusUpdate

	err := json.NewDecoder(r.Body).Decode(&status)
	if err != nil {
		h.logger.Error("Failed to decode shipment status: %v", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.db.UpdateShipmentStatus(context.Background(), status.ID, status.Status)
	if err != nil {
		h.logger.Error("Failed to update shipment status: %v", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handlers) handleUpdateShipmentCarrierStatus(w http.ResponseWriter, r *http.Request) {
	var signal ShipmentCarrierUpdateSignal

	err := json.NewDecoder(r.Body).Decode(&signal)
	if err != nil {
		h.logger.Error("Failed to decode shipment signal: %v", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.temporal.SignalWorkflow(context.Background(),
		ShipmentWorkflowID(r.PathValue("id")), "",
		ShipmentCarrierUpdateSignalName,
		signal,
	)
	if err != nil {
		if _, ok := err.(*serviceerror.NotFound); ok {
			h.logger.Error("Failed to signal shipment workflow: %v", "error", err)
			http.Error(w, "Shipment not found", http.StatusNotFound)
		} else {
			h.logger.Error("Failed to signal shipment workflow: %v", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}
}
