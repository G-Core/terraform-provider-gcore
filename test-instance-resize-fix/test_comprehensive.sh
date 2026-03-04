#!/bin/bash
set -e

echo "=== Comprehensive Instance Resize Testing ==="
echo ""

# Load credentials
echo "Loading credentials..."
source ../.env
export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/.terraformrc"

echo ""
echo "=== TEST 1: Create instance with small flavor ==="
echo "  Initial config: g1-standard-1-2, 10GB volume"
terraform apply -auto-approve -var="flavor_id=g1-standard-1-2" -var="volume_size=10"

# Capture initial state
INSTANCE_ID_1=$(terraform output -raw instance_id)
VOLUME_ID_1=$(terraform output -json instance_volumes | jq -r '.[0].volume_id')

echo ""
echo "✅ Instance created: $INSTANCE_ID_1"
echo "✅ Volume ID populated: $VOLUME_ID_1"

if [ "$VOLUME_ID_1" = "null" ] || [ -z "$VOLUME_ID_1" ]; then
    echo "❌ FAIL: volume_id not populated correctly!"
    exit 1
fi

echo ""
echo "=== TEST 2: Drift Detection (CRITICAL) ==="
echo "  Running terraform plan after apply..."
terraform plan -detailed-exitcode
EXIT_CODE=$?

if [ $EXIT_CODE -eq 0 ]; then
    echo "✅ PASS: No drift detected"
elif [ $EXIT_CODE -eq 2 ]; then
    echo "❌ FAIL: Drift detected - changes found on second plan!"
    terraform plan
    exit 1
else
    echo "❌ FAIL: Terraform plan error"
    exit 1
fi

echo ""
echo "=== TEST 3: Flavor Change (Resize Operation) ==="
echo "  Changing flavor: g1-standard-1-2 → g1-standard-2-4"
echo "  Expected: POST to /changeflavor endpoint, NO replacement"
terraform apply -auto-approve -var="flavor_id=g1-standard-2-4" -var="volume_size=10"

# Verify instance ID unchanged
INSTANCE_ID_2=$(terraform output -raw instance_id)
VOLUME_ID_2=$(terraform output -json instance_volumes | jq -r '.[0].volume_id')

echo ""
if [ "$INSTANCE_ID_1" = "$INSTANCE_ID_2" ]; then
    echo "✅ PASS: Instance updated in-place (ID unchanged)"
    echo "  Before: $INSTANCE_ID_1"
    echo "  After:  $INSTANCE_ID_2"
else
    echo "❌ FAIL: Instance was REPLACED!"
    echo "  Before: $INSTANCE_ID_1"
    echo "  After:  $INSTANCE_ID_2"
    exit 1
fi

if [ "$VOLUME_ID_1" = "$VOLUME_ID_2" ]; then
    echo "✅ PASS: Volume ID stable: $VOLUME_ID_2"
else
    echo "⚠️  WARNING: Volume ID changed (might be expected if volume recreated)"
    echo "  Before: $VOLUME_ID_1"
    echo "  After:  $VOLUME_ID_2"
fi

echo ""
echo "=== TEST 4: Drift Check After Flavor Change ==="
terraform plan -detailed-exitcode -var="flavor_id=g1-standard-2-4" -var="volume_size=10"
if [ $? -eq 0 ]; then
    echo "✅ PASS: No drift after flavor change"
else
    echo "❌ FAIL: Drift detected after flavor change!"
    terraform plan -var="flavor_id=g1-standard-2-4" -var="volume_size=10"
    exit 1
fi

echo ""
echo "=== TEST 5: Volume Extension ==="
echo "  Extending volume: 10GB → 20GB"
echo "  Expected: POST to /extend endpoint, NO replacement"
terraform apply -auto-approve -var="flavor_id=g1-standard-2-4" -var="volume_size=20"

# Verify instance ID unchanged
INSTANCE_ID_3=$(terraform output -raw instance_id)
VOLUME_SIZE_3=$(terraform output -json instance_volumes | jq -r '.[0].size')

if [ "$INSTANCE_ID_2" = "$INSTANCE_ID_3" ]; then
    echo "✅ PASS: Instance still in-place (ID unchanged)"
else
    echo "❌ FAIL: Instance replaced during volume extension!"
    exit 1
fi

if [ "$VOLUME_SIZE_3" = "20" ]; then
    echo "✅ PASS: Volume size updated to 20GB"
else
    echo "❌ FAIL: Volume size not updated correctly (got: $VOLUME_SIZE_3)"
    exit 1
fi

echo ""
echo "=== TEST 6: Final Drift Check ==="
terraform plan -detailed-exitcode -var="flavor_id=g1-standard-2-4" -var="volume_size=20"
if [ $? -eq 0 ]; then
    echo "✅ PASS: No drift after volume extension"
else
    echo "❌ FAIL: Drift after volume extension!"
    terraform plan -var="flavor_id=g1-standard-2-4" -var="volume_size=20"
    exit 1
fi

echo ""
echo "=== ALL TESTS PASSED ✅ ==="
echo ""
echo "Summary:"
echo "  - Instance ID stable across operations: $INSTANCE_ID_3"
echo "  - Flavor changed: g1-standard-1-2 → g1-standard-2-4 ✅"
echo "  - Volume extended: 10GB → 20GB ✅"
echo "  - No drift detected in any test ✅"
echo ""
