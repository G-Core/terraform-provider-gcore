#!/bin/bash
# Comprehensive LB Testing
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Load credentials
set -o allexport
source ../.env
set +o allexport

export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/.terraformrc"
export TF_LOG=DEBUG
export TF_LOG_PATH="test.log"

echo "========================================"
echo "LB COMPREHENSIVE TESTING"
echo "========================================"

# Clean state
rm -f terraform.tfstate* .terraform.lock.hcl
rm -rf .terraform

echo ""
echo "TEST 1: Create minimal LB + Listener (no pool)"
terraform apply -auto-approve -var="create_pool=false"

echo "Checking for drift..."
if terraform plan -detailed-exitcode; then
  echo -e "${GREEN}✅ PASS: No drift${NC}"
else
  echo -e "${RED}❌ FAIL: Drift detected${NC}"
  terraform plan -no-color | head -50
fi

echo ""
echo "TEST 2: Rename LB"
terraform apply -auto-approve -var="lb_name=qa-lb-RENAMED" -var="create_pool=false"

echo "Checking for drift after rename..."
if terraform plan -detailed-exitcode; then
  echo -e "${GREEN}✅ PASS: No drift after rename${NC}"
else
  echo -e "${RED}❌ FAIL: Drift after rename${NC}"
fi

echo ""
echo "TEST 3: Add pool (timing test)"
terraform apply -auto-approve -var="create_pool=true"

POOL_ID=$(terraform show -json | jq -r '.values.root_module.resources[] | select(.type=="gcore_cloud_load_balancer_pool") | .values.id')
if [ -n "$POOL_ID" ] && [ "$POOL_ID" != "null" ]; then
  echo -e "${GREEN}✅ PASS: Pool created: $POOL_ID${NC}"
else
  echo -e "${RED}❌ FAIL: Pool creation failed${NC}"
fi

echo "Checking for drift with pool..."
if terraform plan -detailed-exitcode; then
  echo -e "${GREEN}✅ PASS: No drift with pool${NC}"
else
  echo -e "${RED}❌ FAIL: Drift with pool${NC}"
  terraform plan -no-color | head -80
fi

echo ""
echo "TEST 4: Set timeout values"
terraform apply -auto-approve \
  -var="timeout_client_data=50000" \
  -var="timeout_member_connect=5000" \
  -var="timeout_member_data=50000"

echo "Checking for drift with timeouts..."
if terraform plan -detailed-exitcode; then
  echo -e "${GREEN}✅ PASS: No drift with timeouts${NC}"
else
  echo -e "${RED}❌ FAIL: Drift with timeouts${NC}"
  terraform plan -no-color | head -80
fi

echo ""
echo "TEST 5: Add health monitor"
terraform apply -auto-approve \
  -var='pool_healthmonitor={"type":"HTTP","delay":10,"max_retries":3,"timeout":5}'

echo "Checking for drift with health monitor..."
if terraform plan -detailed-exitcode; then
  echo -e "${GREEN}✅ PASS: No drift with health monitor${NC}"
else
  echo -e "${RED}❌ FAIL: Drift with health monitor${NC}"
  terraform plan -no-color | head -80
fi

echo ""
echo "========================================"
echo "TESTS COMPLETE"
echo "========================================"
echo "To clean up: terraform destroy -auto-approve"
