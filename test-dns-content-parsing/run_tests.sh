#!/bin/bash
set -e

echo "=============================================="
echo "DNS Zone RRSet Content Parsing Tests"
echo "Testing NEW simpler syntax: content = [\"192.168.1.1\"]"
echo "=============================================="

# Load credentials
if [ -f /Users/user/repos/gcore-terraform/.env ]; then
    set -o allexport
    source /Users/user/repos/gcore-terraform/.env
    set +o allexport
    echo "✅ Credentials loaded"
else
    echo "❌ .env not found"
    exit 1
fi

export TF_CLI_CONFIG_FILE=/Users/user/repos/gcore-terraform/.terraformrc

# Build provider first
echo ""
echo "=== Building Provider ==="
cd /Users/user/repos/gcore-terraform
go build -o terraform-provider-gcore .
echo "✅ Provider built"

# Go to test directory
cd /Users/user/repos/gcore-terraform/test-dns-content-parsing
mkdir -p evidence

# =============================================================================
# Test 1: Create all records with NEW simpler syntax
# =============================================================================
echo ""
echo "=== Test 1: Create Records with NEW Syntax ==="
terraform apply -auto-approve 2>&1 | tee evidence/test1_create.log
echo "✅ Test 1: Records created"

# =============================================================================
# Test 2: Drift test - should show NO CHANGES
# =============================================================================
echo ""
echo "=== Test 2: Drift Test (should show no changes) ==="
terraform plan -detailed-exitcode 2>&1 | tee evidence/test2_drift.log
DRIFT_EXIT=$?
if [ $DRIFT_EXIT -eq 0 ]; then
    echo "✅ Test 2 PASS: No drift detected"
else
    echo "❌ Test 2 FAIL: Drift detected (exit code $DRIFT_EXIT)"
fi

# =============================================================================
# Test 3: Check state - content should be UNWRAPPED (no quotes)
# =============================================================================
echo ""
echo "=== Test 3: Check State Format ==="
echo "A record content in state:"
terraform show -json | jq '.values.root_module.resources[] | select(.type=="gcore_dns_zone_rrset" and .name=="test_a") | .values.resource_records[0].content'
echo ""
echo "MX record content in state:"
terraform show -json | jq '.values.root_module.resources[] | select(.type=="gcore_dns_zone_rrset" and .name=="test_mx") | .values.resource_records[0].content'

# =============================================================================
# Test 4: Update TTL (should update in-place, no drift)
# =============================================================================
echo ""
echo "=== Test 4: Update TTL (300 -> 600) ==="
terraform apply -auto-approve -var="ttl_a=600" 2>&1 | tee evidence/test4_update.log
terraform plan -detailed-exitcode -var="ttl_a=600" 2>&1 | tee evidence/test4_drift.log
DRIFT_EXIT=$?
if [ $DRIFT_EXIT -eq 0 ]; then
    echo "✅ Test 4 PASS: Update successful, no drift"
else
    echo "❌ Test 4 FAIL: Drift after update (exit code $DRIFT_EXIT)"
fi

# =============================================================================
# Test 5: Import test
# =============================================================================
echo ""
echo "=== Test 5: Import Test ==="
terraform state rm gcore_dns_zone_rrset.test_a 2>&1 | tee evidence/test5_state_rm.log
terraform import -var="ttl_a=600" gcore_dns_zone_rrset.test_a "maxima.lt/tf-content-test-a.maxima.lt/A" 2>&1 | tee evidence/test5_import.log
terraform plan -detailed-exitcode -var="ttl_a=600" 2>&1 | tee evidence/test5_drift.log
DRIFT_EXIT=$?
if [ $DRIFT_EXIT -eq 0 ]; then
    echo "✅ Test 5 PASS: Import successful, no drift"
else
    echo "❌ Test 5 FAIL: Drift after import (exit code $DRIFT_EXIT)"
fi

# =============================================================================
# Test 6: Cleanup
# =============================================================================
echo ""
echo "=== Test 6: Cleanup ==="
terraform destroy -auto-approve 2>&1 | tee evidence/test6_destroy.log
echo "✅ Test 6: Cleanup complete"

echo ""
echo "=============================================="
echo "Testing Complete!"
echo "Evidence saved to: evidence/"
echo "=============================================="
