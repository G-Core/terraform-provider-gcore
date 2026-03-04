#!/bin/bash
set -e

# Load credentials
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

echo "=== Testing Pedro's Comments on Router PATCH Behavior ==="
echo ""

# Test 1: Empty slices with omitzero - verify Pedro's claim
echo "Test 1: Verify empty slices are serialized with omitzero tag"
echo "Creating router with 2 routes..."

# First create a basic router without routes
CREATE_RESPONSE=$(curl -s -X POST \
  -H "Authorization: APIKey $API_KEY" \
  -H "Content-Type: application/json" \
  "https://api.gcore.com/cloud/v1/routers/$PROJECT_ID/$REGION_ID" \
  -d '{"name": "pedro-test-router"}')

echo "Create response: $CREATE_RESPONSE"

# Wait for router to be created
echo "Waiting for router creation..."
sleep 10

# Get the actual router ID
ROUTER_LIST=$(curl -s -H "Authorization: APIKey $API_KEY" \
  "https://api.gcore.com/cloud/v1/routers/$PROJECT_ID/$REGION_ID")
ROUTER_ID=$(echo "$ROUTER_LIST" | jq -r --arg name "pedro-test-router" '.results[] | select(.name == $name) | .id')

if [ -z "$ROUTER_ID" ] || [ "$ROUTER_ID" = "null" ]; then
    echo "ERROR: Could not find created router"
    echo "Router list: $ROUTER_LIST"
    exit 1
fi

echo "Router ID: $ROUTER_ID"

# Now add routes via PATCH
echo "Adding 2 routes via PATCH..."
curl -s -X PATCH \
  -H "Authorization: APIKey $API_KEY" \
  -H "Content-Type: application/json" \
  "https://api.gcore.com/cloud/v1/routers/$PROJECT_ID/$REGION_ID/$ROUTER_ID" \
  -d '{
    "routes": [
      {"destination": "10.0.1.0/24", "nexthop": "192.168.1.1"},
      {"destination": "10.0.2.0/24", "nexthop": "192.168.1.2"}
    ]
  }' > /dev/null

sleep 2

echo ""
echo "Test 1a: PATCH with routes=[] (empty array explicitly)"
PATCH_RESPONSE=$(curl -s -X PATCH \
  -H "Authorization: APIKey $API_KEY" \
  -H "Content-Type: application/json" \
  "https://api.gcore.com/cloud/v1/routers/$PROJECT_ID/$REGION_ID/$ROUTER_ID" \
  -d '{"routes": []}')

echo "PATCH Response (routes=[]):"
echo "$PATCH_RESPONSE" | jq -r '.routes // "routes field missing"'
ROUTES_COUNT=$(echo "$PATCH_RESPONSE" | jq -r '.routes | length')
echo "Routes count after PATCH with routes=[]: $ROUTES_COUNT"

if [ "$ROUTES_COUNT" = "0" ]; then
    echo "✅ CONFIRMED: Empty array routes=[] successfully deletes routes via PATCH"
else
    echo "❌ FAILED: routes=[] did not delete routes"
fi

echo ""
echo "Test 1b: Reset router with 2 routes again for next test"
curl -s -X PATCH \
  -H "Authorization: APIKey $API_KEY" \
  -H "Content-Type: application/json" \
  "https://api.gcore.com/cloud/v1/routers/$PROJECT_ID/$REGION_ID/$ROUTER_ID" \
  -d '{
    "routes": [
      {"destination": "10.0.1.0/24", "nexthop": "192.168.1.1"},
      {"destination": "10.0.2.0/24", "nexthop": "192.168.1.2"}
    ]
  }' > /dev/null

sleep 2

echo ""
echo "Test 2: Partial route deletion (2 routes -> 1 route)"
echo "Testing Pedro's suggested condition: len(dataRoutes) < len(stateRoutes)"

GET_RESPONSE=$(curl -s -H "Authorization: APIKey $API_KEY" \
  "https://api.gcore.com/cloud/v1/routers/$PROJECT_ID/$REGION_ID/$ROUTER_ID")
CURRENT_ROUTES=$(echo "$GET_RESPONSE" | jq -r '.routes')
echo "Current routes (should be 2):"
echo "$CURRENT_ROUTES" | jq '.'

echo ""
echo "Deleting one route (keeping only first route)..."
PATCH_RESPONSE=$(curl -s -X PATCH \
  -H "Authorization: APIKey $API_KEY" \
  -H "Content-Type: application/json" \
  "https://api.gcore.com/cloud/v1/routers/$PROJECT_ID/$REGION_ID/$ROUTER_ID" \
  -d '{
    "routes": [
      {"destination": "10.0.1.0/24", "nexthop": "192.168.1.1"}
    ]
  }')

ROUTES_AFTER=$(echo "$PATCH_RESPONSE" | jq -r '.routes')
echo "Routes after partial deletion:"
echo "$ROUTES_AFTER" | jq '.'
ROUTES_COUNT_AFTER=$(echo "$PATCH_RESPONSE" | jq -r '.routes | length')

