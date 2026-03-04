#!/bin/bash
set -e

echo "Testing PATCH body for route deletion"

# Load credentials
if [ -f ".env" ]; then
    set -o allexport
    source .env
    set +o allexport
else
    echo "Error: .env file not found"
    exit 1
fi

export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/.terraformrc"

# Create test dir
TEST_DIR="test-patch-debug-$(date +%s)"
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
  name       = "test-patch-network"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "subnet" {
  project_id  = 379987
  region_id   = 76
  name        = "test-patch-subnet"
  cidr        = "192.168.60.0/24"
  network_id  = gcore_cloud_network.network.id
}

variable "include_routes" {
  type    = bool
  default = true
}

resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "test-patch-router"
  external_gateway_info = { enable_snat = true }
  interfaces = [{ subnet_id = gcore_cloud_network_subnet.subnet.id, type = "subnet" }]
  routes = var.include_routes ? [{ destination = "10.0.6.0/24", nexthop = "192.168.60.1" }] : []
}

output "router_id" { value = gcore_cloud_network_router.router.id }
output "router_routes" { value = gcore_cloud_network_router.router.routes }
TFEOF

echo "Creating router with routes..."
terraform apply -auto-approve -var="include_routes=true" 2>&1 | tail -20

ROUTER_ID=$(terraform output -raw router_id)
echo "Router created: $ROUTER_ID"

echo ""
echo "Removing routes and capturing PATCH body..."
echo "=============================================="
terraform apply -auto-approve -var="include_routes=false" 2>&1 | grep -A20 "PATCH REQUEST DETAILS" || echo "No PATCH logs found (check TF_LOG)"

echo ""
echo "Cleaning up..."
terraform destroy -auto-approve > /dev/null 2>&1
cd ..
rm -rf "$TEST_DIR"
echo "Done"
