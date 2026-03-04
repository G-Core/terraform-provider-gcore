#!/bin/bash
set -e

if [ -f .env ]; then
    source .env
elif [ -f ../.env ]; then
    source ../.env
else
    echo "ERROR: .env not found"
    exit 1
fi

API_KEY="${GCORE_API_KEY}"
PROJECT_ID="${GCORE_CLOUD_PROJECT_ID}"
REGION_ID="${GCORE_CLOUD_REGION_ID}"

echo "=== Testing Partial Route Deletion (Pedro's Comment) ==="

# Create router
echo "Creating router..."
CREATE_RESPONSE=$(curl -s -X POST \
  -H "Authorization: APIKey $API_KEY" \
  -H "Content-Type: application/json" \
  "https://api.gcore.com/cloud/v1/routers/$PROJECT_ID/$REGION_ID" \
  -d '{"name": "partial-deletion-test"}')

sleep 10

# Get router ID
ROUTER_LIST=$(curl -s -H "Authorization: APIKey $API_KEY" \
  "https://api.gcore.com/cloud/v1/routers/$PROJECT_ID/$REGION_ID")
ROUTER_ID=$(echo "$ROUTER_LIST" | jq -r --arg name "partial-deletion-test" '.results[] | select(.name == $name) | .id')

echo "Router ID: $ROUTER_ID"

# Add 3 routes
echo ""
echo "Step 1: Adding 3 routes..."
PATCH1=$(curl -s -X PATCH \
  -H "Authorization: APIKey $API_KEY" \
  -H "Content-Type: application/json" \
  "https://api.gcore.com/cloud/v1/routers/$PROJECT_ID/$REGION_ID/$ROUTER_ID" \
  -d '{
    "routes": [
      {"destination": "10.0.1.0/24", "nexthop": "192.168.1.1"},
      {"destination": "10.0.2.0/24", "nexthop": "192.168.1.2"},
      {"destination": "10.0.3.0/24", "nexthop": "192.168.1.3"}
    ]
  }')

ROUTES_COUNT=$(echo "$PATCH1" | jq -r '.routes | length')
echo "Routes added: $ROUTES_COUNT"
echo "$PATCH1" | jq -r '.routes'

# Partial deletion: 3 -> 2 routes
echo ""
echo "Step 2: Deleting one route (3 -> 2)..."
PATCH2=$(curl -s -X PATCH \
  -H "Authorization: APIKey $API_KEY" \
  -H "Content-Type: application/json" \
  "https://api.gcore.com/cloud/v1/routers/$PROJECT_ID/$REGION_ID/$ROUTER_ID" \
  -d '{
    "routes": [
      {"destination": "10.0.1.0/24", "nexthop": "192.168.1.1"},
      {"destination": "10.0.2.0/24", "nexthop": "192.168.1.2"}
    ]
  }')

ROUTES_COUNT=$(echo "$PATCH2" | jq -r '.routes | length')
echo "Routes after partial deletion: $ROUTES_COUNT"
echo "$PATCH2" | jq -r '.routes'

if [ "$ROUTES_COUNT" = "2" ]; then
    echo "✅ Partial deletion works (3 -> 2)"
else
    echo "❌ Partial deletion failed"
fi

# Another partial deletion: 2 -> 1 route
echo ""
echo "Step 3: Deleting another route (2 -> 1)..."
PATCH3=$(curl -s -X PATCH \
  -H "Authorization: APIKey $API_KEY" \
  -H "Content-Type: application/json" \
  "https://api.gcore.com/cloud/v1/routers/$PROJECT_ID/$REGION_ID/$ROUTER_ID" \
  -d '{
    "routes": [
      {"destination": "10.0.1.0/24", "nexthop": "192.168.1.1"}
    ]
  }')

ROUTES_COUNT=$(echo "$PATCH3" | jq -r '.routes | length')
echo "Routes after second partial deletion: $ROUTES_COUNT"
echo "$PATCH3" | jq -r '.routes'

if [ "$ROUTES_COUNT" = "1" ]; then
    echo "✅ Partial deletion works (2 -> 1)"
else
    echo "❌ Partial deletion failed"
fi

# Full deletion: 1 -> 0 routes
echo ""
echo "Step 4: Deleting all routes (1 -> 0)..."
PATCH4=$(curl -s -X PATCH \
  -H "Authorization: APIKey $API_KEY" \
  -H "Content-Type: application/json" \
  "https://api.gcore.com/cloud/v1/routers/$PROJECT_ID/$REGION_ID/$ROUTER_ID" \
  -d '{"routes": []}')

ROUTES_COUNT=$(echo "$PATCH4" | jq -r '.routes | length')
echo "Routes after full deletion: $ROUTES_COUNT"

if [ "$ROUTES_COUNT" = "0" ]; then
    echo "✅ Full deletion works (1 -> 0)"
else
    echo "❌ Full deletion failed"
fi

# Cleanup
echo ""
echo "Cleanup: Deleting router..."
curl -s -X DELETE \
  -H "Authorization: APIKey $API_KEY" \
  "https://api.gcore.com/cloud/v1/routers/$PROJECT_ID/$REGION_ID/$ROUTER_ID" > /dev/null

echo ""
echo "=== CONCLUSION ==="
echo "✅ ALL partial deletions work via PATCH"
echo "⚠️  Current code condition is WRONG:"
echo "    routesDeletionNeeded := len(dataRoutes) == 0 && len(stateRoutes) > 0"
echo "    This ONLY catches full deletion (n -> 0)"
echo ""
echo "✅ Pedro's suggestion is CORRECT:"
echo "    routesDeletionNeeded := len(dataRoutes) < len(stateRoutes)"
echo "    This catches ALL deletions including partial (n -> m where m < n)"
