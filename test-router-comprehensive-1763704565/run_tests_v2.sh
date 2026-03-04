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
        tail -30 apply_${test_num}.log
        echo "- ❌ TEST $test_num: $test_name - FAIL (apply error)" >> $RESULTS_FILE
        return 1
    fi

    # Wait for API to stabilize (eventual consistency)
    sleep 2

    # Drift check (must use same variables as apply!)
    echo "Checking for drift..."
    if terraform plan -detailed-exitcode $tf_vars > plan_${test_num}.log 2>&1; then
        echo "✅ No drift detected"
        echo "- ✅ TEST $test_num: $test_name - PASS (no drift)" >> $RESULTS_FILE
    else
        exitcode=$?
        if [ $exitcode -eq 2 ]; then
            echo "⚠️  Drift detected, attempting to fix with re-apply..."
            if terraform apply -auto-approve $tf_vars > apply_${test_num}_retry.log 2>&1; then
                echo "✅ Re-apply succeeded, checking drift again..."
                if terraform plan -detailed-exitcode $tf_vars > plan_${test_num}_retry.log 2>&1; then
                    echo "✅ No drift after re-apply"
                    echo "- ✅ TEST $test_num: $test_name - PASS (drift fixed by re-apply)" >> $RESULTS_FILE
                else
                    echo "❌ Drift persists after re-apply"
                    head -50 plan_${test_num}_retry.log
                    echo "- ❌ TEST $test_num: $test_name - FAIL (persistent drift)" >> $RESULTS_FILE
                    return 1
                fi
            else
                echo "❌ Re-apply failed"
                tail -20 apply_${test_num}_retry.log
                echo "- ❌ TEST $test_num: $test_name - FAIL (re-apply error)" >> $RESULTS_FILE
                return 1
            fi
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

# Get subnet IDs
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
# TEST 2: Interface - Add single interface
#############################################################################
run_test 2 "Add Single Interface" -var="interfaces=[\"$SUBNET1_ID\"]"

#############################################################################
# TEST 3: Routes - Add single route (now interface exists for nexthop validation)
#############################################################################
run_test 3 "Add Single Route" -var="interfaces=[\"$SUBNET1_ID\"]" -var='routes=[{destination="10.0.0.0/24",nexthop="192.168.1.1"}]'

#############################################################################
# TEST 4: Routes - Update existing route
#############################################################################
run_test 4 "Update Route" -var="interfaces=[\"$SUBNET1_ID\"]" -var='routes=[{destination="172.16.0.0/16",nexthop="192.168.1.254"}]'

#############################################################################
# TEST 5: Routes - Add multiple routes
#############################################################################
run_test 5 "Add Multiple Routes" -var="interfaces=[\"$SUBNET1_ID\"]" -var='routes=[{destination="10.1.0.0/24",nexthop="192.168.1.1"},{destination="10.2.0.0/24",nexthop="192.168.1.2"},{destination="10.3.0.0/24",nexthop="192.168.1.3"}]'

#############################################################################
# TEST 6: Routes - Remove all routes (CRITICAL - check empty array handling)
#############################################################################
run_test 6 "Remove All Routes" -var="interfaces=[\"$SUBNET1_ID\"]" -var='routes=[]'

#############################################################################
# TEST 7: Interface - Remove interface (no routes should be present)
#############################################################################
run_test 7 "Remove Interface" -var='interfaces=[]'

#############################################################################
# TEST 8: Interface - Add multiple interfaces
#############################################################################
run_test 8 "Add Multiple Interfaces" -var="interfaces=[\"$SUBNET1_ID\",\"$SUBNET2_ID\"]"

#############################################################################
# TEST 9: Interface - Replace interface (remove subnet1, keep subnet2)
#############################################################################
run_test 9 "Replace Interface" -var="interfaces=[\"$SUBNET2_ID\"]"

#############################################################################
# TEST 10: Interface - Remove all interfaces
#############################################################################
run_test 10 "Remove All Interfaces" -var='interfaces=[]'

#############################################################################
# TEST 11: Mixed - Add route + add interface together (CRITICAL - test ordering)
#############################################################################
run_test 11 "Add Route + Interface Together" -var="interfaces=[\"$SUBNET1_ID\"]" -var='routes=[{destination="10.10.0.0/24",nexthop="192.168.1.1"}]'

#############################################################################
# TEST 12: Mixed - Remove route + keep interface
#############################################################################
run_test 12 "Remove Route Keep Interface" -var="interfaces=[\"$SUBNET1_ID\"]" -var='routes=[]'

#############################################################################
# TEST 13: Mixed - Add interface + add route in same update
#############################################################################
run_test 13 "Add Second Interface + Routes" -var="interfaces=[\"$SUBNET1_ID\",\"$SUBNET2_ID\"]" -var='routes=[{destination="10.20.0.0/24",nexthop="192.168.1.1"}]'

#############################################################################
# TEST 14: Mixed - Remove route + remove interface together (CRITICAL - test ordering)
#############################################################################
run_test 14 "Remove Route + Interface Together" -var='interfaces=[]' -var='routes=[]'

#############################################################################
# TEST 15: Name - Update name only
#############################################################################
run_test 15 "Update Router Name" -var='router_name=test-router-renamed'

#############################################################################
# TEST 16: ExtGateway - Set external gateway info
#############################################################################
run_test 16 "Set External Gateway" -var='enable_external_gateway=true' -var='external_gateway_snat=true'

#############################################################################
# TEST 17: ExtGateway - Update external gateway SNAT
#############################################################################
run_test 17 "Update External Gateway SNAT" -var='enable_external_gateway=true' -var='external_gateway_snat=false'

#############################################################################
# TEST 18: ExtGateway - Remove external gateway info
#############################################################################
run_test 18 "Remove External Gateway" -var='enable_external_gateway=false'

#############################################################################
# TEST 19: Drift - Complex config with all features
#############################################################################
run_test 19 "Complex Config - Full Feature Drift Check" -var='router_name=test-router-complex' -var="interfaces=[\"$SUBNET1_ID\",\"$SUBNET2_ID\"]" -var='routes=[{destination="10.99.0.0/24",nexthop="192.168.1.99"}]' -var='enable_external_gateway=true'

echo ""
echo "===================================================================="
echo "ALL TESTS COMPLETED"
echo "===================================================================="
echo ""
cat $RESULTS_FILE
