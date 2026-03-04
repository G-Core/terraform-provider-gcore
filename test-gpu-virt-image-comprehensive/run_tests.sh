#!/bin/bash
set -e

cd /Users/user/repos/gcore-terraform/test-gpu-virt-image-comprehensive

export TF_CLI_CONFIG_FILE=/Users/user/repos/gcore-terraform/.terraformrc
export GCORE_API_KEY='21788$1e278ce67b6aa33f178122658b1dd0210d0edff453d348acb9b68bffea6a635b7791925ddda198d5678a4dc20269fe04a263ca92c7e5aa41ea79075f89b66bf6'
export GCORE_CLIENT=3621
export GCORE_CLOUD_PROJECT_ID=379987
export GCORE_CLOUD_REGION_ID=76
export TF_LOG=DEBUG
export TF_LOG_PATH=terraform.log

echo "============================================"
echo "GPU Virtual Cluster Image Comprehensive Test"
echo "============================================"
echo ""

# Test 1: Create resource
echo "=== Test 1: Create Resource ==="
terraform apply -auto-approve
echo "✅ Create completed"
echo ""

# Test 2: Drift detection (no changes on second plan)
echo "=== Test 2: Drift Detection ==="
terraform plan -detailed-exitcode
DRIFT_EXIT=$?
if [ $DRIFT_EXIT -eq 0 ]; then
    echo "✅ No drift detected (exit code 0)"
elif [ $DRIFT_EXIT -eq 2 ]; then
    echo "❌ DRIFT DETECTED (exit code 2) - This is a bug!"
    exit 1
else
    echo "❌ Plan failed (exit code $DRIFT_EXIT)"
    exit 1
fi
echo ""

# Test 3: Verify all outputs
echo "=== Test 3: Verify All Attributes ==="
terraform output -json > outputs.json
echo "Outputs:"
cat outputs.json | jq .
echo "✅ All attributes retrieved"
echo ""

# Test 4: Show full state
echo "=== Test 4: Full State ==="
terraform show
echo ""

# Test 5: Delete and recreate (test destroy)
echo "=== Test 5: Destroy Resource ==="
terraform destroy -auto-approve
echo "✅ Destroy completed"
echo ""

echo "============================================"
echo "All tests passed!"
echo "============================================"
