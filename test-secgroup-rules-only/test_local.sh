#!/bin/bash
set -e

echo "=== Building local provider ==="
cd /Users/user/repos/gcore-terraform
go build -o terraform-provider-gcore

echo "=== Running Terraform with local provider ==="
cd test-secgroup-rules-only

# Set Terraform to use local provider
export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/test-secgroup-rules-only/.terraformrc"

# Source environment variables
if [ -f ../.env ]; then
    source ../.env
fi

echo "=== Terraform Init ==="
terraform init -upgrade

echo "=== Terraform Plan ==="
terraform plan

echo "=== Terraform Apply ==="
terraform apply -auto-approve

echo "=== Check state - should have NO security_group_rules field ==="
terraform show

echo "=== Apply again to check for drift ==="
terraform apply -auto-approve

echo "=== Final state check ==="
terraform show | grep -A 5 "security_group_rules" || echo "✅ NO security_group_rules field in state!"
