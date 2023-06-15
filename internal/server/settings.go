package server

import (
	"time"

	"github.com/devldavydov/gophkeeper/internal/common/nettools"
	gkTLS "github.com/devldavydov/gophkeeper/internal/common/tls"
)

// ServiceSettings represents settings for GophKeeper server.
type ServiceSettings struct {
	GRPCAddress     *nettools.Address
	GRPCServerTLS   *gkTLS.ServerSettings
	DatabaseDsn     string
	ServerSecret    string
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
