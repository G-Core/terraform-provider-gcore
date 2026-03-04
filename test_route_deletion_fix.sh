#!/bin/bash
set -e

echo "=========================================="
echo "Testing Route Deletion Fix"
echo "=========================================="

# Load credentials
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

# Create a test directory
TEST_DIR="test-route-fix-$(date +%s)"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

# Create test configuration
cat > main.tf << 'TFEOF'
terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

resource "gcore_cloud_network" "network" {
  project_id = 379987
  region_id  = 76
  name       = "test-route-fix-network"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "subnet" {
  project_id  = 379987
  region_id   = 76
  name        = "test-route-fix-subnet"
  cidr        = "192.168.50.0/24"
  network_id  = gcore_cloud_network.network.id
}

variable "include_routes" {
  type    = bool
  default = true
}

resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "test-route-fix-router"

  external_gateway_info = {
    enable_snat = true
  }

  interfaces = [{
    subnet_id = gcore_cloud_network_subnet.subnet.id
    type      = "subnet"
  }]

  routes = var.include_routes ? [
    {
      destination = "10.0.5.0/24"
      nexthop     = "192.168.50.1"
    }
  ] : []
}

output "router_id" {
  value = gcore_cloud_network_router.router.id
}

output "router_routes" {
  value = gcore_cloud_network_router.router.routes
}
TFEOF

# Skip terraform init with provider development overrides
echo "✓ Using provider development override, skipping init"

echo ""
echo "Step 1: Create router WITH routes..."
echo "----------------------------------------"
terraform apply -auto-approve -var="include_routes=true"

ROUTER_ID=$(terraform output -raw router_id)
ROUTES_COUNT=$(terraform output -json router_routes | jq 'length')

echo "Router ID: $ROUTER_ID"
echo "Routes count: $ROUTES_COUNT"

if [ "$ROUTES_COUNT" -eq "1" ]; then
    echo "✅ Route created successfully"
else
    echo "❌ Expected 1 route, got $ROUTES_COUNT"
    cd .. && rm -rf "$TEST_DIR"
    exit 1
fi

echo ""
echo "Step 2: Remove routes (set to empty array)..."
echo "----------------------------------------"
terraform apply -auto-approve -var="include_routes=false" 2>&1 | tee /tmp/tf_remove_routes.log

NEW_ROUTER_ID=$(terraform output -raw router_id)
NEW_ROUTES_COUNT=$(terraform output -json router_routes | jq 'length')

echo "Router ID: $NEW_ROUTER_ID"
echo "Routes count: $NEW_ROUTES_COUNT"

# Check if ModifyPlan triggered
if grep -q "ModifyPlan: Forcing routes to empty array for deletion" /tmp/tf_remove_routes.log; then
    echo "✅ ModifyPlan fix triggered"
else
    echo "⚠️  ModifyPlan fix did not trigger (check logs)"
fi

# Verify router not replaced
if [ "$ROUTER_ID" = "$NEW_ROUTER_ID" ]; then
    echo "✅ Router ID unchanged (update in-place)"
else
    echo "❌ Router ID changed (resource was replaced)"
    cd .. && rm -rf "$TEST_DIR"
    exit 1
fi

# Verify routes deleted
if [ "$NEW_ROUTES_COUNT" -eq "0" ]; then
    echo "✅ Routes removed from state"
else
    echo "❌ Expected 0 routes, got $NEW_ROUTES_COUNT"
    cd .. && rm -rf "$TEST_DIR"
    exit 1
fi

echo ""
echo "Step 3: Check for drift..."
echo "----------------------------------------"
if terraform plan -detailed-exitcode; then
    echo "✅ NO DRIFT DETECTED - FIX WORKS!"
else
    EXIT_CODE=$?
    if [ $EXIT_CODE -eq 2 ]; then
        echo "❌ DRIFT DETECTED - FIX FAILED"
        terraform plan
        cd .. && rm -rf "$TEST_DIR"
        exit 1
    fi
fi

echo ""
echo "Step 4: Cleanup..."
echo "----------------------------------------"
terraform destroy -auto-approve

cd ..
rm -rf "$TEST_DIR"

echo ""
echo "=========================================="
echo "✅ ROUTE DELETION FIX VERIFIED"
echo "=========================================="
echo "- Routes created successfully"
echo "- Routes deleted when set to empty array"
echo "- Router updated in-place (not replaced)"
echo "- NO DRIFT after deletion"
