#!/bin/bash

# Load credentials from parent .env file
source ../.env

# Remove quotes from API_KEY if present
API_KEY="${GCORE_API_KEY//\'/}"
PROJECT_ID="${GCORE_CLOUD_PROJECT_ID}"
REGION_ID="${GCORE_CLOUD_REGION_ID}"

echo "=== Listing Load Balancers ==="
echo "Project: $PROJECT_ID"
echo "Region: $REGION_ID"
echo ""

# List all load balancers
LB_LIST=$(curl -s -H "Authorization: APIKey $API_KEY" \
  "https://api.gcore.com/cloud/v1/loadbalancers/$PROJECT_ID/$REGION_ID")

echo "$LB_LIST" | python3 -m json.tool

# Extract load balancer IDs
LB_IDS=$(echo "$LB_LIST" | python3 -c "import sys, json; data = json.load(sys.stdin); print('\n'.join([lb['id'] for lb in data.get('results', [])]))" 2>/dev/null)

if [ -z "$LB_IDS" ]; then
    echo "No load balancers found or error parsing response"
    exit 0
fi

echo ""
echo "=== Found Load Balancers ==="
echo "$LB_IDS"
echo ""
echo "=== Deleting Load Balancers ==="

# Delete each load balancer
for LB_ID in $LB_IDS; do
    echo "Deleting load balancer: $LB_ID"
    RESPONSE=$(curl -s -X DELETE -H "Authorization: APIKey $API_KEY" \
      "https://api.gcore.com/cloud/v1/loadbalancers/$PROJECT_ID/$REGION_ID/$LB_ID")

    echo "Response: $RESPONSE"
    echo ""
done

echo "=== Cleanup Complete ==="
