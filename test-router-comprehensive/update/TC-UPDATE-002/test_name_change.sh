#!/bin/bash
set -e

echo "=========================================="
echo "TC-UPDATE-002: Name Change Test"
echo "=========================================="
echo "This tests that name updates use PATCH, not replacement"
echo ""

cd "$(dirname "$0")"

# Load environment
if [ -f "../../../.env" ]; then
    set -o allexport
    source ../../../.env
    set +o allexport
    echo "✓ Loaded credentials"
else
    echo "✗ Error: .env file not found"
    exit 1
fi

export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/.terraformrc"

# Step 1: Create router with original name
echo ""
echo "Step 1: Creating router with original name..."
echo "----------------------------------------------"
terraform apply -auto-approve -var="router_name=test-router-original"

ROUTER_ID=$(terraform output -raw router_id)
ROUTER_NAME=$(terraform output -raw router_name)
echo "Router ID: $ROUTER_ID"
echo "Router Name: $ROUTER_NAME"

if [ "$ROUTER_NAME" != "test-router-original" ]; then
    echo "❌ Expected name 'test-router-original', got '$ROUTER_NAME'"
    exit 1
fi
echo "✅ Router created with original name"

# Step 2: Update name
echo ""
echo "Step 2: Updating router name..."
echo "--------------------------------"
terraform apply -auto-approve -var="router_name=test-router-updated"

NEW_ROUTER_ID=$(terraform output -raw router_id)
NEW_ROUTER_NAME=$(terraform output -raw router_name)

echo "Router ID: $NEW_ROUTER_ID"
echo "Router Name: $NEW_ROUTER_NAME"

# Verify router not replaced
if [ "$ROUTER_ID" = "$NEW_ROUTER_ID" ]; then
    echo "✅ Router ID unchanged (update in-place via PATCH)"
else
    echo "❌ Router ID changed (resource was replaced)"
    echo "   Old ID: $ROUTER_ID"
    echo "   New ID: $NEW_ROUTER_ID"
    exit 1
fi

# Verify name updated
if [ "$NEW_ROUTER_NAME" = "test-router-updated" ]; then
    echo "✅ Router name successfully updated"
else
    echo "❌ Expected name 'test-router-updated', got '$NEW_ROUTER_NAME'"
    exit 1
fi

# Step 3: Check for drift
echo ""
echo "Step 3: Checking for drift after name change..."
echo "------------------------------------------------"
if terraform plan -detailed-exitcode -var="router_name=test-router-updated"; then
    echo "✅ No drift detected"
else
    EXIT_CODE=$?
    if [ $EXIT_CODE -eq 2 ]; then
        echo "❌ Drift detected after name change"
        terraform plan -var="router_name=test-router-updated"
        exit 1
    fi
fi

echo ""
echo "=========================================="
echo "✅ TC-UPDATE-002 PASSED"
echo "=========================================="
echo "- Router name successfully updated"
echo "- Update used PATCH (router ID unchanged)"
echo "- No drift after update"
