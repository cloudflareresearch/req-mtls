# req-mtls

![GitHub License](https://img.shields.io/github/license/cloudflareresearch/req-mtls)

Repository presenting authentication for orchestrated agents navigating the web. It implements all components required by request mTLS defined by [draft-jhoyla-req-mtls-flag](https://datatracker.ietf.org/doc/draft-jhoyla-req-mtls-flag/), and presents examples.

## Tables of Content

- [Examples](#examples)
- [Usage](#usage)
  - [Prerequisites](#prerequisites)
  - [Connecting to a remote test server](#connecting-to-the-demo-server)
  - [Setup a local environment](#setup-a-local-environment)
  - [Connecting to your local environment](#connecting-to-your-local-environment)
  - [Setup without Docker](#setup-without-docker)
- [Security Considerations](#security-considerations)
- [License](#license)

## Examples

We provide example code in Go in the [cmd](./cmd) folder.
This leverages the Cloudflare Research implementation of request mTLS as part of a Go fork in [cloudflare/go](https://github.com/cloudflare/go).

For convenience, a demonstration is provided on [req-mtls.research.cloudflare.com](https://req-mtls.research.cloudflare.com), which you can test against.

In the assets folder, you can find a [packet capture](./assets/demonstration-capture.pcapng) demonstrating a successful req mTLS connection.
The packet was captured with the [`SSLKEYLOGFILE`](https://everything.curl.dev/usingcurl/tls/sslkeylogfile.html) environment variable set.

## Usage

### Prerequisites

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Make](https://www.gnu.org/software/make/)
- [OpenSSL](https://www.openssl.org/)

### Connecting to a remote test server

The demonstration client and server are both built and run inside docker containers.
You can connect to [req-mtls.research.cloudflare.com](https://req-mtls.research.cloudflare.com) demonstration server by running the following command.

```shell
make client-request
```

This will spin up a test client in a docker container and connect to our demonstration server.

### Setup a local environment

When interacting with Cloudflare demonstration, certificates are already provided.

To do the same with a local demo, you need to generate these certificates yourself.

You can do this with

```shell
make gen-cert
```

This creates a client CA and generate a client leaf cert from it.
The default is to use the [RFC 9500](https://www.rfc-editor.org/rfc/rfc9500) test keys for the client CA key.
This allows anyone to create client certificates that will correctly chain to the client CA in the demo.

### Connecting to your local environment

You can run the local test server with

```shell
make start-test-server
```

This then allows you to run a client request hitting the /debug endpoint of the server

```shell
make client-request-local
```

### Setup without Docker

If you want to run the client and server directly as opposed to in docker you can build them locally by

```shell
make artefacts
```

This will build the client and server binaries and place them in the `bin` directory.

For example to connect to Cloudflare's demonstration site and request mutually authenticated TLS you can run:

```shell
./bin/client \
  -cert "./testdata/certs/client.crt" \
  -key "./testdata/certs/client.key" \
  -connect "https://req-mtls.research.cloudflare.com/debug" \
  -tlsflags 80
```

You can send other flags than the req mTLS flag

```shell
./bin/client \
  -cert "./testdata/certs/client.crt" \
  -key "./testdata/certs/client.key" \
  -connect "https://req-mtls.research.cloudflare.com/debug" \
  -tlsflags 8,170,12
```

Or set a custom root certificate for the server

```shell
./bin/client \
  -cert "./testdata/certs/client.crt" \
  -key "./testdata/certs/client.key" \
  -connect "https://req-mtls.research.cloudflare.com/debug" \
  -server-cacert "./testdata/certs/rootCA.crt" \
  -tlsflags 80
```

## Security Considerations

This software has not been audited. Please use at your sole discretion.

### Cryptographic keys

All keys shared publicly as part of this repository are test keys provided by [RFC 9500](https://www.rfc-editor.org/rfc/rfc9500.html), and formatted in [thibmeu/rfc9500](https://github.com/thibmeu/rfc9500/tree/main).
They MUST NOT be used for production services.

### TLS certificate provision

The repository include [certbot-to-R2](./third_party/certbot-to-r2/README.md), a tool for generating a certificate for a given FQDN and uploading it to a [Cloudflare R2](https://developers.cloudflare.com/r2/) bucket.
This is useful for deploying a demo. We advise against the use of this service in production.

## License

This project is under the Apache 2.0 license.

### Contribution

Unless you explicitly state otherwise, any contribution intentionally submitted for inclusion in the work by you shall be Apache 2.0 licensed as above, without any additional terms or conditions.
