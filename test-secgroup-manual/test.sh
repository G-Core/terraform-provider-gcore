#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "========================================"
echo "Manual Test: Security Group Rules"
echo "========================================"
echo ""

# Load environment
echo "Loading environment variables..."
cd /Users/user/repos/gcore-terraform
set -o allexport
source .env
set +o allexport
cd test-secgroup-manual

echo -e "${GREEN}✓${NC} Environment loaded"
echo ""

# Initialize (skip actual init due to dev overrides)
echo "Initializing Terraform..."
echo -e "${YELLOW}Note:${NC} Skipping terraform init (using dev provider overrides)"
echo -e "${GREEN}✓${NC} Ready to test"
echo ""

# Apply
echo "========================================"
echo "TEST 1: Creating Security Group + Rules"
echo "========================================"
terraform apply -auto-approve

echo ""
echo "========================================"
echo "TEST 2: Verify State"
echo "========================================"
echo "Checking security_group_rules in state..."
STATE_RULES=$(terraform state show gcore_cloud_security_group.test | grep "security_group_rules" | head -1)
echo "$STATE_RULES"

if echo "$STATE_RULES" | grep -q "= \[\]"; then
    echo -e "${GREEN}✓ PASS${NC}: State shows empty list"
else
    echo -e "${RED}✗ FAIL${NC}: State should show empty list"
fi

echo ""
echo "========================================"
echo "TEST 3: Verify API Rule Count"
echo "========================================"
SG_ID=$(terraform output -raw security_group_id)
echo "Security Group ID: $SG_ID"

RULE_COUNT=$(curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
  "https://api.gcore.com/cloud/v1/securitygroups/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID/$SG_ID" | \
  jq '.security_group_rules | length')

echo "API Rule Count: $RULE_COUNT"
echo "Expected: 4 (2 individual + 2 loop rules)"

if [ "$RULE_COUNT" -eq 4 ]; then
    echo -e "${GREEN}✓ PASS${NC}: Correct rule count (no default rules)"
else
    echo -e "${RED}✗ FAIL${NC}: Expected 4 rules, got $RULE_COUNT"
fi

echo ""
echo "========================================"
echo "TEST 4: View Rule Details"
echo "========================================"
curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
  "https://api.gcore.com/cloud/v1/securitygroups/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID/$SG_ID" | \
  jq '.security_group_rules[] | {direction, protocol, port: .port_range_min, description}' | \
  jq -s '.'

echo ""
echo "========================================"
echo "TEST 5: Drift Detection (CRITICAL)"
echo "========================================"
echo "Running terraform plan..."
set +e
terraform plan -detailed-exitcode > /tmp/plan_output.txt 2>&1
PLAN_EXIT=$?
set -e

if [ $PLAN_EXIT -eq 0 ]; then
    echo -e "${GREEN}✓ PASS${NC}: No drift detected (exit code 0)"
elif [ $PLAN_EXIT -eq 2 ]; then
    echo -e "${RED}✗ FAIL${NC}: Drift detected (exit code 2)"
    echo ""
    echo "Plan output:"
    cat /tmp/plan_output.txt
else
    echo -e "${RED}✗ ERROR${NC}: Plan failed (exit code $PLAN_EXIT)"
    cat /tmp/plan_output.txt
fi

echo ""
echo "========================================"
echo "TEST 6: Verify Outputs"
echo "========================================"
terraform output

echo ""
echo "========================================"
echo "All Tests Complete!"
echo "========================================"
echo ""
echo -e "${YELLOW}Note:${NC} Resources are still running. To clean up, run:"
echo "  terraform destroy -auto-approve"
echo ""
