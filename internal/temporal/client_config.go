package temporal

import (
	"crypto/tls"
	"fmt"
	"os"

	"go.temporal.io/sdk/client"
)

// ClientOptionsFromEnv creates a client.Options instance, configures
// it based on environment variables, and returns that instance. It
// supports the following environment variables:
//
//	TEMPORAL_ADDRESS: Host and port (formatted as host:port) of the Temporal Frontend Service
//	TEMPORAL_NAMESPACE: Namespace to be used by the Client
//	TEMPORAL_TLS_CERT: Path to the x509 certificate
//	TEMPORAL_TLS_KEY: Path to the private certificate key
//
// If none of these environment variables are set, this will return a
// ClientOptions instance configured to connect to port 7233 on the
// local machine, without TLS, and using Namespace 'default'
func CreateClientOptionsFromEnv() (client.Options, error) {
	var clientOpts client.Options = client.Options{}

	if isSet("TEMPORAL_TLS_CERT") != isSet("TEMPORAL_TLS_KEY") {
		msg := "client cert and key are both required when using TLS"
		return clientOpts, fmt.Errorf(msg)
	}

	clientOpts.HostPort = getEnvWithFallback("TEMPORAL_ADDRESS", "localhost:7233")
	clientOpts.Namespace = getEnvWithFallback("TEMPORAL_NAMESPACE", "default")

	// Other TLS-related parameters are ignored unless the cert and
	// key paths are specified.
	if isSet("TEMPORAL_TLS_CERT") && isSet("TEMPORAL_TLS_KEY") {
		clientCertPath := os.Getenv("TEMPORAL_TLS_CERT")
		clientKeyPath := os.Getenv("TEMPORAL_TLS_KEY")

		clientCert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
		if err != nil {
			msg := "failed to load client cert and key: %w"
			return client.Options{}, fmt.Errorf(msg, err)
		}

		var tlsConfig = tls.Config{}
		tlsConfig.Certificates = []tls.Certificate{clientCert}

		clientOpts.ConnectionOptions = client.ConnectionOptions{
			TLS: &tlsConfig,
		}
	}

	return clientOpts, nil
}

func isSet(name string) bool {
	_, isSet := os.LookupEnv(name)
	return isSet
}

func getEnvWithFallback(name, defaultValue string) string {
	value, isSet := os.LookupEnv(name)
	if !isSet {
		value = defaultValue
	}

	return value
}
