package main

import (
	"flag"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplicatrionSettingsAdaptFromEnv(t *testing.T) {
	t.Setenv("SERVER_ADDRESS", "127.0.0.1:8888")
	t.Setenv("TLS_CA_CERT", "/tmp/ca.cert")
	t.Setenv("LOG_LEVEL", "DEBUG")
	t.Setenv("LOG_FILE", "log.log")

	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{})
	assert.NoError(t, err)

	cltSettings, err := ClientSettingsAdapt(config)
	assert.NoError(t, err)

	assert.Equal(t, "127.0.0.1:8888", cltSettings.ServerAddress.String())
	assert.Equal(t, "/tmp/ca.cert", cltSettings.TLSCACertPath)
}

func TestApplicationSettingsAdaptFromFlag(t *testing.T) {
	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{
		"-a", "127.0.0.1:8888",
		"-tlscacert", "/tmp/ca.cert",
		"-l", "DEBUG",
		"-f", "log.log",
	})
	assert.NoError(t, err)

	cltSettings, err := ClientSettingsAdapt(config)
	assert.NoError(t, err)

	assert.Equal(t, "127.0.0.1:8888", cltSettings.ServerAddress.String())
	assert.Equal(t, "/tmp/ca.cert", cltSettings.TLSCACertPath)
}

func TestApplicationSettingsWithDefault(t *testing.T) {
	testFlagSet := flag.NewFlagSet("test", flag.ExitOnError)
	config, err := LoadConfig(*testFlagSet, []string{
		"-tlscacert", "/tmp/ca.cert",
	})
	assert.NoError(t, err)

	cltSettings, err := ClientSettingsAdapt(config)
	assert.NoError(t, err)

	assert.Equal(t, "127.0.0.1:8080", cltSettings.ServerAddress.String())
	assert.Equal(t, "/tmp/ca.cert", cltSettings.TLSCACertPath)
}

func TestApplicationSettingsAdaptError(t *testing.T) {
	for i, tt := range []struct {
		flags   []string
		env     map[string]string
		loadErr bool
	}{
		{flags: []string{}, env: map[string]string{}},
		{flags: []string{"-a", ""}, env: map[string]string{}},
		{flags: []string{"-tlscacert", ""}, env: map[string]string{}},
		{flags: []string{}, env: map[string]string{"TLS_CA_CERT": ""}},
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

			_, err = ClientSettingsAdapt(config)
			assert.Error(t, err)
		})
	}
}
