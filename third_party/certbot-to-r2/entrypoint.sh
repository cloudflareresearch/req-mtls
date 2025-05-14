#! /bin/sh
set -e

# Wait for network to be up before fetching
# sleep 60
# while [ -z "$(ip addr show cfeth0 | grep 'inet ' | awk '{print $2}' | cut -d/ -f1)" ]; do
#   echo "The network is not up yet"
#   sleep 3
# done

# Check if the required environment variables are set
if [[ -z "$CLOUDFLARE_API_TOKEN" || -z "$AWS_ACCESS_KEY_ID" || -z "$AWS_SECRET_ACCESS_KEY" || -z "$AWS_DEFAULT_REGION" || -z "$AWS_ENDPOINT_URL" || -z "$AWS_S3_URI" ]]; then
  echo "Error: One or more required environment variables are missing."
  exit 1
fi

# Write Cloudflare configuration in cloudflare.ini
# Taken from https://certbot-dns-cloudflare.readthedocs.io/en/stable/index.html
cat <<EOL > ~/cloudflare.ini
dns_cloudflare_api_token = $CLOUDFLARE_API_TOKEN
EOL

chmod 600 ~/cloudflare.ini

# Issue certificate
certbot certonly -n --agree-tos --email $EMAIL --dns-cloudflare --dns-cloudflare-credentials ~/cloudflare.ini -d "$@"

# Create the .aws directory if it doesn't exist
mkdir -p ~/.aws

# Write the AWS credentials and config files
cat <<EOL > ~/.aws/credentials
[default]
aws_access_key_id = $AWS_ACCESS_KEY_ID
aws_secret_access_key = $AWS_SECRET_ACCESS_KEY
EOL

cat <<EOL > ~/.aws/config
[default]
region = $AWS_DEFAULT_REGION
output = json
endpoint_url = $AWS_ENDPOINT_URL
EOL

# Mandate that server.crt, server.key, and rootCA.crt exist at the AWS S3 URI
aws s3 cp /etc/letsencrypt/live/$@/fullchain.pem "${AWS_S3_URI}/server.crt"
aws s3 cp /etc/letsencrypt/live/$@/privkey.pem "${AWS_S3_URI}/server.key"
