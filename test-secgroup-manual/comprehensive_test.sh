#!/bin/bash
set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo "========================================================================"
echo "COMPREHENSIVE SECURITY GROUP RULES TEST"
echo "========================================================================"
echo ""

# Load environment
cd /Users/user/repos/gcore-terraform
set -o allexport
source .env
set +o allexport
cd test-secgroup-manual

echo -e "${BLUE}=== SCENARIO 1: Initial Create & Verify ===${NC}"
echo ""

# Apply
echo "Step 1: Creating resources..."
terraform apply -auto-approve > /tmp/apply1.log 2>&1
echo -e "${GREEN}✓${NC} Resources created"

# Check state
echo ""
echo "Step 2: Checking state..."
terraform state show gcore_cloud_security_group.test > /tmp/state1.txt
echo -e "${GREEN}✓${NC} State captured"

# Get SG ID
SG_ID=$(terraform output -raw security_group_id)
echo "Security Group ID: $SG_ID"

# Check API
echo ""
echo "Step 3: Verifying API..."
RULE_COUNT=$(curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
  "https://api.gcore.com/cloud/v1/securitygroups/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID/$SG_ID" | \
  jq '.security_group_rules | length')
echo "API Rule Count: $RULE_COUNT"

if [ "$RULE_COUNT" -eq 4 ]; then
    echo -e "${GREEN}✓ PASS${NC}: Correct rule count (no default rules)"
else
    echo -e "${RED}✗ FAIL${NC}: Expected 4 rules, got $RULE_COUNT"
fi

# First drift test
echo ""
echo -e "${BLUE}=== SCENARIO 2: Drift Detection (Immediate) ===${NC}"
echo ""
set +e
terraform plan -detailed-exitcode > /tmp/plan1.txt 2>&1
PLAN1_EXIT=$?
set -e

if [ $PLAN1_EXIT -eq 0 ]; then
    echo -e "${GREEN}✓ PASS${NC}: No drift detected"
else
    echo -e "${RED}✗ FAIL${NC}: Drift detected (exit code $PLAN1_EXIT)"
    echo ""
    echo "Analyzing drift..."

    # Check specific fields
    if grep -q "security_group_rules" /tmp/plan1.txt; then
        echo -e "${RED}  ✗ security_group_rules drift found${NC}"
    else
        echo -e "${GREEN}  ✓ security_group_rules: NO drift${NC}"
    fi

    if grep -q "name.*->.*null" /tmp/plan1.txt; then
        echo -e "${RED}  ✗ name field drift: going to null${NC}"
    fi

    if grep -q "created_at.*known after apply" /tmp/plan1.txt; then
        echo -e "${RED}  ✗ created_at showing as 'known after apply'${NC}"
    fi

    if grep -q "description.*known after apply" /tmp/plan1.txt; then
        echo -e "${RED}  ✗ description showing as 'known after apply'${NC}"
    fi

    if grep -q "revision_number.*known after apply" /tmp/plan1.txt; then
        echo -e "${RED}  ✗ revision_number showing as 'known after apply'${NC}"
    fi
fi

# Apply again (should be no-op if no drift)
echo ""
echo -e "${BLUE}=== SCENARIO 3: Second Apply (Should Be No-Op) ===${NC}"
echo ""
terraform apply -auto-approve > /tmp/apply2.log 2>&1
if grep -q "Apply complete! Resources: 0 added, 0 changed, 0 destroyed" /tmp/apply2.log; then
    echo -e "${GREEN}✓ PASS${NC}: No changes applied (no-op)"
else
    if grep -q "Apply complete! Resources: 0 added, 1 changed, 0 destroyed" /tmp/apply2.log; then
        echo -e "${RED}✗ FAIL${NC}: 1 resource changed (should be no-op)"
    else
        echo -e "${RED}✗ FAIL${NC}: Unexpected changes"
    fi
    tail -20 /tmp/apply2.log
fi

# Third drift test (after second apply)
echo ""
echo -e "${BLUE}=== SCENARIO 4: Drift After Second Apply ===${NC}"
echo ""
set +e
terraform plan -detailed-exitcode > /tmp/plan2.txt 2>&1
PLAN2_EXIT=$?
set -e

if [ $PLAN2_EXIT -eq 0 ]; then
    echo -e "${GREEN}✓ PASS${NC}: No drift after second apply"
else
    echo -e "${RED}✗ FAIL${NC}: Still showing drift (exit code $PLAN2_EXIT)"
fi

