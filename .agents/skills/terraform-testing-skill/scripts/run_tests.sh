#!/bin/bash
# run_tests.sh - Comprehensive test runner for Terraform resources

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test results tracking
TESTS_PASSED=0
TESTS_FAILED=0
FAILED_TESTS=()

# Function to print test header
print_test_header() {
    echo ""
    echo -e "${BLUE}================================================================${NC}"
    echo -e "${BLUE}  TEST: $1${NC}"
    echo -e "${BLUE}================================================================${NC}"
}

# Function to print test result
print_test_result() {
    local test_name=$1
    local result=$2
    
    if [ $result -eq 0 ]; then
        echo -e "${GREEN}✓ PASS: $test_name${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "${RED}✗ FAIL: $test_name${NC}"
        ((TESTS_FAILED++))
        FAILED_TESTS+=("$test_name")
    fi
}

# Function to run drift test
run_drift_test() {
    local resource=$1
    local test_name="TC-DRIFT-001-$resource-no-changes"
    
    print_test_header "$test_name"
    
    # Create test directory
    mkdir -p "tests/drift/$test_name"
    cd "tests/drift/$test_name"
    
    # Copy test configuration
    if [ -f "../../../configs/$resource/minimal.tf" ]; then
        cp "../../../configs/$resource/minimal.tf" main.tf
    else
        echo -e "${RED}Error: No minimal config found for $resource${NC}"
        cd ../../..
        return 1
    fi
    
    # Initialize Terraform
    terraform init -upgrade > /dev/null 2>&1
    
    # Apply configuration
    echo "Creating resource..."
    terraform apply -auto-approve > apply.log 2>&1
    
    # Check for drift
    echo "Checking for drift..."
    terraform plan -detailed-exitcode > plan.log 2>&1
    local exit_code=$?
    
    if [ $exit_code -eq 0 ]; then
        print_test_result "$test_name" 0
    elif [ $exit_code -eq 2 ]; then
        echo -e "${RED}Drift detected! See plan.log for details${NC}"
        print_test_result "$test_name" 1
    else
        echo -e "${RED}Plan failed! See plan.log for details${NC}"
        print_test_result "$test_name" 1
    fi
    
    # Clean up
    terraform destroy -auto-approve > /dev/null 2>&1
    cd ../../..
}

# Function to run update test
run_update_test() {
    local resource=$1
    local field=$2
    local test_name="TC-UPDATE-001-$resource-$field"
    
    print_test_header "$test_name"
    
    # Create test directory
    mkdir -p "tests/update/$test_name"
    cd "tests/update/$test_name"
    
    # Copy initial configuration
    if [ -f "../../../configs/$resource/update_$field.tf" ]; then
        cp "../../../configs/$resource/update_$field.tf" main.tf
    else
        echo -e "${RED}Error: No update config found for $resource/$field${NC}"
        cd ../../..
        return 1
    fi
    
    # Initialize and apply initial configuration
    terraform init -upgrade > /dev/null 2>&1
    terraform apply -auto-approve > apply_initial.log 2>&1
    
    # Get initial resource ID
    local initial_id=$(terraform show -json | jq -r '.values.root_module.resources[0].values.id')
    
    # Apply the update variable
    echo "Updating $field..."
    terraform apply -auto-approve -var="${field}_updated=true" > apply_update.log 2>&1
    
    # Get updated resource ID
    local updated_id=$(terraform show -json | jq -r '.values.root_module.resources[0].values.id')
    
    # Check if resource was recreated
    if [ "$initial_id" = "$updated_id" ]; then
        echo -e "${GREEN}Resource updated in-place (PATCH operation)${NC}"
        
        # Check API calls in log
        if grep -q "PATCH" apply_update.log; then
            print_test_result "$test_name" 0
        else
            echo -e "${YELLOW}Warning: No PATCH operation found in logs${NC}"
            print_test_result "$test_name" 1
        fi
    else
        echo -e "${RED}Resource was recreated (ID changed from $initial_id to $updated_id)${NC}"
        print_test_result "$test_name" 1
    fi
    
    # Clean up
    terraform destroy -auto-approve > /dev/null 2>&1
    cd ../../..
}

