package server

import (
	"time"

	"github.com/devldavydov/gophkeeper/internal/common/nettools"
	gkTLS "github.com/devldavydov/gophkeeper/internal/common/tls"
)

// ServiceSettings represents settings for GophKeeper server.
type ServiceSettings struct {
	// GRPCAddress - listen address of server.
	GRPCAddress *nettools.Address

	// GRPCServerTLS - TLS server settings.
	GRPCServerTLS *gkTLS.ServerSettings

	// DatabaseDsn - database connection string.
	DatabaseDsn string

	// ServerSecret - unique 32 chars string to be used as a key of encryption.
	ServerSecret string

	// ShutdownTimeout - server shitdown timeout.
	ShutdownTimeout time.Duration
}

// NewServiceSettings creates new ServiceSettings object.
func NewServiceSettings(
	grpcAddress *nettools.Address,
	grpcServerTLS *gkTLS.ServerSettings,
	databaseDsn string,
	serverSecret string,
	shutdownTimeout time.Duration,
) *ServiceSettings {
	return &ServiceSettings{
		GRPCAddress:     grpcAddress,
		GRPCServerTLS:   grpcServerTLS,
		DatabaseDsn:     databaseDsn,
		ServerSecret:    serverSecret,
		ShutdownTimeout: shutdownTimeout,
	}
}
