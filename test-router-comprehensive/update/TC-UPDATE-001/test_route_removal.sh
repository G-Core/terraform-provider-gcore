#!/bin/bash
set -e

echo "=========================================="
echo "TC-UPDATE-001: Route Removal Test"
echo "=========================================="
echo "This tests the main bug fix: routes deleted when removed from config"
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

# Step 1: Create router WITH routes
echo ""
echo "Step 1: Creating router WITH routes..."
echo "----------------------------------------"
terraform apply -auto-approve -var="include_routes=true"

ROUTER_ID=$(terraform output -raw router_id)
echo "Router ID: $ROUTER_ID"

ROUTES_COUNT=$(terraform output -json router_routes | jq 'length')
echo "Routes count: $ROUTES_COUNT"

if [ "$ROUTES_COUNT" -eq "1" ]; then
    echo "✅ Route created successfully"
else
    echo "❌ Expected 1 route, got $ROUTES_COUNT"
    exit 1
fi

# Step 2: Remove routes from config
echo ""
echo "Step 2: Removing routes from configuration..."
echo "----------------------------------------"
terraform apply -auto-approve -var="include_routes=false"

NEW_ROUTER_ID=$(terraform output -raw router_id)
NEW_ROUTES_COUNT=$(terraform output -json router_routes | jq 'length')

echo "Router ID: $NEW_ROUTER_ID"
echo "Routes count: $NEW_ROUTES_COUNT"

# Verify router not replaced
if [ "$ROUTER_ID" = "$NEW_ROUTER_ID" ]; then
    echo "✅ Router ID unchanged (update in-place)"
else
    echo "❌ Router ID changed (resource was replaced)"
    exit 1
fi

# Verify routes deleted
if [ "$NEW_ROUTES_COUNT" -eq "0" ]; then
    echo "✅ Routes successfully deleted"
else
    echo "❌ Expected 0 routes, got $NEW_ROUTES_COUNT"
    exit 1
fi

# Step 3: Check for drift
echo ""
echo "Step 3: Checking for drift after route removal..."
echo "----------------------------------------"
if terraform plan -detailed-exitcode; then
    echo "✅ No drift detected"
else
    EXIT_CODE=$?
    if [ $EXIT_CODE -eq 2 ]; then
        echo "❌ Drift detected after route removal"
        terraform plan
        exit 1
    fi
fi

echo ""
echo "=========================================="
echo "✅ TC-UPDATE-001 PASSED"
echo "=========================================="
echo "- Routes successfully added"
echo "- Routes successfully removed when config changed"
echo "- Router updated in-place (not replaced)"
echo "- No drift after update"
