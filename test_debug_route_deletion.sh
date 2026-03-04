#!/bin/bash
set -e

echo "Quick Route Deletion Debug Test"
echo "================================"

# Load credentials
if [ -f ".env" ]; then
    source .env
else
    echo "Error: .env not found"
    exit 1
fi

export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/.terraformrc"

# Use existing router from previous test
ROUTER_ID="79ed794f-7767-4733-b2f7-cd5dc2f88f1e"

# Create simple test config
mkdir -p test-debug-route-del
cd test-debug-route-del

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

resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "test-route-fix-router"
  external_gateway_info = { enable_snat = true }
  interfaces = [{ subnet_id = gcore_cloud_network_subnet.subnet.id, type = "subnet" }]
  routes = []  # Empty array to test route deletion
}
TFEOF

echo "Importing existing router state..."
terraform import gcore_cloud_network.network 6bcead4b-0c17-44a6-8fa0-417ae6274655 2>/dev/null || true
terraform import gcore_cloud_network_subnet.subnet 2778bf67-af24-4622-a61b-b1b18acc0306 2>/dev/null || true
terraform import gcore_cloud_network_router.router $ROUTER_ID 2>/dev/null || true

echo ""
echo "Applying configuration with routes=[]..."
terraform apply -auto-approve 2>&1 | grep -E "(Warning|Error|Modifications|Apply complete)" | head -20

cd ..
