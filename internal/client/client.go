// Package server is the main package for GophKeeper client part.
package client

import (
	"context"
	"fmt"

	"github.com/devldavydov/gophkeeper/internal/client/transport"
	"github.com/devldavydov/gophkeeper/internal/common/model"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
)

// Client represents client application.
type Client struct {
	settings   *Settings
	logger     *logrus.Logger
	tr         transport.Transport
	cltToken   string
	lstSecrets []model.SecretInfo
	//
	uiApp               *tview.Application
	uiPages             *tview.Pages
	wdgLogin            *tview.Form
	wdgCreateUser       *tview.Form
	wdgError            *tview.Form
	wdgUser             *tview.TextView
	wdgLstSecrets       *tview.List
	wdgCreateUserSecret *tview.Form
}

// NewClient creates new Application object.
func NewClient(settngs *Settings, logger *logrus.Logger) *Client {
	return &Client{settings: settngs, logger: logger}
}

func (r *Client) Start(ctx context.Context) error {
	var err error

	// Create Transport
	r.tr, err = transport.NewGrpcTransport(r.settings.ServerAddress, r.settings.TLSCACertPath)
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
