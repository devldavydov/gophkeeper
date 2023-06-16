package tls

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServerSettings(t *testing.T) {
	_, err := NewServerSettings("", "")
	assert.Error(t, err)

	_, err = NewServerSettings(getTLSFile("server-cert.pem"), getTLSFile("server-key.pem"))
	assert.NoError(t, err)
}

func TestServerSettingsLoad(t *testing.T) {
	s, err := NewServerSettings("/tmp/foo", "/tmp/bar")
	assert.NoError(t, err)

	_, err = s.Load()
	assert.Error(t, err)

	s, err = NewServerSettings(getTLSFile("server-cert.pem"), getTLSFile("server-key.pem"))
	assert.NoError(t, err)

	_, err = s.Load()
	assert.NoError(t, err)
}

func TestLoadCACert(t *testing.T) {
	_, err := LoadCACert("/tmp/foo", "test")
	assert.Error(t, err)

	_, err = LoadCACert(getTLSFile("server-key.pem"), "test")
	assert.Error(t, err)

	_, err = LoadCACert(getTLSFile("ca-cert.pem"), "test")
	assert.NoError(t, err)

}

func getTLSFile(fileName string) string {
	_, this, _, _ := runtime.Caller(0)
	tlsRoot := filepath.Join(this, "../../../../tls")
	return filepath.Join(tlsRoot, fileName)
}
