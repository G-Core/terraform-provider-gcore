#!/bin/bash
# Test pool name update (should use PATCH, not replacement)

set -e

echo "==========================================="
echo "TC-UPDATE-001: Pool Name Update Test"
echo "==========================================="
echo ""

# Source credentials
if [ -f "../../.env" ]; then
    source ../../.env
    export GCORE_API_KEY
    export GCORE_CLIENT
    export GCORE_CLOUD_PROJECT_ID
    export GCORE_CLOUD_REGION_ID
fi

# Copy .terraformrc if it exists
if [ -f "../../.terraformrc" ]; then
    cp ../../.terraformrc .
fi

echo "Step 1: Initial apply with pool name 'test-pool-update-01'"
terraform apply -auto-approve

echo ""
echo "Step 2: Update pool name to 'test-pool-update-01-renamed'"
terraform apply -auto-approve -var="pool_name=test-pool-update-01-renamed"

echo ""
echo "Step 3: Check terraform plan output"
PLAN_OUTPUT=$(terraform plan -var="pool_name=test-pool-update-01-renamed" -detailed-exitcode 2>&1 || true)

if echo "$PLAN_OUTPUT" | grep -q "No changes"; then
    echo "✅ PASS: Pool name updated successfully, no further changes"
    exit 0
else
    echo "❌ FAIL: Unexpected changes detected after update"
    terraform plan -var="pool_name=test-pool-update-01-renamed"
    exit 1
fi
