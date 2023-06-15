// Package model represents data structures for GophKeeper.
package model

import (
	"errors"

	gkMsgp "github.com/devldavydov/gophkeeper/internal/common/msgp"
	"github.com/tinylib/msgp/msgp"
)

var (
	ErrUnknownPayload = errors.New("unknown secret payload")
	ErrInvalidPayload = errors.New("invalid secret payload")
)

// SecretType is an enum type for secret types.
type SecretType int

const (
	UnknownSecret SecretType = 0
	CredsSecret   SecretType = 1
	TextSecret    SecretType = 2
	BinarySecret  SecretType = 3
	CardSecret    SecretType = 4
)

func (st SecretType) String() string {
	switch st {
	case CredsSecret:
		return "Credentials"
	case TextSecret:
		return "Text"
	case BinarySecret:
		return "Binary"
	case CardSecret:
		return "Card"
	case UnknownSecret:
		return "Unknown"
	default:
		return "Unknown"
	}
}

// Secret represents data structure for secret.
type Secret struct {
	Type       SecretType
	Name       string
	Meta       string
	Version    int64
	PayloadRaw []byte
}

// GetPayload returns Payload from binary raw.
func (s *Secret) GetPayload() (Payload, error) {
	var decObj msgp.Decodable

	switch s.Type {
	case CredsSecret:
		decObj = &CredsPayload{}
	case TextSecret:
		decObj = &TextPayload{}
	case BinarySecret:
		decObj = &BinaryPayload{}
	case CardSecret:
		decObj = &CardPayload{}
	case UnknownSecret:
		return nil, ErrUnknownPayload
	default:
		return nil, ErrUnknownPayload
	}

	if err := gkMsgp.Deserialize(s.PayloadRaw, decObj); err != nil {
		return nil, err
	}

	payload, ok := decObj.(Payload)
	if !ok || !payload.Valid() {
		return nil, ErrInvalidPayload
	}

	return payload, nil
}
