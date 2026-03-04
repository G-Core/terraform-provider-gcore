#!/bin/bash
set -e

# Load environment
source /Users/user/repos/gcore-terraform/.env
export TF_CLI_CONFIG_FILE=/Users/user/repos/gcore-terraform/.terraformrc
export TF_LOG=TRACE
export TF_LOG_PATH="terraform.log"

TEST_DIR="/Users/user/repos/gcore-terraform/test-existing-volume"
cd "$TEST_DIR"

echo "================================================"
echo "COMPREHENSIVE VOLUME TEST - GCLOUD2-21138 FIX"
echo "================================================"
echo ""

# Clean up any previous state
rm -f terraform.tfstate terraform.tfstate.backup terraform.log 2>/dev/null || true

#############################################
# STEP 1: Create bootable volume from image
#############################################
echo "=== STEP 1: Create bootable volume from image ==="

cat > main.tf << 'EOF'
terraform {
  required_providers {
    gcore = {
      source  = "gcore/gcore"
      version = "9999.0.0"
    }
  }
}

provider "gcore" {}

# Create a bootable volume from Ubuntu image
resource "gcore_cloud_volume" "test_boot_volume" {
  project_id = 379987
  region_id  = 76
  name       = "test-existing-volume-boot"
  source     = "image"
  size       = 10
  type_name  = "ssd_hiiops"
  image_id   = "f84ddba3-7a5a-4199-931a-250e981d16fb"  # ubuntu-25.04-x64
  tags = {
    purpose = "test-existing-volume-fix"
  }
}

output "volume_id" {
  value = gcore_cloud_volume.test_boot_volume.id
}
EOF

# Skip init when using dev overrides
terraform apply -auto-approve

VOLUME_ID=$(terraform output -raw volume_id)
echo ""
echo "✅ Created volume: $VOLUME_ID"
echo ""

# Wait for volume to be available
sleep 5

#############################################
# STEP 2: Test existing-volume source
#############################################
echo "=== STEP 2: Create instance with existing-volume source ==="
echo "This tests the volume_id fix - GCLOUD2-21138"
echo ""

cat > main.tf << EOF
terraform {
  required_providers {
    gcore = {
      source  = "gcore/gcore"
      version = "9999.0.0"
    }
  }
}

provider "gcore" {}

# Keep the volume resource to avoid deletion
resource "gcore_cloud_volume" "test_boot_volume" {
  project_id = 379987
  region_id  = 76
  name       = "test-existing-volume-boot"
  source     = "image"
  size       = 10
  type_name  = "ssd_hiiops"
  image_id   = "f84ddba3-7a5a-4199-931a-250e981d16fb"
  tags = {
    purpose = "test-existing-volume-fix"
  }
}

# Create instance using existing-volume source
resource "gcore_cloud_instance" "test_existing_vol" {
  project_id = 379987
  region_id  = 76
  name       = "test-existing-volume-instance"
  flavor     = "g1-standard-1-2"

  volumes = [
    {
      source     = "existing-volume"
      volume_id  = gcore_cloud_volume.test_boot_volume.id
      boot_index = 0
    }
  ]

  interfaces = [
    {
      type      = "external"
      ip_family = "dual"
    }
  ]
}

output "volume_id" {
  value = gcore_cloud_volume.test_boot_volume.id
}

output "instance_id" {
  value = gcore_cloud_instance.test_existing_vol.id
}
EOF

echo "Applying instance with existing-volume..."
terraform apply -auto-approve

INSTANCE_ID=$(terraform output -raw instance_id)
echo ""
echo "✅ Created instance: $INSTANCE_ID"
echo ""

#############################################
# STEP 3: Verify API request in logs
#############################################
echo "=== STEP 3: Verify volume_id in API request ==="

echo "Checking for volume_id in request body..."
if grep -q '"volume_id"' terraform.log; then
    echo "✅ PASS: Found 'volume_id' in API request"
    grep -o '"volume_id":"[^"]*"' terraform.log | head -3
else
    echo "❌ FAIL: Did not find 'volume_id' in API request"
fi

echo ""
echo "Checking that 'id' is NOT used in volume create request..."
# Look for request body with source=existing-volume and check field name
grep -A5 '"source":"existing-volume"' terraform.log | head -10 || true

#############################################
# STEP 4: Test drift detection (refresh)
#############################################
echo ""
echo "=== STEP 4: Test drift detection ==="
terraform plan -detailed-exitcode && echo "✅ No drift detected" || echo "⚠️  Drift detected"

#############################################
# STEP 5: Cleanup
#############################################
echo ""
echo "=== STEP 5: Cleanup ==="
terraform destroy -auto-approve

echo ""
echo "================================================"
echo "TEST COMPLETE"
echo "================================================"
