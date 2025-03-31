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
	Status string `json:"status" db:"status" bson:"status"`
}

// ShipmentStatsResult holds the stats for the Shipment system.
type ShipmentStatsResult struct {
	WorkerCount int64 `json:"workerCount"`
	Backlog     int64 `json:"backlog"`
}

// Router implements the http.Handler interface for the Shipment API
func Router(client client.Client, db db.DB, logger *slog.Logger) http.Handler {
	r := http.NewServeMux()

	h := handlers{temporal: client, db: db, logger: logger}

	r.HandleFunc("GET /shipments", h.handleListShipments)
	r.HandleFunc("GET /shipments/pending", h.handleListPendingShipments)
	r.HandleFunc("GET /shipments/stats", h.handleGetStats)
	r.HandleFunc("GET /shipments/{id}", h.handleGetShipment)
	r.HandleFunc("POST /shipments/{id}", h.handleUpdateShipmentStatus)
	r.HandleFunc("POST /shipments/{id}/status", h.handleUpdateShipmentCarrierStatus)

	return r
}

func (h *handlers) handleListShipments(w http.ResponseWriter, _ *http.Request) {
	shipments := []db.ShipmentStatus{}

	err := h.db.GetShipments(context.Background(), &shipments)
	if err != nil {
		h.logger.Error("Failed to list shipments: %v", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	list := make([]ListShipmentEntry, len(shipments))
	for i, s := range shipments {
		list[i] = ListShipmentEntry{
			ID:     s.ID,
			Status: s.Status,
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(list); err != nil {
		h.logger.Error("Failed to encode orders: %v", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *handlers) handleListPendingShipments(w http.ResponseWriter, _ *http.Request) {
	shipments := []db.ShipmentStatus{}

	err := h.db.GetPendingShipments(context.Background(), &shipments)
	if err != nil {
		h.logger.Error("Failed to list shipments: %v", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	list := make([]ListShipmentEntry, len(shipments))
	for i, s := range shipments {
		list[i] = ListShipmentEntry{
			ID:     s.ID,
			Status: s.Status,
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(list); err != nil {
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

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(ShipmentStatsResult{
		WorkerCount: int64(workerCount),
		Backlog:     backlog,
	})
	if err != nil {
		h.logger.Error("Failed to encode backlog", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
