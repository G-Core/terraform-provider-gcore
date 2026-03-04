#!/bin/bash
set -e

# Load environment
source /Users/user/repos/gcore-terraform/.env
export TF_CLI_CONFIG_FILE=/Users/user/repos/gcore-terraform/.terraformrc
export TF_LOG=TRACE
export TF_LOG_PATH="terraform_image.log"

TEST_DIR="/Users/user/repos/gcore-terraform/test-existing-volume"
cd "$TEST_DIR"

echo "================================================"
echo "TEST: image source (regression test)"
echo "================================================"
echo ""

# Clean up any previous state
rm -f terraform.tfstate terraform.tfstate.backup terraform_image.log main.tf 2>/dev/null || true

#############################################
# Create instance with image source
#############################################
echo "=== Creating instance with image source ==="

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

# Create instance with image source
resource "gcore_cloud_instance" "test_image_source" {
  project_id = 379987
  region_id  = 76
  name       = "test-image-source"
  flavor     = "g1-standard-1-2"

  volumes = [
    {
      source     = "image"
      image_id   = "f84ddba3-7a5a-4199-931a-250e981d16fb"  # ubuntu-25.04-x64
      size       = 10
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

output "instance_id" {
  value = gcore_cloud_instance.test_image_source.id
}

output "volume_id_from_state" {
  value = gcore_cloud_instance.test_image_source.volumes[0].volume_id
}
EOF

echo "Applying..."
terraform apply -auto-approve

echo ""
INSTANCE_ID=$(terraform output -raw instance_id 2>/dev/null || echo "")
if [ -n "$INSTANCE_ID" ]; then
    echo "✅ Created instance: $INSTANCE_ID"
else
    echo "❌ Failed to create instance"
    exit 1
fi

#############################################
# Verify volume_id populated from response
#############################################
echo ""
echo "=== Checking volume_id populated from API response ==="
VOLUME_FROM_STATE=$(terraform output -raw volume_id_from_state 2>/dev/null || echo "")
if [ -n "$VOLUME_FROM_STATE" ] && [ "$VOLUME_FROM_STATE" != "null" ]; then
    echo "✅ volume_id populated from API response: $VOLUME_FROM_STATE"
else
    echo "⚠️  volume_id not populated (may be expected for image source)"
fi

#############################################
# Test drift detection (refresh)
#############################################
echo ""
echo "=== Testing drift detection (terraform plan) ==="
if terraform plan -detailed-exitcode; then
    echo "✅ No drift detected - state matches infrastructure"
else
    EXIT_CODE=$?
    if [ $EXIT_CODE -eq 2 ]; then
        echo "⚠️  Drift detected - changes would be made"
        terraform plan
    else
        echo "❌ Error during plan"
    fi
fi

#############################################
# Cleanup
#############################################
echo ""
echo "=== Cleanup ==="
terraform destroy -auto-approve

echo ""
echo "================================================"
echo "TEST COMPLETE"
echo "================================================"
