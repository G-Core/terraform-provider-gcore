#!/bin/bash
set -e

echo "=========================================="
echo "Testing Tags Fix - GCLOUD2-20778"
echo "=========================================="
echo ""

# Setup
source ../.env
export TF_CLI_CONFIG_FILE=../.terraformrc

# Step 1: Create LB WITHOUT tags
echo "Step 1: Creating LB without tags..."
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
  flavor     = "lb1-2-4"
  name       = "qa-lb-tags-test-fix"
  # NO TAGS YET
}
EOF

terraform apply -auto-approve
echo ""
echo "✓ LB created without tags"
echo ""

# Check initial state
echo "Current state:"
terraform show | grep -A 5 "tags"
echo ""

# Step 2: Add tags to config
echo "Step 2: Adding tags to configuration..."
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
  flavor     = "lb1-2-4"
  name       = "qa-lb-tags-test-fix"

  # ADDING TAGS
  tags = {
    "qa"          = "load-balancer"
    "environment" = "test"
  }
}
EOF

echo ""
echo "Step 3: Applying tags change (this is where the error would occur)..."
echo "=========================================="
if terraform apply -auto-approve 2>&1 | tee apply_result.log; then
    echo ""
    echo "✅✅✅ SUCCESS! Tags applied without error ✅✅✅"
    echo ""
    echo "Final state:"
    terraform show | grep -A 10 "tags"
    echo ""
    echo "🎉 BUG FIXED! The tags inconsistency error is RESOLVED!"
    exit 0
else
    echo ""
    echo "❌ FAILED: Error occurred"
    if grep -q "tags_v2.*appeared" apply_result.log; then
        echo "🐛 BUG STILL EXISTS: tags_v2 inconsistency error"
    else
        echo "Different error - check apply_result.log"
    fi
    exit 1
fi
