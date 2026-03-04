#!/bin/bash
set -e

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

# Configure environment
export TF_CLI_CONFIG_FILE=/Users/user/repos/gcore-terraform/.terraformrc
unset HTTP_PROXY HTTPS_PROXY

# Initialize test results
RESULTS_FILE="test_results.md"
echo "# Router Comprehensive Test Results" > $RESULTS_FILE
echo "Date: $(date)" >> $RESULTS_FILE
echo "" >> $RESULTS_FILE

# Helper function to run test
run_test() {
    local test_num=$1
    local test_name=$2
    shift 2
    local tf_vars="$@"

    echo "======================================================================="
    echo "TEST $test_num: $test_name"
    echo "======================================================================="
    echo "Variables: $tf_vars"

    # Apply changes
    echo "Applying..."
    if terraform apply -auto-approve $tf_vars > apply_${test_num}.log 2>&1; then
        echo "✅ Apply succeeded"
    else
        echo "❌ Apply failed"
        tail -20 apply_${test_num}.log
        echo "- ❌ TEST $test_num: $test_name - FAIL (apply error)" >> $RESULTS_FILE
        return 1
    fi

    # Drift check
    echo "Checking for drift..."
    if terraform plan -detailed-exitcode > plan_${test_num}.log 2>&1; then
        echo "✅ No drift detected"
        echo "- ✅ TEST $test_num: $test_name - PASS (no drift)" >> $RESULTS_FILE
    else
        exitcode=$?
        if [ $exitcode -eq 2 ]; then
            echo "❌ Drift detected!"
            head -50 plan_${test_num}.log
            echo "- ❌ TEST $test_num: $test_name - FAIL (drift detected)" >> $RESULTS_FILE
            return 1
        else
            echo "❌ Plan error!"
            tail -20 plan_${test_num}.log
            echo "- ❌ TEST $test_num: $test_name - FAIL (plan error)" >> $RESULTS_FILE
            return 1
        fi
    fi

    # Save state snapshot
    terraform show -json > evidence/state_test_${test_num}.json 2>/dev/null || true

    echo ""
    return 0
}

echo "===================================================================="
echo "COMPREHENSIVE ROUTER TESTING - 18 Test Scenarios"
echo "===================================================================="
echo ""

# Get subnet IDs after initial creation
SUBNET1_ID=$(terraform output -raw subnet1_id)
SUBNET2_ID=$(terraform output -raw subnet2_id)
echo "Subnet1 ID: $SUBNET1_ID"
echo "Subnet2 ID: $SUBNET2_ID"
echo ""

#############################################################################
# TEST 1: Drift - Minimal router (already created, just check drift)
#############################################################################
run_test 1 "Drift Detection - Minimal Router"

#############################################################################
# TEST 2: Routes - Add single route
#############################################################################
run_test 2 "Add Single Route" -var='routes=[{destination="10.0.0.0/24",nexthop="192.168.1.1"}]'

#############################################################################
# TEST 3: Routes - Update existing route
#############################################################################
run_test 3 "Update Route" -var='routes=[{destination="172.16.0.0/16",nexthop="192.168.1.254"}]'

#############################################################################
# TEST 4: Routes - Remove all routes (CRITICAL - check empty array)
#############################################################################
run_test 4 "Remove All Routes" -var='routes=[]'

#############################################################################
# TEST 5: Routes - Add multiple routes
#############################################################################
run_test 5 "Add Multiple Routes" -var='routes=[{destination="10.1.0.0/24",nexthop="192.168.1.1"},{destination="10.2.0.0/24",nexthop="192.168.1.2"},{destination="10.3.0.0/24",nexthop="192.168.1.3"}]'

# Clear routes for interface tests
echo "Clearing routes..."
terraform apply -auto-approve -var='routes=[]' > /dev/null 2>&1

#############################################################################
# TEST 6: Interface - Add single interface
#############################################################################
run_test 6 "Add Single Interface" -var="interfaces=[\"$SUBNET1_ID\"]"

#############################################################################
# TEST 7: Interface - Remove interface
#############################################################################
run_test 7 "Remove Interface" -var='interfaces=[]'

#############################################################################
# TEST 8: Interface - Add multiple interfaces
#############################################################################
run_test 8 "Add Multiple Interfaces" -var="interfaces=[\"$SUBNET1_ID\",\"$SUBNET2_ID\"]"

#############################################################################
# TEST 9: Interface - Replace interface
#############################################################################
run_test 9 "Replace Interface" -var="interfaces=[\"$SUBNET2_ID\"]"

# Clear interfaces for mixed tests
echo "Clearing interfaces..."
terraform apply -auto-approve -var='interfaces=[]' > /dev/null 2>&1

#############################################################################
# TEST 10: Mixed - Add route + add interface (CRITICAL - check order)
#############################################################################
run_test 10 "Add Route + Interface Together" -var="interfaces=[\"$SUBNET1_ID\"]" -var='routes=[{destination="10.10.0.0/24",nexthop="192.168.1.1"}]'

#############################################################################
# TEST 11: Mixed - Remove route + remove interface (CRITICAL - check order)
#############################################################################
run_test 11 "Remove Route + Interface Together" -var='interfaces=[]' -var='routes=[]'

#############################################################################
# TEST 12: Mixed - Add interface + remove route
#############################################################################
# First set up state with route but no interface
echo "Setting up for TEST 12..."
terraform apply -auto-approve -var='routes=[{destination="10.20.0.0/24",nexthop="192.168.1.1"}]' > /dev/null 2>&1
run_test 12 "Add Interface + Remove Route" -var="interfaces=[\"$SUBNET1_ID\"]" -var='routes=[]'

# Clear for next tests
echo "Clearing interfaces..."
terraform apply -auto-approve -var='interfaces=[]' > /dev/null 2>&1

#############################################################################
# TEST 13: Name - Update name only
#############################################################################
run_test 13 "Update Router Name" -var='router_name=test-router-renamed'

#############################################################################
# TEST 14: ExtGateway - Set external gateway info
#############################################################################
run_test 14 "Set External Gateway" -var='enable_external_gateway=true' -var='external_gateway_snat=true'

#############################################################################
# TEST 15: ExtGateway - Update external gateway info
#############################################################################
run_test 15 "Update External Gateway SNAT" -var='enable_external_gateway=true' -var='external_gateway_snat=false'

#############################################################################
# TEST 16: ExtGateway - Remove external gateway info
#############################################################################
run_test 16 "Remove External Gateway" -var='enable_external_gateway=false'

#############################################################################
# TEST 17: Interface Change Only
#############################################################################
run_test 17 "Interface Change Only" -var="interfaces=[\"$SUBNET1_ID\"]"

#############################################################################
# TEST 18: Drift - Complex config with all features
#############################################################################
run_test 18 "Complex Config - Full Feature Drift Check" -var='router_name=test-router-complex' -var="interfaces=[\"$SUBNET1_ID\",\"$SUBNET2_ID\"]" -var='routes=[{destination="10.99.0.0/24",nexthop="192.168.1.99"}]' -var='enable_external_gateway=true'

echo ""
echo "===================================================================="
echo "ALL TESTS COMPLETED"
echo "===================================================================="
echo ""
cat $RESULTS_FILE
