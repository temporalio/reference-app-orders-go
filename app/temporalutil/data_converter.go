package temporalutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	commonpb "go.temporal.io/api/common/v1"
	"go.temporal.io/sdk/converter"
)

const (
	// metadataEncodingEncrypted identifies payloads encoded with an encrypted binary format
	metadataEncodingEncrypted = "binary/encrypted"
	// metadataEncryptionKeyID identifies the key used to encrypt a payload
	metadataEncryptionKeyID = "encryption-key-id"
)

// Codec provides methods for encrypting and decrypting payload data
type Codec struct {
	EncryptionKeyID string
}

// this function simulates the retrieval of an encryption key (identified by
// the provided key ID) from secure storage, such as a key management server
func (e *Codec) retrieveKey(keyID string) (key []byte, err error) {
	if keyID == "" {
		return nil, fmt.Errorf("key retrieval failed due to empty identifier")
	}

	// Simulate key retrieval by using a hash function to generate
	// a 256-bit value that will be consistent for a given key ID
	h := sha256.Sum256([]byte(keyID))

	return h[:], nil
}

// NewEncryptionDataConverter creates and returns a DataConverter instance that
// wraps the default DataConverter with a CodecDataConverter that uses encryption
// to protect the confidentiality of payload data. This instance will encrypt data
// using a key associated with the specified encryption key ID.
func NewEncryptionDataConverter(underlying converter.DataConverter, encryptionKeyID string) converter.DataConverter {
	codecs := []converter.PayloadCodec{
		&Codec{EncryptionKeyID: encryptionKeyID},
	}

	return converter.NewCodecDataConverter(underlying, codecs...)
}

// Encode implements the Encode method defined by the converter.PayloadCodec interface
func (e *Codec) Encode(payloads []*commonpb.Payload) ([]*commonpb.Payload, error) {
	result := make([]*commonpb.Payload, len(payloads))
	for i, payload := range payloads {
		unencryptedData, err := payload.Marshal()
		if err != nil {
			return payloads, err
		}

		key, err := e.retrieveKey(e.EncryptionKeyID)
		if err != nil {
			return payloads, err
		}

		encryptedData, err := encrypt(unencryptedData, key)
		if err != nil {
			return payloads, err
		}

		result[i] = &commonpb.Payload{
			Metadata: map[string][]byte{
				converter.MetadataEncoding: []byte(metadataEncodingEncrypted),
				metadataEncryptionKeyID:    []byte(e.EncryptionKeyID),
			},
			Data: encryptedData,
		}
	}

	return result, nil
}

// Decode implements the Decode method defined by the converter.PayloadCodec interface
func (e *Codec) Decode(payloads []*commonpb.Payload) ([]*commonpb.Payload, error) {
	result := make([]*commonpb.Payload, len(payloads))
	for i, payload := range payloads {
		payloadFormatID := string(payload.Metadata[converter.MetadataEncoding])

		// Skip decryption for any payload not using our encrypted format
		if payloadFormatID != metadataEncodingEncrypted {
			result[i] = payload
			continue
		}

		encryptedData := payload.Data

		keyID, ok := payload.Metadata[metadataEncryptionKeyID]
		if !ok {
			return payloads, fmt.Errorf("encryption key id missing from metadata")
		}

		key, err := e.retrieveKey(string(keyID))
		if err != nil {
			return payloads, err
		}

		decryptedData, err := decrypt(encryptedData, key)
		if err != nil {
			return payloads, err
		}

		result[i] = &commonpb.Payload{}
		err = result[i].Unmarshal(decryptedData)
		if err != nil {
			return payloads, err
		}
	}

	return result, nil
}

// Uses AES to encrypt the provided block of data using with the provided key
func encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	encrypted := aesgcm.Seal(nonce, nonce, data, nil)
	return encrypted, nil
}

// Uses AES to decrypt the provided block of data using the provided key
func decrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesgcm.NonceSize()
	nonce, encrypted := data[:nonceSize], data[nonceSize:]
	return aesgcm.Open(nil, nonce, encrypted, nil)
}
