#!/bin/bash
# Run a single test case

set -e

TEST_DIR=$1

if [ -z "$TEST_DIR" ]; then
    echo "Usage: $0 <test-directory>"
    echo "Example: $0 drift/TC-DRIFT-001-lb-no-changes"
    exit 1
fi

echo "==========================================="
echo "Running test: $TEST_DIR"
echo "==========================================="
echo ""

cd "$TEST_DIR"

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

# Skip init when using provider development overrides
if [ -f ".terraformrc" ]; then
    echo "Step 1: Skipping terraform init (using provider development override)"
else
    echo "Step 1: Terraform init"
    terraform init -upgrade
fi

# First apply
echo ""
echo "Step 2: First terraform apply"
terraform apply -auto-approve

# Check for drift
echo ""
echo "Step 3: Check for configuration drift"
echo ""

if terraform plan -detailed-exitcode; then
    echo ""
    echo "✅ PASS: No drift detected"
    exit 0
else
    EXIT_CODE=$?
    if [ $EXIT_CODE -eq 2 ]; then
        echo ""
        echo "❌ FAIL: Drift detected"
        terraform plan
        exit 1
    else
        echo ""
        echo "❌ ERROR: terraform plan failed"
        exit 1
    fi
fi
