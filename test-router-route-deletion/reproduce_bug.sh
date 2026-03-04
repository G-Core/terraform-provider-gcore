#!/bin/bash

# This script reproduces the router route deletion bug
# GCLOUD2-21144: Routes are not deleted when removed from Terraform config

set -e

echo "=========================================="
echo "Router Route Deletion Bug Reproduction"
echo "=========================================="
echo ""

cd "$(dirname "$0")"

# Load environment variables from .env
if [ -f "../.env" ]; then
    set -o allexport
    source "../.env"
    set +o allexport
    echo "✓ Loaded credentials from .env"
    echo "  API Key: ${GCORE_API_KEY:0:10}..."
    echo "  Project ID: $GCORE_CLOUD_PROJECT_ID"
    echo "  Region ID: $GCORE_CLOUD_REGION_ID"
else
    echo "✗ Error: .env file not found in parent directory"
    exit 1
fi
echo ""

# Step 1: Create initial router with interface (no routes)
echo ""
echo "Step 1: Creating router with interface (no routes)..."
echo "----------------------------------------------"
cat > main.tf << 'EOF'
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
  name       = "qa-terr-router-network"
}

resource "gcore_cloud_network_subnet" "sb" {
  project_id  = 379987
  region_id   = 76
  name        = "qa-terr-router-subnet"
  cidr        = "192.168.0.0/24"
  network_id  = gcore_cloud_network.network.id
}

resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "qa-terr-route-bug"
  external_gateway_info = {
    enable_snat = true
  }
  interfaces = [{
    subnet_id = gcore_cloud_network_subnet.sb.id
    type      = "subnet"
  }]
}
EOF

# No need for terraform init with dev_overrides
echo "Running terraform apply..."
terraform apply -auto-approve

echo ""
echo "✓ Router created"
echo ""
echo "Current state (no routes):"
terraform state show gcore_cloud_network_router.router | grep -A 5 "routes" || echo "  (no routes field or empty)"

# Step 2: Add route to router
echo ""
echo "Step 2: Adding route to router..."
echo "----------------------------------------------"
cat > main.tf << 'EOF'
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
  name       = "qa-terr-router-network"
}

resource "gcore_cloud_network_subnet" "sb" {
  project_id  = 379987
  region_id   = 76
  name        = "qa-terr-router-subnet"
  cidr        = "192.168.0.0/24"
  network_id  = gcore_cloud_network.network.id
}

resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "qa-terr-route-bug"
  external_gateway_info = {
    enable_snat = true
  }
  interfaces = [{
    subnet_id = gcore_cloud_network_subnet.sb.id
    type      = "subnet"
  }]
  routes = [{
    destination = "10.0.3.0/24"
    nexthop     = "192.168.0.1"
  }]
}
EOF

terraform apply -auto-approve

echo ""
echo "✓ Route added"
echo ""
echo "Current state (with route):"
terraform state show gcore_cloud_network_router.router | grep -A 10 "routes"

# Step 3: Remove route from configuration
echo ""
echo "Step 3: Removing route from configuration..."
echo "----------------------------------------------"
cat > main.tf << 'EOF'
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
  name       = "qa-terr-router-network"
}

resource "gcore_cloud_network_subnet" "sb" {
  project_id  = 379987
  region_id   = 76
  name        = "qa-terr-router-subnet"
  cidr        = "192.168.0.0/24"
  network_id  = gcore_cloud_network.network.id
}

resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "qa-terr-route-bug"
  external_gateway_info = {
    enable_snat = true
  }
  interfaces = [{
    subnet_id = gcore_cloud_network_subnet.sb.id
    type      = "subnet"
  }]
  # routes removed - should delete the route
}
EOF

echo ""
echo "Running terraform apply to remove route..."
terraform apply -auto-approve

echo ""
echo "=========================================="
echo "BUG CHECK"
echo "=========================================="
echo ""
echo "Terraform state after removing routes:"
terraform state show gcore_cloud_network_router.router | grep -A 10 "routes"

echo ""
echo "Expected: routes = []"
echo ""
echo "If routes still exist in state, THE BUG IS CONFIRMED!"
echo ""
echo "=========================================="
echo ""

read -p "Press Enter to cleanup resources..." -t 5 || echo ""

echo ""
echo "Cleaning up resources..."
terraform destroy -auto-approve

echo ""
echo "✓ Cleanup complete"
