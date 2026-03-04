#!/bin/bash
# Comprehensive LB Testing Script
# Tests all issues reported by Kirill in GCLOUD2-20778

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Load credentials
if [ -f ../.env ]; then
  set -o allexport
  source ../.env
  set +o allexport
  echo -e "${GREEN}✅ Credentials loaded${NC}"
else
  echo -e "${RED}❌ .env file not found${NC}"
  exit 1
fi

# Set Terraform config
export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/.terraformrc"
export TF_LOG=DEBUG
export TF_LOG_PATH="comprehensive_test.log"

echo "========================================="
echo "COMPREHENSIVE LB TESTING"
echo "Testing GCLOUD2-20778 issues"
echo "========================================="

# Clean up any previous state
rm -f terraform.tfstate* .terraform.lock.hcl
rm -rf .terraform

echo ""
echo "========================================="
echo "TEST 1: DRIFT DETECTION - Minimal Config"
echo "Issue: LB rename causes drift"
echo "Issue: Listener shows fields as 'known after apply'"
echo "========================================="
terraform init
terraform apply -auto-approve \
  -var="lb_name=qa-lb-drift-test" \
  -var="create_pool=false"

echo ""
echo "Running second plan to check for drift..."
if terraform plan -detailed-exitcode; then
  echo -e "${GREEN}✅ PASS: No drift detected${NC}"
else
  EXIT_CODE=$?
  if [ $EXIT_CODE -eq 2 ]; then
    echo -e "${RED}❌ FAIL: Drift detected on listener/LB${NC}"
    terraform plan -no-color | tee drift_detected.txt
  fi
fi

echo ""
echo "========================================="
echo "TEST 2: LB RENAME"
echo "Issue: Renaming LB with listener causes drift"
echo "========================================="
terraform apply -auto-approve -var="lb_name=qa-lb-RENAMED"

echo "Checking for drift after rename..."
if terraform plan -detailed-exitcode; then
  echo -e "${GREEN}✅ PASS: No drift after rename${NC}"
else
  echo -e "${RED}❌ FAIL: Drift detected after rename${NC}"
  terraform plan -no-color | tee rename_drift.txt
fi

echo ""
echo "========================================="
echo "TEST 3: TAG OPERATIONS"
echo "Issue: Tag removal not working"
echo "========================================="
echo "Adding tags..."
terraform apply -auto-approve -var='lb_tags=["test-tag-1","test-tag-2"]'

echo "Verifying tags were added..."
terraform show -json | jq '.values.root_module.resources[] | select(.type=="gcore_cloud_load_balancer") | .values.tags_v2'

echo "Removing tags..."
terraform apply -auto-approve -var='lb_tags=[]'

echo "Verifying tags were removed..."
TAGS=$(terraform show -json | jq -r '.values.root_module.resources[] | select(.type=="gcore_cloud_load_balancer") | .values.tags_v2 | length')
if [ "$TAGS" -eq 0 ]; then
  echo -e "${GREEN}✅ PASS: Tags removed successfully${NC}"
else
  echo -e "${RED}❌ FAIL: Tags still present: $TAGS${NC}"
fi

echo ""
echo "========================================="
echo "TEST 4: POOL CREATION WITH LISTENER"
echo "Issue: Pool created before listener ready"
echo "========================================="
echo "Creating pool with listener..."
terraform apply -auto-approve -var="create_pool=true"

POOL_STATUS=$(terraform show -json | jq -r '.values.root_module.resources[] | select(.type=="gcore_cloud_load_balancer_pool") | .values.operating_status')
echo "Pool operating status: $POOL_STATUS"

if [ "$POOL_STATUS" = "ONLINE" ] || [ "$POOL_STATUS" = "DEGRADED" ] || [ "$POOL_STATUS" = "NO_MONITOR" ]; then
  echo -e "${GREEN}✅ PASS: Pool created successfully${NC}"
else
  echo -e "${RED}❌ FAIL: Pool in unexpected status: $POOL_STATUS${NC}"
fi

echo ""
echo "========================================="
echo "TEST 5: TIMEOUT FIELDS (COMPUTED_OPTIONAL)"
echo "Issue: Listener timeout fields cause drift"
echo "========================================="
echo "Setting timeout values..."
terraform apply -auto-approve \
  -var="timeout_client_data=50000" \
  -var="timeout_member_connect=5000" \
  -var="timeout_member_data=50000"

