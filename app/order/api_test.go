package order_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/temporalio/reference-app-orders-go/app/order"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/mocks"
)

func TestCreateOrderSetsReceivedAt(t *testing.T) {
	for _, test := range []struct{ name, body string }{
		{
			name: "receipt time omitted",
			body: `{"id":"order-1","customerId":"customer-1","items":[{"sku":"test","quantity":1}]}`,
		},
		{
			name: "client receipt time replaced",
			body: `{"id":"order-1","customerId":"customer-1","items":[{"sku":"test","quantity":1}],"receivedAt":"2000-01-01T00:00:00Z"}`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			c := mocks.NewClient(t)
			var receivedAt time.Time
			c.On("ExecuteWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Run(func(args mock.Arguments) {
					options := args.Get(1).(client.StartWorkflowOptions)
					require.Equal(t, order.OrderWorkflowID("order-1"), options.ID)
					require.Equal(t, order.TaskQueue, options.TaskQueue)
					input := args.Get(3).(*order.OrderInput)
					require.Equal(t, "order-1", input.ID)
					require.Equal(t, "customer-1", input.CustomerID)
					require.Equal(t, []*order.Item{{SKU: "test", Quantity: 1}}, input.Items)
					receivedAt = input.ReceivedAt
				}).Return(nil, nil).Once()

			request := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(test.body))
			response := httptest.NewRecorder()
			before := time.Now().UTC()
			order.Router(c, nil, slog.Default()).ServeHTTP(response, request)
			after := time.Now().UTC()

			require.Equal(t, http.StatusCreated, response.Code)
			require.Equal(t, "/orders/order-1", response.Header().Get("Location"))
			require.False(t, receivedAt.Before(before))
			require.False(t, receivedAt.After(after))
			require.Equal(t, time.UTC, receivedAt.Location())
		})
	}
}
