package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v7"
	"github.com/devldavydov/gophkeeper/internal/common/cipher"
	"github.com/devldavydov/gophkeeper/internal/common/nettools"
	gkTLS "github.com/devldavydov/gophkeeper/internal/common/tls"
	"github.com/devldavydov/gophkeeper/internal/server"
)

var errInvalidSettings = errors.New("invalid settings")

const (
	_defaultConfigGRPCAddress       = "127.0.0.1:8080"
	_defaultConfigGRPCServerTLSCert = ""
	_defaultConfigGRPCServerTLSKey  = ""
	_defaultConfigDatabaseDsn       = ""
	_defaultConfigLogLevel          = "INFO"
	_defaultConfigServerSecret      = "GophKeeperSupaSecretKeyForCrypto" //nolint:gosec // OK
	_defaultConfigShutdownTimeout   = 10 * time.Second
)

// Config is a command line/env server configuration options.
type Config struct {
	// GRPCAddress - listen address of server.
	// env: "GRPC_ADDRESS", flag: "a".
	GRPCAddress string `env:"GRPC_ADDRESS"`

	// GRPCServerTLSCert - TLS certificate of server.
	// env: "GRPC_SERVER_TLS_CERT", flag: "tlscert".
	GRPCServerTLSCert string `env:"GRPC_SERVER_TLS_CERT"`

	// GRPCServerTLSKey - TLS certificate key of server.
	// env: "GRPC_SERVER_TLS_KEY", flag: "tlskey".
	GRPCServerTLSKey string `env:"GRPC_SERVER_TLS_KEY"`

	// DatabaseDsn - database connection string.
	// env: "DATABASE_DSN", flag: "d".
	DatabaseDsn string `env:"DATABASE_DSN"`

	// LogLevel - logging level.
	// env: "LOG_LEVEL", flag: "l".
	LogLevel string `env:"LOG_LEVEL"`

	// ServerSecret - unique 32 chars string to be used as a key of encryption.
	// env: "SERVER_SECRET", flag: "s".
	ServerSecret string `env:"SERVER_SECRET"`

	// ShutdownTimeout - server shitdown timeout.
	// env: "SHUTDOWN_TIMEOUT", flag: "t".
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT"`
}

// LoadConfig loads server configuration from flags/env.
func LoadConfig(flagSet flag.FlagSet, flags []string) (*Config, error) {
	var err error
	config := &Config{}

	// Flags
	flagSet.StringVar(&config.GRPCAddress, "a", _defaultConfigGRPCAddress, "gRPC address")
	flagSet.StringVar(&config.GRPCServerTLSCert, "tlscert", _defaultConfigGRPCServerTLSCert, "gRPC server certificate")
	flagSet.StringVar(&config.GRPCServerTLSKey, "tlskey", _defaultConfigGRPCServerTLSKey, "gRPC server certificate key")
	flagSet.StringVar(&config.DatabaseDsn, "d", _defaultConfigDatabaseDsn, "database dsn")
	flagSet.StringVar(&config.LogLevel, "l", _defaultConfigLogLevel, "log level")
	flagSet.StringVar(&config.ServerSecret, "s", _defaultConfigServerSecret, "server secret")
	flagSet.DurationVar(&config.ShutdownTimeout, "t", _defaultConfigShutdownTimeout, "server shutdown timeout")

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

// ServiceSettingsAdapt adapts flag/en configuration to server service settings internal format.
// Returns error "errInvalidSettings" in case of invalid configuration.
func ServiceSettingsAdapt(config *Config) (*server.ServiceSettings, error) {
	grpcAddress, err := nettools.NewAddress(config.GRPCAddress)
	if err != nil {
		return nil, err
	}

	grpcServerTLS, err := gkTLS.NewServerSettings(config.GRPCServerTLSCert, config.GRPCServerTLSKey)
	if err != nil {
		return nil, err
	}

	if config.DatabaseDsn == "" {
		return nil, errInvalidSettings
	}

	if len(config.ServerSecret) != cipher.AESKeyLength {
		return nil, errInvalidSettings
	}

	serverSettings := server.NewServiceSettings(
		grpcAddress,
		grpcServerTLS,
		config.DatabaseDsn,
		config.ServerSecret,
		config.ShutdownTimeout,
	)
	return serverSettings, nil
}
