// Packages storage contains functions to work with persistent storage.
package storage

// Storage is an interface to store users and secrets in persistent storage.
type Storage interface {
	Close()
}
