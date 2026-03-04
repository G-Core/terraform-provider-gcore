#!/bin/bash
# Test script for gcore_gpu_virtual_image with old provider
# JIRA: GCLOUD2-21134
#
# Required environment variables:
#   GCORE_API_KEY - Gcore API token
#   GCORE_CLOUD_PROJECT_ID - Project ID (default: 379987)
#   GCORE_CLOUD_REGION_ID - Region ID (default: 76)

set -e

# Use local terraformrc without dev_overrides to get the real old provider
export TF_CLI_CONFIG_FILE="$(pwd)/../.terraformrc.old"
echo "Using TF_CLI_CONFIG_FILE=$TF_CLI_CONFIG_FILE (no dev overrides)"

# Load credentials from parent .env if exists
if [ -f "../.env" ]; then
    echo "Loading credentials from ../.env"
    set -o allexport
    source "../.env"
    set +o allexport
fi

# Validate required vars
if [ -z "$GCORE_API_KEY" ]; then
    echo "ERROR: GCORE_API_KEY environment variable is required"
    exit 1
fi

PROJECT_ID="${GCORE_CLOUD_PROJECT_ID:-379987}"
REGION_ID="${GCORE_CLOUD_REGION_ID:-76}"

echo "=== Test GPU Virtual Image (Old Provider) ==="
echo "Project ID: $PROJECT_ID"
echo "Region ID: $REGION_ID"
echo ""

# Initialize terraform with old provider
echo "=== Initializing Terraform ==="
terraform init -upgrade

# Create terraform.tfvars
cat > terraform.tfvars <<EOF
gcore_api_token = "$GCORE_API_KEY"
project_id      = $PROJECT_ID
region_id       = $REGION_ID
EOF

# Plan
echo ""
echo "=== Terraform Plan ==="
terraform plan

# Ask for confirmation
echo ""
read -p "Proceed with apply? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborted."
    exit 0
fi

# Apply
echo ""
echo "=== Terraform Apply ==="
terraform apply -auto-approve

# Show state
echo ""
echo "=== Terraform State ==="
terraform show

# Copy files for JIRA
echo ""
echo "=== Saving files for JIRA ticket GCLOUD2-21134 ==="
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
mkdir -p "jira_artifacts_$TIMESTAMP"
cp main.tf "jira_artifacts_$TIMESTAMP/"
cp terraform.tfstate "jira_artifacts_$TIMESTAMP/" 2>/dev/null || echo "No state file yet"
cp terraform.tfvars.example "jira_artifacts_$TIMESTAMP/"
echo "Files saved to jira_artifacts_$TIMESTAMP/"

echo ""
echo "=== Done ==="
echo "Image ID: $(terraform output -raw image_id 2>/dev/null || echo 'N/A')"
