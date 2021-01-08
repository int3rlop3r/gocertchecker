package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"
)

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
	fmt.Println("connecting to:", addr)
	conn, err := tls.Dial("tcp", addr, &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't connect to host '%s', err: %s\n", addr, err)
		return
	}

	err = conn.VerifyHostname(host)
	if err != nil {
		fmt.Fprintf(os.Stderr, "hostname doesn't match cert, err: %s\n", err)
		return
	} else {
		fmt.Printf("hostname '%s' verified\n", host)
	}
	state := conn.ConnectionState()
	fmt.Println(conn.ConnectionState().Version)
	for _, pCert := range state.PeerCertificates {
		fmt.Println("Issuer:", pCert.Issuer)
	}
}
