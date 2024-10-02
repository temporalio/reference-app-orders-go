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
	"github.com/temporalio/reference-app-orders-go/app/server"
	"github.com/temporalio/reference-app-orders-go/app/shipment"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.temporal.io/sdk/mocks"
)

func TestShipmentUpdate(t *testing.T) {
	ctx := context.Background()
	c := mocks.NewClient(t)

	c.On("SignalWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	mongoDBContainer, err := mongodb.Run(ctx, "mongo:6")
	require.NoError(t, err)
	defer mongoDBContainer.Terminate(context.Background())

	port, err := mongoDBContainer.MappedPort(context.Background(), "27017/tcp")
	require.NoError(t, err)

	mc, err := mongo.Connect(context.Background(), options.Client().ApplyURI(fmt.Sprintf("mongodb://localhost:%s", port.Port())))
	require.NoError(t, err)

	db := mc.Database("testdb")

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
