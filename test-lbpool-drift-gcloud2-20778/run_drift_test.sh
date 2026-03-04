#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "================================================"
echo "GCLOUD2-20778 Drift Detection Test"
echo "================================================"
echo ""

# Load credentials
if [ -f ../.env ]; then
    source ../.env
    echo "✓ Loaded credentials from .env"
else
    echo "❌ .env file not found"
    exit 1
fi

# Set terraform config
export TF_CLI_CONFIG_FILE=/Users/user/repos/gcore-terraform/.terraformrc

echo ""
echo "================================================"
echo "STEP 1: Creating LB + Listener (WITHOUT pool)"
echo "================================================"
echo ""

# Create config without pool
cat > main.tf << 'EOF'
terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

data "gcore_cloud_projects" "my_projects" {
  name = "default"
}

locals {
  project_id = [for p in data.gcore_cloud_projects.my_projects.items : p.id]
}

data "gcore_cloud_region" "rg" {
  region_id = 76
}

resource "gcore_cloud_load_balancer" "lb" {
  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id
  flavor     = "lb1-1-2"
  name       = "qa-drift-test-lb"
}

resource "gcore_cloud_load_balancer_listener" "ls" {
  project_id       = local.project_id[0]
  region_id        = data.gcore_cloud_region.rg.id
  load_balancer_id = gcore_cloud_load_balancer.lb.id
  name             = "qa-ls-name"
  protocol         = "HTTP"
  protocol_port    = 80

  # This user_list is key to reproducing the bug
  user_list = [
    {
      username           = "testuser"
      encrypted_password = "$5$isRr.HJ1IrQP38.m$oViu3DJOpUG2ZsjCBtbITV3mqpxxbZfyWJojLPNSPO5"
    }
  ]
}
EOF

echo "Applying configuration (LB + Listener)..."
terraform apply -auto-approve > step1_apply.log 2>&1

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ LB and Listener created successfully${NC}"
else
    echo -e "${RED}❌ Failed to create resources${NC}"
    cat step1_apply.log
    exit 1
fi

echo ""
echo "Checking for drift after initial creation..."
terraform plan > step1_plan.log 2>&1

if grep -q "No changes" step1_plan.log; then
    echo -e "${GREEN}✓ No drift detected after initial creation${NC}"
else
    echo -e "${YELLOW}⚠ Unexpected drift after initial creation${NC}"
    cat step1_plan.log
fi

echo ""
echo "================================================"
echo "STEP 2: Adding Pool to Configuration"
echo "================================================"
echo ""

sleep 2

# Add pool to config
cat >> main.tf << 'EOF'

# Adding pool - this is where drift might appear
resource "gcore_cloud_load_balancer_pool" "lb_pool" {
  project_id   = local.project_id[0]
  region_id    = data.gcore_cloud_region.rg.id
  lb_algorithm = "LEAST_CONNECTIONS"
  name         = "pool-drift-test"
  protocol     = "HTTP"
  listener_id  = gcore_cloud_load_balancer_listener.ls.id
}
EOF

echo "Applying pool addition..."
terraform apply -auto-approve > step2_apply.log 2>&1

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Pool created and attached successfully${NC}"
else
    echo -e "${RED}❌ Failed to create pool${NC}"
    cat step2_apply.log
    exit 1
fi

echo ""
echo "================================================"
echo "STEP 3: CRITICAL TEST - Checking for Drift"
echo "================================================"
echo ""

sleep 2

terraform plan > step3_plan.log 2>&1

echo ""
echo "================================================"
echo "RESULTS"
echo "================================================"
echo ""

if grep -q "No changes" step3_plan.log; then
    echo -e "${GREEN}✅ SUCCESS: No drift detected!${NC}"
    echo ""
    echo "The bug is FIXED. The listener does not show spurious drift"
    echo "after adding a pool, even with user_list configured."
    echo ""
else
    echo -e "${RED}❌ DRIFT DETECTED: Bug is present!${NC}"
    echo ""
    echo "Terraform detected changes in the listener:"
    grep -A 30 "gcore_cloud_load_balancer_listener" step3_plan.log || cat step3_plan.log
    echo ""
    echo "This indicates the bug from GCLOUD2-20778 is still present."
    echo ""
fi

echo "================================================"
echo "Test logs saved:"
echo "  - step1_apply.log (initial creation)"
echo "  - step1_plan.log  (drift check after creation)"
echo "  - step2_apply.log (pool addition)"
echo "  - step3_plan.log  (final drift check)"
echo "================================================"
echo ""

# Ask about cleanup
read -p "Do you want to destroy the test resources? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo ""
    echo "Cleaning up resources..."
    terraform destroy -auto-approve > cleanup.log 2>&1
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Resources destroyed successfully${NC}"
    else
        echo -e "${RED}❌ Cleanup failed - check cleanup.log${NC}"
    fi
fi
