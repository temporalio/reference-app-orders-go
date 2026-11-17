//go:build integration

package test

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/temporalio/reference-app-orders-go/app/config"
	"github.com/temporalio/reference-app-orders-go/app/db"
	"github.com/temporalio/reference-app-orders-go/app/order"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/worker"
)

func TestOrderInsertRecoversAfterWorkerRestart(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 45*time.Second)
	defer cancel()
	temporal, err := testsuite.StartDevServer(ctx, testsuite.DevServerOptions{})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, temporal.Stop()) })
	c := temporal.Client()

	t.Chdir(t.TempDir())
	store := db.CreateDB(config.AppConfig{})
	require.NoError(t, store.Connect(ctx))
	t.Cleanup(func() { require.NoError(t, store.Close()) })
	require.NoError(t, store.Setup())

	committed := make(chan struct{})
	firstResponseLost := make(chan struct{})
	var insertAttempts atomic.Int32
	router := order.Router(c, store, slog.Default())
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/orders/order-1/insert" && insertAttempts.Add(1) == 1 {
			response := httptest.NewRecorder()
			router.ServeHTTP(response, r)
			if response.Code != http.StatusOK {
				http.Error(w, response.Body.String(), response.Code)
				return
			}
			// Keep the response from the first worker after SQLite commits. The worker
			// stops and its activity context cancels the connection before it sees OK.
			close(committed)
			<-r.Context().Done()
			close(firstResponseLost)
			return
		}
		router.ServeHTTP(w, r)
	}))
	t.Cleanup(api.Close)

	startWorker := func(identity string) worker.Worker {
		w := worker.New(c, order.TaskQueue, worker.Options{
			Identity: identity, WorkerStopTimeout: 50 * time.Millisecond,
		})
		w.RegisterWorkflow(order.Order)
		activities := &order.Activities{OrderURL: api.URL}
		w.RegisterActivity(activities.InsertOrder)
		w.RegisterActivity(activities.UpdateOrderStatus)
		// Keep the test at the order persistence boundary by returning no shipments.
		w.RegisterActivityWithOptions(func(_ context.Context, _ *order.ReserveItemsInput) (*order.ReserveItemsResult, error) {
			return &order.ReserveItemsResult{}, nil
		}, activity.RegisterOptions{Name: "ReserveItems"})
		require.NoError(t, w.Start())
		return w
	}
	first := startWorker("orders-original")
	firstStopped := false
	t.Cleanup(func() {
		if !firstStopped {
			first.Stop()
		}
	})

	run, err := c.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		ID: order.OrderWorkflowID("order-1"), TaskQueue: order.TaskQueue, WorkflowTaskTimeout: 2 * time.Second,
	}, order.Order, &order.OrderInput{
		ID: "order-1", CustomerID: "customer-1", Items: []*order.Item{{SKU: "shared-sku", Quantity: 1}},
	})
	require.NoError(t, err)
	select {
	case <-committed:
	case <-ctx.Done():
		t.Fatal("the order did not commit before the worker restart")
	}
	var beforeRestart []db.OrderStatus
	require.NoError(t, store.GetOrders(ctx, &beforeRestart))
	require.Len(t, beforeRestart, 1)
	require.Equal(t, order.OrderStatusPending, beforeRestart[0].Status)

	firstStopped = true
	first.Stop()
	select {
	case <-firstResponseLost:
	case <-ctx.Done():
		t.Fatal("stopping the worker did not cancel the request before acknowledgement")
	}
	second := startWorker("orders-replacement")
	t.Cleanup(second.Stop)
	var result order.OrderResult
	require.NoError(t, run.Get(ctx, &result))
	require.Equal(t, order.OrderStatusCompleted, result.Status)
	require.EqualValues(t, 2, insertAttempts.Load())

	var afterRestart []db.OrderStatus
	require.NoError(t, store.GetOrders(ctx, &afterRestart))
	require.Len(t, afterRestart, 1)
	require.Equal(t, beforeRestart[0].ID, afterRestart[0].ID)
	require.Equal(t, beforeRestart[0].ReceivedAt, afterRestart[0].ReceivedAt)
	require.Equal(t, order.OrderStatusCompleted, afterRestart[0].Status)
}