# Function to run corner case tests
run_corner_case_test() {
    local resource=$1
    local case=$2
    local test_name="TC-CORNER-001-$resource-$case"
    
    print_test_header "$test_name"
    
    # Create test directory
    mkdir -p "tests/corner_cases/$test_name"
    cd "tests/corner_cases/$test_name"
    
    # Copy corner case configuration
    if [ -f "../../../configs/$resource/corner_$case.tf" ]; then
        cp "../../../configs/$resource/corner_$case.tf" main.tf
    else
        echo -e "${YELLOW}No corner case config for $resource/$case, skipping${NC}"
        cd ../../..
        return 0
    fi
    
    # Run the test
    terraform init -upgrade > /dev/null 2>&1
    terraform apply -auto-approve > apply.log 2>&1
    local apply_result=$?
    
    if [ $apply_result -eq 0 ]; then
        # Check for drift after apply
        terraform plan -detailed-exitcode > plan.log 2>&1
        local plan_result=$?
        
        if [ $plan_result -eq 0 ]; then
            print_test_result "$test_name" 0
        else
            echo -e "${RED}Corner case has drift or errors${NC}"
            print_test_result "$test_name" 1
        fi
    else
        echo -e "${RED}Failed to apply corner case configuration${NC}"
        print_test_result "$test_name" 1
    fi
    
    # Clean up
    terraform destroy -auto-approve > /dev/null 2>&1 || true
    cd ../../..
}

# Function to generate test report
generate_report() {
    local resource=$1
    local report_file="tests/TEST_REPORT_$(date +%Y%m%d_%H%M%S).md"
    
    cat > "$report_file" << EOF
# Test Report for $resource

**Date:** $(date)
**Resource:** $resource
**Provider Version:** $(terraform version -json | jq -r '.provider_selections."registry.terraform.io/gcore/gcore"' 2>/dev/null || echo "dev-override")

## Summary

- **Total Tests:** $((TESTS_PASSED + TESTS_FAILED))
- **Passed:** $TESTS_PASSED ✅
- **Failed:** $TESTS_FAILED ❌
- **Success Rate:** $(( TESTS_PASSED * 100 / (TESTS_PASSED + TESTS_FAILED) ))%

## Test Results

### Passed Tests
EOF
    
    if [ $TESTS_PASSED -gt 0 ]; then
        echo "✅ All drift tests passed" >> "$report_file"
    fi
    
    if [ ${#FAILED_TESTS[@]} -gt 0 ]; then
        echo "" >> "$report_file"
        echo "### Failed Tests" >> "$report_file"
        for test in "${FAILED_TESTS[@]}"; do
            echo "- ❌ $test" >> "$report_file"
        done
    fi
    
    echo "" >> "$report_file"
    echo "## Recommendations" >> "$report_file"
    
    if [ $TESTS_FAILED -gt 0 ]; then
        echo "- Review failed test logs in tests/ directory" >> "$report_file"
        echo "- Check drift issues with UnmarshalComputed usage" >> "$report_file"
        echo "- Verify computed_optional field configurations" >> "$report_file"
    else
        echo "- All tests passed! Resource is ready for production." >> "$report_file"
    fi
    
    echo ""
    echo -e "${GREEN}Test report saved to: $report_file${NC}"
}

# Main execution
main() {
    local resource=${1:-router}
    
    echo -e "${GREEN}Starting comprehensive tests for: $resource${NC}"
    
    # Setup environment
    source scripts/set_env.sh
    
    # Create test directories
    mkdir -p tests/{drift,update,corner_cases,reports}
    
    # Run drift tests
    echo -e "\n${YELLOW}Running Drift Tests...${NC}"
    run_drift_test "$resource"
    
    # Run update tests
    echo -e "\n${YELLOW}Running Update Tests...${NC}"
    run_update_test "$resource" "name"
    
    # Run corner case tests
    echo -e "\n${YELLOW}Running Corner Case Tests...${NC}"
    run_corner_case_test "$resource" "empty_arrays"
    run_corner_case_test "$resource" "null_fields"
    
    # Generate report
    echo -e "\n${YELLOW}Generating Test Report...${NC}"
    generate_report "$resource"
    
    # Final summary
    echo ""
    echo -e "${BLUE}================================================================${NC}"
    echo -e "${BLUE}  FINAL RESULTS${NC}"
    echo -e "${BLUE}================================================================${NC}"
    echo -e "Tests Passed: ${GREEN}$TESTS_PASSED${NC}"
    echo -e "Tests Failed: ${RED}$TESTS_FAILED${NC}"
    
    if [ $TESTS_FAILED -eq 0 ]; then
        echo -e "\n${GREEN}✅ ALL TESTS PASSED!${NC}"
        exit 0
    else
        echo -e "\n${RED}❌ SOME TESTS FAILED - Review logs for details${NC}"
        exit 1
    fi
}

# Check if resource name provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 <resource_name>"
    echo "Example: $0 router"
    exit 1
fi

# Run main function
main "$@"