package order

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/sdk/client"
)

// TaskQueue is the default task queue for the Order system.
const TaskQueue = "orders"

// StatusQuery is the name of the query to use to fetch an Order's status.
const StatusQuery = "status"

// OrderWorkflowID returns the workflow ID for an Order.
func OrderWorkflowID(id string) string {
	return "Order:" + id
}

// OrderIDFromWorkflowID returns the ID for an Order from a WorkflowID.
func OrderIDFromWorkflowID(id string) string {
	return strings.TrimPrefix(id, "Order:")
}

// Item represents an item being ordered.
// All fields are required.
type Item struct {
	SKU      string `json:"sku"`
	Quantity int32  `json:"quantity"`
}

// OrderInput is the input for an Order workflow.
type OrderInput struct {
	ID         string  `json:"id"`
	CustomerID string  `json:"customerId"`
	Items      []*Item `json:"items"`
}

// OrderStatus holds the status of an Order workflow.
type OrderStatus struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customerId" db:"customer_id"`
	ReceivedAt time.Time `json:"receivedAt" db:"received_at"`

	Status string `json:"status"`

	Fulfillments []*Fulfillment `json:"fulfillments"`
}

// OrderStatusUpdate is used to update an Order's status.
type OrderStatusUpdate struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

const (
	// OrderStatusPending is the status of a pending Order.
	OrderStatusPending = "pending"

	// OrderStatusProcessing is the status of a processing Order.
	OrderStatusProcessing = "processing"

	// OrderStatusCustomerActionRequired is the status of an Order that requires customer action.
	OrderStatusCustomerActionRequired = "customerActionRequired"

	// OrderStatusCompleted is the status of a completed Order.
	OrderStatusCompleted = "completed"

	// OrderStatusCancelled is the status of a cancelled Order.
	OrderStatusCancelled = "cancelled"
)

// ListOrderEntry is an entry in the Order list.
type ListOrderEntry struct {
	ID         string    `json:"id"`
	Status     string    `json:"status"`
	ReceivedAt time.Time `json:"receivedAt" db:"received_at"`
}

