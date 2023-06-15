package model

import (
	"fmt"
	"testing"

	gkMsgp "github.com/devldavydov/gophkeeper/internal/common/msgp"
	"github.com/stretchr/testify/assert"
	"github.com/tinylib/msgp/msgp"
)

func TestSecretTypeString(t *testing.T) {
	assert.Equal(t, "Unknown", UnknownSecret.String())
	assert.Equal(t, "Unknown", SecretType(100).String())
	assert.Equal(t, "Credentials", CredsSecret.String())
	assert.Equal(t, "Text", TextSecret.String())
	assert.Equal(t, "Binary", BinarySecret.String())
	assert.Equal(t, "Card", CardSecret.String())
}

func TestPayloadValid(t *testing.T) {
	for i, tt := range []struct {
		payload Payload
		valid   bool
	}{
		{payload: NewCredsPayload("foo", "bar"), valid: true},
		{payload: &CredsPayload{Login: "foo", Password: "bar"}, valid: false},
		{payload: NewTextPayload("foo"), valid: true},
		{payload: &TextPayload{Data: "foo"}, valid: false},
		{payload: NewBinaryPayload([]byte("foo")), valid: true},
		{payload: &BinaryPayload{Data: []byte("foo")}, valid: false},
		{payload: NewCardPayload("2202", "foo", "11/26", "777"), valid: true},
		{payload: &CardPayload{
			CardNum:    "2202",
			CardHolder: "foo",
			ValidThru:  "11/26",
			CVV:        "777",
		}, valid: false},
	} {
		tt := tt
		t.Run(fmt.Sprintf("Run %d", i), func(t *testing.T) {
			assert.Equal(t, tt.valid, tt.payload.Valid())
		})
	}
}

func TestPayloadString(t *testing.T) {
	for i, tt := range []struct {
		payload Payload
		str     string
	}{
		{
			payload: NewCredsPayload("foo", "bar"),
			str:     "Login=foo Password=bar",
		},
		{
			payload: NewTextPayload("foo"),
			str:     "Text=foo",
		},
		{
			payload: NewBinaryPayload([]byte("foobar")),
			str:     "Data=666f6f626172",
		},
		{
			payload: NewCardPayload("2202", "foo", "11/26", "777"),
			str:     "CardNum=2202 CardHolder=foo ValidThru=11/26 CVV=777",
		},
	} {
		tt := tt
		t.Run(fmt.Sprintf("Run %d", i), func(t *testing.T) {
			assert.Equal(t, tt.str, tt.payload.String())
		})
	}
}

func TestPayloadMsgpConversion(t *testing.T) {
	for i, tt := range []struct {
		input  msgp.Encodable
		output msgp.Decodable
	}{
		{input: NewCredsPayload("foo", "bar"), output: &CredsPayload{}},
		{input: NewTextPayload("foobar"), output: &TextPayload{}},
		{input: NewBinaryPayload([]byte("foobar")), output: &BinaryPayload{}},
		{input: NewCardPayload("2032", "foo", "11/26", "777"), output: &CardPayload{}},
	} {
		tt := tt
		t.Run(fmt.Sprintf("Run %d", i), func(t *testing.T) {
			assert.NoError(t, gkMsgp.SerDe(tt.input, tt.output))
			assert.Equal(t, tt.input, tt.output)
			assert.True(t, tt.output.(Payload).Valid())
		})
	}
}

func TestPayloadMsgpInvalidConversion(t *testing.T) {
	input := NewCredsPayload("foo", "bar")
	output := &CardPayload{}

	assert.NoError(t, gkMsgp.SerDe(input, output))
	assert.False(t, output.Valid())
}

func TestSecretGetPayload(t *testing.T) {
	for i, tt := range []struct {
		input      msgp.Encodable
		secretType SecretType
		isErr      bool
	}{
		{input: NewCredsPayload("foo", "bar"), secretType: CredsSecret},
		{input: NewTextPayload("foo"), secretType: TextSecret},
		{input: NewBinaryPayload([]byte("foo")), secretType: BinarySecret},
		{input: NewCardPayload("2202", "foo", "11/26", "777"), secretType: CardSecret},
		{input: NewCredsPayload("foo", "bar"), secretType: UnknownSecret, isErr: true},
		{input: NewCredsPayload("foo", "bar"), secretType: SecretType(100), isErr: true},
		{input: NewCredsPayload("foo", "bar"), secretType: BinarySecret, isErr: true},
	} {
		tt := tt
		t.Run(fmt.Sprintf("Run %d", i), func(t *testing.T) {
			data, err := gkMsgp.Serialize(tt.input)
			assert.NoError(t, err)

			secret := &Secret{Type: tt.secretType, PayloadRaw: data}
			payload, err := secret.GetPayload()
			if tt.isErr {
				assert.Error(t, err)
			} else {
				expPayload, ok := tt.input.(Payload)
				assert.True(t, ok)
				assert.Equal(t, expPayload.String(), payload.String())
			}
		})
	}
}
