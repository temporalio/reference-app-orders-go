package shipment

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/sdk/client"
)

// TaskQueue is the default task queue for the Shipment system.
const TaskQueue = "shipments"

// StatusQuery is the name of the query to use to fetch a Shipment's status.
const StatusQuery = "status"

// ShipmentCollection is the name of the MongoDB collection to use for Shipment data.
const ShipmentCollection = "shipments"

// ShipmentWorkflowID returns the workflow ID for a Shipment.
func ShipmentWorkflowID(id string) string {
	return "Shipment:" + id
}

// ShipmentIDFromWorkflowID returns the ID for a Shipment from a WorkflowID.
func ShipmentIDFromWorkflowID(id string) string {
	return strings.TrimPrefix(id, "Shipment:")
}

type handlers struct {
	temporal  client.Client
	shipments *mongo.Collection
	logger    *slog.Logger
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
	ID     string `json:"id" bson:"id"`
	Status string `json:"status" bson:"status"`
}

// Router implements the http.Handler interface for the Shipment API
func Router(client client.Client, db *mongo.Database, logger *slog.Logger) http.Handler {
	r := http.NewServeMux()

	shipments := db.Collection(ShipmentCollection)

	h := handlers{temporal: client, shipments: shipments, logger: logger}

	r.HandleFunc("GET /shipments", h.handleListShipments)
	r.HandleFunc("GET /shipments/{id}", h.handleGetShipment)
	r.HandleFunc("POST /shipments/{id}", h.handleUpdateShipmentStatus)
	r.HandleFunc("POST /shipments/{id}/status", h.handleUpdateShipmentCarrierStatus)

	return r
}

func (h *handlers) handleListShipments(w http.ResponseWriter, _ *http.Request) {
	shipments := []ListShipmentEntry{}

	res, err := h.shipments.Find(context.Background(), nil)
	if err != nil {
		h.logger.Error("Failed to list shipments: %v", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = res.All(context.TODO(), &shipments)
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

	_, err = h.shipments.UpdateOne(
		context.Background(),
		bson.M{"id": status.ID},
		bson.M{"$set": bson.M{"status": status.Status}, "$setOnInsert": bson.M{"booked_at": time.Now().UTC()}},
		options.Update().SetUpsert(true),
	)
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
