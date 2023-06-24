// Package server is the main package for GophKeeper server part.
package server

import (
	"context"
	"fmt"
	"net"

	"github.com/devldavydov/gophkeeper/internal/server/storage"
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

// Start - starts GophKeeper server service.
func (s *Service) Start(ctx context.Context) error {
	// Storage
	stg, err := storage.NewPgStorage(s.settings.DatabaseDsn, s.logger)
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}
	defer stg.Close()

	// TLS
	tlsCredentials, err := s.settings.GRPCServerTLS.Load()
	if err != nil {
		return fmt.Errorf("failed to load TLS: %w", err)
	}

	// gRPC server
	listen, err := net.Listen("tcp", s.settings.GRPCAddress.String())
	if err != nil {
		return fmt.Errorf("failed to setup gRPC listen: %w", err)
	}

	grpcSrv, _ := NewGrpcServer(stg, tlsCredentials, []byte(s.settings.ServerSecret), s.logger)

	errChan := make(chan error)
	go func(ch chan error) {
		s.logger.Infof("gRPC service started on [%s]", s.settings.GRPCAddress.String())
		ch <- grpcSrv.Serve(listen)
	}(errChan)

	select {
	case err = <-errChan:
		return fmt.Errorf("gRPC service exited with err: %w", err)
	case <-ctx.Done():
		s.logger.Infof("gRPC service context canceled")

		grpcSrv.GracefulStop()

		s.logger.Info("gRPC service finished")
		return nil
	}
}
