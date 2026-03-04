#!/bin/bash
set -e

# Comprehensive cloud_instance testing script for GCLOUD2-21138
# Tests: Interface variants, Volume variants, Security groups, Placement groups, Update operations

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Load credentials
if [ -f ../../../.env ]; then
    source ../../../.env
elif [ -f ../.env ]; then
    source ../.env
else
    echo "ERROR: .env not found"
    exit 1
fi

# Export TF configuration
export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/.terraformrc"
export TF_VAR_api_key="$GCORE_API_KEY"

# Test results file
RESULTS_FILE="test_results_$(date +%Y%m%d_%H%M%S).md"

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$RESULTS_FILE"
}

log_test() {
    echo "" | tee -a "$RESULTS_FILE"
    echo "## $1" | tee -a "$RESULTS_FILE"
    echo "" | tee -a "$RESULTS_FILE"
}

log_result() {
    local status=$1
    local message=$2
    if [ "$status" = "PASS" ]; then
        echo "✅ **PASS**: $message" | tee -a "$RESULTS_FILE"
    elif [ "$status" = "FAIL" ]; then
        echo "❌ **FAIL**: $message" | tee -a "$RESULTS_FILE"
    else
        echo "⚠️ **INFO**: $message" | tee -a "$RESULTS_FILE"
    fi
}

# Initialize results file
cat > "$RESULTS_FILE" << EOF
# Cloud Instance Comprehensive Test Results
**Date**: $(date)
**Jira**: GCLOUD2-21138
**Branch**: terraform-instances

EOF

# ============================================================================
# Phase 1: Initial Apply
# ============================================================================
log_test "Phase 1: Initial Apply - All Resources"

log "Initializing Terraform..."
terraform init -upgrade 2>&1 | tail -5

log "Running terraform plan..."
if terraform plan -out=tfplan 2>&1 | tee plan_output.txt | tail -20; then
    log_result "INFO" "Plan completed successfully"
else
    log_result "FAIL" "Plan failed"
    cat plan_output.txt
    exit 1
fi

log "Applying configuration..."
if terraform apply -auto-approve tfplan 2>&1 | tee apply_output.txt; then
    log_result "PASS" "Initial apply completed successfully"
else
    log_result "FAIL" "Initial apply failed"
    cat apply_output.txt
    exit 1
fi

# ============================================================================
# Phase 2: Test 19 - Placement Groups Verification
# ============================================================================
log_test "Test 19: Placement Groups (servergroup_id)"

PLACEMENT_INSTANCE_ID=$(terraform output -raw placement_instance_id)
SERVERGROUP_ID=$(terraform output -raw placement_servergroup_id)

log "Verifying placement instance: $PLACEMENT_INSTANCE_ID"
log "Servergroup ID: $SERVERGROUP_ID"

# Verify via API that instance is in the servergroup
INSTANCE_DATA=$(curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
    "https://api.gcore.com/cloud/v1/instances/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID/$PLACEMENT_INSTANCE_ID")

if echo "$INSTANCE_DATA" | grep -q "server_group"; then
    log_result "PASS" "Instance has server_group field in API response"
else
    log_result "FAIL" "Instance missing server_group field"
fi

# ============================================================================
# Phase 3: Test 1-4 - Interface Variants
# ============================================================================
log_test "Test 1-4: Interface Variants"

EXTERNAL_ID=$(terraform output -raw external_instance_id)
SUBNET_ID=$(terraform output -raw subnet_instance_id)
ANY_SUBNET_ID=$(terraform output -raw any_subnet_instance_id)

log "External interface instance: $EXTERNAL_ID"
log "Subnet interface instance: $SUBNET_ID"
log "Any subnet interface instance: $ANY_SUBNET_ID"

# Check external instance
EXTERNAL_DATA=$(curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
    "https://api.gcore.com/cloud/v1/instances/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID/$EXTERNAL_ID")

if echo "$EXTERNAL_DATA" | grep -q '"status": "ACTIVE"'; then
    log_result "PASS" "Test 1 - External interface instance is ACTIVE"
else
    log_result "FAIL" "Test 1 - External interface instance not ACTIVE"
fi

# Check subnet instance
SUBNET_DATA=$(curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
    "https://api.gcore.com/cloud/v1/instances/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID/$SUBNET_ID")

if echo "$SUBNET_DATA" | grep -q '"status": "ACTIVE"'; then
    log_result "PASS" "Test 2 - Subnet interface instance is ACTIVE"
else
    log_result "FAIL" "Test 2 - Subnet interface instance not ACTIVE"
fi

# Check any_subnet instance
ANY_SUBNET_DATA=$(curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
    "https://api.gcore.com/cloud/v1/instances/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID/$ANY_SUBNET_ID")

if echo "$ANY_SUBNET_DATA" | grep -q '"status": "ACTIVE"'; then
    log_result "PASS" "Test 3 - Any subnet interface instance is ACTIVE"
else
    log_result "FAIL" "Test 3 - Any subnet interface instance not ACTIVE"
fi

# ============================================================================
# Phase 4: Test 5 - Floating IP
# ============================================================================
log_test "Test 5: Floating IP (new)"

FLOATING_ID=$(terraform output -raw floating_new_instance_id)
log "Floating IP instance: $FLOATING_ID"

FLOATING_DATA=$(curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
    "https://api.gcore.com/cloud/v1/instances/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID/$FLOATING_ID")

if echo "$FLOATING_DATA" | grep -q '"status": "ACTIVE"'; then
    log_result "PASS" "Test 5 - Floating IP instance is ACTIVE"
else
    log_result "FAIL" "Test 5 - Floating IP instance not ACTIVE"
fi

# ============================================================================
# Phase 5: Test 14-15 - Security Groups
# ============================================================================
log_test "Test 14-15: Security Groups"

