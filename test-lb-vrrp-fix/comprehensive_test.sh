#!/bin/bash
set -e

# Comprehensive test for all LB resources with mitmproxy
echo "========================================="
echo "Comprehensive LB Test with mitmproxy"
echo "========================================="

# Setup
cd /Users/user/repos/gcore-terraform
set -o allexport
source .env
set +o allexport
cd test-lb-vrrp-fix

export TF_CLI_CONFIG_FILE=/Users/user/repos/gcore-terraform/.terraformrc

# Setup proxy for mitmproxy
export HTTP_PROXY="http://127.0.0.1:9092"
export HTTPS_PROXY="http://127.0.0.1:9092"
export NO_PROXY="registry.terraform.io,releases.hashicorp.com"

# SSL cert bundle (if exists)
if [ -f ../ca-bundle.pem ]; then
    export SSL_CERT_FILE="../ca-bundle.pem"
    export REQUESTS_CA_BUNDLE="../ca-bundle.pem"
fi

echo ""
echo "=== Test 1: Create all resources ==="
terraform apply -auto-approve

LB_ID=$(terraform output -raw lb_id)
LISTENER_ID=$(terraform output -raw listener_id)
POOL_ID=$(terraform output -raw pool_id)
VRRP_COUNT=$(terraform output -raw vrrp_ips_count)

echo "LB ID: $LB_ID"
echo "Listener ID: $LISTENER_ID"
echo "Pool ID: $POOL_ID"
echo "vrrp_ips count: $VRRP_COUNT"
echo ""

echo "=== Test 2: Drift check after creation ==="
if terraform plan -detailed-exitcode; then
    echo "✅ PASS: No drift after creation"
else
    EXIT_CODE=$?
    if [ $EXIT_CODE -eq 2 ]; then
        echo "❌ FAIL: Drift detected after creation"
        exit 1
    fi
fi
echo ""

echo "=== Test 3: Update LB name (PATCH + GET pattern) ==="
terraform apply -auto-approve -var="lb_name=test-lb-renamed"

LB_ID_AFTER=$(terraform output -raw lb_id)
VRRP_COUNT_AFTER=$(terraform output -raw vrrp_ips_count)

if [ "$LB_ID" = "$LB_ID_AFTER" ]; then
    echo "✅ PASS: LB updated in-place"
else
    echo "❌ FAIL: LB was recreated"
    exit 1
fi

if [ "$VRRP_COUNT" = "$VRRP_COUNT_AFTER" ]; then
    echo "✅ PASS: vrrp_ips preserved ($VRRP_COUNT elements)"
else
    echo "❌ FAIL: vrrp_ips changed"
    exit 1
fi
echo ""

echo "=== Test 4: Drift check after LB update ==="
if terraform plan -detailed-exitcode -var="lb_name=test-lb-renamed"; then
    echo "✅ PASS: No drift after LB update"
else
    EXIT_CODE=$?
    if [ $EXIT_CODE -eq 2 ]; then
        echo "❌ FAIL: Drift detected after LB update"
        terraform plan -var="lb_name=test-lb-renamed"
        exit 1
    fi
fi
echo ""

echo "=== Test 5: Update Listener (UpdateAndPoll pattern) ==="
terraform apply -auto-approve \
    -var="lb_name=test-lb-renamed" \
    -var="listener_name=test-listener-renamed" \
    -var="timeout_client_data=60000"

LISTENER_ID_AFTER=$(terraform output -raw listener_id)
if [ "$LISTENER_ID" = "$LISTENER_ID_AFTER" ]; then
    echo "✅ PASS: Listener updated in-place"
else
    echo "❌ FAIL: Listener was recreated"
    exit 1
fi
echo ""

echo "=== Test 6: Drift check after Listener update ==="
if terraform plan -detailed-exitcode \
    -var="lb_name=test-lb-renamed" \
    -var="listener_name=test-listener-renamed" \
    -var="timeout_client_data=60000"; then
    echo "✅ PASS: No drift after Listener update"
else
    EXIT_CODE=$?
    if [ $EXIT_CODE -eq 2 ]; then
        echo "❌ FAIL: Drift detected after Listener update"
        exit 1
    fi
fi
echo ""

