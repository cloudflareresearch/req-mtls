// Copyright (c) 2025 Cloudflare, Inc.
// Licensed under the Apache 2.0 license found in the LICENSE file or at:
//     https://opensource.org/licenses/Apache-2.0

package main

import (
	"crypto/x509"
	_ "embed"
	"fmt"
	"log"
	"net/http"
)

func debugHandler(clientRootPool *x509.CertPool) func(http.ResponseWriter, *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(resp, "Client sent TLS Flags: %v\n", req.TLS.PeerTLSFlags)

		fmt.Fprintf(resp, "Mutually supported TLS Flags: %v\n", req.TLS.AgreedTLSFlags)

		fmt.Fprintf(resp, "Client cert requested: %v\n", req.TLS.RequestClientCert)

		if len(req.TLS.PeerCertificates) != 0 {
			fmt.Fprintf(resp, "Peer certificate received: true\n")
			fmt.Fprintf(resp, "Cert received: %s\n", req.TLS.PeerCertificates[0].DNSNames[0])
			cert := req.TLS.PeerCertificates[0]
			opts := x509.VerifyOptions{
				Roots:     clientRootPool,
				KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			}
			chains, err := cert.Verify(opts)
			if err != nil {
				log.Print(err)
			}

			fmt.Fprintf(resp, "Cert validated: %v\n", len(chains) != 0)
		} else {
			fmt.Fprintf(resp, "Peer certificate received: false\n")
		}
	}
}

//go:embed index.html
var indexHTMLTemplate string

func indexHandler() func(resp http.ResponseWriter, req *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(resp, indexHTMLTemplate)
	}
}
