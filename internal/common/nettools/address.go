// Package nettools provide utility functions to work with network.
package nettools

import (
	"errors"
	"fmt"
	"net"
	"strconv"
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
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, ErrAddressFormat
	}

	nPort, err := strconv.Atoi(port)
	if err != nil {
		return nil, ErrAddressFormat
	}

	return &Address{Host: host, Port: nPort}, nil
}

func (a Address) String() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}
