package jbase

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// #region StringBuffer
type StringBuffer struct {
	Value    string
	Encoding StringEncoding
}
type StringEncoding string

const (
	Utf8   StringEncoding = "utf8"
	Base64 StringEncoding = "base64"
	Hex    StringEncoding = "hex"
)

func NewStringBuffer(value string, encoding StringEncoding) *StringBuffer {
	return &StringBuffer{
		Value:    value,
		Encoding: encoding,
	}
}
func (sb *StringBuffer) GetBinary() (binary []byte) {
	binary, err := sb.GetBinaryWithError()
	if err != nil {
		panic(err)
	}
	return
}

func (sb *StringBuffer) GetBinaryWithError() (binary []byte, err error) {
	if sb.Encoding == Utf8 {
		binary = []byte(sb.Value)
	} else if sb.Encoding == Base64 {
		binary, err = base64.StdEncoding.DecodeString(sb.Value)
	} else if sb.Encoding == Hex {
		binary, err = hex.DecodeString(sb.Value)
	} else {
		err = fmt.Errorf("invalid StringEncoding: %s", sb.Encoding)
	}
	return
}

func (sb *StringBuffer) UnmarshalJSON(data []byte) error {
	sb.Encoding = Hex
	json.Unmarshal(data, &sb.Value)
	return nil
}

func (sb StringBuffer) MarshalJSON() ([]byte, error) {
	return ([]byte)(fmt.Sprintf("%q", sb.Value)), nil
}

func (sb *StringBuffer) AsJsBuffer() string {
	return fmt.Sprintf("Buffer.from(%q, %q)", sb.Value, sb.Encoding)
}

// #endregion

// #region NewStringBuffer
type HexStringBuffer struct {
	StringBuffer
}

func NewHexStringBuffer(value string) *HexStringBuffer {
	return &HexStringBuffer{
		StringBuffer: StringBuffer{
			Value:    value,
			Encoding: Hex,
		},
	}
}
func CreateHexStringBuffer(binary []byte) *HexStringBuffer {
	return NewHexStringBuffer(hex.EncodeToString(binary))
}

// #endregion
// #region Base64StringBuffer

type Base64StringBuffer struct {
	StringBuffer
}

func NewBase64StringBuffer(value string) *Base64StringBuffer {
	return &Base64StringBuffer{
		StringBuffer: StringBuffer{
			Value:    value,
			Encoding: Base64,
		},
	}
}
func CreateBase64StringBuffer(binary []byte) *Base64StringBuffer {
	return NewBase64StringBuffer(base64.RawStdEncoding.EncodeToString(binary))
}
func (sb *Base64StringBuffer) UnmarshalJSON(data []byte) error {
	sb.Encoding = Base64
	sb.Value = string(data)
	return nil
}

// #endregion

// #region Utf8StringBuffer
type Utf8StringBuffer struct {
	StringBuffer
}

func NewUtf8StringBuffer(value string) *Utf8StringBuffer {
	return &Utf8StringBuffer{
		StringBuffer: StringBuffer{
			Value:    value,
			Encoding: Utf8,
		},
	}
}
func CreateUtf8StringBuffer(binary []byte) *Utf8StringBuffer {
	return NewUtf8StringBuffer(string(binary))
}
func (sb *Utf8StringBuffer) UnmarshalJSON(data []byte) error {
	sb.Encoding = Utf8
	sb.Value = string(data)
	return nil
}

// #endregion
