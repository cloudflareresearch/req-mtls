#! /bin/sh
set -e

# Wait for network to be up before fetching
sleep 60
while [ -z "$(ip addr show cfeth0 | grep 'inet ' | awk '{print $2}' | cut -d/ -f1)" ]; do
  echo "The network is not up yet"
  sleep 3
done

# Check if the required environment variables are set
if [[ -z "$AWS_ACCESS_KEY_ID" || -z "$AWS_SECRET_ACCESS_KEY" || -z "$AWS_DEFAULT_REGION" || -z "$AWS_ENDPOINT_URL" ]]; then
  echo "Error: One or more required environment variables are missing."
  exit 1
fi

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
aws s3 cp "$AWS_S3_URI" /home/server --recursive

# Run command
exec "$@"
