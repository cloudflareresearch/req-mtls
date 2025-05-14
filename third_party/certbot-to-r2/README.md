# certbot-to-r2

This service generate a certificate for a given FQDN, and uploads it to an R2 bucket.

## Requirement

* Docker

## Run

You need to provision the following environment variables. A sample .env file is provided in [.env.sample](./.env.sample).

| Variable                | Description                                                                   |
|:------------------------|:------------------------------------------------------------------------------|
| `CLOUDFLARE_API_TOKEN`  | Cloudflare API TOKEN allowing to Edit DNS                                     |
| `EMAIL`                 | Email cerbot bill use for cert issue. Has to be registered with let's encrypt |
| `AWS_ACCESS_KEY_ID`     | R2 S3 API access key ID                                                       |
| `AWS_SECRET_ACCESS_KEY` | R2 S3 API access secret                                                       |
| `AWS_DEFAULT_REGION`    | Jurisdiction, probably "auto"                                                 |
| `AWS_ENDPOINT_URL`      | Cloudflare R2 S3 URL                                                          |
| `AWS_S3_URI`            | s3:// URL which is the prefix you want to use to store your certificate       |

Put these in a .env file.

You can then provision the keys using the following command, replacing example.com with the domain you want to provision a certificate for.

```shell
docker build -t certbot . && \
docker run --env-file .env certbot "example.com"
```
