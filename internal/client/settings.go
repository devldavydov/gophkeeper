package client

import "github.com/devldavydov/gophkeeper/internal/common/nettools"

// ApplicationSettings rerpesents client application settings.
type ApplicationSettings struct {
	ServerAddress *nettools.Address
	TLSCACertPath string
}

// NewApplicationSettings creates new ApplicationSettings object.
func NewApplicationSettings(serverAddress *nettools.Address, tlsCACertPath string) *ApplicationSettings {
	return &ApplicationSettings{ServerAddress: serverAddress, TLSCACertPath: tlsCACertPath}
}
