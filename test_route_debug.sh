#!/bin/bash
set -e

# Load environment
if [ -f ".env" ]; then
    set -o allexport
    source .env
    set +o allexport
    echo "✓ Loaded credentials"
else
    echo "✗ Error: .env file not found"
    exit 1
fi

export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/.terraformrc"
export TF_LOG=DEBUG

cd test-router-comprehensive/update/TC-UPDATE-001

echo "Step 1: Creating router WITH routes..."
terraform apply -auto-approve -var="include_routes=true" 2>&1 | tee /tmp/terraform_create.log

echo ""
echo "Step 2: Removing routes..."
terraform apply -auto-approve -var="include_routes=false" 2>&1 | tee /tmp/terraform_update.log

echo ""
echo "Checking logs for debug info..."
grep -i "modifyplan routes check" /tmp/terraform_update.log || echo "No ModifyPlan logs found"
grep -i "update needsupdate check" /tmp/terraform_update.log || echo "No Update logs found"  
grep -i "patch body" /tmp/terraform_update.log || echo "No PATCH body logs found"
