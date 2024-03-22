package shipmentapi_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	api "github.com/temporalio/orders-reference-app-go/api/shipmentapi"
	"github.com/temporalio/orders-reference-app-go/internal/shipmentapi"
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
		shipmentapi.ShipmentUpdateSignal{
			Status: shipmentapi.ShipmentStatusDispatched,
		},
	)
}
