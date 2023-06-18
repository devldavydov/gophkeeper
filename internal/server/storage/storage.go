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

	CreateSecret(ctx context.Context, userID int64, secret model.Secret) error

	Ping(ctx context.Context) bool
	Close()
}