echo "Checking for drift after setting timeouts..."
if terraform plan -detailed-exitcode; then
  echo -e "${GREEN}✅ PASS: No drift with timeout values${NC}"
else
  echo -e "${RED}❌ FAIL: Drift with timeout values${NC}"
fi

echo "Clearing timeout values (set to null)..."
terraform apply -auto-approve \
  -var="timeout_client_data=null" \
  -var="timeout_member_connect=null" \
  -var="timeout_member_data=null"

echo "Checking for drift after clearing timeouts..."
if terraform plan -detailed-exitcode; then
  echo -e "${GREEN}✅ PASS: No drift after clearing timeouts${NC}"
else
  echo -e "${RED}❌ FAIL: Drift after clearing timeouts${NC}"
  terraform plan -no-color | tee timeout_drift.txt
fi

echo ""
echo "========================================="
echo "TEST 6: HEALTH MONITOR WITH COMPUTED FIELDS"
echo "========================================="
echo "Adding health monitor with only required fields..."
terraform apply -auto-approve \
  -var='pool_healthmonitor={"type":"HTTP","delay":10,"max_retries":3,"timeout":5}'

echo "Checking computed fields (http_method, max_retries_down)..."
HM=$(terraform show -json | jq '.values.root_module.resources[] | select(.type=="gcore_cloud_load_balancer_pool") | .values.healthmonitor')
echo "Health monitor: $HM"

echo "Checking for drift..."
if terraform plan -detailed-exitcode; then
  echo -e "${GREEN}✅ PASS: No drift with health monitor${NC}"
else
  echo -e "${RED}❌ FAIL: Drift with health monitor${NC}"
fi

echo ""
echo "========================================="
echo "TEST 7: POOL MEMBERS INLINE"
echo "========================================="
# Get subnet ID for members
SUBNET_ID=$(terraform show -json | jq -r '.values.root_module.resources[] | select(.type=="gcore_cloud_subnet") | .values.id')

echo "Adding pool members..."
terraform apply -auto-approve \
  -var="pool_members=[{address=\"10.0.1.10\",protocol_port=80,subnet_id=\"$SUBNET_ID\"},{address=\"10.0.1.11\",protocol_port=80,subnet_id=\"$SUBNET_ID\"}]"

echo "Checking for drift..."
if terraform plan -detailed-exitcode; then
  echo -e "${GREEN}✅ PASS: No drift with members${NC}"
else
  echo -e "${RED}❌ FAIL: Drift with members${NC}"
fi

echo ""
echo "========================================="
echo "TEST 8: CLEARING POOL MEMBERS"
echo "========================================="
echo "Clearing all members..."
terraform apply -auto-approve -var="pool_members=[]"

MEMBER_COUNT=$(terraform show -json | jq '.values.root_module.resources[] | select(.type=="gcore_cloud_load_balancer_pool") | .values.members | length')
if [ "$MEMBER_COUNT" -eq 0 ]; then
  echo -e "${GREEN}✅ PASS: Members cleared${NC}"
else
  echo -e "${RED}❌ FAIL: Members still present: $MEMBER_COUNT${NC}"
fi

echo ""
echo "========================================="
echo "TEST 9: SEPARATE POOL MEMBER RESOURCE"
echo "========================================="
echo "Creating separate pool member..."
terraform apply -auto-approve \
  -var="create_separate_member=true" \
  -var="separate_member_address=10.0.1.20"

MEMBER_ID=$(terraform show -json | jq -r '.values.root_module.resources[] | select(.type=="gcore_cloud_load_balancer_pool_member") | .values.id')
if [ -n "$MEMBER_ID" ] && [ "$MEMBER_ID" != "null" ]; then
  echo -e "${GREEN}✅ PASS: Separate member created: $MEMBER_ID${NC}"
else
  echo -e "${RED}❌ FAIL: Separate member not created${NC}"
fi

echo ""
echo "========================================="
echo "TEST SUMMARY"
echo "========================================="
echo "Review the results above"
echo "Logs saved to: comprehensive_test.log"
echo ""
echo "To clean up resources:"
echo "  terraform destroy -auto-approve"
