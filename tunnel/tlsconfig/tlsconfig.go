// Package tlsconfig provides convenience functions for configuring TLS connections from the
// command line.
package tlsconfig

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"
)

// Config is the user provided parameters to create a tls.Config
type TLSParams struct {
	Cert                 string
	Key                  string
	GetCertificate       *CertReloader
	GetClientCertificate *CertReloader
	ClientCAs            []string
	RootCAs              []string
	ServerName           string
	CurvePreferences     []tls.CurveID
	MinVersion           uint16 // min tls version. If zero, TLS1.0 is defined as minimum.
	MaxVersion           uint16 // max tls version. If zero, last TLS version is used defined as limit (currently TLS1.3)
}

func (p TLSParams) String() string {
	return fmt.Sprintf(
		`TLSParams{
	Cert: %q,
	Key: %q,
	GetCertificate: %s,
	GetClientCertificate: %s,
	ClientCAs: %s,
	RootCAs: %s,
	ServerName: %q,
	CurvePreferences: %v,
	MinVersion: %d,
	MaxVersion: %d,
}`,
		p.Cert,
		p.Key,
		p.GetCertificate.String(),
		p.GetClientCertificate.String(),
		strings.Join(p.ClientCAs, ", "),
		strings.Join(p.RootCAs, ", "),
		p.ServerName,
		p.CurvePreferences,
		p.MinVersion,
		p.MaxVersion,
	)
}

// GetTLSConfig returns a TLS configuration according to the GetTLSConfig set by the user.
func GetTLSConfig(params *TLSParams) (*tls.Config, error) {
	tlsconfig := &tls.Config{}
	if params.Cert != "" && params.Key != "" {
		cert, err := tls.LoadX509KeyPair(params.Cert, params.Key)
		if err != nil {
			return nil, fmt.Errorf("Error parsing X509 key pair: %v", err)
		}
		tlsconfig.Certificates = []tls.Certificate{cert}
		// BuildNameToCertificate parses Certificates and builds NameToCertificate from common name
		// and SAN fields of leaf certificates
		tlsconfig.BuildNameToCertificate()
	}

	if params.GetCertificate != nil {
		// GetCertificate is called when client supplies SNI info or Certificates is empty.
		// Order of retrieving certificate is GetCertificate, NameToCertificate and lastly first element of Certificates
		tlsconfig.GetCertificate = params.GetCertificate.Cert
	}

	if params.GetClientCertificate != nil {
		// GetClientCertificate is called when using an HTTP client library and mTLS is required.
		tlsconfig.GetClientCertificate = params.GetClientCertificate.ClientCert
	}

	if len(params.ClientCAs) > 0 {
		// set of root certificate authorities that servers use if required to verify a client certificate
		// by the policy in ClientAuth
		clientCAs, err := LoadCert(params.ClientCAs)
		if err != nil {
			return nil, fmt.Errorf("error loading client CAs: %v", err)
		}
		tlsconfig.ClientCAs = clientCAs
		// server's policy for TLS Client Authentication. Default is no client cert
		tlsconfig.ClientAuth = tls.RequireAndVerifyClientCert
	}

	if len(params.RootCAs) > 0 {
		rootCAs, err := LoadCert(params.RootCAs)
		if err != nil {
			return nil, fmt.Errorf("Error loading root CAs: %v", err)
		}
		tlsconfig.RootCAs = rootCAs
	}

	if params.ServerName != "" {
		tlsconfig.ServerName = params.ServerName
	}

	if len(params.CurvePreferences) > 0 {
		tlsconfig.CurvePreferences = params.CurvePreferences
	} else {
		// Cloudflare optimize CurveP256
		tlsconfig.CurvePreferences = []tls.CurveID{tls.CurveP256}
	}

	tlsconfig.MinVersion = params.MinVersion
	tlsconfig.MaxVersion = params.MaxVersion

	return tlsconfig, nil
}

// LoadCert creates a CertPool containing all certificates in a PEM-format file.
func LoadCert(certPaths []string) (*x509.CertPool, error) {
	ca := x509.NewCertPool()
	for _, certPath := range certPaths {
		caCert, err := os.ReadFile(certPath)
		if err != nil {
			return nil, fmt.Errorf("error reading certificate %s: %v", certPath, err)
		}
		if !ca.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("Error parsing certificate %s: %v", certPath, err)
		}
	}
	return ca, nil
}
