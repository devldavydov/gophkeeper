// Packages storage contains functions to work with persistent storage.
package storage

import (
	"context"

	"github.com/devldavydov/gophkeeper/internal/common/model"
)

// Storage is an interface to store users and secrets in persistent storage.
type Storage interface {
	CreateUser(ctx context.Context, login, password string) (int64, error)
	FindUser(ctx context.Context, login string) (int64, string, error)

	CreateSecret(ctx context.Context, userID int64, secret *model.Secret) error
	GetSecret(ctx context.Context, userID int64, name string) (*model.Secret, error)
	GetAllSecrets(ctx context.Context, userID int64) ([]model.SecretInfo, error)
	DeleteSecret(ctx context.Context, userID int64, name string) error
	DeleteAllSecrets(ctx context.Context) error
	UpdateSecret(ctx context.Context, userID int64, name string, update *model.SecretUpdate) error

	Ping(ctx context.Context) bool
	Close()
}
