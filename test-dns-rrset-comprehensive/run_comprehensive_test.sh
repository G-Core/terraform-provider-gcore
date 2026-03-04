#!/bin/bash
set -e

echo "=== DNS Zone RRSet Comprehensive Test Suite ==="
echo "Date: $(date)"
echo ""

# Load credentials
if [ -f /Users/user/repos/gcore-terraform/.env ]; then
    set -o allexport
    source /Users/user/repos/gcore-terraform/.env
    set +o allexport
    echo "Credentials loaded"
else
    echo "ERROR: .env file not found"
    exit 1
fi

# Set terraform config
export TF_CLI_CONFIG_FILE=/Users/user/repos/gcore-terraform/.terraformrc
export TF_LOG=DEBUG
export TF_LOG_PATH=terraform_test.log

# Initialize terraform
echo ""
echo "=== Initializing Terraform ==="
terraform init -upgrade

# Clean up any existing record first
echo ""
echo "=== Cleanup: Removing any existing test record ==="
API_KEY="$GCORE_API_KEY"
curl -s -X DELETE -H "Authorization: APIKey $API_KEY" \
    "https://api.gcore.com/dns/v2/zones/maxima.lt/tf-test-fixed.maxima.lt/A" || true

echo ""
echo "=== TEST 1: Create A Record ==="
terraform apply -auto-approve 2>&1 | tee test1_apply.log

echo ""
echo "=== TEST 1: Drift Check (should show no changes) ==="
terraform plan -detailed-exitcode 2>&1 | tee test1_drift.log
DRIFT_EXIT=$?
if [ $DRIFT_EXIT -eq 0 ]; then
    echo "DRIFT TEST 1 PASSED: No changes detected"
elif [ $DRIFT_EXIT -eq 2 ]; then
    echo "DRIFT TEST 1 FAILED: Changes detected!"
else
    echo "DRIFT TEST 1 ERROR: Exit code $DRIFT_EXIT"
fi

echo ""
echo "=== TEST 2: Update TTL (300 -> 600) ==="
terraform apply -auto-approve -var="ttl=600" 2>&1 | tee test2_apply.log

echo ""
echo "=== TEST 2: Verify in-place update (check plan output for 'Update' not 'Replace') ==="
# Check if resource was replaced or updated
if grep -q "must be replaced" test2_apply.log; then
    echo "TEST 2 FAILED: Resource was replaced instead of updated"
else
    echo "TEST 2 PASSED: Resource updated in-place"
fi

echo ""
echo "=== TEST 2: Drift Check after TTL update ==="
terraform plan -detailed-exitcode 2>&1 | tee test2_drift.log
DRIFT_EXIT=$?
if [ $DRIFT_EXIT -eq 0 ]; then
    echo "DRIFT TEST 2 PASSED: No changes detected"
elif [ $DRIFT_EXIT -eq 2 ]; then
    echo "DRIFT TEST 2 FAILED: Changes detected!"
else
    echo "DRIFT TEST 2 ERROR: Exit code $DRIFT_EXIT"
fi

echo ""
echo "=== TEST 3: Update resource_records (add second IP) ==="
terraform apply -auto-approve -var="ttl=600" \
    -var='resource_records=[{content=["\"192.168.1.100\""],enabled=true},{content=["\"192.168.1.101\""],enabled=true}]' 2>&1 | tee test3_apply.log

echo ""
echo "=== TEST 3: Drift Check after resource_records update ==="
terraform plan -detailed-exitcode \
    -var="ttl=600" \
    -var='resource_records=[{content=["\"192.168.1.100\""],enabled=true},{content=["\"192.168.1.101\""],enabled=true}]' 2>&1 | tee test3_drift.log
DRIFT_EXIT=$?
if [ $DRIFT_EXIT -eq 0 ]; then
    echo "DRIFT TEST 3 PASSED: No changes detected"
elif [ $DRIFT_EXIT -eq 2 ]; then
    echo "DRIFT TEST 3 FAILED: Changes detected!"
else
    echo "DRIFT TEST 3 ERROR: Exit code $DRIFT_EXIT"
fi

echo ""
echo "=== TEST 4: Delete resource ==="
terraform destroy -auto-approve 2>&1 | tee test4_destroy.log

echo ""
echo "=== Verify deletion via API ==="
RESPONSE=$(curl -s -H "Authorization: APIKey $API_KEY" \
    "https://api.gcore.com/dns/v2/zones/maxima.lt/tf-test-fixed.maxima.lt/A")
if echo "$RESPONSE" | grep -q "not found\|404"; then
    echo "DELETE TEST PASSED: Record no longer exists"
else
    echo "DELETE TEST WARNING: Record may still exist"
    echo "$RESPONSE"
fi

echo ""
echo "=== Test Suite Complete ==="