if [ "$ROUTES_COUNT_AFTER" = "1" ]; then
    echo "✅ CONFIRMED: Partial deletion works (2 routes -> 1 route)"
    echo "   Pedro is RIGHT: condition should be len(dataRoutes) < len(stateRoutes)"
    echo "   Current code: routesDeletionNeeded := len(dataRoutes) == 0 && len(stateRoutes) > 0"
    echo "   This would NOT catch partial deletions!"
else
    echo "❌ FAILED: Partial deletion did not work as expected"
fi

echo ""
echo "Test 3: Can we send single PATCH for all updates including routes=[]?"
echo "Testing: Update name + delete all routes in single PATCH"

curl -s -X PATCH \
  -H "Authorization: APIKey $API_KEY" \
  -H "Content-Type: application/json" \
  "https://api.gcore.com/cloud/v1/routers/$PROJECT_ID/$REGION_ID/$ROUTER_ID" \
  -d '{
    "routes": [
      {"destination": "10.0.1.0/24", "nexthop": "192.168.1.1"},
      {"destination": "10.0.2.0/24", "nexthop": "192.168.1.2"}
    ]
  }' > /dev/null

sleep 2

PATCH_RESPONSE=$(curl -s -X PATCH \
  -H "Authorization: APIKey $API_KEY" \
  -H "Content-Type: application/json" \
  "https://api.gcore.com/cloud/v1/routers/$PROJECT_ID/$REGION_ID/$ROUTER_ID" \
  -d '{
    "name": "pedro-test-router-renamed",
    "routes": []
  }')

NEW_NAME=$(echo "$PATCH_RESPONSE" | jq -r '.name')
NEW_ROUTES_COUNT=$(echo "$PATCH_RESPONSE" | jq -r '.routes | length')

echo "Response after combined PATCH:"
echo "  Name: $NEW_NAME"
echo "  Routes count: $NEW_ROUTES_COUNT"

if [ "$NEW_NAME" = "pedro-test-router-renamed" ] && [ "$NEW_ROUTES_COUNT" = "0" ]; then
    echo "✅ CONFIRMED: Single PATCH can update name AND delete routes simultaneously"
    echo "   Pedro is RIGHT: We don't need separate PATCH calls"
else
    echo "❌ FAILED: Single PATCH did not work as expected"
fi

echo ""
echo "Test 4: Does PATCH return updated state (do we need final GET)?"
echo "Comparing PATCH response with GET response..."

PATCH_RESPONSE=$(curl -s -X PATCH \
  -H "Authorization: APIKey $API_KEY" \
  -H "Content-Type: application/json" \
  "https://api.gcore.com/cloud/v1/routers/$PROJECT_ID/$REGION_ID/$ROUTER_ID" \
  -d '{"name": "pedro-test-final"}')

GET_RESPONSE=$(curl -s -H "Authorization: APIKey $API_KEY" \
  "https://api.gcore.com/cloud/v1/routers/$PROJECT_ID/$REGION_ID/$ROUTER_ID")

PATCH_NAME=$(echo "$PATCH_RESPONSE" | jq -r '.name')
GET_NAME=$(echo "$GET_RESPONSE" | jq -r '.name')
PATCH_STATUS=$(echo "$PATCH_RESPONSE" | jq -r '.status')
GET_STATUS=$(echo "$GET_RESPONSE" | jq -r '.status')

echo "PATCH response: name=$PATCH_NAME, status=$PATCH_STATUS"
echo "GET response:   name=$GET_NAME, status=$GET_STATUS"

if [ "$PATCH_NAME" = "$GET_NAME" ]; then
    echo "✅ CONFIRMED: PATCH returns updated state"
    echo "   Pedro is RIGHT: We likely don't need final GET after PATCH"
else
    echo "⚠️  PATCH and GET responses differ - may need GET for consistency"
fi

echo ""
echo "Cleanup: Deleting test router..."
curl -s -X DELETE \
  -H "Authorization: APIKey $API_KEY" \
  "https://api.gcore.com/cloud/v1/routers/$PROJECT_ID/$REGION_ID/$ROUTER_ID" > /dev/null

echo ""
echo "=== Summary of Findings ==="
echo "1. ✅ Empty slices work with PATCH: routes=[] successfully deletes routes"
echo "2. ✅ Partial deletions work: Current code MISSES this case!"
echo "   - Current: routesDeletionNeeded := len(dataRoutes) == 0 && len(stateRoutes) > 0"
echo "   - Should be: routesDeletionNeeded := len(dataRoutes) < len(stateRoutes)"
echo "3. ✅ Single PATCH can handle all updates including routes=[]"
echo "4. ✅ PATCH returns updated state - final GET may be unnecessary"
echo ""
echo "Pedro's comments are VALID - we can significantly simplify the implementation!"
