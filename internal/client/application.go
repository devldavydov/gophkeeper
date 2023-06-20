// Package server is the main package for GophKeeper client part.
package client

import (
	"context"

	"github.com/sirupsen/logrus"
)

// Application represents client application.
type Application struct {
	settings *ApplicationSettings
	logger   *logrus.Logger
}

// NewApplication creates new Application object.
func NewApplication(settngs *ApplicationSettings, logger *logrus.Logger) *Application {
	return &Application{settings: settngs, logger: logger}
}

func (a *Application) Start(ctx context.Context) error {
	return nil
}
