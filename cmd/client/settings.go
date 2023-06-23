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

// Config is a command line/env client configuration options.
// Options:
//   - ServerAddress - address of server to connect.
//     env: "SERVER_ADDRESS", flag: "a".
//   - CACert - TLS Certification Authority certificate file.
//     env: "TLS_CA_CERT", flag: "tlscacert".
//   - LogLevel - logging level.
//     env: "LOG_LEVEL", flag: "l".
//   - LogFile - file to log to.
//     env: "LOG_FILE", flag: "f".
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

// ClientSettingsAdapt adapts flag/env configuration settings to client settings internal format.
// Returns error "errInvalidSettings" in case of invalid configuration.
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
