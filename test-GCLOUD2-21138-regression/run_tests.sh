#!/bin/bash

# GCLOUD2-21138 Regression Test Script
# Tests all bug fixes for cloud_instance resource

set -e

export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/.terraformrc"
export HTTP_PROXY="http://127.0.0.1:9092"
export HTTPS_PROXY="http://127.0.0.1:9092"
export SSL_CERT_FILE="$HOME/.mitmproxy/mitmproxy-ca-cert.pem"

# Load env vars
set -a
source ../.env
set +a

ARTIFACTS_DIR="./artifacts"
mkdir -p "$ARTIFACTS_DIR"

# Helper function to save state and output
save_artifact() {
    local name=$1
    local step=$2
    echo "=== Saving artifact: $name - $step ===" | tee -a "$ARTIFACTS_DIR/test_log.txt"
    terraform show -json > "$ARTIFACTS_DIR/${name}_${step}_state.json" 2>/dev/null || true
    terraform output -json > "$ARTIFACTS_DIR/${name}_${step}_output.json" 2>/dev/null || true
}

# Helper function to run terraform command and log
run_tf() {
    local cmd=$1
    local name=$2
    echo "=== Running: terraform $cmd ($name) ===" | tee -a "$ARTIFACTS_DIR/test_log.txt"
    terraform $cmd 2>&1 | tee -a "$ARTIFACTS_DIR/${name}.log"
    return ${PIPESTATUS[0]}
}

echo "========================================" | tee "$ARTIFACTS_DIR/test_log.txt"
echo "GCLOUD2-21138 Regression Tests Started" | tee -a "$ARTIFACTS_DIR/test_log.txt"
echo "Date: $(date)" | tee -a "$ARTIFACTS_DIR/test_log.txt"
echo "========================================" | tee -a "$ARTIFACTS_DIR/test_log.txt"

# Copy terraform config for reference
cp main.tf "$ARTIFACTS_DIR/main.tf"

# TEST 1: Basic Instance Creation
echo ""
echo "=== TEST 1: Basic Instance Creation ===" | tee -a "$ARTIFACTS_DIR/test_log.txt"
run_tf "apply -auto-approve" "test1_create"
save_artifact "test1" "after_create"

# Capture instance ID for later tests
INSTANCE_ID=$(terraform output -raw instance_id)
PLACEMENT_GROUP_ID=$(terraform output -raw placement_group_id)
SUBNET_ID=$(terraform output -raw subnet_id)
echo "Instance ID: $INSTANCE_ID" | tee -a "$ARTIFACTS_DIR/test_log.txt"
echo "Placement Group ID: $PLACEMENT_GROUP_ID" | tee -a "$ARTIFACTS_DIR/test_log.txt"
echo "Subnet ID: $SUBNET_ID" | tee -a "$ARTIFACTS_DIR/test_log.txt"

# Check for drift
echo ""
echo "=== TEST 1b: Check for drift after creation ===" | tee -a "$ARTIFACTS_DIR/test_log.txt"
run_tf "plan -detailed-exitcode" "test1b_drift_check" || {
    if [ $? -eq 2 ]; then
        echo "DRIFT DETECTED after creation!" | tee -a "$ARTIFACTS_DIR/test_log.txt"
    fi
}

# TEST 2: Flavor Change (should NOT replace instance)
echo ""
echo "=== TEST 2: Flavor Change ===" | tee -a "$ARTIFACTS_DIR/test_log.txt"
run_tf "apply -auto-approve -var='flavor=g1-standard-2-4'" "test2_flavor_change"
save_artifact "test2" "after_flavor_change"

# Verify instance ID is the same
NEW_INSTANCE_ID=$(terraform output -raw instance_id)
if [ "$INSTANCE_ID" = "$NEW_INSTANCE_ID" ]; then
    echo "SUCCESS: Instance ID unchanged after flavor change" | tee -a "$ARTIFACTS_DIR/test_log.txt"
else
    echo "FAILURE: Instance was replaced during flavor change!" | tee -a "$ARTIFACTS_DIR/test_log.txt"
fi

# TEST 3: Tags Update
echo ""
echo "=== TEST 3: Tags Update ===" | tee -a "$ARTIFACTS_DIR/test_log.txt"
run_tf "apply -auto-approve -var='flavor=g1-standard-2-4' -var='tags={env=\"production\",project=\"regression\",new_tag=\"added\"}'" "test3_tags_add"
save_artifact "test3" "after_tags_add"

# TEST 4: Tags Delete (remove one tag)
echo ""
echo "=== TEST 4: Tags Delete ===" | tee -a "$ARTIFACTS_DIR/test_log.txt"
run_tf "apply -auto-approve -var='flavor=g1-standard-2-4' -var='tags={env=\"production\"}'" "test4_tags_delete"
save_artifact "test4" "after_tags_delete"

