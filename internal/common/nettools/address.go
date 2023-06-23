// Package nettools provide utility functions to work with network.
package nettools

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrAddressFormat = errors.New("wrong address format")

// Address - utility struct for holding address pair.
type Address struct {
	Host string
	Port int
}

// NewAddress creates new Address object.
//
// In case of invalid input argument, returns ErrAddressFormat error.
func NewAddress(addr string) (*Address, error) {
	parts := strings.Split(addr, ":")
	if len(parts) != 2 { //nolint:gomnd // No magic numbers
		return nil, ErrAddressFormat
	}

	host := parts[0]
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, ErrAddressFormat
	}

	return &Address{Host: host, Port: port}, nil
}

func (a Address) String() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}
