#!/bin/bash

set -e

echo "===================================================================="
echo "Test Case: Router Route Deletion Bug Reproduction"
echo "Issue: Routes are not deleted when removed from Terraform config"
echo "===================================================================="

cd "$(dirname "$0")"

# Load environment variables
if [ -f ../.env ]; then
    source ../.env
else
    echo "Error: .env file not found in parent directory"
    exit 1
fi

echo ""
echo "Step 1: Initialize Terraform"
echo "----------------------------"
terraform init

echo ""
echo "Step 2: Create router with interface (no routes)"
echo "-------------------------------------------------"
terraform apply -auto-approve

ROUTER_ID=$(terraform state show gcore_cloud_network_router.router | grep -E "^\s+id\s+=" | awk '{print $3}' | tr -d '"')
echo "Router ID: $ROUTER_ID"

echo ""
echo "Step 3: Verify router has no routes"
echo "------------------------------------"
terraform state show gcore_cloud_network_router.router | grep -A 5 "routes"

echo ""
echo "Step 4: Add route to router configuration"
echo "------------------------------------------"
cat > main_with_route.tf << 'EOF'
terraform {
  required_providers {
    gcore = {
      source = "local.gcore.com/repo/gcore"
    }
  }
}

data "gcore_cloud_projects" "my_projects" {}

data "gcore_cloud_region" "rg" {
  name = "Luxembourg-2"
}

locals {
  project_id = [for v in data.gcore_cloud_projects.my_projects.projects : v.id if v.name == "qa-terraform"]
}

resource "gcore_cloud_network" "network" {
  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id
  name       = "qa-terr-router-network"
}

resource "gcore_cloud_network_subnet" "sb" {
  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id
  name       = "qa-terr-router-subnet"
  cidr       = "192.168.0.0/24"
  network_id = gcore_cloud_network.network.id
}

resource "gcore_cloud_network_router" "router" {
  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id
  name       = "qa-terr-rename"
  external_gateway_info = {
    enable_snat = true
  }
  interfaces = [
    {
      subnet_id = gcore_cloud_network_subnet.sb.id
      type      = "subnet"
    }
  ]
  routes = [{
    destination = "10.0.3.0/24"
    nexthop     = "192.168.0.1"
  }]
}
EOF

cp main_with_route.tf main.tf
terraform apply -auto-approve

echo ""
echo "Step 5: Verify router now has the route"
echo "----------------------------------------"
terraform state show gcore_cloud_network_router.router | grep -A 10 "routes"

echo ""
echo "Step 6: Remove route from configuration"
echo "----------------------------------------"
cat > main_without_route.tf << 'EOF'
terraform {
  required_providers {
    gcore = {
      source = "local.gcore.com/repo/gcore"
    }
  }
}

data "gcore_cloud_projects" "my_projects" {}

data "gcore_cloud_region" "rg" {
  name = "Luxembourg-2"
}

locals {
  project_id = [for v in data.gcore_cloud_projects.my_projects.projects : v.id if v.name == "qa-terraform"]
}

resource "gcore_cloud_network" "network" {
  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id
  name       = "qa-terr-router-network"
}

resource "gcore_cloud_network_subnet" "sb" {
  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id
  name       = "qa-terr-router-subnet"
  cidr       = "192.168.0.0/24"
  network_id = gcore_cloud_network.network.id
}

resource "gcore_cloud_network_router" "router" {
  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id
  name       = "qa-terr-rename"
  external_gateway_info = {
    enable_snat = true
  }
  interfaces = [
    {
      subnet_id = gcore_cloud_network_subnet.sb.id
      type      = "subnet"
    }
  ]
}
EOF

cp main_without_route.tf main.tf
terraform apply -auto-approve

echo ""
echo "Step 7: Check if route was actually deleted"
echo "--------------------------------------------"
echo "Checking Terraform state:"
terraform state show gcore_cloud_network_router.router | grep -A 10 "routes"

echo ""
echo "===================================================================="
echo "BUG: If the route still appears in state or on the actual router,"
echo "then the route deletion is not working properly."
echo "Expected: routes = []"
echo "===================================================================="

echo ""
echo "Cleaning up resources..."
terraform destroy -auto-approve

echo ""
echo "Test complete!"
