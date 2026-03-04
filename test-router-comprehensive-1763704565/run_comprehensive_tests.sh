#!/bin/bash
set -e

# Configure environment
export TF_CLI_CONFIG_FILE=/Users/user/repos/gcore-terraform/.terraformrc
export HTTP_PROXY="http://127.0.0.1:9092"
export HTTPS_PROXY="http://127.0.0.1:9092"
export NO_PROXY="registry.terraform.io,releases.hashicorp.com"
export SSL_CERT_FILE="/Users/user/repos/gcore-terraform/ca-bundle.pem"
export REQUESTS_CA_BUNDLE="/Users/user/repos/gcore-terraform/ca-bundle.pem"

# Load credentials
set -o allexport
source /Users/user/repos/gcore-terraform/.env
set +o allexport

# Initialize test results
RESULTS_FILE="test_results.md"
echo "# Router Comprehensive Test Results" > $RESULTS_FILE
echo "Date: $(date)" >> $RESULTS_FILE
echo "" >> $RESULTS_FILE

# Helper function to run test
run_test() {
    local test_num=$1
    local test_name=$2
    local tf_vars=$3
    local check_cmd=$4

    echo "======================================================================="
    echo "TEST $test_num: $test_name"
    echo "======================================================================="
    echo "Variables: $tf_vars"

    # Apply changes
    echo "Applying..."
    terraform apply -auto-approve $tf_vars > /dev/null 2>&1

    # Run custom check if provided
    if [ -n "$check_cmd" ]; then
        eval "$check_cmd"
    fi

    # Drift check
    echo "Checking for drift..."
    if terraform plan -detailed-exitcode > /dev/null 2>&1; then
        echo "✅ No drift detected"
        echo "- ✅ TEST $test_num: $test_name - PASS (no drift)" >> $RESULTS_FILE
    else
        echo "❌ Drift detected!"
        terraform plan | head -20
        echo "- ❌ TEST $test_num: $test_name - FAIL (drift detected)" >> $RESULTS_FILE
    fi

    # Save MITM capture for this test
    cp ../flow.mitm evidence/flow_test_${test_num}.mitm 2>/dev/null || true

    echo ""
}

# Save state snapshot helper
save_state() {
    local test_num=$1
    terraform show -json > evidence/state_test_${test_num}.json 2>/dev/null || true
}

echo "==================================================================="
echo "COMPREHENSIVE ROUTER TESTING - 18 Test Scenarios"
echo "==================================================================="
echo ""

# Get subnet IDs after initial creation
SUBNET1_ID=""
SUBNET2_ID=""

#############################################################################
# TEST 1: Drift - Minimal router
#############################################################################
run_test 1 "Drift Detection - Minimal Router" "-var='router_name=test-router-minimal'"
ROUTER_ID=$(terraform output -raw router_id)
SUBNET1_ID=$(terraform output -raw subnet1_id)
SUBNET2_ID=$(terraform output -raw subnet2_id)
echo "Router ID: $ROUTER_ID"
echo "Subnet1 ID: $SUBNET1_ID"
echo "Subnet2 ID: $SUBNET2_ID"
save_state 1

#############################################################################
# TEST 2: Routes - Add single route
#############################################################################
run_test 2 "Add Single Route" \
    "-var='routes=[{destination=\"10.0.0.0/24\",nexthop=\"192.168.1.1\"}]'"
save_state 2

#############################################################################
# TEST 3: Routes - Update existing route
#############################################################################
run_test 3 "Update Route" \
    "-var='routes=[{destination=\"172.16.0.0/16\",nexthop=\"192.168.1.254\"}]'"
save_state 3

#############################################################################
# TEST 4: Routes - Remove all routes (CRITICAL - check empty array)
#############################################################################
run_test 4 "Remove All Routes" "-var='routes=[]'"
save_state 4

#############################################################################
# TEST 5: Routes - Add multiple routes
#############################################################################
run_test 5 "Add Multiple Routes" \
    "-var='routes=[{destination=\"10.1.0.0/24\",nexthop=\"192.168.1.1\"},{destination=\"10.2.0.0/24\",nexthop=\"192.168.1.2\"},{destination=\"10.3.0.0/24\",nexthop=\"192.168.1.3\"}]'"