// ShipmentStatus holds the status of a Shipment.
type ShipmentStatus struct {
	ID string `json:"id"`

	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// PaymentStatus holds the status of a Payment.
type PaymentStatus struct {
	SubTotal int32 `json:"subTotal"`
	Tax      int32 `json:"tax"`
	Shipping int32 `json:"shipping"`
	Total    int32 `json:"total"`

	Status string `json:"status"`
}

const (
	// PaymentStatusPending is the status of a pending payment.
	PaymentStatusPending = "pending"

	// PaymentStatusSuccess is the status of a successful payment.
	PaymentStatusSuccess = "success"

	// PaymentStatusFailed is the status of a failed payment.
	PaymentStatusFailed = "failed"
)

// Fulfillment holds a set of items that will be delivered in one shipment (due to location and stock level).
type Fulfillment struct {
	// OrderID is the ID of the order that this fulfillment is part of.
	orderID string

	// CustomerID is the ID of the customer that this fulfillment is for.
	customerID string

	// ID is an identifier for the fulfillment
	ID string `json:"id"`

	// Items is the set of the items that will be part of this shipment.
	Items []*Item `json:"items"`

	// Location is the address for carrier pickup.
	Location string `json:"location,omitempty"`

	// Status is the status of the fulfillment, one of "unavailable", "pending", "processing", "dispatched", "delivered", "failed".
	Status string `json:"status"`

	// PaymentStatus is the status of the payment for this fulfillment.
	Payment *PaymentStatus `json:"payment,omitempty"`

	// ShipmentStatus is the status of the shipment for this fulfillment.
	Shipment *ShipmentStatus `json:"shipment,omitempty"`
}

const (
	// FulfillmentStatusUnavailable is the status of an unavailable Fulfillment.
	FulfillmentStatusUnavailable = "unavailable"

	// FulfillmentStatusPending is the status of a pending Fulfillment.
	FulfillmentStatusPending = "pending"

	// FulfillmentStatusProcessing is the status of a processing Fulfillment.
	FulfillmentStatusProcessing = "processing"

	// FulfillmentStatusCompleted is the status of a processing Fulfillment.
	FulfillmentStatusCompleted = "completed"

	// FulfillmentStatusCancelled is the status of a cancelled Fulfillment.
	FulfillmentStatusCancelled = "cancelled"

	// FulfillmentStatusFailed is the status of a failed Fulfillment.
	FulfillmentStatusFailed = "failed"
)

// CustomerActionSignalName is the name of the signal used to send customer actions.
const CustomerActionSignalName = "CustomerAction"

// CustomerActionSignal is the signal sent to the Fulfillment workflow to indicate a customer action.
type CustomerActionSignal struct {
	Action string `json:"action"`
}

const (
	// CustomerActionCancel is the action to cancel a Fulfillment.
	CustomerActionCancel = "cancel"

	// CustomerActionAmend is the action to amend a Fulfillment.
	CustomerActionAmend = "amend"
)

// OrderResult is the result of an Order workflow.
type OrderResult struct {
	Status string `json:"status"`
}

type handlers struct {
	temporal client.Client
	db       *sqlx.DB
}

// SetupDB creates the necessary tables in the database.
func SetupDB(db *sqlx.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS orders (
		id TEXT PRIMARY KEY,
		customer_id TEXT NOT NULL,
		received_at TIMESTAMP NOT NULL,
		status TEXT NOT NULL
	);

	CREATE INDEX IF NOT EXISTS orders_received_at ON orders(received_at DESC);
	`)
	if err != nil {
		return fmt.Errorf("failed to create orders table: %w", err)
	}

	return nil
}

// RunServer runs a Order API HTTP server on the given port.
func RunServer(ctx context.Context, port int, client client.Client, db *sqlx.DB) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		Handler: Router(client, db),
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
func Router(client client.Client, db *sqlx.DB) *mux.Router {
	r := mux.NewRouter()

	h := handlers{temporal: client, db: db}

	r.HandleFunc("/orders", h.handleCreateOrder).Methods("POST")
	r.HandleFunc("/orders", h.handleListOrders).Methods("GET")
	r.HandleFunc("/orders/{id}", h.handleGetOrder).Methods("GET")
	r.HandleFunc("/orders/{id}/status", h.handleUpdateOrderStatus).Methods("POST")
	r.HandleFunc("/orders/{id}/action", h.handleCustomerAction).Methods("POST")

	return r
}

func (h *handlers) handleListOrders(w http.ResponseWriter, _ *http.Request) {
	orders := []ListOrderEntry{}

	err := h.db.Select(&orders, `SELECT id, status, received_at FROM orders ORDER BY received_at DESC`)
	if err != nil {
		log.Printf("Failed to list orders: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(orders)
	if err != nil {
		log.Printf("Failed to encode orders: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *handlers) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var input OrderInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Printf("Failed to decode order input: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.temporal.ExecuteWorkflow(context.Background(),
		client.StartWorkflowOptions{
			TaskQueue:             TaskQueue,
			ID:                    OrderWorkflowID(input.ID),
			WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE,
		},
		Order,
		&input,
	)
	if err != nil {
		log.Printf("Failed to start order workflow: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	status := &OrderStatus{
		ID:         input.ID,
		CustomerID: input.CustomerID,
		ReceivedAt: time.Now().UTC(),
		Status:     OrderStatusPending,
	}

	_, err = h.db.NamedExec(`INSERT OR IGNORE INTO orders (id, customer_id, received_at, status) VALUES (:id, :customer_id, :received_at, :status)`, status)
	if err != nil {
		log.Printf("Failed to record workflow status: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/orders/"+input.ID)
	w.WriteHeader(http.StatusCreated)
}

func (h *handlers) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var status OrderStatus

	q, err := h.temporal.QueryWorkflow(r.Context(),
		OrderWorkflowID(vars["id"]), "",
		StatusQuery,
	)
	if err != nil {
		if _, ok := err.(*serviceerror.NotFound); ok {
			http.Error(w, "Order not found", http.StatusNotFound)
		} else {
			log.Printf("Failed to query order workflow: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := q.Get(&status); err != nil {
		log.Printf("Failed to get order query result: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(status); err != nil {
		log.Printf("Failed to encode order status: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *handlers) handleUpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	var status OrderStatusUpdate

	err := json.NewDecoder(r.Body).Decode(&status)
	if err != nil {
		log.Printf("Failed to decode order status: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.db.NamedExec(`UPDATE orders SET status = :status WHERE id = :id`, status)
	if err != nil {
		log.Printf("Failed to update order status: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handlers) handleCustomerAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var signal CustomerActionSignal

	err := json.NewDecoder(r.Body).Decode(&signal)
	if err != nil {
		log.Printf("Failed to decode customer action signal: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.temporal.SignalWorkflow(context.Background(),
		OrderWorkflowID(vars["id"]), "",
		CustomerActionSignalName,
		signal,
	)
	if err != nil {
		if _, ok := err.(*serviceerror.NotFound); ok {
			log.Printf("Failed to signal order workflow: %v", err)
			http.Error(w, "Order not found", http.StatusNotFound)
		} else {
			log.Printf("Failed to signal order workflow: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}
}
