package temporal

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// verifies that the client options has the default address and Namespace
func TestCreateClientOptionsFromEnvDefaults(t *testing.T) {
	os.Unsetenv("TEMPORAL_ADDRESS")
	os.Unsetenv("TEMPORAL_NAMESPACE")
	got, err := CreateClientOptionsFromEnv()

	require.Nil(t, err)
	require.Equal(t, "localhost:7233", got.HostPort)
	require.Equal(t, "default", got.Namespace)
}

// verifies that the client options pulls the address from the env var
func TestCreateClientOptionsFromEnvAddressOnly(t *testing.T) {
	address := "darkstar.example.com:7233"
	t.Setenv("TEMPORAL_ADDRESS", address)
	os.Unsetenv("TEMPORAL_NAMESPACE")
	got, err := CreateClientOptionsFromEnv()

	require.Nil(t, err)
	require.Equal(t, "darkstar.example.com:7233", got.HostPort)
}

// verifies that CreateClientOptionsFromEnv fails when the cert is
// provide but the key is not
func TestCreateClientOptionsFromEnvCertButNoKeyFails(t *testing.T) {
	certPath := "/tmp/example.pem"
	t.Setenv("TEMPORAL_TLS_CERT", certPath)
	os.Unsetenv("TEMPORAL_TLS_KEY")
	_, err := CreateClientOptionsFromEnv()

	require.NotNil(t, err)
}

func TestCreateClientOptionsFromEnvWithKeyAndCert(t *testing.T) {
	keyPath := "/tmp/example.key"
	t.Setenv("TEMPORAL_TLS_KEY", keyPath)
	os.Unsetenv("TEMPORAL_TLS_CERT")
	_, err := CreateClientOptionsFromEnv()

	require.NotNil(t, err)
}

// verifies that isSet returns true if the env var is set
func TestIsSetPositiveCase(t *testing.T) {
	key := "EXAMPLE_ENV_VAR_FOR_TEST"
	t.Setenv(key, "example-value")

	require.True(t, isSet(key))
}

// verifies that isSet returns false if the env var is not set
func TestIsSetNegativeCase(t *testing.T) {
	key := "EXAMPLE_ENV_VAR_FOR_TEST"
	os.Unsetenv(key) // just to be sure it's not set

	require.False(t, isSet(key))
}

// verifies that getEnvWithFallback returns env var value if set
func TestGetEnvWithFallbackPositiveCase(t *testing.T) {
	key := "EXAMPLE_ENV_VAR_FOR_TEST"
	val := "example-value"
	t.Setenv(key, val)

	require.Equal(t, val, getEnvWithFallback(key, "nope"))
}

// verifies that getEnvWithFallback returns default value if env var unset
func TestGetEnvWithFallbackNegativeCase(t *testing.T) {
	key := "EXAMPLE_ENV_VAR_FOR_TEST"
	os.Unsetenv(key) // just to be sure it's not set

	def := "default-value"

	require.Equal(t, def, getEnvWithFallback(key, def))
}
