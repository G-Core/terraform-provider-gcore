#!/bin/bash
set -e

echo "=== Test: Adding Route + Interface Together ==="
echo "This tests if we need to attach interface BEFORE patching route"
echo ""

if [ -f .env ]; then
    source .env
else
    echo "ERROR: .env not found"
    exit 1
fi

TEST_DIR="test-add-route-interface-$(date +%s)"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

export TF_CLI_CONFIG_FILE=/Users/user/repos/gcore-terraform/.terraformrc

# Create initial config: router with ONE interface, NO routes
cat > main.tf <<EOF
terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

resource "gcore_cloud_network" "test" {
  project_id = ${GCORE_CLOUD_PROJECT_ID}
  region_id  = ${GCORE_CLOUD_REGION_ID}
  name       = "test-add-route-iface-net"
}

resource "gcore_cloud_network_subnet" "subnet1" {
  project_id = ${GCORE_CLOUD_PROJECT_ID}
  region_id  = ${GCORE_CLOUD_REGION_ID}
  name       = "subnet1"
  cidr       = "192.168.1.0/24"
  network_id = gcore_cloud_network.test.id
}

resource "gcore_cloud_network_subnet" "subnet2" {
  project_id = ${GCORE_CLOUD_PROJECT_ID}
  region_id  = ${GCORE_CLOUD_REGION_ID}
  name       = "subnet2"
  cidr       = "192.168.2.0/24"
  network_id = gcore_cloud_network.test.id
}

resource "gcore_cloud_network_router" "test" {
  project_id = ${GCORE_CLOUD_PROJECT_ID}
  region_id  = ${GCORE_CLOUD_REGION_ID}
  name       = "test-add-route-iface"

  external_gateway_info = {
    enable_snat = true
  }

  # Start with ONLY subnet1
  interfaces = [
    {
      type      = "subnet"
      subnet_id = gcore_cloud_network_subnet.subnet1.id
    }
  ]

  # NO routes initially
}

output "router_id" {
  value = gcore_cloud_network_router.test.id
}

output "subnet2_gateway" {
  value = gcore_cloud_network_subnet.subnet2.gateway_ip
}
EOF

echo "Step 1: Create router with ONE interface, NO routes"
terraform apply -auto-approve > /dev/null
ROUTER_ID=$(terraform output -raw router_id)
SUBNET2_GW=$(terraform output -raw subnet2_gateway)
echo "Created router: $ROUTER_ID"
echo "Subnet2 gateway: $SUBNET2_GW"
echo ""

# Update: Add BOTH subnet2 interface AND a route that uses it
echo "Step 2: Update - Add subnet2 interface AND route that references it"
cat > main.tf <<EOF
terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

resource "gcore_cloud_network" "test" {
  project_id = ${GCORE_CLOUD_PROJECT_ID}
  region_id  = ${GCORE_CLOUD_REGION_ID}
  name       = "test-add-route-iface-net"
}

resource "gcore_cloud_network_subnet" "subnet1" {
  project_id = ${GCORE_CLOUD_PROJECT_ID}
  region_id  = ${GCORE_CLOUD_REGION_ID}
  name       = "subnet1"
  cidr       = "192.168.1.0/24"
  network_id = gcore_cloud_network.test.id
}

resource "gcore_cloud_network_subnet" "subnet2" {
  project_id = ${GCORE_CLOUD_PROJECT_ID}
  region_id  = ${GCORE_CLOUD_REGION_ID}
  name       = "subnet2"
  cidr       = "192.168.2.0/24"
  network_id = gcore_cloud_network.test.id
}

resource "gcore_cloud_network_router" "test" {
  project_id = ${GCORE_CLOUD_PROJECT_ID}
  region_id  = ${GCORE_CLOUD_REGION_ID}
  name       = "test-add-route-iface"

  external_gateway_info = {
    enable_snat = true
  }

  # Add subnet2 interface
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

  # Add route that uses subnet2's gateway as nexthop
  routes = [
    {
      destination = "10.0.0.0/24"
      nexthop     = "$SUBNET2_GW"
    }
  ]
}

output "router_id" {
  value = gcore_cloud_network_router.test.id
}
EOF

echo "Applying update with DEBUG logging..."
TF_LOG=DEBUG terraform apply -auto-approve 2>&1 | tee apply.log

# Check if it succeeded
if [ $? -eq 0 ]; then
  echo ""
  echo "✅ SUCCESS: Update worked!"
  echo ""

  # Analyze the order of API calls
  echo "=== Analyzing API Call Order ==="

  # Check what happened first: PATCH or attach?
  PATCH_LINE=$(grep -n "PATCH.*routers/$ROUTER_ID\"" apply.log | head -1 | cut -d: -f1)
  ATTACH_LINE=$(grep -n "POST.*routers/$ROUTER_ID/attach_subnet" apply.log | head -1 | cut -d: -f1)

  if [ -n "$PATCH_LINE" ] && [ -n "$ATTACH_LINE" ]; then
    if [ "$PATCH_LINE" -lt "$ATTACH_LINE" ]; then
      echo "Order: PATCH (line $PATCH_LINE) → ATTACH (line $ATTACH_LINE)"
      echo "⚠️  WARNING: PATCH happened before ATTACH"
      echo "   If this worked, API doesn't validate nexthops against attached interfaces"
    else
      echo "Order: ATTACH (line $ATTACH_LINE) → PATCH (line $PATCH_LINE)"
      echo "✅ ATTACH happened before PATCH (correct order)"
    fi
  else
    echo "Could not determine order from logs"
  fi
else
  echo ""
  echo "❌ FAILED: Update failed!"
  echo "This confirms we need to attach interfaces BEFORE patching routes"
fi

echo ""
echo "Step 3: Cleanup"
terraform destroy -auto-approve > /dev/null 2>&1

cd ..
rm -rf "$TEST_DIR"

echo "Test complete"
