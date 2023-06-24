// Package model represents data structures for GophKeeper.
package model

import (
	"errors"

	gkMsgp "github.com/devldavydov/gophkeeper/internal/common/msgp"
	"github.com/tinylib/msgp/msgp"
)

var (
	ErrUnknownPayload    = errors.New("unknown secret payload")
	ErrInvalidPayload    = errors.New("invalid secret payload")
	ErrInvalidSecretType = errors.New("invalid secret type")
)

// SecretType is an enum type for secret types.
type SecretType int32

const (
	UnknownSecret SecretType = iota
	CredsSecret
	TextSecret
	BinarySecret
	CardSecret
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

// ValidSecretType checks validity of given secret type.
//
// Returns nil or ErrInvalidSecretType error.
func ValidSecretType(st SecretType) error {
	switch st {
	case CredsSecret, TextSecret, BinarySecret, CardSecret:
		return nil
	case UnknownSecret:
		return ErrInvalidSecretType
	default:
		return ErrInvalidSecretType
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

// SecretInfo represents short information about Secret. Used in list of secrets.
type SecretInfo struct {
	Type    SecretType
	Name    string
	Version int64
}

// SecretUpdate represents Secret fields to update.
type SecretUpdate struct {
	Meta          string
	Version       int64
	PayloadRaw    []byte
	UpdatePayload bool
}

// GetPayload returns Payload from binary raw.
//
// In case of unknown secret type returns ErrUnknownPayload.
//
// In case if MSGP deserialization error, returns specific Msgp error.
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
