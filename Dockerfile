FROM golang:1.24 AS go-builder
ADD https://github.com/cloudflare/go.git#37bc41c6ff79507200a315b72834fce6ca427a7e /go
WORKDIR /go/src
RUN ./make.bash

FROM alpine AS builder
RUN apk add git
COPY --from=go-builder /go/ /go/
COPY cmd ./cmd
COPY go.mod go.mod
COPY go.sum go.sum
ENV GOROOT=/go
RUN /go/bin/go build -o /client ./cmd/client
RUN /go/bin/go build -o /server ./cmd/server

FROM scratch AS artefacts
COPY --from=builder /client /client
COPY --from=builder /server /server

FROM alpine AS client
RUN  addgroup -S clientg && adduser -S client -G clientg -h /home/client
COPY --from=builder /client /usr/local/bin/client
USER client
WORKDIR /home/client
ENTRYPOINT ["/usr/local/bin/client"]

FROM alpine AS server
RUN  addgroup -S serverg && adduser -S server -G serverg -h /home/server
COPY --from=builder /server /usr/local/bin/server
RUN apk add libcap-setcap
RUN setcap CAP_NET_BIND_SERVICE=+eip /usr/local/bin/server
USER server
WORKDIR /home/server
ENTRYPOINT ["/usr/local/bin/server"]

FROM server
USER root
RUN apk add aws-cli
USER server
COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
CMD ["/usr/local/bin/server", "--cert", "/home/server/server.crt", "--key", "/home/server/server.key", "--client-cacert", "/home/server/rootCA.crt", "--listen-addr", "0.0.0.0:443"]
