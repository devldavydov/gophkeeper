// Package transport provides different transport mechanisms to connect with server.
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
	ErrUserPermissionDenied      = errors.New("user with login/password permission denied")
	ErrSecretAlreadyExists       = errors.New("secret already exists")
	ErrSecretNotFound            = errors.New("secret not found")
	ErrSecretOutdated            = errors.New("secret outdated")
	ErrSecretPayloadSizeExceeded = errors.New("secret payload size exceeded")
	ErrSecretInvalid             = errors.New("invalid secret")
)

// Transport is a common interface to connect with server.
// Containts description of user and secret methods.
type Transport interface {
	// Create new user on server.
	UserCreate(userLogin, userPassword string) (string, error)
	// Login existing user on server.
	UserLogin(userLogin, userPassword string) (string, error)

	// Retreive users's secret list.
	SecretGetList(token string) ([]model.SecretInfo, error)
	// Retrieve user secret.
	SecretGet(token, name string) (*model.Secret, error)
	// Create user secret.
	SecretCreate(token string, secret *model.Secret) error
	// Update user secret.
	SecretUpdate(token, name string, updSecret *model.SecretUpdate) error
	// Delete user secret.
	SecretDelete(token, name string) error
}
