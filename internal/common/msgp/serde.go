// Package msgp contains utils for MSGP library.
package msgp

import (
	"bytes"

	"github.com/tinylib/msgp/msgp"
)

// Serialize input object to bytes.
//
// Returns serialized byte slice or Msgp error.
func Serialize(input msgp.Encodable) ([]byte, error) {
	var buf bytes.Buffer

	msgpW := msgp.NewWriter(&buf)
	err := input.EncodeMsg(msgpW)
	if err != nil {
		return nil, err
	}
	msgpW.Flush()

	return buf.Bytes(), nil
}

// Deserialize bytes to output object.
//
// In case of success argument ouput contains deserialized object.
//
// In case of error returns MSGP error.
func Deserialize(data []byte, output msgp.Decodable) error {
	buf := bytes.NewBuffer(data)

	msgpR := msgp.NewReader(buf)
	err := output.DecodeMsg(msgpR)
	if err != nil {
		return err
	}

	return nil
}

// SerDe - serialize and deserialize.
//
// In case of success argument ouput contains deserialized object.
//
// In case of error returns MSGP error.
func SerDe(input msgp.Encodable, output msgp.Decodable) error {
	data, err := Serialize(input)
	if err != nil {
		return err
	}

	return Deserialize(data, output)
}
