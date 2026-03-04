#!/bin/bash
set -e

# Load environment variables
if [ -f ../.env ]; then
    set -o allexport
    source ../.env
    set +o allexport
    echo "Loaded environment from ../.env"
fi

# Set Terraform config to use local provider
export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/.terraformrc"

echo "=========================================="
echo "Test: Security Group Name Drift & Rules Nulling"
echo "=========================================="
echo ""

# Clean up any existing state
rm -f terraform.tfstate terraform.tfstate.backup .terraform.lock.hcl
rm -rf .terraform

echo "Step 1: Initialize Terraform"
echo "----------------------------"
# Note: With dev_overrides, init may fail but we can skip it
terraform init -input=false 2>&1 || echo "(init failed but continuing with dev_overrides)"
echo ""

echo "Step 2: Apply #1 - Initial Creation"
echo "------------------------------------"
terraform apply -auto-approve -input=false
echo ""

echo "Checking outputs after Apply #1:"
terraform output
echo ""

echo "Step 3: Plan #1 - Check for drift"
echo "----------------------------------"
terraform plan -detailed-exitcode -input=false && PLAN1_EXIT=$? || PLAN1_EXIT=$?
echo "Plan exit code: $PLAN1_EXIT (0=no changes, 2=changes detected)"
echo ""

echo "Step 4: Apply #2 - Should be no-op if no drift"
echo "-----------------------------------------------"
terraform apply -auto-approve -input=false
echo ""

echo "Checking outputs after Apply #2:"
terraform output
echo ""

echo "Step 5: Plan #2 - Verify no drift"
echo "----------------------------------"
terraform plan -detailed-exitcode -input=false && PLAN2_EXIT=$? || PLAN2_EXIT=$?
echo "Plan exit code: $PLAN2_EXIT (0=no changes, 2=changes detected)"
echo ""

echo "Step 6: Apply #3 - Final stability check"
echo "-----------------------------------------"
terraform apply -auto-approve -input=false
echo ""

echo "Checking outputs after Apply #3:"
terraform output
echo ""

echo "Step 7: Final Plan - Should be clean"
echo "-------------------------------------"
terraform plan -detailed-exitcode -input=false && PLAN3_EXIT=$? || PLAN3_EXIT=$?
echo "Plan exit code: $PLAN3_EXIT (0=no changes, 2=changes detected)"
echo ""

echo "=========================================="
echo "Test Results Summary"
echo "=========================================="
echo "Plan #1 exit code: $PLAN1_EXIT"
echo "Plan #2 exit code: $PLAN2_EXIT"
echo "Plan #3 exit code: $PLAN3_EXIT"
echo ""

if [ "$PLAN1_EXIT" -eq 0 ] && [ "$PLAN2_EXIT" -eq 0 ] && [ "$PLAN3_EXIT" -eq 0 ]; then
    echo "SUCCESS: No drift detected across 3 applies"
else
    echo "ISSUE: Drift detected - check output above"
fi

echo ""
echo "State inspection:"
echo "-----------------"
terraform state show gcore_cloud_security_group.test | grep -E "^[[:space:]]*(name|security_group_rules)" | head -20
