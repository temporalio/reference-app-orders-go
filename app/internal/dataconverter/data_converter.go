package dataconverter

import (
	"fmt"

	commonpb "go.temporal.io/api/common/v1"
	"go.temporal.io/sdk/converter"
)

const (
	// MetadataEncodingEncrypted identifies payloads encoded with an encrypted binary format
	MetadataEncodingEncrypted = "binary/encrypted"
	// MetadataEncryptionKeyID identifies the key used to encrypt a payload
	MetadataEncryptionKeyID   = "encryption-key-id"
)

// DataConverter wraps an underlying DataConverter with a CodecDataConverter that uses encryption to protect the confidentiality of payload data
type DataConverter struct {
	parent converter.DataConverter
	converter.DataConverter
	options DataConverterOptions
}

// DataConverterOptions holds options related to DataConverter configuration
type DataConverterOptions struct {
	EncryptionKeyID string
}

// Codec provides methods for encrypting and decrypting payload data
type Codec struct {
	EncryptionKeyID string
}

// this function simulates the retrieval of an encryption key (identified by
// the provided key ID) from secure storage, such as a key management server
func (e *Codec) getKey(keyID string) (key []byte, err error) {
	if keyID == "" {
		return nil, fmt.Errorf("key retrieval failed due to empty identifier")
	}
	return []byte("trivial-key-for-example-use-only"), nil
}

// NewEncryptionDataConverter creates and returns an instance of a DataConverter that wraps the default DataConverter with a CodecDataConverter that uses encryption to protect the confidentiality of payload data
func NewEncryptionDataConverter(dataConverter converter.DataConverter, options DataConverterOptions) *DataConverter {
	codecs := []converter.PayloadCodec{
		&Codec{EncryptionKeyID: options.EncryptionKeyID},
	}

	return &DataConverter{
		parent:        dataConverter,
		DataConverter: converter.NewCodecDataConverter(dataConverter, codecs...),
		options:       options,
	}
}

// Encode implements the Encode method defined by the converter.PayloadCodec interface
func (e *Codec) Encode(payloads []*commonpb.Payload) ([]*commonpb.Payload, error) {
	result := make([]*commonpb.Payload, len(payloads))
	for i, payload := range payloads {
		unencryptedData, err := payload.Marshal()
		if err != nil {
			return payloads, err
		}

		key, err := e.getKey(e.EncryptionKeyID)
		if err != nil {
			return payloads, err
		}

		encryptedData, err := encrypt(unencryptedData, key)
		if err != nil {
			return payloads, err
		}

		result[i] = &commonpb.Payload{
			Metadata: map[string][]byte{
				converter.MetadataEncoding: []byte(MetadataEncodingEncrypted),
				MetadataEncryptionKeyID:    []byte(e.EncryptionKeyID),
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
		if payloadFormatID != MetadataEncodingEncrypted {
			result[i] = payload
			continue
		}

		encryptedData := payload.Data

		keyID, ok := payload.Metadata[MetadataEncryptionKeyID]
		if !ok {
			return payloads, fmt.Errorf("encryption key id missing from metadata")
		}

		key, err := e.getKey(string(keyID))
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
