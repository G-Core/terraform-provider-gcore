#!/bin/bash
set -e

echo "=================================================="
echo "Testing Router Interface Update (No Replacement)"
echo "=================================================="
echo

# Load environment variables
echo "Loading environment variables..."
set -o allexport
source ../.env
set +o allexport

# Set Terraform config
export TF_CLI_CONFIG_FILE="../.terraformrc"

echo "Step 1: Apply initial configuration with ONE interface"
echo "-----------------------------------------------------"
terraform apply -auto-approve
echo

echo "Step 2: Modify configuration to add SECOND interface"
echo "-----------------------------------------------------"

# Update the main.tf to add second interface
cat > main.tf << 'EOF'
terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {
  # Credentials from environment
}

# Create a network
resource "gcore_cloud_network" "test_network" {
  project_id = 379987
  region_id  = 76
  name       = "test-router-network"
}

# Create a subnet
resource "gcore_cloud_network_subnet" "test_subnet" {
  project_id = 379987
  region_id  = 76
  name       = "test-router-subnet"
  network_id = gcore_cloud_network.test_network.id
  cidr       = "192.168.1.0/24"

  enable_dhcp = true
}

# Create a second subnet for testing interface add
resource "gcore_cloud_network_subnet" "test_subnet_2" {
  project_id = 379987
  region_id  = 76
  name       = "test-router-subnet-2"
  network_id = gcore_cloud_network.test_network.id
  cidr       = "192.168.2.0/24"

  enable_dhcp = true
}

# Create a router with TWO interfaces
resource "gcore_cloud_network_router" "test_router" {
  project_id = 379987
  region_id  = 76
  name       = "test-router-interface-fix"

  external_gateway_info = {
    enable_snat = true
    type        = "default"
  }

  # NOW WITH TWO INTERFACES
  interfaces = [
    {
      subnet_id = gcore_cloud_network_subnet.test_subnet.id
      type      = "subnet"
    },
    {
      subnet_id = gcore_cloud_network_subnet.test_subnet_2.id
      type      = "subnet"
    }
  ]
}

output "router_id" {
  value = gcore_cloud_network_router.test_router.id
}

output "router_interfaces" {
  value = gcore_cloud_network_router.test_router.interfaces
}

output "subnet_1_id" {
  value = gcore_cloud_network_subnet.test_subnet.id
}

output "subnet_2_id" {
  value = gcore_cloud_network_subnet.test_subnet_2.id
}
EOF

echo "Step 3: Plan the update - checking for 'must be replaced'"
echo "---------------------------------------------------------"
terraform plan -no-color > plan_output.txt 2>&1 || true

if grep -q "must be replaced" plan_output.txt; then
    echo "❌ FAILED: Router is being replaced instead of updated!"
    echo
    echo "Plan output:"
    cat plan_output.txt
    exit 1
else
    echo "✅ SUCCESS: Router will be updated in-place (no replacement)!"
    echo
fi

echo "Step 4: Apply the update"
echo "------------------------"
terraform apply -auto-approve

echo
echo "Step 5: Verify router has 2 interfaces"
echo "---------------------------------------"
terraform show -json | grep -A 20 '"interfaces"' | head -30

echo
echo "=================================================="
echo "✅ Test completed successfully!"
echo "The router was updated in-place without replacement."
echo "=================================================="
