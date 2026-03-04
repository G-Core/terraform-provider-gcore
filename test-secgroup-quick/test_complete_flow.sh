#!/bin/bash
set -e

echo "=========================================="
echo "Testing Security Group Rules with Registry"
echo "=========================================="

source ../.env
export TF_CLI_CONFIG_FILE=/Users/user/repos/gcore-terraform/.terraformrc

echo ""
echo "=== Apply 1: Create resources ==="
terraform apply -auto-approve

SG_ID=$(terraform output -raw security_group_id)
echo "Security group ID: $SG_ID"

echo ""
echo "=== Check API rule count ==="
RULE_COUNT=$(curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
  "https://api.gcore.com/cloud/v1/securitygroups/379987/76/$SG_ID" | jq '.security_group_rules | length')
echo "Total rules in API: $RULE_COUNT"
echo "Expected: 42 (3 user + 39 backend)"

echo ""
echo "=== Update description to trigger cleanup ==="
sed -i.bak 's/Testing rule drift/Testing rule drift - updated/' main.tf

echo ""
echo "=== Apply 2: Trigger backend rule cleanup ==="
terraform apply -auto-approve 2>&1 | grep -E "(Plan:|Apply complete)" || true

echo ""
echo "=== Check API rule count after cleanup ==="
RULE_COUNT=$(curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
  "https://api.gcore.com/cloud/v1/securitygroups/379987/76/$SG_ID" | jq '.security_group_rules | length')
echo "Total rules in API: $RULE_COUNT"
echo "Expected: 3 (all user rules preserved, including one without description)"

echo ""
echo "=== Check that rule without description still exists ==="
curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
  "https://api.gcore.com/cloud/v1/securitygroups/379987/76/$SG_ID" | \
  jq '.security_group_rules[] | select(.port_range_min == 8080) | {id: .id, port: .port_range_min, description: .description}'

echo ""
echo "=== Apply 3: Verify no further changes ==="
terraform apply -auto-approve 2>&1 | grep -E "(No changes|Apply complete)" || true

echo ""
echo "=== Restore original config ==="
mv main.tf.bak main.tf

echo ""
echo "=========================================="
echo "Test Complete!"
echo "=========================================="
