// Packages storage contains functions to work with persistent storage.
package storage

import "context"

// Storage is an interface to store users and secrets in persistent storage.
type Storage interface {
	CreateUser(ctx context.Context, login, password string) (int, error)
	FindUser(ctx context.Context, login string) (int, string, error)

	Ping(ctx context.Context) bool
	Close()
}
