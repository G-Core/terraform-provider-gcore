#!/bin/bash
set -e

cd /Users/user/repos/gcore-terraform
source .env

SG_ID="5969403e-c10f-4864-9085-70df4ba18af7"
RULE_ID="0f028d9f-7eba-4697-a0d8-3351b34bcd54"

echo "========================================================================"
echo "TESTING SECURITY GROUP RULE PUT ENDPOINT BEHAVIOR"
echo "========================================================================"
echo ""
echo "Security Group ID: $SG_ID"
echo "Rule ID (HTTP rule): $RULE_ID"
echo ""

echo "=== STEP 1: Get rule details BEFORE update ==="
curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
  "https://api.gcore.com/cloud/v1/securitygroups/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID/$SG_ID" \
  | jq ".security_group_rules[] | select(.id == \"$RULE_ID\")" > /tmp/rule_before.json

echo "Rule BEFORE update:"
cat /tmp/rule_before.json | jq '{id, description, direction, protocol, port_range_min, port_range_max, ethertype, security_group_id, revision_number, updated_at}'
echo ""

ORIGINAL_ID=$(cat /tmp/rule_before.json | jq -r '.id')
ORIGINAL_DESC=$(cat /tmp/rule_before.json | jq -r '.description')
ORIGINAL_REV=$(cat /tmp/rule_before.json | jq -r '.revision_number')
echo "  Original ID: $ORIGINAL_ID"
echo "  Original Description: $ORIGINAL_DESC"
echo "  Original Revision: $ORIGINAL_REV"
echo ""

echo "=== STEP 2: Update rule using PUT endpoint ==="
echo "Changing description from '$ORIGINAL_DESC' to 'UPDATED VIA CURL PUT'"
echo ""

PUT_RESPONSE=$(curl -s -X PUT \
  -H "Authorization: APIKey $GCORE_API_KEY" \
  -H "Content-Type: application/json" \
  "https://api.gcore.com/cloud/v1/securitygrouprules/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID/$RULE_ID" \
  -d '{
    "security_group_id": "'$SG_ID'",
    "direction": "ingress",
    "ethertype": "IPv4",
    "protocol": "tcp",
    "port_range_min": 80,
    "port_range_max": 80,
    "remote_ip_prefix": "0.0.0.0/0",
    "description": "UPDATED VIA CURL PUT"
  }')

echo "$PUT_RESPONSE" > /tmp/rule_put_response.json
echo "PUT Response:"
echo "$PUT_RESPONSE" | jq '{id, description, direction, protocol, port_range_min, port_range_max, security_group_id, revision_number, updated_at}'
echo ""

NEW_ID=$(echo "$PUT_RESPONSE" | jq -r '.id')
NEW_DESC=$(echo "$PUT_RESPONSE" | jq -r '.description')
NEW_REV=$(echo "$PUT_RESPONSE" | jq -r '.revision_number')
echo "  New ID: $NEW_ID"
echo "  New Description: $NEW_DESC"
echo "  New Revision: $NEW_REV"
echo ""

echo "=== STEP 3: Verify rule via security group GET ==="
sleep 1
curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
  "https://api.gcore.com/cloud/v1/securitygroups/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID/$SG_ID" \
  > /tmp/sg_after_update.json

echo "Searching for ORIGINAL_ID ($ORIGINAL_ID) in security group rules:"
cat /tmp/sg_after_update.json | jq ".security_group_rules[] | select(.id == \"$ORIGINAL_ID\")" > /tmp/search_original.json
if [ -s /tmp/search_original.json ]; then
  echo "  ✓ FOUND - Rule with original ID still exists"
  cat /tmp/search_original.json | jq '{id, description}'
else
  echo "  ✗ NOT FOUND - Original ID no longer exists in security group"
fi
echo ""

echo "Searching for NEW_ID ($NEW_ID) in security group rules:"
cat /tmp/sg_after_update.json | jq ".security_group_rules[] | select(.id == \"$NEW_ID\")" > /tmp/search_new.json
if [ -s /tmp/search_new.json ]; then
  echo "  ✓ FOUND - Rule with new ID exists"
  cat /tmp/search_new.json | jq '{id, description}'
else
  echo "  ✗ NOT FOUND - New ID not found in security group"
fi
echo ""

echo "========================================================================"
echo "ANALYSIS"
echo "========================================================================"
echo ""

if [ "$ORIGINAL_ID" = "$NEW_ID" ]; then
  echo "✓ ID PRESERVED: The rule ID remained the same ($ORIGINAL_ID)"
  echo "  This indicates an in-place UPDATE operation"
else
  echo "✗ ID CHANGED: The rule ID changed!"
  echo "  Before: $ORIGINAL_ID"
  echo "  After:  $NEW_ID"
  echo "  This indicates a REPLACE operation (delete old + create new)"
fi
echo ""

if [ "$ORIGINAL_REV" != "$NEW_REV" ]; then
  echo "✓ Revision number incremented: $ORIGINAL_REV → $NEW_REV"
else
  echo "  Revision number unchanged: $ORIGINAL_REV"
fi
echo ""

echo "Description changed: '$ORIGINAL_DESC' → '$NEW_DESC'"
echo ""

echo "========================================================================"
echo "Files saved:"
echo "  /tmp/rule_before.json - Rule state before PUT"
echo "  /tmp/rule_put_response.json - PUT endpoint response"
echo "  /tmp/sg_after_update.json - Security group state after update"
echo "========================================================================"
