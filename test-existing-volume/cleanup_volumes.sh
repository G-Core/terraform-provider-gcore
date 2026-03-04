#!/bin/bash
source /Users/user/repos/gcore-terraform/.env

echo "=== All Volumes ==="
curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
  "https://api.gcore.com/cloud/v1/volumes/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID" | \
  jq -r '.[] | "\(.id) | \(.name) | status=\(.status) | attached_to=\(.attachments[0].instance_id // "none")"'

echo ""
echo "=== Volumes to clean (available, test-related) ==="
VOLUMES_TO_DELETE=$(curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
  "https://api.gcore.com/cloud/v1/volumes/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID" | \
  jq -r '.[] | select(.status == "available" and (.name | test("test-"; "i"))) | .id')

echo "$VOLUMES_TO_DELETE"

echo ""
read -p "Delete these volumes? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    for vol_id in $VOLUMES_TO_DELETE; do
        echo "Deleting $vol_id..."
        curl -s -X DELETE -H "Authorization: APIKey $GCORE_API_KEY" \
          "https://api.gcore.com/cloud/v1/volumes/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID/$vol_id"
        echo ""
    done
fi
