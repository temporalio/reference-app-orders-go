package order

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/temporalio/orders-reference-app-go/app/internal/temporalutil"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

// TaskQueue is the default task queue for the Order system.
const TaskQueue = "orders"

// StatusQuery is the name of the query to use to fetch an Order's status.
const StatusQuery = "status"

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
	ID           string         `json:"id"`
	CustomerID   string         `json:"customerId"`
	Fulfillments []*Fulfillment `json:"fulfillments"`
}

// ListOrderEntry is an entry in the Order list.
type ListOrderEntry struct {
	ID        string    `json:"id"`
	StartedAt time.Time `json:"startedAt"`
}

// ShipmentStatus holds the status of a ShipmentStatus.
type ShipmentStatus struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// PaymentStatus holds the status of a PaymentStatus.
type PaymentStatus struct {
	// Status is the status of the payment.
	Status   string `json:"status"`
	SubTotal int32  `json:"subTotal"`
	Tax      int32  `json:"tax"`
	Shipping int32  `json:"shipping"`
	Total    int32  `json:"total"`
}

// Fulfillment holds a set of items that will be delivered in one shipment (due to location and stock level).
type Fulfillment struct {
	// ID is an identifier for the fulfillment
	ID string `json:"id"`
	// Items is the set of the items that will be part of this shipment.
	Items []*Item `json:"items"`

	// Payment stores details of the payment
	Payment *PaymentStatus `json:"payment"`
	// Shipment stores details of the shipment
	Shipment *ShipmentStatus `json:"shipment"`

	// location is the address for courier pickup (the warehouse).
	Location string
}

// OrderResult is the result of an Order workflow.
type OrderResult struct{}

type handlers struct {
	temporal client.Client
}

// RunServer runs a Order API HTTP server on the given port.
func RunServer(ctx context.Context, port int) error {
	clientOptions, err := temporalutil.CreateClientOptionsFromEnv()
	if err != nil {
		return fmt.Errorf("failed to create client options: %v", err)
	}

	c, err := client.Dial(clientOptions)
	if err != nil {
		return fmt.Errorf("client error: %v", err)
	}
	defer c.Close()

	srv := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		Handler: Router(c),
	}

	fmt.Printf("Listening on http://127.0.0.1:%d\n", port)

	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()

	select {
	case <-ctx.Done():
		srv.Close()
	case err = <-errCh:
		return err
	}

	return nil
}

// Router implements the http.Handler interface for the Billing API
func Router(c client.Client) *mux.Router {
	r := mux.NewRouter()
	h := handlers{temporal: c}

	r.HandleFunc("/orders", h.handleCreateOrder).Methods("POST")
	r.HandleFunc("/orders", h.handleListOrders).Methods("GET")
	r.HandleFunc("/orders/{id}", h.handleGetOrder)

	return r
}

func (h *handlers) handleListOrders(w http.ResponseWriter, r *http.Request) {
	orders := []ListOrderEntry{}
	var nextPageToken []byte

	for {
		resp, err := h.temporal.ListWorkflow(r.Context(), &workflowservice.ListWorkflowExecutionsRequest{
			PageSize:      10,
			NextPageToken: nextPageToken,
			Query:         "WorkflowType='Order' AND ExecutionStatus='Running'",
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, e := range resp.Executions {
			orders = append(orders, ListOrderEntry{ID: e.GetExecution().GetWorkflowId(), StartedAt: e.GetStartTime().AsTime()})
		}

		if len(resp.NextPageToken) == 0 {
			break
		}

		nextPageToken = resp.NextPageToken
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *handlers) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var input OrderInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println("Error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.temporal.ExecuteWorkflow(context.Background(),
		client.StartWorkflowOptions{
			TaskQueue: TaskQueue,
			ID:        input.ID,
		},
		Order,
		&input,
	)
	if err != nil {
		log.Println("Error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *handlers) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var status OrderStatus

	q, err := h.temporal.QueryWorkflow(r.Context(),
		vars["id"], "",
		StatusQuery,
	)
	if err != nil {
		switch err.(type) {
		case *serviceerror.NotFound:
			http.Error(w, "Order not found", http.StatusNotFound)
		default:
			log.Println("Error: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := q.Get(&status); err != nil {
		log.Println("Error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(status); err != nil {
		log.Println("Error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
