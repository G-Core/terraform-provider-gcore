#!/bin/bash
set -e

# Load credentials
set -o allexport
source /Users/user/repos/gcore-terraform/.env
set +o allexport

# Set Terraform config
export TF_CLI_CONFIG_FILE=/Users/user/repos/gcore-terraform/.terraformrc
export TF_LOG=DEBUG
export TF_LOG_PATH="terraform_flavor_change.log"

echo "=== Testing Flavor Change (g1-standard-1-2 -> g1-standard-2-4) ==="

# Get instance ID before change
BEFORE_ID=$(terraform output -raw instance_id)
echo "Instance ID before: $BEFORE_ID"

# Apply flavor change
terraform apply -auto-approve -var="flavor=g1-standard-2-4"

# Get instance ID after change
AFTER_ID=$(terraform output -raw instance_id)
echo "Instance ID after: $AFTER_ID"

if [ "$BEFORE_ID" = "$AFTER_ID" ]; then
    echo "✅ PASS: Flavor changed in-place (same instance ID)"
else
    echo "❌ FAIL: Instance was recreated (different ID)"
    exit 1
fi

# Drift check
terraform plan -detailed-exitcode
if [ $? -eq 0 ]; then
    echo "✅ No drift after flavor change"
else
    echo "❌ Drift detected after flavor change"
fi
