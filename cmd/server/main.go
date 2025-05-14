// Copyright (c) 2025 Cloudflare, Inc.
// Licensed under the Apache 2.0 license found in the LICENSE file or at:
//     https://opensource.org/licenses/Apache-2.0

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/cloudflare/certinel/fswatcher" // allow to watch for certificate file modification automatically
	"github.com/oklog/run"
)

var (
	certPath   = flag.String("cert", "", "path to the server certificate, i.e. server.crt")
	keyPath    = flag.String("key", "", "path to certificate private key, i.e. server.key")
	caCertPath = flag.String("client-cacert", "", "path to root CA certificate to validate the client certificate, i.e. rootCA.crt")
	addr       = flag.String("listen-addr", "localhost:8775", "listening address")
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	flag.Parse()
	cacertPEM, err := os.ReadFile(*caCertPath)
	if err != nil {
		log.Fatalf("Couldn't read cert 1")
	}

	block, rest := pem.Decode(cacertPEM)
	if len(rest) != 0 {
		log.Fatalf("excess data in cert")
	}

	cAcert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatalf("Couldn't parse cert")
	}

	clientRootPool := x509.NewCertPool()
	clientRootPool.AddCert(cAcert)

	certinel, err := fswatcher.New(*certPath, *keyPath)
	if err != nil {
		log.Fatalf("Couldn't load cert")
	}
	g := run.Group{}
	{
		g.Add(func() error {
			return certinel.Start(ctx)
		}, func(err error) {
			cancel()
		})
	}
	{
		config := &tls.Config{
			GetCertificate:    certinel.GetCertificate,
			ClientCAs:         clientRootPool,
			TLSFlagsSupported: []tls.TLSFlag{tls.ExperimentalFlagSupportMTLS},
			NextProtos:        []string{"h2"},
		}
		conn, err := tls.Listen("tcp", *addr, config)
		if err != nil {
			log.Fatalf("Couldn't listen on %s", *addr)
		}
		mux := http.NewServeMux()
		mux.HandleFunc("/debug", debugHandler(clientRootPool))
		mux.HandleFunc("/", indexHandler())
		g.Add(
			func() error {
				return http.Serve(conn, mux)
			},
			func(err error) {
				conn.Close()
			},
		)
	}

	if err := g.Run(); err != nil {
		log.Fatalf("err='%s'", err)
	}
}
