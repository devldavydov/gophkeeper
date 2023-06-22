// Package transport contains functions to connect with server in different ways.
package transport

import (
	"errors"

	"github.com/devldavydov/gophkeeper/internal/common/model"
)

var (
	ErrInternalServerError       = errors.New("internal server error")
	ErrInternalError             = errors.New("internal error")
	ErrUserAlreadyExists         = errors.New("user already exists")
	ErrUserInvalidCredentials    = errors.New("user invalid credentials")
	ErrUserNotFound              = errors.New("user not found")
	ErrUserLoginFailed           = errors.New("user wrong login/password")
	ErrSecretAlreadyExists       = errors.New("secret already exists")
	ErrSecretNotFound            = errors.New("secret not found")
	ErrSecretOutdated            = errors.New("secret outdated")
	ErrSecretPayloadSizeExceeded = errors.New("secret payload size exceeded")
	ErrSecretInvalid             = errors.New("invalid secret")
)

// Transport is an interface to connect with server.
type Transport interface {
	UserCreate(userLogin, userPassword string) (string, error)
	UserLogin(userLogin, userPassword string) (string, error)

	SecretGetList(token string) ([]model.SecretInfo, error)
	SecretGet(token, name string) (*model.Secret, error)
	SecretCreate(token string, secret *model.Secret) error
	SecretUpdate(token, name string, updSecret *model.SecretUpdate) error
	SecretDelete(token, name string) error
}
