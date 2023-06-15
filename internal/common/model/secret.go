// Package model represents data structures for GophKeeper.
package model

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

func (s *Secret) GetPayload() Payload {
	return nil
}
