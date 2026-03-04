#!/bin/bash
set -e

# Load credentials
set -o allexport
source /Users/user/repos/gcore-terraform/.env
set +o allexport

# Set Terraform config
export TF_CLI_CONFIG_FILE=/Users/user/repos/gcore-terraform/.terraformrc

# Set proxy
export HTTP_PROXY="http://127.0.0.1:9092"
export HTTPS_PROXY="http://127.0.0.1:9092"
export NO_PROXY="registry.terraform.io,releases.hashicorp.com"
export SSL_CERT_FILE="/Users/user/.mitmproxy/mitmproxy-ca-cert.pem"

# Set debug logging
export TF_LOG=DEBUG
export TF_LOG_PATH="terraform_test.log"

echo "=== Environment Setup ==="
echo "GCORE_API_KEY set: $(echo $GCORE_API_KEY | head -c 10)..."
echo "TF_CLI_CONFIG_FILE: $TF_CLI_CONFIG_FILE"
echo "HTTP_PROXY: $HTTP_PROXY"

# Skip terraform init when using dev_overrides
# echo "=== Terraform Init ==="
# terraform init
echo "Skipping terraform init (dev_overrides in effect)"

echo ""
echo "=== Terraform Apply ==="
terraform apply -auto-approve

echo ""
echo "=== Terraform Plan (Drift Check) ==="
terraform plan -detailed-exitcode
DRIFT_EXIT=$?

if [ $DRIFT_EXIT -eq 0 ]; then
    echo "✅ No drift detected"
elif [ $DRIFT_EXIT -eq 2 ]; then
    echo "❌ DRIFT DETECTED - Changes would be made"
    exit 1
else
    echo "❌ Plan failed"
    exit 1
fi

echo ""
echo "=== Instance Details ==="
terraform output

echo ""
echo "=== Test Complete ==="
