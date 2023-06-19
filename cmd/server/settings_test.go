package main

import (
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServiceSettingsAdaptFromEnv(t *testing.T) {
	t.Setenv("GRPC_ADDRESS", "127.0.0.1:8080")
	t.Setenv("GRPC_SERVER_TLS_CERT", "/tmp/tls.cert")
	t.Setenv("GRPC_SERVER_TLS_KEY", "/tmp/tls.key")
	t.Setenv("DATABASE_DSN", "postgre:5432")
	t.Setenv("LOG_LEVEL", "DEBUG")
	t.Setenv("SERVER_SECRET", "asuperstrong32bitpasswordgohere!")
	t.Setenv("SHUTDOWN_TIMEOUT", "1s")

	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{})
	assert.NoError(t, err)

	serviceSettings, err := ServiceSettingsAdapt(config)
	assert.NoError(t, err)

	assert.Equal(t, "127.0.0.1:8080", serviceSettings.GRPCAddress.String())
	assert.Equal(t, "/tmp/tls.cert", serviceSettings.GRPCServerTLS.ServerCertPath)
	assert.Equal(t, "/tmp/tls.key", serviceSettings.GRPCServerTLS.ServerKeyPath)
	assert.Equal(t, "postgre:5432", serviceSettings.DatabaseDsn)
	assert.Equal(t, "asuperstrong32bitpasswordgohere!", serviceSettings.ServerSecret)
	assert.Equal(t, 1*time.Second, serviceSettings.ShutdownTimeout)
}

func TestServiceSettingsAdaptFromFlag(t *testing.T) {
	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{
		"-a", "127.0.0.1:8080",
		"-tlscert", "/tmp/tls.cert",
		"-tlskey", "/tmp/tls.key",
		"-d", "postgre:5432",
		"-l", "DEBUG",
		"-s", "asuperstrong32bitpasswordgohere!",
		"-t", "1s",
	})
	assert.NoError(t, err)

	serviceSettings, err := ServiceSettingsAdapt(config)
	assert.NoError(t, err)

	assert.Equal(t, "127.0.0.1:8080", serviceSettings.GRPCAddress.String())
	assert.Equal(t, "/tmp/tls.cert", serviceSettings.GRPCServerTLS.ServerCertPath)
	assert.Equal(t, "/tmp/tls.key", serviceSettings.GRPCServerTLS.ServerKeyPath)
	assert.Equal(t, "postgre:5432", serviceSettings.DatabaseDsn)
	assert.Equal(t, "asuperstrong32bitpasswordgohere!", serviceSettings.ServerSecret)
	assert.Equal(t, 1*time.Second, serviceSettings.ShutdownTimeout)
}

func TestServiceSettingsAdaptWithDefault(t *testing.T) {
	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{
		"-tlscert", "/tmp/tls.cert",
		"-tlskey", "/tmp/tls.key",
		"-d", "postgre:5432",
	})
	assert.NoError(t, err)

	serviceSettings, err := ServiceSettingsAdapt(config)
	assert.NoError(t, err)

	assert.Equal(t, "127.0.0.1:8080", serviceSettings.GRPCAddress.String())
	assert.Equal(t, "/tmp/tls.cert", serviceSettings.GRPCServerTLS.ServerCertPath)
	assert.Equal(t, "/tmp/tls.key", serviceSettings.GRPCServerTLS.ServerKeyPath)
	assert.Equal(t, "postgre:5432", serviceSettings.DatabaseDsn)
	assert.Equal(t, "GophKeeperSupaSecretKeyForCrypto", serviceSettings.ServerSecret)
	assert.Equal(t, 10*time.Second, serviceSettings.ShutdownTimeout)
}

func TestServiceSettingsAdaptError(t *testing.T) {
	for i, tt := range []struct {
		flags   []string
		env     map[string]string
		loadErr bool
	}{
		{flags: []string{}, env: map[string]string{}},
		{flags: []string{"-a", ""}, env: map[string]string{}},
		{flags: []string{"-s", "123"}, env: map[string]string{}},
		{
			flags: []string{"-d", "", "-tlscert", "/tmp/tls.cert", "-tlskey", "/tmp/tls.key"},
			env:   map[string]string{},
		},
		{
			flags: []string{},
			env: map[string]string{
				"SHUTDOWN_TIMEOUT": "foobar",
			},
			loadErr: true},
	} {
		tt := tt
		t.Run(fmt.Sprintf("Run %d", i), func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
			config, err := LoadConfig(*testFlagSet, tt.flags)
			if tt.loadErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			_, err = ServiceSettingsAdapt(config)
			assert.Error(t, err)
		})
	}
}
