Q ?= @

.PHONY: build
build: build-client build-server

.PHONY: build-client
build-client:
	$(Q)docker compose build client

.PHONY: build-server
build-server:
	$(Q)docker compose build server

.PHONY: gen-cert
gen-cert:
	$(Q)openssl req -x509 -new -nodes -key testdata/certs/rootCA.key -sha256 -days 3650 -out testdata/certs/rootCA.crt -config testdata/certs/rootCA.cnf
	$(Q)openssl req -new -key testdata/certs/client.key -out testdata/certs/client.csr -config testdata/certs/client.cnf
	$(Q)openssl x509 -req -in testdata/certs/client.csr -CA testdata/certs/rootCA.crt -CAkey testdata/certs/rootCA.key -CAcreateserial -out testdata/certs/client.crt -days 3650 -sha256 -extfile testdata/certs/client.cnf -extensions v3_req
	$(Q)openssl req -new -key testdata/certs/server.key -out testdata/certs/server.csr -config testdata/certs/server.cnf
	$(Q)openssl x509 -req -in testdata/certs/server.csr -CA testdata/certs/rootCA.crt -CAkey testdata/certs/rootCA.key -CAcreateserial -out testdata/certs/server.crt -days 3650 -sha256 -extfile testdata/certs/server.cnf -extensions v3_req

.PHONY: client-request-local
client-request-local:
	$(Q)docker compose run client -cert /testdata/certs/client.crt -key /testdata/certs/client.key -server-cacert /testdata/certs/rootCA.crt -tlsflags 80 --connect https://localhost:8775/debug

.PHONY: client-request
client-request:
	$(Q)docker compose run client -cert /testdata/certs/client.crt -key /testdata/certs/client.key -tlsflags 80 --connect https://req-mtls.research.cloudflare.com/debug

.PHONY: start-test-server
start-test-server:
	$(Q)docker compose run server -cert /testdata/certs/server.crt -key /testdata/certs/server.key -client-cacert /testdata/certs/rootCA.crt

.PHONY: artefacts
artefacts:
	$(Q)docker buildx build --output type=local,dest=./bin --target=artefacts .
