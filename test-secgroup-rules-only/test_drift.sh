#!/bin/bash
set -e

echo "========================================="
echo "Testing Security Group Rules - NO DRIFT"
echo "========================================="

# Build the local provider
echo ""
echo "=== Step 1: Build LOCAL provider ==="
cd /Users/user/repos/gcore-terraform
go build -o terraform-provider-gcore
echo "✅ Provider built successfully"

# Setup test directory
echo ""
echo "=== Step 2: Setup test directory ==="
cd test-secgroup-rules-only

# Clean up any previous state
rm -f terraform.tfstate terraform.tfstate.backup .terraform.lock.hcl
rm -rf .terraform/

# Source environment variables
if [ -f ../.env ]; then
    source ../.env
    echo "✅ Environment variables loaded"
else
    echo "❌ .env file not found"
    exit 1
fi

# Set Terraform to use local provider
export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/test-secgroup-rules-only/.terraformrc"
echo "✅ Using local provider override"
echo "⚠️  Skipping terraform init (not needed with dev_overrides)"

# First apply
echo ""
echo "=== Step 4: First terraform apply ==="
terraform apply -auto-approve

# Get the security group ID
SG_ID=$(terraform output -raw security_group_id)
echo ""
echo "✅ Created security group: $SG_ID"

# Check API to see how many rules exist
echo ""
echo "=== Step 5: Check rules via API ==="
RULE_COUNT=$(curl -s -H "Authorization: APIKey ${GCORE_API_KEY}" \
  "https://api.gcore.com/cloud/v1/securitygroups/${SG_ID}?project_id=379987&region_id=76" | \
  jq '.security_group_rules | length')
echo "Rules in backend: $RULE_COUNT"

# Second apply - should show NO CHANGES
echo ""
echo "=== Step 6: Second terraform apply (checking for drift) ==="
terraform plan -detailed-exitcode > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✅ NO DRIFT DETECTED - Test PASSED!"
else
    echo "❌ DRIFT DETECTED - Test FAILED!"
    echo ""
    echo "Terraform plan output:"
    terraform plan
    exit 1
fi

# Third apply - final verification
echo ""
echo "=== Step 7: Third terraform apply (final check) ==="
terraform apply -auto-approve

# Check that rules weren't deleted
echo ""
echo "=== Step 8: Verify rules still exist ==="
RULE_COUNT_AFTER=$(curl -s -H "Authorization: APIKey ${GCORE_API_KEY}" \
  "https://api.gcore.com/cloud/v1/securitygroups/${SG_ID}?project_id=379987&region_id=76" | \
  jq '.security_group_rules | length')
echo "Rules in backend after multiple applies: $RULE_COUNT_AFTER"

if [ "$RULE_COUNT" -eq "$RULE_COUNT_AFTER" ]; then
    echo "✅ Rules preserved - no unexpected deletions"
else
    echo "❌ Rule count changed: $RULE_COUNT -> $RULE_COUNT_AFTER"
    exit 1
fi

# Check state file - should NOT have security_group_rules
echo ""
echo "=== Step 9: Check Terraform state ==="
if terraform show | grep -q "security_group_rules"; then
    echo "❌ State contains security_group_rules field (should NOT be there)"
    terraform show | grep -A 10 security_group_rules
    exit 1
else
    echo "✅ State does NOT contain security_group_rules field"
fi

echo ""
echo "========================================="
echo "✅ ALL TESTS PASSED!"
echo "========================================="
echo ""
echo "Summary:"
echo "- Security group created: $SG_ID"
echo "- Backend has $RULE_COUNT_AFTER default rules"
echo "- NO drift detected across multiple applies"
echo "- State does NOT track backend rules"
echo ""