# Verify instance ID is the same
NEW_INSTANCE_ID=$(terraform output -raw instance_id)
if [ "$INSTANCE_ID" = "$NEW_INSTANCE_ID" ]; then
    echo "SUCCESS: Instance ID unchanged after tag changes" | tee -a "$ARTIFACTS_DIR/test_log.txt"
else
    echo "FAILURE: Instance was replaced during tag changes!" | tee -a "$ARTIFACTS_DIR/test_log.txt"
fi

# TEST 5: Servergroup Add (should NOT replace instance)
echo ""
echo "=== TEST 5: Servergroup Add ===" | tee -a "$ARTIFACTS_DIR/test_log.txt"
run_tf "apply -auto-approve -var='flavor=g1-standard-2-4' -var='tags={env=\"production\"}' -var='servergroup_id=$PLACEMENT_GROUP_ID'" "test5_servergroup_add"
save_artifact "test5" "after_servergroup_add"

# Verify instance ID is the same
NEW_INSTANCE_ID=$(terraform output -raw instance_id)
if [ "$INSTANCE_ID" = "$NEW_INSTANCE_ID" ]; then
    echo "SUCCESS: Instance ID unchanged after servergroup add" | tee -a "$ARTIFACTS_DIR/test_log.txt"
else
    echo "FAILURE: Instance was replaced during servergroup add!" | tee -a "$ARTIFACTS_DIR/test_log.txt"
fi

# TEST 6: Servergroup Remove
echo ""
echo "=== TEST 6: Servergroup Remove ===" | tee -a "$ARTIFACTS_DIR/test_log.txt"
run_tf "apply -auto-approve -var='flavor=g1-standard-2-4' -var='tags={env=\"production\"}'" "test6_servergroup_remove"
save_artifact "test6" "after_servergroup_remove"

# Verify instance ID is the same
NEW_INSTANCE_ID=$(terraform output -raw instance_id)
if [ "$INSTANCE_ID" = "$NEW_INSTANCE_ID" ]; then
    echo "SUCCESS: Instance ID unchanged after servergroup remove" | tee -a "$ARTIFACTS_DIR/test_log.txt"
else
    echo "FAILURE: Instance was replaced during servergroup remove!" | tee -a "$ARTIFACTS_DIR/test_log.txt"
fi

# TEST 7: Import test (remove from state, import, check drift)
echo ""
echo "=== TEST 7: Import and Drift Test ===" | tee -a "$ARTIFACTS_DIR/test_log.txt"
PROJECT_ID=$(terraform output -json | jq -r '.instance_id.value' | cut -d'/' -f1 || echo "")
REGION_ID=$(terraform output -json | jq -r '.instance_id.value' | cut -d'/' -f2 || echo "")

# Get project and region from the state
terraform state show gcore_cloud_instance.test > "$ARTIFACTS_DIR/test7_state_before_import.txt" 2>/dev/null || true

# Remove instance from state (keep infrastructure)
terraform state rm gcore_cloud_instance.test 2>&1 | tee -a "$ARTIFACTS_DIR/test7_import.log"

# Import the instance back using the stored ID
terraform import -var='flavor=g1-standard-2-4' -var='tags={env="production"}' gcore_cloud_instance.test "$INSTANCE_ID" 2>&1 | tee -a "$ARTIFACTS_DIR/test7_import.log"
save_artifact "test7" "after_import"

# Check for drift after import
echo ""
echo "=== TEST 7b: Check for drift after import ===" | tee -a "$ARTIFACTS_DIR/test_log.txt"
terraform plan -detailed-exitcode -var='flavor=g1-standard-2-4' -var='tags={env="production"}' 2>&1 | tee "$ARTIFACTS_DIR/test7b_drift_plan.log"
DRIFT_EXIT=$?
if [ $DRIFT_EXIT -eq 0 ]; then
    echo "SUCCESS: No drift after import" | tee -a "$ARTIFACTS_DIR/test_log.txt"
elif [ $DRIFT_EXIT -eq 2 ]; then
    echo "FAILURE: Drift detected after import!" | tee -a "$ARTIFACTS_DIR/test_log.txt"
fi

# TEST 8: Cleanup
echo ""
echo "=== TEST 8: Cleanup ===" | tee -a "$ARTIFACTS_DIR/test_log.txt"
run_tf "destroy -auto-approve -var='flavor=g1-standard-2-4' -var='tags={env=\"production\"}'" "test8_destroy"

echo ""
echo "========================================" | tee -a "$ARTIFACTS_DIR/test_log.txt"
echo "GCLOUD2-21138 Regression Tests Complete" | tee -a "$ARTIFACTS_DIR/test_log.txt"
echo "Date: $(date)" | tee -a "$ARTIFACTS_DIR/test_log.txt"
echo "========================================" | tee -a "$ARTIFACTS_DIR/test_log.txt"

echo ""
echo "Artifacts saved to: $ARTIFACTS_DIR"
ls -la "$ARTIFACTS_DIR"
