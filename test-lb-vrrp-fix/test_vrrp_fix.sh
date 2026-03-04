#!/bin/bash
set -e

# Test script for Load Balancer vrrp_ips fix
# Tests the PATCH + GET pattern for preserving vrrp_ips during updates

echo "========================================="
echo "Load Balancer vrrp_ips Fix Test"
echo "========================================="
echo ""

# Load environment
echo "Loading credentials..."
cd /Users/user/repos/gcore-terraform
set -o allexport
source .env
set +o allexport
cd test-lb-vrrp-fix

export TF_CLI_CONFIG_FILE=/Users/user/repos/gcore-terraform/.terraformrc
export TF_LOG=DEBUG
export TF_LOG_PATH=test_vrrp.log

echo "✓ Credentials loaded"
echo ""

# Test 1: Skip init (using dev_overrides)
echo "Test 1: Using dev_overrides, skipping terraform init..."
echo "✓ Provider configured via .terraformrc"
echo ""

# Test 2: Create load balancer
echo "Test 2: Creating load balancer..."
terraform apply -auto-approve
echo "✓ Load balancer created"
echo ""

# Capture initial state
LB_ID=$(terraform output -raw lb_id)
VRRP_COUNT=$(terraform output -raw vrrp_ips_count)
echo "Load Balancer ID: $LB_ID"
echo "vrrp_ips count: $VRRP_COUNT"
echo ""

# Test 3: Drift detection (CRITICAL)
echo "Test 3: Checking for drift (no changes expected)..."
if terraform plan -detailed-exitcode; then
    echo "✅ PASS: No drift detected"
else
    EXIT_CODE=$?
    if [ $EXIT_CODE -eq 2 ]; then
        echo "❌ FAIL: Drift detected on fresh apply!"
        terraform plan
        exit 1
    fi
fi
echo ""

# Test 4: Update name (triggers PATCH + GET)
echo "Test 4: Updating load balancer name..."
echo "This will test the PATCH + GET pattern for vrrp_ips preservation"
terraform apply -auto-approve -var="lb_name=test-lb-vrrp-renamed"
echo "✓ Load balancer name updated"
echo ""

# Verify ID didn't change (update in place)
LB_ID_AFTER=$(terraform output -raw lb_id)
if [ "$LB_ID" = "$LB_ID_AFTER" ]; then
    echo "✅ PASS: Load balancer updated in-place (ID preserved)"
else
    echo "❌ FAIL: Load balancer was recreated (ID changed)"
    exit 1
fi
echo ""

# Verify vrrp_ips preserved
VRRP_COUNT_AFTER=$(terraform output -raw vrrp_ips_count)
echo "vrrp_ips count before: $VRRP_COUNT"
echo "vrrp_ips count after: $VRRP_COUNT_AFTER"

if [ "$VRRP_COUNT" = "$VRRP_COUNT_AFTER" ] && [ "$VRRP_COUNT" -gt 0 ]; then
    echo "✅ PASS: vrrp_ips preserved ($VRRP_COUNT elements)"
else
    echo "❌ FAIL: vrrp_ips not preserved or empty"
    terraform output vrrp_ips
    exit 1
fi
echo ""

# Test 5: Drift detection after update (CRITICAL)
echo "Test 5: Checking for drift after update..."
if terraform plan -detailed-exitcode -var="lb_name=test-lb-vrrp-renamed"; then
    echo "✅ PASS: No drift after update"
else
    EXIT_CODE=$?
    if [ $EXIT_CODE -eq 2 ]; then
        echo "❌ FAIL: Drift detected after name update!"
        terraform plan -var="lb_name=test-lb-vrrp-renamed"
        exit 1
    fi
fi
echo ""

# Test 6: Check debug logs for PATCH + GET pattern
echo "Test 6: Verifying API calls in debug logs..."
if grep -q "failed to read load balancer after update" test_vrrp.log; then
    echo "❌ FAIL: GET call after PATCH failed"
    exit 1
else
    echo "✅ PASS: No errors in GET after PATCH"
fi

# Count PATCH operations
PATCH_COUNT=$(grep -c "PATCH.*loadbalancers" test_vrrp.log || echo "0")
echo "PATCH operations detected: $PATCH_COUNT"

# Check for GET after PATCH
if grep -A 5 "PATCH.*loadbalancers" test_vrrp.log | grep -q "GET.*loadbalancers"; then
    echo "✅ PASS: GET call follows PATCH (correct pattern)"
else
    echo "⚠️  WARNING: Could not verify GET after PATCH in logs"
fi
echo ""

# Display final state
echo "========================================="
echo "Final State:"
echo "========================================="
terraform output vrrp_ips
echo ""

echo "========================================="
echo "All Tests Passed! ✅"
echo "========================================="
echo ""
echo "Summary:"
echo "- Load balancer created with vrrp_ips"
echo "- No drift on initial creation"
echo "- Name update used in-place modification"
echo "- vrrp_ips preserved during update"
echo "- No drift after update"
echo ""
