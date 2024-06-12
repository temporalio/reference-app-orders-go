package shipment_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/temporalio/reference-app-orders-go/app/server"
	"github.com/temporalio/reference-app-orders-go/app/shipment"
	"go.temporal.io/sdk/mocks"
	_ "modernc.org/sqlite"
)

func TestShipmentUpdate(t *testing.T) {
	c := mocks.NewClient(t)

	c.On("SignalWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	db, err := sqlx.Open("sqlite", ":memory:")
	require.NoError(t, err)
	err = server.SetupDB(db)
	require.NoError(t, err)

	logger := slog.Default()

	r := shipment.Router(c, db, logger)
	req, err := http.NewRequest("POST", "/shipments/test/status", strings.NewReader(`{"status":"dispatched"}`))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, rr.Body.String(), "")

	c.AssertCalled(t,
		"SignalWorkflow",
		mock.Anything,
		shipment.ShipmentWorkflowID("test"), "",
		"ShipmentCarrierUpdate",
		shipment.ShipmentCarrierUpdateSignal{
			Status: shipment.ShipmentStatusDispatched,
		},
	)
}
