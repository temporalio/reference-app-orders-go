package shipment

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	"github.com/temporalio/orders-reference-app-go/app/internal/temporalutil"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/sdk/client"
)

type handlers struct {
	temporal client.Client
}

func Server(port int) error {
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
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		Handler: Router(c),
	}

	fmt.Printf("Listening on http://0.0.0.0:%d\n", port)

	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	select {
	case <-sigCh:
		srv.Close()
	case err = <-errCh:
		return err
	}

	return nil
}

// Router implements the http.Handler interface for the Shipment API
func Router(c client.Client) *mux.Router {
	r := mux.NewRouter()
	h := handlers{temporal: c}

	r.HandleFunc("/shipments/{id}/status", h.handleShipmentStatus)

	return r
}

func (h *handlers) handleShipmentStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var signal ShipmentUpdateSignal

	err := json.NewDecoder(r.Body).Decode(&signal)
	if err != nil {
		log.Println("Error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.temporal.SignalWorkflow(context.Background(),
		vars["id"], "",
		ShipmentUpdateSignalName,
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
