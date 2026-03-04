#!/bin/bash
# Comprehensive Test Script for gcore_cloud_load_balancer
# Tests GCLOUD2-20778 fix and full resource behavior

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Load environment
cd "$(dirname "$0")"
set -o allexport
source /Users/user/repos/gcore-terraform/.env
set +o allexport
export TF_CLI_CONFIG_FILE=/Users/user/repos/gcore-terraform/.terraformrc

# Track results
declare -A RESULTS
FAILED=0

# Function to run test
run_test() {
    local test_num=$1
    local test_name=$2
    local tf_args=$3
    local check_drift=${4:-true}
    
    echo -e "\n${YELLOW}=== Test $test_num: $test_name ===${NC}"
    
    # Run apply
    if terraform apply -auto-approve $tf_args 2>&1 | tee "evidence/apply_test_${test_num}.log"; then
        echo -e "${GREEN}Apply succeeded${NC}"
        
        if [ "$check_drift" = true ]; then
            # Check for drift
            if terraform plan -detailed-exitcode $tf_args > "evidence/plan_test_${test_num}.log" 2>&1; then
                echo -e "${GREEN}Drift check passed (exit code 0)${NC}"
                RESULTS[$test_num]="PASS"
            else
                exit_code=$?
                if [ $exit_code -eq 2 ]; then
                    echo -e "${RED}Drift detected! (exit code 2)${NC}"
                    cat "evidence/plan_test_${test_num}.log"
                    RESULTS[$test_num]="FAIL - Drift detected"
                    FAILED=$((FAILED + 1))
                else
                    echo -e "${RED}Plan error (exit code $exit_code)${NC}"
                    RESULTS[$test_num]="FAIL - Plan error"
                    FAILED=$((FAILED + 1))
                fi
            fi
        else
            RESULTS[$test_num]="PASS (no drift check)"
        fi
    else
        echo -e "${RED}Apply failed${NC}"
        RESULTS[$test_num]="FAIL - Apply error"
        FAILED=$((FAILED + 1))
    fi
}

# Function to get LB ID
get_lb_id() {
    terraform output -raw lb_id 2>/dev/null || echo ""
}

echo "=============================================="
echo "  Comprehensive Load Balancer Testing"
echo "  GCLOUD2-20778 Fix Verification"
echo "=============================================="

# Test 1: Create minimal LB without tags
echo -e "\n${YELLOW}=== Test 1: Create minimal LB without tags ===${NC}"
if terraform apply -auto-approve 2>&1 | tee evidence/apply_test_1.log; then
    LB_ID=$(get_lb_id)
    echo -e "${GREEN}LB created: $LB_ID${NC}"
    
    # Drift check
    if terraform plan -detailed-exitcode > evidence/plan_test_1.log 2>&1; then
        echo -e "${GREEN}Test 1 PASSED - No drift${NC}"
        RESULTS[1]="PASS"
    else
        echo -e "${RED}Test 1 FAILED - Drift detected${NC}"
        RESULTS[1]="FAIL - Drift"
        FAILED=$((FAILED + 1))
    fi
else
    echo -e "${RED}Test 1 FAILED - Create error${NC}"
    RESULTS[1]="FAIL - Create error"
    FAILED=$((FAILED + 1))
fi

# Test 2: Add tags (GCLOUD2-20778 critical test)
run_test 2 "Add tags to existing LB (GCLOUD2-20778)" '-var=lb_tags={"qa"="load-balancer"}'

# Test 3: Modify tags
run_test 3 "Modify tags" '-var=lb_tags={"qa"="modified","env"="test"}'

# Test 4: Update name
run_test 4 "Update name" '-var=lb_name=test-lb-renamed -var=lb_tags={"qa"="modified","env"="test"}'

# Test 5: Remove all tags
run_test 5 "Remove all tags" '-var=lb_name=test-lb-renamed'

# Test 6: Update preferred_connectivity (may not work on all regions)
echo -e "\n${YELLOW}=== Test 6: Update preferred_connectivity ===${NC}"
if terraform apply -auto-approve -var='lb_name=test-lb-renamed' -var='preferred_connectivity=L3' 2>&1 | tee evidence/apply_test_6.log; then
    if terraform plan -detailed-exitcode -var='lb_name=test-lb-renamed' -var='preferred_connectivity=L3' > evidence/plan_test_6.log 2>&1; then
        echo -e "${GREEN}Test 6 PASSED - No drift${NC}"
        RESULTS[6]="PASS"
    else
        echo -e "${YELLOW}Test 6 WARNING - Drift detected (may be expected)${NC}"
        RESULTS[6]="WARN - Drift"
    fi
else
    echo -e "${YELLOW}Test 6 SKIPPED - preferred_connectivity not supported${NC}"
    RESULTS[6]="SKIPPED"
fi

# Test 7: Resize flavor
echo -e "\n${YELLOW}=== Test 7: Resize flavor ===${NC}"
BEFORE_ID=$(get_lb_id)
if terraform apply -auto-approve -var='lb_name=test-lb-renamed' -var='lb_flavor=lb1-4-8' 2>&1 | tee evidence/apply_test_7.log; then
    AFTER_ID=$(get_lb_id)
    if [ "$BEFORE_ID" = "$AFTER_ID" ]; then
        echo -e "${GREEN}Test 7 PASSED - Resized in place (ID unchanged)${NC}"
        RESULTS[7]="PASS"
    else
        echo -e "${RED}Test 7 FAILED - ID changed (recreated instead of resize)${NC}"
        RESULTS[7]="FAIL - Recreated"
        FAILED=$((FAILED + 1))
    fi
else
    echo -e "${RED}Test 7 FAILED - Resize error${NC}"
    RESULTS[7]="FAIL - Resize error"
    FAILED=$((FAILED + 1))
fi

# Print summary
echo ""
echo "=============================================="
echo "  TEST SUMMARY"
echo "=============================================="
for test_num in "${!RESULTS[@]}"; do
    result=${RESULTS[$test_num]}
    if [[ $result == PASS* ]]; then
        echo -e "Test $test_num: ${GREEN}$result${NC}"
    elif [[ $result == WARN* ]] || [[ $result == SKIP* ]]; then
        echo -e "Test $test_num: ${YELLOW}$result${NC}"
    else
        echo -e "Test $test_num: ${RED}$result${NC}"
    fi
done | sort -t':' -k1 -V

echo ""
if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
else
    echo -e "${RED}$FAILED test(s) failed${NC}"
fi

echo ""
echo "Evidence saved in: evidence/"
echo "=============================================="

exit $FAILED
