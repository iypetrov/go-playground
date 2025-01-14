package tlsconfig

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"sync"
)

const (
	OriginCAPoolFlag = "origin-ca-pool"
	CaCertFlag       = "cacert"
)

// CertReloader can load and reload a TLS certificate from a particular filepath.
// Hooks into tls.Config's GetCertificate to allow a TLS server to update its certificate without restarting.
type CertReloader struct {
	sync.Mutex
	certificate *tls.Certificate
	certPath    string
	keyPath     string
}

func (cr *CertReloader) String() string {
	return fmt.Sprintf(
		`CertReloader{
	CertPath: %q,
	KeyPath: %q,
	CertificateLoaded: %t,
}`,
		cr.certPath,
		cr.keyPath,
		cr.certificate != nil,
	)
}

// NewCertReloader makes a CertReloader. It loads the cert during initialization to make sure certPath and keyPath are valid
func NewCertReloader(certPath, keyPath string) (*CertReloader, error) {
	cr := new(CertReloader)
	cr.certPath = certPath
	cr.keyPath = keyPath
	if err := cr.LoadCert(); err != nil {
		return nil, err
	}
	return cr, nil
}

// Cert returns the TLS certificate most recently read by the CertReloader.
// This method works as a direct utility method for tls.Config#Cert.
func (cr *CertReloader) Cert(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	cr.Lock()
	defer cr.Unlock()
	return cr.certificate, nil
}

// ClientCert returns the TLS certificate most recently read by the CertReloader.
// This method works as a direct utility method for tls.Config#ClientCert.
func (cr *CertReloader) ClientCert(certRequestInfo *tls.CertificateRequestInfo) (*tls.Certificate, error) {
	cr.Lock()
	defer cr.Unlock()
	return cr.certificate, nil
}

// LoadCert loads a TLS certificate from the CertReloader's specified filepath.
// Call this after writing a new certificate to the disk (e.g. after renewing a certificate)
func (cr *CertReloader) LoadCert() error {
	cr.Lock()
	defer cr.Unlock()

	cert, err := tls.LoadX509KeyPair(cr.certPath, cr.keyPath)

	// Keep the old certificate if there's a problem reading the new one.
	if err != nil {
		log.Printf("Error parsing X509 key pair: %v", err)
		return err
	}
	cr.certificate = &cert
	return nil
}

func LoadOriginCA(originCAPoolFilename string) (*x509.CertPool, error) {
	var originCustomCAPool []byte

	if originCAPoolFilename != "" {
		var err error
		originCustomCAPool, err = os.ReadFile(originCAPoolFilename)
		if err != nil {
			return nil, fmt.Errorf("unable to read the file %s for --%s: %v", originCAPoolFilename, OriginCAPoolFlag, err)
		}
	}

	originCertPool, err := loadOriginCertPool(originCustomCAPool)
	if err != nil {
		return nil, fmt.Errorf("error loading the certificate pool: %v", err)
	}

	return originCertPool, nil
}

func LoadCustomOriginCA(originCAFilename string) (*x509.CertPool, error) {
	// First, obtain the system certificate pool
	certPool, err := x509.SystemCertPool()
	if err != nil {
		certPool = x509.NewCertPool()
	}

	// Next, append the Cloudflare CAs into the system pool
	cfRootCA, err := GetCloudflareRootCA()
	if err != nil {
		return nil, fmt.Errorf("could not append Cloudflare Root CAs to cloudflared certificate pool: %v", err)
	}
	for _, cert := range cfRootCA {
		certPool.AddCert(cert)
	}

	if originCAFilename == "" {
		return certPool, nil
	}

	customOriginCA, err := os.ReadFile(originCAFilename)
	if err != nil {
		return nil, fmt.Errorf("unable to read the file %s: %v", originCAFilename, err)
	}

	if !certPool.AppendCertsFromPEM(customOriginCA) {
		return nil, fmt.Errorf("error appending custom CA to cert pool")
	}
	return certPool, nil
}

func CreateTunnelConfig(serverName string) (*tls.Config, error) {
	var rootCAs []string

	userConfig := &TLSParams{RootCAs: rootCAs, ServerName: serverName}
	tlsConfig, err := Config(userConfig)
	if err != nil {
		return nil, err
	}

	if tlsConfig.RootCAs == nil {
		rootCAPool, err := x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("unable to get x509 system cert pool: %v", err)
		}
		cfRootCA, err := GetCloudflareRootCA()
		if err != nil {
			return nil, fmt.Errorf("could not append Cloudflare Root CAs to cloudflared certificate pool: %v", err)
		}
		for _, cert := range cfRootCA {
			rootCAPool.AddCert(cert)
		}
		tlsConfig.RootCAs = rootCAPool
	}

	if tlsConfig.ServerName == "" && !tlsConfig.InsecureSkipVerify {
		return nil, fmt.Errorf("either ServerName or InsecureSkipVerify must be specified in the tls.Config")
	}
	return tlsConfig, nil
}

func loadOriginCertPool(originCAPoolPEM []byte) (*x509.CertPool, error) {
	// Get the global pool
	certPool, err := loadGlobalCertPool()
	if err != nil {
		return nil, err
	}

	// Then, add any custom origin CA pool the user may have passed
	if originCAPoolPEM != nil {
		if !certPool.AppendCertsFromPEM(originCAPoolPEM) {
			log.Println("could not append the provided origin CA to the cloudflared certificate pool")
		}
	}

	return certPool, nil
}

func loadGlobalCertPool() (*x509.CertPool, error) {
	// First, obtain the system certificate pool
	certPool, err := x509.SystemCertPool()
	if err != nil {
		certPool = x509.NewCertPool()
	}

	// Next, append the Cloudflare CAs into the system pool
	cfRootCA, err := GetCloudflareRootCA()
	if err != nil {
		return nil, fmt.Errorf("could not append Cloudflare Root CAs to cloudflared certificate pool: %v", err)
	}
	for _, cert := range cfRootCA {
		certPool.AddCert(cert)
	}

	// Finally, add the Hello certificate into the pool (since it's self-signed)
	helloCert, err := GetHelloCertificateX509()
	if err != nil {
		return nil, fmt.Errorf("could not append Hello server certificate to cloudflared certificate pool: %v", err)
	}
	certPool.AddCert(helloCert)

	return certPool, nil
}
