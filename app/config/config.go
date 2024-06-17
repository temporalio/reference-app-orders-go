package config

import (
	"fmt"
	"os"
	"strconv"
)

// AppConfig is a struct that holds the configuration for the Order/Shipment/Fraud/Billing system.
type AppConfig struct {
	BindOnIP     string
	DataDir      string
	BillingPort  int32
	BillingURL   string
	OrderPort    int32
	OrderURL     string
	ShipmentPort int32
	ShipmentURL  string
	FraudPort    int32
	FraudURL     string
}

// ServiceHostPort returns the host:port for a given service.
func (c *AppConfig) ServiceHostPort(service string) (string, error) {
	var port int32

	switch service {
	case "billing":
		port = c.BillingPort
	case "fraud":
		port = c.FraudPort
	case "order":
		port = c.OrderPort
	case "shipment":
		port = c.ShipmentPort
	default:
		return "", fmt.Errorf("unknown service: %s", service)
	}

	return fmt.Sprintf("%s:%d", c.BindOnIP, port), nil
}

// AppConfigFromEnv creates an AppConfig from environment variables.
func AppConfigFromEnv() (AppConfig, error) {
	conf := AppConfig{
		BindOnIP:     "127.0.0.1",
		DataDir:      "./",
		BillingPort:  8081,
		BillingURL:   "http://127.0.0.1:8081",
		OrderPort:    8082,
		OrderURL:     "http://127.0.0.1:8082",
		ShipmentPort: 8083,
		ShipmentURL:  "http://127.0.0.1:8083",
		FraudPort:    8084,
		FraudURL:     "http://127.0.0.1:8084",
	}

	if ip := os.Getenv("BIND_ON_IP"); ip != "" {
		conf.BindOnIP = ip
	}

	if p := os.Getenv("DATA_DIR"); p != "" {
		conf.DataDir = p
	}

	if p := os.Getenv("BILLING_API_URL"); p != "" {
		conf.BillingURL = p
	}

	if p := os.Getenv("BILLING_API_PORT"); p != "" {
		v, err := strconv.Atoi(p)
		if err != nil {
			return conf, err
		}
		conf.BillingPort = int32(v)
	}

	if p := os.Getenv("ORDER_API_URL"); p != "" {
		conf.OrderURL = p
	}

	if p := os.Getenv("ORDER_API_PORT"); p != "" {
		v, err := strconv.Atoi(p)
		if err != nil {
			return conf, err
		}
		conf.OrderPort = int32(v)
	}

	if p := os.Getenv("SHIPMENT_API_URL"); p != "" {
		conf.ShipmentURL = p
	}

	if p := os.Getenv("SHIPMENT_API_PORT"); p != "" {
		v, err := strconv.Atoi(p)
		if err != nil {
			return conf, err
		}
		conf.ShipmentPort = int32(v)
	}

	if p := os.Getenv("FRAUD_API_URL"); p != "" {
		conf.FraudURL = p
	}

	if p := os.Getenv("FRAUD_API_PORT"); p != "" {
		v, err := strconv.Atoi(p)
		if err != nil {
			return conf, err
		}
		conf.FraudPort = int32(v)
	}

	return conf, nil
}