save_state 5

# Clear routes for interface tests
terraform apply -auto-approve -var='routes=[]' > /dev/null 2>&1

#############################################################################
# TEST 6: Interface - Add single interface
#############################################################################
run_test 6 "Add Single Interface" \
    "-var='interfaces=[\"$SUBNET1_ID\"]'"
save_state 6

#############################################################################
# TEST 7: Interface - Remove interface
#############################################################################
run_test 7 "Remove Interface" "-var='interfaces=[]'"
save_state 7

#############################################################################
# TEST 8: Interface - Add multiple interfaces
#############################################################################
run_test 8 "Add Multiple Interfaces" \
    "-var='interfaces=[\"$SUBNET1_ID\",\"$SUBNET2_ID\"]'"
save_state 8

#############################################################################
# TEST 9: Interface - Replace interface
#############################################################################
run_test 9 "Replace Interface" \
    "-var='interfaces=[\"$SUBNET2_ID\"]'"
save_state 9

# Clear interfaces for mixed tests
terraform apply -auto-approve -var='interfaces=[]' > /dev/null 2>&1

#############################################################################
# TEST 10: Mixed - Add route + add interface (CRITICAL - check order)
#############################################################################
run_test 10 "Add Route + Interface Together" \
    "-var='interfaces=[\"$SUBNET1_ID\"]' -var='routes=[{destination=\"10.10.0.0/24\",nexthop=\"192.168.1.1\"}]'"
save_state 10

#############################################################################
# TEST 11: Mixed - Remove route + remove interface (CRITICAL - check order)
#############################################################################
run_test 11 "Remove Route + Interface Together" \
    "-var='interfaces=[]' -var='routes=[]'"
save_state 11

#############################################################################
# TEST 12: Mixed - Add interface + remove route
#############################################################################
# First set up state with route but no interface
terraform apply -auto-approve -var='routes=[{destination=\"10.20.0.0/24\",nexthop=\"192.168.1.1\"}]' > /dev/null 2>&1
run_test 12 "Add Interface + Remove Route" \
    "-var='interfaces=[\"$SUBNET1_ID\"]' -var='routes=[]'"
save_state 12

# Clear for next tests
terraform apply -auto-approve -var='interfaces=[]' > /dev/null 2>&1

#############################################################################
# TEST 13: Name - Update name only
#############################################################################
run_test 13 "Update Router Name" "-var='router_name=test-router-renamed'"
save_state 13

#############################################################################
# TEST 14: ExtGateway - Set external gateway info
#############################################################################
run_test 14 "Set External Gateway" \
    "-var='enable_external_gateway=true' -var='external_gateway_snat=true'"
save_state 14

#############################################################################
# TEST 15: ExtGateway - Update external gateway info
#############################################################################
run_test 15 "Update External Gateway SNAT" \
    "-var='enable_external_gateway=true' -var='external_gateway_snat=false'"
save_state 15

#############################################################################
# TEST 16: ExtGateway - Remove external gateway info
#############################################################################
run_test 16 "Remove External Gateway" \
    "-var='enable_external_gateway=false'"
save_state 16

#############################################################################
# TEST 17: Empty PATCH Check - Only change interfaces (no PATCH should be sent)
#############################################################################
run_test 17 "Interface Change Only (No PATCH)" \
    "-var='interfaces=[\"$SUBNET1_ID\"]'"
save_state 17

#############################################################################
# TEST 18: Drift - Complex config with all features
#############################################################################
run_test 18 "Complex Config - Full Feature Drift Check" \
    "-var='router_name=test-router-complex' -var='interfaces=[\"$SUBNET1_ID\",\"$SUBNET2_ID\"]' -var='routes=[{destination=\"10.99.0.0/24\",nexthop=\"192.168.1.99\"}]' -var='enable_external_gateway=true'"
save_state 18

echo ""
echo "==================================================================="
echo "ALL TESTS COMPLETED"
echo "==================================================================="
echo ""
cat $RESULTS_FILE
