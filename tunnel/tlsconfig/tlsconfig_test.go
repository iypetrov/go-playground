package tlsconfig

import (
	"crypto/tls"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// testcert.pem and testcert2.pem are Generated using `openssl req -newkey rsa:512 -nodes -x509 -days 3650`
const (
	testcertCommonName = "localhost"
)

func TestGetFromEmptyConfig(t *testing.T) {
	tlsParams := &TLSParams{}

	tlsConfig, err := GetTLSConfig(tlsParams)
	fmt.Printf(tlsConfig)
	assert.NoError(t, err)
	assert.Empty(t, tlsConfig.Certificates)

	assert.Empty(t, tlsConfig.NameToCertificate)

	assert.Nil(t, tlsConfig.ClientCAs)
	assert.Equal(t, tls.NoClientCert, tlsConfig.ClientAuth)

	assert.Nil(t, tlsConfig.RootCAs)

	assert.Len(t, tlsConfig.CurvePreferences, 1)
	assert.Equal(t, tls.CurveP256, tlsConfig.CurvePreferences[0])
}

func TestGetConfig(t *testing.T) {
	cert, err := tls.LoadX509KeyPair("testcert.pem", "testkey.pem")
	assert.NoError(t, err)

	params := &TLSParams{
		Cert:             "testcert.pem",
		Key:              "testkey.pem",
		ClientCAs:        []string{"testcert.pem", "testcert2.pem"},
		RootCAs:          []string{"testcert.pem", "testcert2.pem"},
		ServerName:       "test",
		CurvePreferences: []tls.CurveID{tls.CurveP384},
	}
	tlsConfig, err := GetTLSConfig(params)
	assert.NoError(t, err)
	assert.Len(t, tlsConfig.Certificates, 1)
	assert.Equal(t, cert, tlsConfig.Certificates[0])

	assert.Equal(t, cert, *tlsConfig.NameToCertificate[testcertCommonName])

	assert.NotNil(t, tlsConfig.ClientCAs)
	assert.Equal(t, tls.RequireAndVerifyClientCert, tlsConfig.ClientAuth)

	assert.NotNil(t, tlsConfig.RootCAs)

	assert.Len(t, tlsConfig.CurvePreferences, 1)
	assert.Equal(t, tls.CurveP384, tlsConfig.CurvePreferences[0])
}

func TestCertReloader(t *testing.T) {
	expectedCert, err := tls.LoadX509KeyPair("testcert.pem", "testkey.pem")
	assert.NoError(t, err)

	certReloader, err := NewCertReloader("testcert.pem", "testkey.pem")
	assert.NoError(t, err)

	chi := &tls.ClientHelloInfo{ServerName: testcertCommonName}
	cert, err := certReloader.Cert(chi)
	assert.NoError(t, err)
	assert.Equal(t, expectedCert, *cert)

	params := &TLSParams{
		GetCertificate: certReloader,
	}
	tlsConfig, err := GetTLSConfig(params)
	assert.NoError(t, err)

	cert, err = tlsConfig.GetCertificate(chi)
	assert.NoError(t, err)
	assert.Equal(t, expectedCert, *cert)
}
