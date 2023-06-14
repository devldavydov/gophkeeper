// Package server is the main package for GophKeeper server part.
package server

import (
	"time"

	"github.com/sirupsen/logrus"
)

// Service represents GophKeeper server.
type Service struct {
	settings        *ServiceSettings
	logger          *logrus.Logger
	shutdownTimeout time.Duration
}

// NewService creates new Service object.
func NewService(settings *ServiceSettings, shutdownTimeout time.Duration, logger *logrus.Logger) *Service {
	return &Service{settings: settings, logger: logger, shutdownTimeout: shutdownTimeout}
}