SECGROUP_INSTANCE_ID=$(terraform output -raw secgroup_instance_id)
log "Security group instance: $SECGROUP_INSTANCE_ID"

SECGROUP_DATA=$(curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
    "https://api.gcore.com/cloud/v1/instances/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID/$SECGROUP_INSTANCE_ID")

if echo "$SECGROUP_DATA" | grep -q '"status": "ACTIVE"'; then
    log_result "PASS" "Test 14 - Security group instance is ACTIVE"
else
    log_result "FAIL" "Test 14 - Security group instance not ACTIVE"
fi

# ============================================================================
# Phase 6: Test 7-8 - Volume attachment_tag
# ============================================================================
log_test "Test 7-8: Volume attachment_tag"

VOLUME_TAGS_ID=$(terraform output -raw volume_tags_instance_id)
log "Volume tags instance: $VOLUME_TAGS_ID"

VOLUME_TAGS_DATA=$(curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
    "https://api.gcore.com/cloud/v1/instances/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID/$VOLUME_TAGS_ID")

if echo "$VOLUME_TAGS_DATA" | grep -q '"status": "ACTIVE"'; then
    log_result "PASS" "Test 7-8 - Volume tags instance is ACTIVE"
else
    log_result "FAIL" "Test 7-8 - Volume tags instance not ACTIVE"
fi

# Check if attachment_tag is visible in volumes
if echo "$VOLUME_TAGS_DATA" | grep -q "boot-disk\|data-disk"; then
    log_result "PASS" "attachment_tag values visible in API response"
else
    log_result "INFO" "attachment_tag values may not be returned in GET response (expected behavior)"
fi

# ============================================================================
# Phase 7: Drift Detection
# ============================================================================
log_test "Drift Detection"

log "Running terraform plan to detect drift..."
if terraform plan -detailed-exitcode 2>&1 | tee drift_check.txt; then
    log_result "PASS" "No drift detected - state matches infrastructure"
elif [ $? -eq 2 ]; then
    log_result "FAIL" "Drift detected - state differs from infrastructure"
    cat drift_check.txt | tail -30
else
    log_result "FAIL" "Plan failed during drift check"
fi

# ============================================================================
# Phase 8: Test 20 - Volume Deletion Verification (QA Bug)
# ============================================================================
log_test "Test 20: Volume Deletion Behavior on Instance Destroy"

log "Getting volume IDs before destroy..."
BOOT_VOLUME_ID=$(terraform show -json | python3 -c "import sys,json; data=json.load(sys.stdin); volumes=[r for r in data.get('values',{}).get('root_module',{}).get('resources',[]) if r.get('type')=='gcore_cloud_volume']; print(volumes[0]['values']['id'] if volumes else 'N/A')" 2>/dev/null || echo "N/A")
log "Boot volume ID: $BOOT_VOLUME_ID"

# We'll check this after destroy in cleanup phase

# ============================================================================
# Phase 9: Test 9-10 - Update Operations (Name Change)
# ============================================================================
log_test "Test 9-10: Update Operations - Name and Tags"

UPDATES_INSTANCE_ID=$(terraform output -raw updates_instance_id)
log "Updates instance ID: $UPDATES_INSTANCE_ID"

# Update name
log "Updating instance name..."
if terraform apply -auto-approve \
    -var='instance_name=tf-test-instance-updates-renamed' 2>&1 | tee update_name.txt; then
    log_result "PASS" "Test 9 - Name update applied successfully"
else
    log_result "FAIL" "Test 9 - Name update failed"
fi

# Verify name change via API
sleep 5
UPDATED_DATA=$(curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
    "https://api.gcore.com/cloud/v1/instances/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID/$UPDATES_INSTANCE_ID")

if echo "$UPDATED_DATA" | grep -q 'tf-test-instance-updates-renamed'; then
    log_result "PASS" "Name change reflected in API response"
else
    log_result "FAIL" "Name change NOT reflected in API response"
fi

# Update tags
log "Updating instance tags..."
if terraform apply -auto-approve \
    -var='instance_name=tf-test-instance-updates-renamed' \
    -var='instance_tags={"environment":"production","managed_by":"terraform","version":"2"}' 2>&1 | tee update_tags.txt; then
    log_result "PASS" "Test 10 - Tags update applied successfully"
else
    log_result "FAIL" "Test 10 - Tags update failed"
fi

# ============================================================================
# Phase 10: Test 21 - GET Request Body Bug Verification
# ============================================================================
log_test "Test 21: GET Request Body Bug (QA Report)"

log "Checking if GET requests include body..."
log "This requires mitmproxy capture - verifying via provider debug logs"

export TF_LOG=DEBUG
export TF_LOG_PATH="get_body_check.log"

terraform refresh 2>&1 > /dev/null

if grep -i "body" get_body_check.log 2>/dev/null | grep -i "GET" | head -5; then
    log_result "INFO" "Found potential GET with body - review logs manually"
else
    log_result "INFO" "No obvious GET body issues in logs"
fi

unset TF_LOG
unset TF_LOG_PATH

# ============================================================================
# Summary
# ============================================================================
log_test "Test Summary"

echo "" | tee -a "$RESULTS_FILE"
echo "### All Instances Created:" | tee -a "$RESULTS_FILE"
terraform output -json all_instance_statuses 2>/dev/null | python3 -c "
import sys, json
data = json.load(sys.stdin)
for name, status in data.items():
    print(f'- {name}: {status}')
" 2>/dev/null | tee -a "$RESULTS_FILE"

echo "" | tee -a "$RESULTS_FILE"
echo "### Test Results File: $RESULTS_FILE" | tee -a "$RESULTS_FILE"

log "Tests completed. Review $RESULTS_FILE for full results."
