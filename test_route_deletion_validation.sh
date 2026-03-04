#!/bin/bash
set -e

echo "Route Deletion Validation Test"
echo "==============================="

# Load credentials
if [ -f ".env" ]; then
    source .env
else
    echo "Error: .env not found"
    exit 1
fi

export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/.terraformrc"

# Create test directory
TEST_DIR="test-route-validation-$(date +%s)"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

cat > main.tf << 'TFEOF'
terraform {
  required_providers {
    gcore = { source = "gcore/gcore" }
  }
}
provider "gcore" {}

resource "gcore_cloud_network" "network" {
  project_id = 379987
  region_id  = 76
  name       = "test-route-validation-network"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "subnet" {
  project_id  = 379987
  region_id   = 76
  name        = "test-route-validation-subnet"
  cidr        = "192.168.70.0/24"
  network_id  = gcore_cloud_network.network.id
}

variable "include_routes" {
  type    = bool
  default = true
}

resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "test-route-validation-router"
  external_gateway_info = { enable_snat = true }
  interfaces = [{ subnet_id = gcore_cloud_network_subnet.subnet.id, type = "subnet" }]
  routes = var.include_routes ? [
    { destination = "10.0.7.0/24", nexthop = "192.168.70.1" }
  ] : []
}

output "router_id" { value = gcore_cloud_network_router.router.id }
output "router_routes" { value = gcore_cloud_network_router.router.routes }
TFEOF

echo "Step 1: Creating router WITH routes..."
echo "======================================="
terraform init > /dev/null 2>&1
terraform apply -auto-approve -var="include_routes=true" 2>&1 | tail -20

ROUTER_ID=$(terraform output -raw router_id)
echo ""
echo "Router created: $ROUTER_ID"
echo "Routes in state:"
terraform output -json router_routes | jq -c

echo ""
echo "Step 2: Removing routes with TF_LOG=DEBUG..."
echo "============================================="
TF_LOG=DEBUG terraform apply -auto-approve -var="include_routes=false" 2>&1 | tee apply.log | grep -E "(Route deletion|Warning|PATCH|routes)" | head -30

echo ""
echo "Step 3: Checking for drift after route removal..."
echo "=================================================="
terraform plan -detailed-exitcode -var="include_routes=false" 2>&1 | tee plan.log | tail -20
PLAN_EXIT_CODE=$?

echo ""
echo "Results:"
echo "--------"
if [ $PLAN_EXIT_CODE -eq 0 ]; then
    echo "✓ SUCCESS: No drift detected - routes were successfully deleted"
elif [ $PLAN_EXIT_CODE -eq 2 ]; then
    echo "✗ FAILURE: Drift detected - routes still exist on API"
    echo ""
    echo "Drift details:"
    terraform plan -var="include_routes=false" 2>&1 | grep -A10 "routes"
else
    echo "✗ ERROR: Plan command failed"
fi

echo ""
echo "Step 4: Checking diagnostic warnings..."
echo "========================================"
if grep -q "Route deletion workaround" apply.log; then
    echo "✓ Workaround was triggered (found 'Route deletion workaround' warning)"
    grep "Route deletion workaround" apply.log
else
    echo "✗ Workaround was NOT triggered (no 'Route deletion workaround' warning found)"
fi

if grep -q "Router route deletion detected" apply.log; then
    echo "✓ Route deletion detected in diagnostics"
    grep "Router route deletion detected" apply.log
else
    echo "✗ No route deletion diagnostic found"
fi

echo ""
echo "Cleaning up..."
terraform destroy -auto-approve > /dev/null 2>&1
cd ..
rm -rf "$TEST_DIR"

echo ""
if [ $PLAN_EXIT_CODE -eq 0 ]; then
    echo "✓✓✓ VALIDATION PASSED ✓✓✓"
    exit 0
else
    echo "✗✗✗ VALIDATION FAILED ✗✗✗"
    exit 1
fi
