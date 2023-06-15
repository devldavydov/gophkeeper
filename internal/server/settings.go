package server

import "time"

// ServiceSettings represents settings for GophKeeper server.
type ServiceSettings struct {
	DatabaseDsn     string
	ServerSecret    string
	ShutdownTimeout time.Duration
}

// NewServiceSettings creates new ServiceSettings object.
func NewServiceSettings(
	databaseDsn string,
	serverSecret string,
	shutdownTimeout time.Duration,
) *ServiceSettings {
	return &ServiceSettings{
		DatabaseDsn:     databaseDsn,
		ServerSecret:    serverSecret,
		ShutdownTimeout: shutdownTimeout,
	}
}
