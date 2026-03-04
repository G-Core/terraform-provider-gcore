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
echo "TEST: existing-volume source - GCLOUD2-21138 FIX"
echo "================================================"
echo ""

# Use an existing available bootable volume
EXISTING_VOLUME_ID="b2ea3a8a-a493-4726-a759-14b977768413"
echo "Using existing bootable volume: $EXISTING_VOLUME_ID"
echo ""

# Clean up any previous state
rm -f terraform.tfstate terraform.tfstate.backup terraform.log main.tf 2>/dev/null || true

#############################################
# Create instance with existing-volume source
#############################################
echo "=== Creating instance with existing-volume source ==="

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

# Create instance using existing-volume source
resource "gcore_cloud_instance" "test_existing_vol" {
  project_id = 379987
  region_id  = 76
  name       = "test-existing-volume-fix"
  flavor     = "g1-standard-1-2"

  volumes = [
    {
      source     = "existing-volume"
      volume_id  = "$EXISTING_VOLUME_ID"
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
  value = gcore_cloud_instance.test_existing_vol.id
}

output "volume_id_from_state" {
  value = gcore_cloud_instance.test_existing_vol.volumes[0].volume_id
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
# Verify API request in logs
#############################################
echo ""
echo "=== Verifying volume_id in API request ==="

echo "Checking for 'volume_id' in create request body..."
if grep -E '"volume_id"\s*:\s*"[^"]+"' terraform.log | head -5; then
    echo "✅ PASS: Found 'volume_id' field in API request"
else
    echo "❌ FAIL: Did not find 'volume_id' in API request"
fi

echo ""
echo "Checking request body with source=existing-volume..."
grep -A20 'Request Body' terraform.log | grep -A15 '"source":"existing-volume"' | head -20 || echo "Could not extract request body"

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
# Test state contains volume_id
#############################################
echo ""
echo "=== Checking state for volume_id ==="
VOLUME_FROM_STATE=$(terraform output -raw volume_id_from_state 2>/dev/null || echo "")
if [ "$VOLUME_FROM_STATE" = "$EXISTING_VOLUME_ID" ]; then
    echo "✅ volume_id correctly stored in state: $VOLUME_FROM_STATE"
else
    echo "⚠️  volume_id in state: '$VOLUME_FROM_STATE' (expected: $EXISTING_VOLUME_ID)"
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
