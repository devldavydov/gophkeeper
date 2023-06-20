// Package server is the main package for GophKeeper client part.
package client

import (
	"context"
	"fmt"
	"time"

	"github.com/devldavydov/gophkeeper/internal/common/model"
	gkTLS "github.com/devldavydov/gophkeeper/internal/common/tls"
	pb "github.com/devldavydov/gophkeeper/internal/grpc"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const _serverRequestTimeout = 15 * time.Second

// Client represents client application.
type Client struct {
	settings   *Settings
	logger     *logrus.Logger
	gClt       pb.GophKeeperServiceClient
	cltToken   string
	lstSecrets []model.SecretInfo
	//
	uiApp         *tview.Application
	uiPages       *tview.Pages
	wdgLogin      *tview.Form
	wdgCreateUser *tview.Form
	wdgError      *tview.Form
	wdgUser       *tview.TextView
	wdgLstSecrets *tview.List
}

// NewClient creates new Application object.
func NewClient(settngs *Settings, logger *logrus.Logger) *Client {
	return &Client{settings: settngs, logger: logger}
}

func (r *Client) Start(ctx context.Context) error {
	var err error

	// Create gRPC client
	r.gClt, err = r.createGrpcClient()
	if err != nil {
		return err
	}

	// Start UI application
	r.createUIApplication()

	errChan := make(chan error)
	go func(ch chan error) {
		r.logger.Infof("client application started")
		ch <- r.uiApp.Run()
	}(errChan)

	select {
	case err = <-errChan:
		if err != nil {
			return fmt.Errorf("client exited with err: %w", err)
		}
		r.logger.Info("client application finished")
		return nil
	case <-ctx.Done():
		r.logger.Info("client application context canceled")
		r.logger.Info("client application finished")
		return nil
	}
}

func (r *Client) createGrpcClient() (pb.GophKeeperServiceClient, error) {
	tlsCredentials, err := gkTLS.LoadCACert(r.settings.TLSCACertPath, "")
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(r.settings.ServerAddress.String(), grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		return nil, err
	}

	return pb.NewGophKeeperServiceClient(conn), nil
}
