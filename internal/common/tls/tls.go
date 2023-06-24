// Package tls contains utils for crypto TLS.
package tls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials"
)

// ServerSettings represents TLS server settings.
type ServerSettings struct {
	// Server certificate file path
	ServerCertPath string

	// Server certificate key file path
	ServerKeyPath string
}

// NewServerSettings creates new ServerSettings object.
//
// Input arguments - certificate and key file paths.
func NewServerSettings(certPath, keyPath string) (*ServerSettings, error) {
	if certPath == "" || keyPath == "" {
		return nil, errors.New("invalid TLS server settings")
	}

	return &ServerSettings{ServerCertPath: certPath, ServerKeyPath: keyPath}, nil
}

// Load creates new TransportCredentials from certificate and key file paths.
//
// In case of invalid format returns specific TLS error.
func (t *ServerSettings) Load() (credentials.TransportCredentials, error) {
	serverCert, err := tls.LoadX509KeyPair(t.ServerCertPath, t.ServerKeyPath)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}

// LoadCACert loads Certificate Authority certificate from file.
func LoadCACert(caCertPath string, customServerName string) (credentials.TransportCredentials, error) {
	pemServerCA, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	config := &tls.Config{
		MinVersion: tls.VersionTLS12,
		RootCAs:    certPool,
	}
	if customServerName != "" {
		config.ServerName = customServerName
	}

	return credentials.NewTLS(config), nil
}
