FROM alpine
RUN apk add certbot certbot-dns-cloudflare aws-cli
COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