# Update a rule
echo ""
echo -e "${BLUE}=== SCENARIO 5: Update Rule Description ===${NC}"
echo ""
cat > test_update.tf << 'EOF'
resource "gcore_cloud_security_group_rule" "test_update" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = 379987
  region_id  = 76

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = 9000
  port_range_max   = 9000
  remote_ip_prefix = "0.0.0.0/0"
  description      = "Test rule for updates"
}
EOF

terraform apply -auto-approve > /tmp/apply3.log 2>&1
echo -e "${GREEN}✓${NC} Rule added"

# Verify rule count increased
NEW_RULE_COUNT=$(curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
  "https://api.gcore.com/cloud/v1/securitygroups/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID/$SG_ID" | \
  jq '.security_group_rules | length')
echo "New API Rule Count: $NEW_RULE_COUNT (expected: 5)"

if [ "$NEW_RULE_COUNT" -eq 5 ]; then
    echo -e "${GREEN}✓ PASS${NC}: Rule added successfully"
else
    echo -e "${RED}✗ FAIL${NC}: Expected 5 rules, got $NEW_RULE_COUNT"
fi

# Remove the test rule
echo ""
echo "Removing test rule..."
rm test_update.tf
terraform apply -auto-approve > /tmp/apply4.log 2>&1
echo -e "${GREEN}✓${NC} Rule removed"

# Verify rule count back to original
FINAL_RULE_COUNT=$(curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
  "https://api.gcore.com/cloud/v1/securitygroups/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID/$SG_ID" | \
  jq '.security_group_rules | length')
echo "Final API Rule Count: $FINAL_RULE_COUNT (expected: 4)"

if [ "$FINAL_RULE_COUNT" -eq 4 ]; then
    echo -e "${GREEN}✓ PASS${NC}: Rule removed successfully"
else
    echo -e "${RED}✗ FAIL${NC}: Expected 4 rules, got $FINAL_RULE_COUNT"
fi

# Summary
echo ""
echo "========================================================================"
echo "TEST SUMMARY"
echo "========================================================================"
echo ""
echo "Key Findings:"
echo ""
echo "1. Default Rule Deletion:"
if [ "$RULE_COUNT" -eq 4 ]; then
    echo -e "   ${GREEN}✓ WORKING${NC}: Only 4 user rules (no 39 defaults)"
else
    echo -e "   ${RED}✗ BROKEN${NC}: Found $RULE_COUNT rules (expected 4)"
fi

echo ""
echo "2. security_group_rules Drift:"
if ! grep -q "security_group_rules" /tmp/plan1.txt; then
    echo -e "   ${GREEN}✓ FIXED${NC}: No drift on security_group_rules field"
else
    echo -e "   ${RED}✗ ISSUE${NC}: Drift detected on security_group_rules"
fi

echo ""
echo "3. Computed Fields Drift:"
if [ $PLAN1_EXIT -eq 0 ]; then
    echo -e "   ${GREEN}✓ FIXED${NC}: No drift on computed fields"
else
    if grep -q "known after apply" /tmp/plan1.txt; then
        echo -e "   ${RED}✗ ISSUE${NC}: Computed fields showing 'known after apply'"
        echo "   Affected fields:"
        grep -o "[a-z_]*.*known after apply" /tmp/plan1.txt | head -5 | sed 's/^/     - /'
    fi
fi

echo ""
echo "4. Top-Level name Field:"
if grep -q "name.*->.*null" /tmp/plan1.txt; then
    echo -e "   ${RED}✗ ISSUE${NC}: name field drifting to null"
else
    echo -e "   ${GREEN}✓ FIXED${NC}: No drift on name field"
fi

echo ""
echo "5. Rule Add/Remove:"
if [ "$NEW_RULE_COUNT" -eq 5 ] && [ "$FINAL_RULE_COUNT" -eq 4 ]; then
    echo -e "   ${GREEN}✓ WORKING${NC}: Rules can be added and removed"
else
    echo -e "   ${RED}✗ ISSUE${NC}: Problem with rule management"
fi

echo ""
echo "========================================================================"
echo ""
echo -e "${YELLOW}Logs saved to /tmp/ for analysis:${NC}"
echo "  - /tmp/apply1.log (initial create)"
echo "  - /tmp/plan1.txt (first drift check)"
echo "  - /tmp/apply2.log (second apply)"
echo "  - /tmp/plan2.txt (drift after second apply)"
echo "  - /tmp/state1.txt (state dump)"
echo ""
