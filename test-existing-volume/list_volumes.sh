#!/bin/bash
source /Users/user/repos/gcore-terraform/.env

echo "=== All Volumes ==="
curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
  "https://api.gcore.com/cloud/v1/volumes/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID" | \
  jq -r '.results[] | "\(.id) | \(.name) | status=\(.status) | bootable=\(.bootable)"'

echo ""
echo "=== Available (unattached) volumes ==="
curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
  "https://api.gcore.com/cloud/v1/volumes/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID" | \
  jq -r '.results[] | select(.status == "available") | "\(.id) | \(.name) | bootable=\(.bootable)"'
