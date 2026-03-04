#!/bin/bash
set -e

# Load environment variables
if [ -f /Users/user/repos/gcore-terraform/.env ]; then
    echo "=== Loading credentials ==="
    set -o allexport
    source /Users/user/repos/gcore-terraform/.env
    set +o allexport
else
    echo "ERROR: .env not found"
    exit 1
fi

export TF_CLI_CONFIG_FILE=/Users/user/repos/gcore-terraform/.terraformrc

echo "=== Planning (create volume + instance) ==="
echo "(Skipping terraform init - using dev overrides)"
terraform plan -var "api_key=$GCORE_API_KEY"

echo ""
echo "=== Applying ==="
terraform apply -auto-approve -var "api_key=$GCORE_API_KEY"

echo ""
echo "=== Verifying state ==="
terraform show

echo ""
echo "=== Testing plan (should show no changes) ==="
terraform plan -var "api_key=$GCORE_API_KEY" -detailed-exitcode || {
    exitcode=$?
    if [ $exitcode -eq 2 ]; then
        echo "ERROR: Plan shows changes (drift detected)"
        exit 1
    elif [ $exitcode -eq 1 ]; then
        echo "ERROR: Terraform plan failed"
        exit 1
    fi
}

echo ""
echo "=== Test passed! Cleaning up ==="
terraform destroy -auto-approve -var "api_key=$GCORE_API_KEY"

echo ""
echo "=== All tests passed ==="
