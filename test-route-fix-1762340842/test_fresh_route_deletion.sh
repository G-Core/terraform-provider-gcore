#!/bin/bash
set -e

echo "Fresh Route Deletion Test with Diagnostics"
echo "==========================================="

# Load credentials
source .env 2>/dev/null || { echo "No .env"; exit 1; }
export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/.terraformrc"

# Create fresh test directory
rm -rf test-fresh-route-del
mkdir test-fresh-route-del
cd test-fresh-route-del

# Create config
cat > main.tf << 'TFEOF'
terraform {
  required_providers { gcore = { source = "gcore/gcore" } }
}
provider "gcore" {}

resource "gcore_cloud_network" "net" {
  project_id = 379987
  region_id  = 76
  name       = "test-fresh-route-net"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "sub" {
  project_id  = 379987
  region_id   = 76
  name        = "test-fresh-route-sub"
  cidr        = "192.168.70.0/24"
  network_id  = gcore_cloud_network.net.id
}

variable "routes" {
  default = [{destination = "10.0.7.0/24", nexthop = "192.168.70.1"}]
}

resource "gcore_cloud_network_router" "rtr" {
  project_id = 379987
  region_id  = 76
  name       = "test-fresh-route-rtr"
  external_gateway_info = { enable_snat = true }
  interfaces = [{ subnet_id = gcore_cloud_network_subnet.sub.id, type = "subnet" }]
  routes = var.routes
}

output "router_id" { value = gcore_cloud_network_router.rtr.id }
output "routes" { value = gcore_cloud_network_router.rtr.routes }
TFEOF

echo "Step 1: Creating router with 1 route..."
terraform apply -auto-approve > /dev/null 2>&1

ROUTER_ID=$(terraform output -raw router_id)
echo "Router created: $ROUTER_ID"
echo "Routes: $(terraform output -json routes | jq -c '.')"

echo ""
echo "Step 2: Removing routes (setting to empty array)..."
terraform apply -auto-approve -var='routes=[]' 2>&1 | tee /tmp/fresh_route_test.log

echo ""
echo "Step 3: Checking for warnings in output..."
grep -i "warning.*route" /tmp/fresh_route_test.log || echo "No route warnings found"

cd ..
