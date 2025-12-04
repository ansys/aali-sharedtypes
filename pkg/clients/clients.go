package clients

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/ansys/aali-sharedtypes/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// GetHttpClientWithCert creates an HTTP client configured with TLS using a custom SSL certificate.
//
// Returns:
//   - httpClient: Pointer to http.Client configured with TLS.
//   - err: an error message if the setup fails.
func GetHttpClient() (httpClient *http.Client, err error) {
	if config.GlobalConfig.USE_SSL {
		// attach custom certificate to HTTP client
		tlsConfig, err := GetTlsConfigWithCert()
		if err != nil {
			return nil, fmt.Errorf("failed to get TLS config with cert: %v", err)
		}

		transport := &http.Transport{
			TLSClientConfig: tlsConfig,
		}

		httpClient = &http.Client{
			Transport: transport,
		}
	} else {
		httpClient = &http.Client{}
	}

	return httpClient, nil
}

// GetGrpcDialOptions creates gRPC dial options with custom dialing logic and transport credentials based on the scheme.
//
// Parameters:
//   - scheme: A string indicating the connection scheme ("http" or "https").
//
// Returns:
//   - options: A slice of grpc.DialOption configured for the connection.
//   - err: an error message if the setup fails.
func GetGrpcDialOptions(scheme string) (options []grpc.DialOption, err error) {
	// Add custom dialer with IPv4 first, fallback to IPv6
	options = append(options, grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
		d := &net.Dialer{}

		// Try IPv4 first
		conn, err := d.DialContext(ctx, "tcp4", addr)
		if err == nil {
			return conn, nil
		}

		// Fall back to IPv6 if IPv4 fails
		return d.DialContext(ctx, "tcp6", addr)
	}))

	// Set up transport credentials based on the scheme
	if scheme == "https" {
		// Set up a secure connection
		var tlsConfig *tls.Config
		if config.GlobalConfig.USE_SSL {
			tlsConfig, err = GetTlsConfigWithCert()
			if err != nil {
				return nil, fmt.Errorf("unable to set up TLS config with custom certificate: %v", err)
			}
		}
		creds := credentials.NewTLS(tlsConfig)
		options = append(options, grpc.WithTransportCredentials(creds))
	} else {
		// Set up an insecure connection
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	return options, nil
}

// getTlsConfigWithCert sets up a TLS configuration using a custom SSL certificate.
//
// Returns:
//   - tlsConfig: Pointer to tls.Config configured with the custom certificate.
//   - err: an error message if the setup fails.
func GetTlsConfigWithCert() (tlsConfig *tls.Config, err error) {
	certPool, err := GetCertPool()
	if err != nil {
		return nil, fmt.Errorf("failed to get cert pool: %v", err)
	}

	tlsConfig = &tls.Config{
		RootCAs: certPool,
	}

	return tlsConfig, nil
}

// GetCertPool reads the SSL certificate from the configured file and creates a certificate pool.
//
// Returns:
//   - certPool: Pointer to x509.CertPool containing the loaded certificate.
//   - err: an error message if the setup fails.
func GetCertPool() (certPool *x509.CertPool, err error) {
	certPEM, err := os.ReadFile(config.GlobalConfig.SSL_CERT_PUBLIC_KEY_FILE)
	if err != nil {
		return nil, fmt.Errorf("failed to read SSL certificate public key file: %v", err)
	}

	certPool = x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(certPEM) {
		return nil, fmt.Errorf("failed to append certificate to CA pool")
	}

	return certPool, nil
}
