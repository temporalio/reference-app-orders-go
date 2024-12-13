//go:build integration

package shipment_test

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/temporalio/reference-app-orders-go/app/config"
	"github.com/temporalio/reference-app-orders-go/app/db"
	"github.com/temporalio/reference-app-orders-go/app/shipment"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.temporal.io/sdk/mocks"
)

func TestShipmentUpdate(t *testing.T) {
	ctx := context.Background()
	c := mocks.NewClient(t)

	c.On("SignalWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mongoDBContainer, err := mongodb.Run(ctx, "mongo:6")
	require.NoError(t, err)
	defer mongoDBContainer.Terminate(ctx)

	port, err := mongoDBContainer.MappedPort(ctx, "27017/tcp")
	require.NoError(t, err)

	uri := fmt.Sprintf("mongodb://localhost:%s", port.Port())

	config := config.AppConfig{MongoURL: uri}

	db := db.CreateDB(config)
	require.NoError(t, db.Connect(ctx))
	require.NoError(t, db.Setup())

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
