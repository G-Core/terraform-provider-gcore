#!/bin/bash
set -e

echo "=========================================="
echo "Router Drift Detection Tests"
echo "=========================================="
echo ""

# Load environment
if [ -f "../../.env" ]; then
    set -o allexport
    source ../../.env
    set +o allexport
    echo "✓ Loaded credentials from .env"
else
    echo "✗ Error: .env file not found"
    exit 1
fi

export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/.terraformrc"

TESTS=("TC-DRIFT-001-baseline" "TC-DRIFT-002-with-routes")
PASSED=0
FAILED=0

for TEST in "${TESTS[@]}"; do
    echo ""
    echo "=========================================="
    echo "Running: $TEST"
    echo "=========================================="

    # Create test directory if needed
    mkdir -p "$TEST"
    cp "${TEST}.tf" "${TEST}/main.tf"
    cd "$TEST"

    # Apply
    echo ""
    echo "Step 1: Creating resources..."
    if ! terraform apply -auto-approve; then
        echo "❌ FAIL: terraform apply failed"
        FAILED=$((FAILED + 1))
        cd ..
        continue
    fi

    echo ""
    echo "Step 2: Checking for drift..."
    if terraform plan -detailed-exitcode; then
        echo ""
        echo "✅ PASS: $TEST - No drift detected"
        PASSED=$((PASSED + 1))
    else
        EXIT_CODE=$?
        if [ $EXIT_CODE -eq 2 ]; then
            echo ""
            echo "❌ FAIL: $TEST - Drift detected"
            terraform plan
            FAILED=$((FAILED + 1))
        else
            echo ""
            echo "❌ ERROR: $TEST - terraform plan failed"
            FAILED=$((FAILED + 1))
        fi
    fi

    # Cleanup
    echo ""
    echo "Cleaning up..."
    terraform destroy -auto-approve > /dev/null 2>&1 || true

    cd ..
done

echo ""
echo "=========================================="
echo "Drift Test Results"
echo "=========================================="
echo "PASSED: $PASSED"
echo "FAILED: $FAILED"
echo ""

if [ $FAILED -gt 0 ]; then
    exit 1
fi
