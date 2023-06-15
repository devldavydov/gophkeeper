// Package msgp contains utils for MSGP library.
package msgp

import (
	"bytes"

	"github.com/tinylib/msgp/msgp"
)

func Serde(input msgp.Encodable, output msgp.Decodable) error {
	var buf bytes.Buffer

	msgpW := msgp.NewWriter(&buf)
	err := input.EncodeMsg(msgpW)
	if err != nil {
		return err
	}
	msgpW.Flush()

	msgpR := msgp.NewReader(&buf)
	err = output.DecodeMsg(msgpR)
	if err != nil {
		return err
	}

	return nil
}
