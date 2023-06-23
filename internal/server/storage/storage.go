// Packages storage contains functions to work with persistent storage.
package storage

import (
	"context"

	"github.com/devldavydov/gophkeeper/internal/common/model"
)

// Storage is an interface to store users and secrets in persistent storage.
type Storage interface {
	// Create new user in storage.
	CreateUser(ctx context.Context, login, password string) (int64, error)
	// Find user in storage.
	FindUser(ctx context.Context, login string) (int64, string, error)

	// Create user's secret.
	CreateSecret(ctx context.Context, userID int64, secret *model.Secret) error
	// Get user's secret.
	GetSecret(ctx context.Context, userID int64, name string) (*model.Secret, error)
	// Get all user's secrets.
	GetAllSecrets(ctx context.Context, userID int64) ([]model.SecretInfo, error)
	// Delete user's secret.
	DeleteSecret(ctx context.Context, userID int64, name string) error
	// Delete all user's secret.
	DeleteAllSecrets(ctx context.Context, userID int64) error
	// Update user's secret.
	UpdateSecret(ctx context.Context, userID int64, name string, update *model.SecretUpdate) error

	// Check storage availability.
	Ping(ctx context.Context) bool
	// Close connection with storage.
	Close()
}
