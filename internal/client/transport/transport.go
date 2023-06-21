// Package transport contains functions to connect with server in different ways.
package transport

import (
	"errors"

	"github.com/devldavydov/gophkeeper/internal/common/model"
)

var (
	ErrInternalServerError    = errors.New("internal server error")
	ErrInternalError          = errors.New("internal error")
	ErrUserAlreadyExists      = errors.New("user already exists")
	ErrUserInvalidCredentials = errors.New("user invalid credentials")
	ErrUserNotFound           = errors.New("user not found")
	ErrUserLoginFailed        = errors.New("user wrong login/password")
	ErrSecretAlreadyExists    = errors.New("secret already exists")
)

// Transport is an interface to connect with server.
type Transport interface {
	UserCreate(userLogin, userPassword string) (string, error)
	UserLogin(userLogin, userPassword string) (string, error)

	SecretGetList(token string) ([]model.SecretInfo, error)
	SecretCreate(token string, secret *model.Secret, payload model.Payload) error
}
