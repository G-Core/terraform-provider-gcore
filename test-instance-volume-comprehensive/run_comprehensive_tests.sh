#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="/Users/user/repos/gcore-terraform"
EVIDENCE_DIR="$SCRIPT_DIR/evidence"

echo "=============================================="
echo "Instance + Volume Comprehensive Test Suite"
echo "=============================================="
echo ""

# Load credentials
if [ -f "$REPO_ROOT/.env" ]; then
    echo "=== Loading credentials ==="
    set -o allexport
    source "$REPO_ROOT/.env"
    set +o allexport
    echo "✓ Credentials loaded"
else
    echo "ERROR: .env not found at $REPO_ROOT/.env"
    exit 1
fi

export TF_CLI_CONFIG_FILE="$REPO_ROOT/.terraformrc"
export TF_LOG=DEBUG
export TF_LOG_PATH="$EVIDENCE_DIR/terraform.log"

cd "$SCRIPT_DIR"

# Clean any existing state
rm -f terraform.tfstate terraform.tfstate.backup 2>/dev/null || true
rm -f "$EVIDENCE_DIR"/*.log "$EVIDENCE_DIR"/*.json 2>/dev/null || true

echo ""
echo "=============================================="
echo "TEST 1: Create Boot Volume + Instance"
echo "=============================================="

echo "=== Planning ==="
terraform plan -var "api_key=$GCORE_API_KEY" -no-color 2>&1 | tee "$EVIDENCE_DIR/test1_plan.log"

echo ""
echo "=== Applying ==="
terraform apply -auto-approve -var "api_key=$GCORE_API_KEY" -no-color 2>&1 | tee "$EVIDENCE_DIR/test1_apply.log"

echo ""
echo "=== Saving state snapshot ==="
cp terraform.tfstate "$EVIDENCE_DIR/test1_state.json"

INSTANCE_ID=$(terraform output -raw instance_id)
BOOT_VOLUME_ID=$(terraform output -raw boot_volume_id)
echo "✓ Instance ID: $INSTANCE_ID"
echo "✓ Boot Volume ID: $BOOT_VOLUME_ID"

echo ""
echo "=============================================="
echo "TEST 2: Drift Detection (No Changes Expected)"
echo "=============================================="

echo "=== Running plan -detailed-exitcode ==="
set +e
terraform plan -var "api_key=$GCORE_API_KEY" -detailed-exitcode -no-color 2>&1 | tee "$EVIDENCE_DIR/test2_drift.log"
DRIFT_EXIT_CODE=$?
set -e

if [ $DRIFT_EXIT_CODE -eq 0 ]; then
    echo ""
    echo "✓ TEST 2 PASSED: No drift detected (exit code 0)"
elif [ $DRIFT_EXIT_CODE -eq 2 ]; then
    echo ""
    echo "✗ TEST 2 FAILED: Drift detected (exit code 2)"
    echo "See $EVIDENCE_DIR/test2_drift.log for details"
    # Don't exit - continue with other tests
else
    echo ""
    echo "✗ TEST 2 ERROR: Terraform plan failed (exit code $DRIFT_EXIT_CODE)"
fi

echo ""
echo "=============================================="
echo "TEST 3: Add Data Volume (Multiple Volumes)"
echo "=============================================="

echo "=== Planning with data volume ==="
terraform plan -var "api_key=$GCORE_API_KEY" -var "test_data_volume=true" -no-color 2>&1 | tee "$EVIDENCE_DIR/test3_plan.log"

echo ""
echo "=== Applying with data volume ==="
terraform apply -auto-approve -var "api_key=$GCORE_API_KEY" -var "test_data_volume=true" -no-color 2>&1 | tee "$EVIDENCE_DIR/test3_apply.log"

echo ""
echo "=== Saving state snapshot ==="
cp terraform.tfstate "$EVIDENCE_DIR/test3_state.json"

DATA_VOLUME_ID=$(terraform output -raw data_volume_id 2>/dev/null || echo "none")
echo "✓ Data Volume ID: $DATA_VOLUME_ID"

echo ""
echo "=== Drift check after adding volume ==="
set +e
terraform plan -var "api_key=$GCORE_API_KEY" -var "test_data_volume=true" -detailed-exitcode -no-color 2>&1 | tee "$EVIDENCE_DIR/test3_drift.log"
DRIFT_EXIT_CODE=$?
set -e

if [ $DRIFT_EXIT_CODE -eq 0 ]; then
    echo "✓ TEST 3 PASSED: No drift after adding data volume"
elif [ $DRIFT_EXIT_CODE -eq 2 ]; then
    echo "✗ TEST 3 DRIFT: Changes detected after adding data volume"
fi

echo ""
echo "=============================================="
echo "TEST 4: Import Test"
echo "=============================================="

# Save current state
cp terraform.tfstate "$EVIDENCE_DIR/test4_before_import.json"

# Remove instance from state
echo "=== Removing instance from state ==="
terraform state rm gcore_cloud_instance.test 2>&1 | tee "$EVIDENCE_DIR/test4_state_rm.log"

# Import the instance back
echo ""
echo "=== Importing instance ==="
terraform import -var "api_key=$GCORE_API_KEY" -var "test_data_volume=true" \
    "gcore_cloud_instance.test" "379987/76/$INSTANCE_ID" 2>&1 | tee "$EVIDENCE_DIR/test4_import.log"

echo ""
echo "=== Post-import state ==="
terraform show -no-color 2>&1 | tee "$EVIDENCE_DIR/test4_show.log"
cp terraform.tfstate "$EVIDENCE_DIR/test4_after_import.json"

echo ""
echo "=== Import drift check ==="
set +e
terraform plan -var "api_key=$GCORE_API_KEY" -var "test_data_volume=true" -detailed-exitcode -no-color 2>&1 | tee "$EVIDENCE_DIR/test4_drift.log"
IMPORT_DRIFT_CODE=$?
set -e

if [ $IMPORT_DRIFT_CODE -eq 0 ]; then
    echo "✓ TEST 4 PASSED: Import successful, no drift"
elif [ $IMPORT_DRIFT_CODE -eq 2 ]; then
    echo "⚠ TEST 4 WARNING: Import shows drift (may need manual volume config)"
fi

echo ""
echo "=============================================="
echo "TEST 5: Cleanup (Destroy)"
echo "=============================================="

echo "=== Destroying resources ==="
terraform destroy -auto-approve -var "api_key=$GCORE_API_KEY" -var "test_data_volume=true" -no-color 2>&1 | tee "$EVIDENCE_DIR/test5_destroy.log"

echo ""
echo "=== Final state ==="
terraform show -no-color 2>&1 | tee "$EVIDENCE_DIR/test5_final.log"

echo ""
echo "=============================================="
echo "TEST SUMMARY"
echo "=============================================="
echo "Evidence stored in: $EVIDENCE_DIR"
echo ""
echo "Test Results:"
echo "  Test 1 (Create): Check test1_apply.log"
echo "  Test 2 (Drift):  Exit code was $DRIFT_EXIT_CODE (0=pass, 2=fail)"
echo "  Test 3 (Multi):  Check test3_apply.log"
echo "  Test 4 (Import): Exit code was $IMPORT_DRIFT_CODE"
echo "  Test 5 (Destroy): Check test5_destroy.log"
echo ""
echo "=============================================="
echo "Tests completed!"
echo "=============================================="
