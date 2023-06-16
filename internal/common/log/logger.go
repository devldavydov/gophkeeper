// Package log contains util functions for logging.
package log

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// NewLogger creates logger.
func NewLogger(logLevel string) (*logrus.Logger, error) {
	logger := logrus.New()
	logLvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, fmt.Errorf("wrong LOG_LEVEL: %w", err)
	}

	logger.SetLevel(logLvl)
	return logger, nil
}
