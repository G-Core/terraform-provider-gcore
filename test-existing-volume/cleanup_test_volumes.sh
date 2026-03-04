#!/bin/bash
source /Users/user/repos/gcore-terraform/.env

echo "=== Deleting available test volumes ==="

# Get available test volumes and delete them
curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
  "https://api.gcore.com/cloud/v1/volumes/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID" | \
  jq -r '.results[] | select(.status == "available" and (.name | test("test-"; "i"))) | .id' | \
while read vol_id; do
    if [ -n "$vol_id" ]; then
        echo "Deleting $vol_id..."
        curl -s -X DELETE -H "Authorization: APIKey $GCORE_API_KEY" \
          "https://api.gcore.com/cloud/v1/volumes/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID/$vol_id"
        echo ""
        sleep 2
    fi
done

echo ""
echo "=== Remaining volumes ==="
curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
  "https://api.gcore.com/cloud/v1/volumes/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID" | \
  jq -r '.results[] | "\(.id) | \(.name) | status=\(.status)"'
