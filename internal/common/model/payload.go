package model

//go:generate msgp -tests=false

import (
	"crypto/sha256"
	"fmt"
)

// Payload interface represents secret valuable information.
type Payload interface {
	fmt.Stringer
	Valid() bool
}

// CredsPayload represents user login/password pair.
type CredsPayload struct {
	Hash     string
	Login    string
	Password string
}

var _ Payload = (*CredsPayload)(nil)

// NewCredsPayload creates new CredsPayload object.
func NewCredsPayload(login, password string) *CredsPayload {
	cp := &CredsPayload{Login: login, Password: password}
	cp.Hash = cp.hash()
	return cp
}

func (cp *CredsPayload) String() string {
	return fmt.Sprintf("Login=%s Password=%s", cp.Login, cp.Password)
}

func (cp *CredsPayload) Valid() bool {
	return cp.Hash == cp.hash()
}

func (cp *CredsPayload) hash() string {
	return sha256sum([]byte(cp.Login + cp.Password))
}

// TextPayload represents arbitrary text data.
type TextPayload struct {
	Hash string
	Data string
}

var _ Payload = (*TextPayload)(nil)

// NewTextPayload creates new TextPayload object.
func NewTextPayload(data string) *TextPayload {
	tp := &TextPayload{Data: data}
	tp.Hash = tp.hash()
	return tp
}

func (tp *TextPayload) String() string {
	return fmt.Sprintf("Text=%s", tp.Data)
}

func (tp *TextPayload) Valid() bool {
	return tp.Hash == tp.hash()
}

func (tp *TextPayload) hash() string {
	return sha256sum([]byte(tp.Data))
}

// BinaryPayload represents arbitrary binary data.
type BinaryPayload struct {
	Hash string
	Data []byte
}

var _ Payload = (*BinaryPayload)(nil)

// NewBinaryPayload creates new BinaryPayload object.
func NewBinaryPayload(data []byte) *BinaryPayload {
	bp := &BinaryPayload{Data: data}
	bp.Hash = bp.hash()
	return bp
}

func (bp *BinaryPayload) String() string {
	return fmt.Sprintf("Data=%x", bp.Data)
}

func (bp *BinaryPayload) Valid() bool {
	return bp.Hash == bp.hash()
}

func (bp *BinaryPayload) hash() string {
	return sha256sum(bp.Data)
}

// CardPayload represents card data.
type CardPayload struct {
	Hash       string
	CardNum    string
	CardHolder string
	ValidThru  string
	CVV        string
}

var _ Payload = (*CardPayload)(nil)

// NewCardPayload creates new CardPayload object.
func NewCardPayload(cardNum, cardHolder, validThru, cvv string) *CardPayload {
	cp := &CardPayload{CardNum: cardNum, CardHolder: cardHolder, ValidThru: validThru, CVV: cvv}
	cp.Hash = cp.hash()
	return cp
}

func (cp *CardPayload) String() string {
	return fmt.Sprintf(
		"CardNum=%s CardHolder=%s ValidThru=%s CVV=%s",
		cp.CardNum,
		cp.CardHolder,
		cp.ValidThru,
		cp.CVV)
}

func (cp *CardPayload) Valid() bool {
	return cp.Hash == cp.hash()
}

func (cp *CardPayload) hash() string {
	return sha256sum([]byte(cp.CardNum + cp.CardHolder + cp.ValidThru + cp.CVV))
}

func sha256sum(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}
