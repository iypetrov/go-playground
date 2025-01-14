package main

import (
	"fmt"
)

type Protocol int64

type TLSSettings struct {
	ServerName string
	NextProtos []string
}

// func createTunnelConfig(ctx context.Context, serverName string) (*tls.Config, error) {
// 	var rootCAs []string
//
// 	userConfig := &TLSParameters{RootCAs: rootCAs, ServerName: serverName}
// 	tlsConfig, err := GetConfig(userConfig)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if tlsConfig.RootCAs == nil {
// 		rootCAPool, err := x509.SystemCertPool()
// 		if err != nil {
// 			return nil, fmt.Errorf("unable to get x509 system cert pool")
// 		}
// 		cfRootCA, err := GetCloudflareRootCA()
// 		if err != nil {
// 			return nil, fmt.Errorf("could not append Cloudflare Root CAs to cloudflared certificate pool")
// 		}
// 		for _, cert := range cfRootCA {
// 			rootCAPool.AddCert(cert)
// 		}
// 		tlsConfig.RootCAs = rootCAPool
// 	}
//
// 	if tlsConfig.ServerName == "" && !tlsConfig.InsecureSkipVerify {
// 		return nil, fmt.Errorf("either ServerName or InsecureSkipVerify must be specified in the tls.Config")
// 	}
// 	return tlsConfig, nil
// }

func main() {
	// ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	// defer cancel()

	// tlsSettings := &TLSSettings{
	// 	// server name to establish http2 connection with edge
	// 	ServerName: "h2.cftunnel.com",
	// }

	// edgeTLSConfig, err := createTunnelConfig(ctx, tlsSettings.ServerName)
	// if err != nil {
	// 	fmt.Printf("unable to create TLS config to connect with edge: %v", err)
	// }

	// if len(tlsSettings.NextProtos) > 0 {
	// 	edgeTLSConfig.NextProtos = tlsSettings.NextProtos
	// }

	// var dialer net.Dialer
	// edgeConn, err := dialer.DialContext(ctx, "tcp", "localhost:12345")
	// if err != nil {
	// 	fmt.Printf("failed to dial: %v", err)
	// }
	// defer edgeConn.Close()

	// tlsEdgeConn := tls.Client(edgeConn, tlsConfig)
	// tlsEdgeConn.SetDeadline(time.Now().Add(timeout))

	// if _, err := conn.Write([]byte("Hello, World!")); err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println("hello tunnel")
}
