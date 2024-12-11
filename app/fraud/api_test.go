package fraud_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/temporalio/reference-app-orders-go/app/fraud"
)

func TestMaintenanceMode(t *testing.T) {

	logger := slog.Default()

	r := fraud.Router(logger)

	req, err := http.NewRequest("POST", "/check", strings.NewReader(`{"customer_id":"1","charge":100}`))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	require.Equal(t, rr.Code, http.StatusOK)

	req, err = http.NewRequest("POST", "/maintenance", nil)
	require.NoError(t, err)

	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	require.Equal(t, rr.Code, http.StatusOK)

	req, err = http.NewRequest("POST", "/check", strings.NewReader(`{"customer_id":"1","charge":100}`))
	require.NoError(t, err)

	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusServiceUnavailable, rr.Code)

	req, err = http.NewRequest("POST", "/reset", nil)
	require.NoError(t, err)

	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)

	req, err = http.NewRequest("POST", "/check", strings.NewReader(`{"customer_id":"1","charge":100}`))
	require.NoError(t, err)

	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	require.Equal(t, rr.Code, http.StatusOK)
}
