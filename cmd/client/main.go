// Copyright (c) 2025 Cloudflare, Inc.
// Licensed under the Apache 2.0 license found in the LICENSE file or at:
//     https://opensource.org/licenses/Apache-2.0

package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	certPath       = flag.String("server-cacert", "", "path to root CA certificate to validate the server certificate, i.e. rootCA.crt")
	clientCertPath = flag.String("cert", "", "path to client certificate")
	clientKeyPath  = flag.String("key", "", "path to client key")
	addr           = flag.String("connect", "https://localhost:8775/", "address to connect to")
	tlsflags       = flag.String("tlsflags", "", "a comma seperated list of TLS Flags")
)

func main() {
	flag.Parse()
	var rootPool *x509.CertPool
	var err error
	if *certPath != "" {
		certPEM, err := os.ReadFile(*certPath)
		if err != nil {
			log.Fatalf("Couldn't read cert")
		}

		block, rest := pem.Decode(certPEM)
		if len(rest) != 0 {
			log.Fatalf("excess data in cert")
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			log.Fatalf("Couldn't parse cert")
		}

		rootPool = x509.NewCertPool()
		rootPool.AddCert(cert)
	} else {
		rootPool, err = x509.SystemCertPool()
		if err != nil {
			log.Fatalf("Couldn't load system root pool")
		}
	}

	var keylog io.Writer
	if os.Getenv("SSLKEYLOGFILE") != "" {
		keylog, err = os.OpenFile(os.Getenv("SSLKEYLOGFILE"), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			log.Fatalf("couldn't write to sslkeylog: %v", err)
		}
	}

	certs := make([]tls.Certificate, 0, 1)
	if *clientCertPath != "" && *clientKeyPath != "" {
		clientCert, err := tls.LoadX509KeyPair(*clientCertPath, *clientKeyPath)
		if err != nil {
			log.Fatalf("couldn't load client cert: %v", err)
		}
		certs = append(certs, clientCert)
	} else {
		log.Fatalf("couldn't load client cert")
	}

	flags := make([]tls.TLSFlag, 0)
	if *tlsflags != "" {
		flagStrings := strings.Split(*tlsflags, ",")
		for _, flagString := range flagStrings {
			flag, err := strconv.ParseUint(flagString, 0, 16)
			if err != nil {
				log.Fatalf("Invalid TLS Flag: %s", flagString)
			}
			flags = append(flags, tls.TLSFlag(uint16(flag)))
		}
	}

	config := &tls.Config{
		TLSFlagsSupported: flags,
		RootCAs:           rootPool,
		KeyLogWriter:      keylog,
		Certificates:      certs,
		NextProtos:        []string{"h2"},
	}

	req, err := http.NewRequest(http.MethodGet, *addr, nil)
	if err != nil {
		log.Fatal("Could not create request")
	}

	trans := http.Transport{TLSClientConfig: config, ForceAttemptHTTP2: true}
	resp, err := trans.RoundTrip(req)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Couldn't read response body: %v", err)
	}

	fmt.Printf("%s\n", string(respBody))
}
