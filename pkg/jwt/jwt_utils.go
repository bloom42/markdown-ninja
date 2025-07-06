package jwt

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

// BytesBase64RawUrl is a simple []byte that encodes to base64RawUrl when marshaling to JSONs
type BytesBase64RawUrl []byte

var BytesBase64RawUrlNUll = []byte("null")

func (b BytesBase64RawUrl) String() string {
	return base64.RawURLEncoding.EncodeToString(b)
}

func (b BytesBase64RawUrl) MarshalJSON() ([]byte, error) {
	if b == nil {
		return BytesBase64RawUrlNUll, nil
	}

	buffer := bytes.NewBuffer(make([]byte, 0, (2 + len(b)*2)))
	buffer.WriteRune('"')
	buffer.WriteString(b.String())
	buffer.WriteRune('"')
	return buffer.Bytes(), nil
}

func (b *BytesBase64RawUrl) UnmarshalJSON(data []byte) (err error) {
	if data == nil || bytes.Equal(data, BytesBase64RawUrlNUll) {
		return nil
	}

	data = bytes.Trim(data, `"`)
	decodedData, err := base64.RawURLEncoding.DecodeString(string(data))
	if err != nil {
		return
	}

	*b = decodedData
	return nil
}

func (b BytesBase64RawUrl) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(b.String())
}

func (b *BytesBase64RawUrl) UnmarshalCBOR(data []byte) (err error) {
	if data == nil {
		return nil
	}

	decodedData, err := base64.RawURLEncoding.DecodeString(string(data))
	if err != nil {
		return
	}

	*b = decodedData
	return nil
}

func (b *BytesBase64RawUrl) Scan(val any) error {
	switch v := val.(type) {
	case []byte:
		*b = v
		return nil
	case nil:
		return nil
	default:
		return fmt.Errorf("BytesBase64RawUrl.Scan: Unsupported type: %T", v)
	}
}
