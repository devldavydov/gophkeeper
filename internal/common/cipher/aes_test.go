package cipher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAES(t *testing.T) {
	key := []byte("asuperstrong32bitpasswordgohere!")
	plainData := []byte("Hello world!!!")

	cipherData, err := AESEncrpyt(plainData, key)
	assert.NoError(t, err)

	decrData, err := AESDecrypt(cipherData, key)
	assert.NoError(t, err)

	assert.Equal(t, plainData, decrData)
}

func TestAESErrors(t *testing.T) {
	_, err := AESEncrpyt([]byte("test"), []byte("test"))
	assert.ErrorIs(t, err, ErrWrongKeyLength)

	_, err = AESDecrypt([]byte("test"), []byte("test"))
	assert.ErrorIs(t, err, ErrWrongKeyLength)

	_, err = AESDecrypt([]byte("test"), []byte("asuperstrong32bitpasswordgohere!"))
	assert.ErrorIs(t, err, ErrCipherDataTooShort)
}
