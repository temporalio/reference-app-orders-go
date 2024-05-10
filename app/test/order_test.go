package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/temporalio/orders-reference-app-go/app/order"
	"github.com/temporalio/orders-reference-app-go/app/server"
	"github.com/temporalio/orders-reference-app-go/app/shipment"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/testsuite"
)

func postJSON(url string, input interface{}) (*http.Response, error) {
	jsonInput, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("unable to encode input: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonInput))
	if err != nil {
		return nil, fmt.Errorf("unable to build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	return client.Do(req)
}

func getJSON(url string, result interface{}) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to build request: %w", err)
	}

	client := http.DefaultClient
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode >= 200 && r.StatusCode < 300 {
		err = json.NewDecoder(r.Body).Decode(result)
		return r, err
	}

	message, _ := io.ReadAll(r.Body)

	return r, fmt.Errorf("%s: %s", http.StatusText(r.StatusCode), message)
}

func Test_Order(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	s, err := testsuite.StartDevServer(ctx, testsuite.DevServerOptions{
		ClientOptions: &client.Options{},
		EnableUI:      true,
		ExtraArgs:     []string{"--dynamic-config-value", "system.forceSearchAttributesCacheRefreshOnRead=true"},
	})
	require.NoError(t, err)

	var (
		c client.Client
	)

	options := client.Options{
		HostPort:  s.FrontendHostPort(),
		Namespace: "default",
	}

	c, err = client.Dial(options)
	require.NoError(t, err)
	defer c.Close()

	err = shipment.EnsureValidTemporalEnv(ctx, options)
	require.NoError(t, err)

	go func() {
		_ = server.RunServer(ctx, c)
	}()

	res, err := postJSON("http://localhost:8083/orders", &order.OrderInput{
		ID:         "order123",
		CustomerID: "customer123",
		Items: []*order.Item{
			{SKU: "Adidas Classic", Quantity: 1},
			{SKU: "Nike Air", Quantity: 2},
		},
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode)

	assert.EventuallyWithT(t, func(c *assert.CollectT) {
		var o order.OrderStatus
		res, err = getJSON("http://localhost:8083/orders/order123", &o)
		require.NoError(t, err)

		assert.Equal(c, "customerActionRequired", o.Status)
	}, 3*time.Second, 100*time.Millisecond)

	res, err = postJSON("http://localhost:8083/orders/order123/action", &order.CustomerActionSignal{
		Action: "amend",
	})
	require.Equal(t, http.StatusOK, res.StatusCode)

	assert.EventuallyWithT(t, func(c *assert.CollectT) {
		var o order.OrderStatus
		res, err := getJSON("http://localhost:8083/orders/order123", &o)
		require.NoError(t, err)

		require.Equal(c, http.StatusOK, res.StatusCode)
		assert.NotNil(c, o.Fulfillments[0].Shipment)
	}, 3*time.Second, 100*time.Millisecond)

	var o order.OrderStatus
	res, err = getJSON("http://localhost:8083/orders/order123", &o)
	require.NoError(t, err)

	for _, f := range o.Fulfillments {
		res, err := postJSON("http://localhost:8081/shipments/"+f.Shipment.ID+"/status", &shipment.ShipmentCarrierUpdateSignal{Status: "delivered"})
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.NoError(t, err)
	}

	assert.EventuallyWithT(t, func(c *assert.CollectT) {
		var o order.OrderStatus
		res, err = getJSON("http://localhost:8083/orders/order123", &o)
		require.NoError(t, err)

		require.Equal(c, http.StatusOK, res.StatusCode)
		assert.Equal(c, "completed", o.Status)
	}, 3*time.Second, 100*time.Millisecond)

	cancel()

	err = s.Stop()
	assert.NoError(t, err)
}
