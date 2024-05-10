package shipment_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/temporalio/orders-reference-app-go/app/shipment"
	"go.temporal.io/sdk/mocks"
)

func TestShipmentUpdate(t *testing.T) {
	c := mocks.NewClient(t)

	c.On("SignalWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	r := shipment.Router(c)
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
