// Package server is the main package for GophKeeper server part.
package server

import (
	"context"

	"github.com/sirupsen/logrus"
)

// Service represents GophKeeper server.
type Service struct {
	settings *ServiceSettings
	logger   *logrus.Logger
}

// NewService creates new Service object.
func NewService(settings *ServiceSettings, logger *logrus.Logger) *Service {
	return &Service{settings: settings, logger: logger}
}

func (s *Service) Start(ctx context.Context) error {
	return nil
}
