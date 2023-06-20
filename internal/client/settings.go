package client

import "github.com/devldavydov/gophkeeper/internal/common/nettools"

// Settings rerpesents client application settings.
type Settings struct {
	ServerAddress *nettools.Address
	TLSCACertPath string
}

// NewSettings creates new Settings object.
func NewSettings(serverAddress *nettools.Address, tlsCACertPath string) *Settings {
	return &Settings{ServerAddress: serverAddress, TLSCACertPath: tlsCACertPath}
}
