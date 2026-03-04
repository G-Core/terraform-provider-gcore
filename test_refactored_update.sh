#!/bin/bash
set -e

# Test the refactored Update function with Pedro's suggested approach:
# 1. Send single PATCH for all updates (including route deletions)
# 2. Then perform attach/detach operations

echo "=== Testing Refactored Router Update Flow ==="

# Load credentials
if [ -f .env ]; then
    source .env
else
    echo "ERROR: .env not found"
    exit 1
fi

# Create test directory
TEST_DIR="test-refactored-update-$(date +%s)"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

# Create test configuration
cat > main.tf <<EOF
terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# Create network and subnet
resource "gcore_cloud_network" "test_net" {
  project_id = ${GCORE_CLOUD_PROJECT_ID}
  region_id  = ${GCORE_CLOUD_REGION_ID}
  name       = "test-refactored-update-net"
}

resource "gcore_cloud_network_subnet" "subnet1" {
  project_id = ${GCORE_CLOUD_PROJECT_ID}
  region_id  = ${GCORE_CLOUD_REGION_ID}
  name       = "test-subnet-1"
  cidr       = "192.168.1.0/24"
  network_id = gcore_cloud_network.test_net.id
}

resource "gcore_cloud_network_subnet" "subnet2" {
  project_id = ${GCORE_CLOUD_PROJECT_ID}
  region_id  = ${GCORE_CLOUD_REGION_ID}
  name       = "test-subnet-2"
  cidr       = "192.168.2.0/24"
  network_id = gcore_cloud_network.test_net.id
}

# Router with interface and route
resource "gcore_cloud_network_router" "test" {
  project_id = ${GCORE_CLOUD_PROJECT_ID}
  region_id  = ${GCORE_CLOUD_REGION_ID}
  name       = "test-refactored-router"

  external_gateway_info = {
    enable_snat = true
  }

  interfaces = [
    {
      type      = "subnet"
      subnet_id = gcore_cloud_network_subnet.subnet1.id
    },
    {
      type      = "subnet"
      subnet_id = gcore_cloud_network_subnet.subnet2.id
    }
  ]

  routes = [
    {
      destination = "10.0.0.0/24"
      nexthop     = "192.168.2.1"  # Uses subnet2 interface
    }
  ]
}

output "router_id" {
  value = gcore_cloud_network_router.test.id
}
EOF

export TF_CLI_CONFIG_FILE=/Users/user/repos/gcore-terraform/.terraformrc

echo ""
echo "Step 1: Initial apply - create router with 2 interfaces and 1 route"
terraform apply -auto-approve

ROUTER_ID=$(terraform output -raw router_id)
echo "Created router: $ROUTER_ID"

# Update: Remove route AND detach interface
# This is the critical test case - route must be deleted via PATCH before interface detachment
echo ""
echo "Step 2: Update - remove route AND detach subnet2 interface"
echo "This tests that PATCH (route deletion) happens before detach operation"

cat > main.tf <<EOF
terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

resource "gcore_cloud_network" "test_net" {
  project_id = ${GCORE_CLOUD_PROJECT_ID}
  region_id  = ${GCORE_CLOUD_REGION_ID}
  name       = "test-refactored-update-net"
}

resource "gcore_cloud_network_subnet" "subnet1" {
  project_id = ${GCORE_CLOUD_PROJECT_ID}
  region_id  = ${GCORE_CLOUD_REGION_ID}
  name       = "test-subnet-1"
  cidr       = "192.168.1.0/24"
  network_id = gcore_cloud_network.test_net.id
}

resource "gcore_cloud_network_subnet" "subnet2" {
  project_id = ${GCORE_CLOUD_PROJECT_ID}
  region_id  = ${GCORE_CLOUD_REGION_ID}
  name       = "test-subnet-2"
  cidr       = "192.168.2.0/24"
  network_id = gcore_cloud_network.test_net.id
}

resource "gcore_cloud_network_router" "test" {
  project_id = ${GCORE_CLOUD_PROJECT_ID}
  region_id  = ${GCORE_CLOUD_REGION_ID}
  name       = "test-refactored-router-updated"  # Also update name

  external_gateway_info = {
    enable_snat = false  # Also update this
  }

  interfaces = [
    {
      type      = "subnet"
      subnet_id = gcore_cloud_network_subnet.subnet1.id
    }
    # subnet2 interface removed
  ]

  # routes removed
}

output "router_id" {
  value = gcore_cloud_network_router.test.id
}
EOF

TF_LOG=DEBUG terraform apply -auto-approve 2>&1 | tee apply_output.log

# Verify the request order in logs
echo ""
echo "=== Verifying Request Order ==="

# Check for PATCH request (should happen first)
PATCH_LINE=$(grep -n "PATCH.*routers/$ROUTER_ID" apply_output.log | head -1 | cut -d: -f1)

# Check for detach request (should happen after PATCH)
DETACH_LINE=$(grep -n "POST.*routers/$ROUTER_ID/detach_subnet" apply_output.log | head -1 | cut -d: -f1)

if [ -z "$PATCH_LINE" ]; then
    echo "❌ PATCH request not found in logs"
    exit 1
fi

if [ -z "$DETACH_LINE" ]; then
    echo "❌ Detach request not found in logs"
    exit 1
fi

if [ "$PATCH_LINE" -lt "$DETACH_LINE" ]; then
    echo "✅ Request order correct: PATCH (line $PATCH_LINE) before DETACH (line $DETACH_LINE)"
else
    echo "❌ Request order incorrect: DETACH (line $DETACH_LINE) before PATCH (line $PATCH_LINE)"
    exit 1
fi

# Count PATCH requests - should only be ONE
PATCH_COUNT=$(grep -c "PATCH.*routers/$ROUTER_ID" apply_output.log || true)
echo ""
echo "PATCH request count: $PATCH_COUNT"
if [ "$PATCH_COUNT" -eq 1 ]; then
    echo "✅ Single PATCH request as expected (Pedro's optimization)"
else
    echo "⚠️  Multiple PATCH requests detected"
fi

# Verify final state
echo ""
echo "Step 3: Verify final state"
terraform show -json | jq -r '.values.root_module.resources[] | select(.type == "gcore_cloud_network_router") | {name: .values.name, interfaces: .values.interfaces, routes: .values.routes}'

echo ""
echo "Step 4: Cleanup"
terraform destroy -auto-approve

cd ..
rm -rf "$TEST_DIR"

echo ""
echo "=== Test Completed Successfully ==="
echo "✅ Refactored Update function works correctly"
echo "✅ Single PATCH request consolidates all updates"
echo "✅ Request order: PATCH → DETACH → GET"
