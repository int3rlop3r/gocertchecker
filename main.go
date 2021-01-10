package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"os"
)

func sslVersionStr(v uint16) (string, error) {
	switch v {
	case tls.VersionTLS10:
		return "TLS 1.0", nil
	case tls.VersionTLS11:
		return "TLS 1.1", nil
	case tls.VersionTLS12:
		return "TLS 1.2", nil
	case tls.VersionTLS13:
		return "TLS 1.3", nil
	case tls.VersionSSL30:
		return "TLS 3.0", nil
	default:
		return "", fmt.Errorf("%x not a valid ssl version")

	}
}

func main() {
	var port string
	flag.StringVar(&port, "port", "443", "Port no. of the remote service")
	flag.Parse()
	host := flag.Arg(0)
	if host == "" {
		fmt.Fprintln(os.Stderr, "hostname cannot be blank")
		return
	}

	var config tls.Config
	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("Connecting to: '%s'. ", addr)
	conn, err := tls.Dial("tcp", addr, &config)
	if _, ok := err.(x509.UnknownAuthorityError); ok {
		config.InsecureSkipVerify = true
		fmt.Printf("Unverified hostname '%s'\n", host)
		conn, err = tls.Dial("tcp", addr, &config)
	} else if err == nil {
		fmt.Printf("Verified hostname '%s'\n", host)
	}

	// catches both errors
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't connect to host '%s', err: %s\n", addr, err)
		return
	}

	state := conn.ConnectionState()
	version, err := sslVersionStr(state.Version)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("SSL/TLS version:", version)
	}

	for _, pCert := range state.PeerCertificates {
		fmt.Println("Issuer:", pCert.Issuer)
	}
}
