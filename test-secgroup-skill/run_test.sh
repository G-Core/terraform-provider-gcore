#!/bin/bash
set -e

# Load credentials
cd /Users/user/repos/gcore-terraform
set -o allexport
source .env
set +o allexport

cd test-secgroup-skill

echo "=== TC-01: Create security group + rules ==="
terraform apply -auto-approve

echo ""
echo "=== TC-02: Drift Test - Second Plan ==="
terraform plan -detailed-exitcode
EXIT_CODE=$?
if [ $EXIT_CODE -eq 0 ]; then
  echo "✅ TC-02 PASS: No drift detected"
else
  echo "❌ TC-02 FAIL: Drift detected (exit code: $EXIT_CODE)"
fi

echo ""
echo "=== Check API: Security Group Rules ==="
SG_ID=$(terraform output -raw security_group_id)
echo "Security Group ID: $SG_ID"
RULE_COUNT=$(curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
  "https://api.gcore.com/cloud/v1/securitygroups/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID/$SG_ID" | \
  jq '.security_group_rules | length')
echo "API Rule Count: $RULE_COUNT"
echo "Expected: 6 (3 basic rules + 3 loop rules)"

echo ""
echo "=== TC-05: Update Security Group Description ==="
terraform apply -auto-approve -var='description=Updated description'

echo ""
echo "=== Cleanup ==="
terraform destroy -auto-approve

echo ""
echo "=== Test Complete! ==="
