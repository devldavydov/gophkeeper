// Package cipher contains util functions to work with cryptography.
package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

var ErrCipherDataTooShort = errors.New("cipher data too short")
var ErrWrongKeyLength = errors.New("key length should be 32 bytes")

const AESKeyLength = 32

// AESEncrypt encrypts plainData with secretKey.
//
// Accespts plain data and key. Returns encrypted data or error.
func AESEncrpyt(plainData []byte, secretKey []byte) ([]byte, error) {
	if len(secretKey) != AESKeyLength {
		return nil, ErrWrongKeyLength
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}

	cipherData := make([]byte, aes.BlockSize+len(plainData))

	iv := cipherData[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherData[aes.BlockSize:], plainData)

	return cipherData, nil
}

// AESDecrypt decrypts cipherData with secretKey.
//
// Accepts encrypted and key. Returns plain data or error.
func AESDecrypt(cipherData []byte, secretKey []byte) ([]byte, error) {
	if len(secretKey) != AESKeyLength {
		return nil, ErrWrongKeyLength
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}

	if len(cipherData) < aes.BlockSize {
		return nil, ErrCipherDataTooShort
	}

	iv := cipherData[:aes.BlockSize]
	cipherData = cipherData[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherData, cipherData)

	return cipherData, nil
}
