#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Load credentials
if [ -f ../.env ]; then
    source ../.env
else
    echo "ERROR: .env not found"
    exit 1
fi

export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/.terraformrc"
export TF_VAR_api_key="$GCORE_API_KEY"

echo "=== Starting cleanup ==="

# Get volume IDs before destroy
echo "Getting volume IDs before destroy..."
VOLUME_IDS=$(terraform show -json 2>/dev/null | python3 -c "
import sys,json
try:
    data=json.load(sys.stdin)
    resources=data.get('values',{}).get('root_module',{}).get('resources',[])
    for r in resources:
        if r.get('type')=='gcore_cloud_volume':
            print(r['values']['id'])
except:
    pass
" 2>/dev/null)

echo "Volume IDs: $VOLUME_IDS"

echo "Destroying all resources..."
terraform destroy -auto-approve

echo ""
echo "=== Test 20 Verification: Volume Deletion Behavior ==="
echo "Checking if volumes still exist after instance destroy..."
sleep 5

for vol_id in $VOLUME_IDS; do
    if [ -n "$vol_id" ]; then
        RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -H "Authorization: APIKey $GCORE_API_KEY" \
            "https://api.gcore.com/cloud/v1/volumes/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID/$vol_id")
        HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE:" | cut -d':' -f2)

        if [ "$HTTP_CODE" = "200" ]; then
            echo "Volume $vol_id still EXISTS after instance destroy (expected with new provider)"
        elif [ "$HTTP_CODE" = "404" ]; then
            echo "Volume $vol_id was DELETED with instance (old provider behavior)"
        else
            echo "Volume $vol_id check returned HTTP $HTTP_CODE"
        fi
    fi
done

echo ""
echo "=== Cleanup complete ==="