echo "=== Test 7: Update Pool (UpdateAndPoll pattern) ==="
terraform apply -auto-approve \
    -var="lb_name=test-lb-renamed" \
    -var="listener_name=test-listener-renamed" \
    -var="timeout_client_data=60000" \
    -var="pool_name=test-pool-renamed" \
    -var="health_check_path=/status"

POOL_ID_AFTER=$(terraform output -raw pool_id)
if [ "$POOL_ID" = "$POOL_ID_AFTER" ]; then
    echo "✅ PASS: Pool updated in-place"
else
    echo "❌ FAIL: Pool was recreated"
    exit 1
fi
echo ""

echo "=== Test 8: Drift check after Pool update ==="
if terraform plan -detailed-exitcode \
    -var="lb_name=test-lb-renamed" \
    -var="listener_name=test-listener-renamed" \
    -var="timeout_client_data=60000" \
    -var="pool_name=test-pool-renamed" \
    -var="health_check_path=/status"; then
    echo "✅ PASS: No drift after Pool update"
else
    EXIT_CODE=$?
    if [ $EXIT_CODE -eq 2 ]; then
        echo "❌ FAIL: Drift detected after Pool update"
        exit 1
    fi
fi
echo ""

echo "=== Test 9: Add Pool Member (AddAndPoll pattern) ==="
terraform apply -auto-approve \
    -var="lb_name=test-lb-renamed" \
    -var="listener_name=test-listener-renamed" \
    -var="timeout_client_data=60000" \
    -var="pool_name=test-pool-renamed" \
    -var="health_check_path=/status" \
    -var="with_member=true"

MEMBER_ID=$(terraform output -raw member_id)
echo "Member ID: $MEMBER_ID"
echo ""

echo "=== Test 10: Drift check after Member creation ==="
if terraform plan -detailed-exitcode \
    -var="lb_name=test-lb-renamed" \
    -var="listener_name=test-listener-renamed" \
    -var="timeout_client_data=60000" \
    -var="pool_name=test-pool-renamed" \
    -var="health_check_path=/status" \
    -var="with_member=true"; then
    echo "✅ PASS: No drift after Member creation"
else
    EXIT_CODE=$?
    if [ $EXIT_CODE -eq 2 ]; then
        echo "❌ FAIL: Drift detected after Member creation"
        exit 1
    fi
fi
echo ""

echo "=== Test 11: Update Pool Member weight (UpdateAndPoll pattern) ==="
terraform apply -auto-approve \
    -var="lb_name=test-lb-renamed" \
    -var="listener_name=test-listener-renamed" \
    -var="timeout_client_data=60000" \
    -var="pool_name=test-pool-renamed" \
    -var="health_check_path=/status" \
    -var="with_member=true" \
    -var="member_weight=5"

MEMBER_ID_AFTER=$(terraform output -raw member_id)
if [ "$MEMBER_ID" = "$MEMBER_ID_AFTER" ]; then
    echo "✅ PASS: Member updated in-place"
else
    echo "❌ FAIL: Member was recreated"
    exit 1
fi
echo ""

echo "=== Test 12: Final drift check ==="
if terraform plan -detailed-exitcode \
    -var="lb_name=test-lb-renamed" \
    -var="listener_name=test-listener-renamed" \
    -var="timeout_client_data=60000" \
    -var="pool_name=test-pool-renamed" \
    -var="health_check_path=/status" \
    -var="with_member=true" \
    -var="member_weight=5"; then
    echo "✅ PASS: No drift after all updates"
else
    EXIT_CODE=$?
    if [ $EXIT_CODE -eq 2 ]; then
        echo "❌ FAIL: Final drift detected"
        exit 1
    fi
fi
echo ""

echo "========================================="
echo "All Tests Passed! ✅"
echo "========================================="
echo ""
echo "Summary:"
echo "- LB: vrrp_ips preserved during PATCH + GET"
echo "- Listener: UpdateAndPoll successful"
echo "- Pool: UpdateAndPoll successful"
echo "- Member: AddAndPoll/UpdateAndPoll successful"
echo "- No drift detected after any operation"
echo ""
echo "Check mitmproxy captures with:"
echo "  mitmdump -r flow.mitm -n --set flow_detail=2"
