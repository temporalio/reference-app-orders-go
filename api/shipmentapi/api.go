package shipmentapi

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/temporalio/orders-reference-app-go/internal/shipmentapi"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/sdk/client"
)

type handlers struct {
	temporal client.Client
}

func Router(c client.Client) *mux.Router {
	r := mux.NewRouter()
	h := handlers{temporal: c}

	r.HandleFunc("/shipments/{id}/status", h.handleShipmentStatus)

	return r
}

func (h *handlers) handleShipmentStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var signal shipmentapi.ShipmentUpdateSignal

	err := json.NewDecoder(r.Body).Decode(&signal)
	if err != nil {
		log.Println("Error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.temporal.SignalWorkflow(context.Background(),
		vars["id"], "",
		shipmentapi.ShipmentUpdateSignalName,
		signal,
	)
	if err != nil {
		switch err.(type) {
		case *serviceerror.NotFound:
			http.Error(w, "Shipment not found", http.StatusNotFound)
		default:
			log.Println("Error: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}
}
