package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v7"
	"github.com/devldavydov/gophkeeper/internal/client"
	"github.com/devldavydov/gophkeeper/internal/common/nettools"
)

var errInvalidSettings = errors.New("invalid settings")

const (
	_defaultConfigServerAddress = "127.0.0.1:8080"
	_defaultConfigCACert        = ""
	_defaultConfigLogLevel      = "INFO"
	_defaultConfigLogFile       = "client.log"
)

// Config represents command line/env client configuration options.
type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	CACert        string `env:"TLS_CA_CERT"`
	LogLevel      string `env:"LOG_LEVEL"`
	LogFile       string `env:"LOG_FILE"`
}

// LoadConfig loads server configuration from flags/env.
func LoadConfig(flagSet flag.FlagSet, flags []string) (*Config, error) {
	var err error
	config := &Config{}

	// Flags
	flagSet.StringVar(&config.ServerAddress, "a", _defaultConfigServerAddress, "server address")
	flagSet.StringVar(&config.CACert, "tlscacert", _defaultConfigCACert, "CA certificate")
	flagSet.StringVar(&config.LogLevel, "l", _defaultConfigLogLevel, "log level")
	flagSet.StringVar(&config.LogFile, "f", _defaultConfigLogFile, "log file")

	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flagSet.PrintDefaults()
	}

	_ = flagSet.Parse(flags)

	// Check env
	if err = env.Parse(config); err != nil {
		return nil, err
	}

	return config, nil
}

func ClientSettingsAdapt(config *Config) (*client.Settings, error) {
	serverAddress, err := nettools.NewAddress(config.ServerAddress)
	if err != nil {
		return nil, err
	}

	if config.CACert == "" {
		return nil, errInvalidSettings
	}

	return client.NewSettings(serverAddress, config.CACert), nil
}
