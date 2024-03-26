package api_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/temporalio/orders-reference-app-go/cmd/shipment-api/api"
	"github.com/temporalio/orders-reference-app-go/workflows"
	"go.temporal.io/sdk/mocks"
)

func TestShipmentUpdate(t *testing.T) {
	c := mocks.NewClient(t)

	c.On("SignalWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	r := api.Router(c)
	req, err := http.NewRequest("POST", "/shipments/test/status", strings.NewReader(`{"status":1}`))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, rr.Body.String(), "")

	c.AssertCalled(t,
		"SignalWorkflow",
		mock.Anything,
		"test", "",
		"ShipmentUpdate",
		workflows.ShipmentUpdateSignal{
			Status: workflows.ShipmentStatusDispatched,
		},
	)
}
