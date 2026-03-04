#!/bin/bash
source /Users/user/repos/gcore-terraform/.env

echo "=== Available Volumes ==="
curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
  "https://api.gcore.com/cloud/v1/volumes/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID" | \
  jq -r 'if type == "array" then .[] | select(.status == "available") | "\(.id) | \(.name) | \(.size)GB | bootable=\(.bootable)" else "No volumes or unexpected format" end' 2>/dev/null | head -10

echo ""
echo "=== Available Images (public Ubuntu) ==="
curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
  "https://api.gcore.com/cloud/v1/images/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID?visibility=public" | \
  jq -r '.results[] | select(.os_distro == "ubuntu") | "\(.id) | \(.display_name // .name) | \(.min_disk)GB min"' 2>/dev/null | head -5

echo ""
echo "=== Networks ==="
curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
  "https://api.gcore.com/cloud/v1/networks/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID" | \
  jq -r '.results[] | "\(.id) | \(.name) | external=\(.external)"' 2>/dev/null | head -5
