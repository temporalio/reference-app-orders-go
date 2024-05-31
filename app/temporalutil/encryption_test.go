package temporalutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// verifies that data created by the encrypt function is
// compatible with the decrypt function
func TestEncryptAndDecrypt(t *testing.T) {
	original := []byte("The crow flies at midnight")
	key := []byte("c555af41f0f17e7bd4bdab1e9e6f7873")

	encrypted, err := encrypt(original, key)
	require.NoError(t, err)

	decrypted, err := decrypt(encrypted, key)
	require.NoError(t, err)

	require.NotEqual(t, original, encrypted)
	require.Equal(t, original, decrypted)
}
