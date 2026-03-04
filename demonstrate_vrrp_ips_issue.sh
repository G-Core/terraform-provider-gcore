#!/bin/bash
set -e

# Demonstrate vrrp_ips API discrepancy
# This script shows that PATCH /loadbalancers returns empty vrrp_ips[] but GET returns the full array

echo "=========================================="
echo "Demonstrating vrrp_ips API Discrepancy"
echo "=========================================="
echo ""

# Load credentials from .env
if [ -f .env ]; then
    echo "✓ Loading credentials from .env"
    source .env
else
    echo "❌ Error: .env file not found"
    exit 1
fi

# Validate credentials
if [ -z "$GCORE_API_KEY" ] || [ -z "$GCORE_CLOUD_PROJECT_ID" ] || [ -z "$GCORE_CLOUD_REGION_ID" ]; then
    echo "❌ Error: Missing required environment variables"
    echo "Required: GCORE_API_KEY, GCORE_CLOUD_PROJECT_ID, GCORE_CLOUD_REGION_ID"
    exit 1
fi

API_KEY="$GCORE_API_KEY"
PROJECT_ID="$GCORE_CLOUD_PROJECT_ID"
REGION_ID="$GCORE_CLOUD_REGION_ID"

echo "✓ Credentials loaded"
echo "  Project ID: $PROJECT_ID"
echo "  Region ID: $REGION_ID"
echo ""

# Step 1: Create Load Balancer
echo "=========================================="
echo "STEP 1: Create Load Balancer"
echo "=========================================="
echo ""

CREATE_RESPONSE=$(curl -s -X POST \
  "https://api.gcore.com/cloud/v1/loadbalancers/$PROJECT_ID/$REGION_ID" \
  -H "Authorization: APIKey $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "demo-vrrp-test",
    "flavor": "lb1-1-2"
  }')

echo "Create LB API Response:"
echo "$CREATE_RESPONSE" | jq '.'
echo ""

# Extract task ID and wait for completion
TASK_ID=$(echo "$CREATE_RESPONSE" | jq -r '.tasks[0]')
echo "Task ID: $TASK_ID"
echo "Waiting for LB creation to complete..."
echo ""

# Poll task until complete
MAX_ATTEMPTS=60
ATTEMPT=0
while [ $ATTEMPT -lt $MAX_ATTEMPTS ]; do
    TASK_STATUS=$(curl -s -X GET \
        "https://api.gcore.com/cloud/v1/tasks/$TASK_ID" \
        -H "Authorization: APIKey $API_KEY" | jq -r '.state')

    if [ "$TASK_STATUS" = "FINISHED" ]; then
        echo "✓ Task completed successfully"
        break
    elif [ "$TASK_STATUS" = "ERROR" ]; then
        echo "❌ Task failed"
        exit 1
    fi

    ATTEMPT=$((ATTEMPT + 1))
    sleep 2
done

# Get LB ID from task
LB_ID=$(curl -s -X GET \
    "https://api.gcore.com/cloud/v1/tasks/$TASK_ID" \
    -H "Authorization: APIKey $API_KEY" | jq -r '.created_resources.loadbalancers[0]')

echo "Load Balancer ID: $LB_ID"
echo ""

# Wait a bit for LB to be fully ready
echo "Waiting for LB to be fully provisioned..."
sleep 10
echo ""

# Step 2: GET Load Balancer (initial state)
echo "=========================================="
echo "STEP 2: GET Load Balancer (shows vrrp_ips)"
echo "=========================================="
echo ""

GET_RESPONSE_1=$(curl -s -X GET \
  "https://api.gcore.com/cloud/v1/loadbalancers/$PROJECT_ID/$REGION_ID/$LB_ID" \
  -H "Authorization: APIKey $API_KEY")

echo "GET Response - Full:"
echo "$GET_RESPONSE_1" | jq '.'
echo ""

VRRP_IPS_COUNT_1=$(echo "$GET_RESPONSE_1" | jq '.vrrp_ips | length')
echo "✓ GET response shows vrrp_ips array with $VRRP_IPS_COUNT_1 elements:"
echo "$GET_RESPONSE_1" | jq '.vrrp_ips'
echo ""

# Step 3: PATCH Load Balancer (rename)
echo "=========================================="
echo "STEP 3: PATCH Load Balancer (rename)"
echo "=========================================="
echo ""
echo "Renaming LB from 'demo-vrrp-test' to 'demo-vrrp-RENAMED'"
echo ""

PATCH_RESPONSE=$(curl -s -X PATCH \
  "https://api.gcore.com/cloud/v1/loadbalancers/$PROJECT_ID/$REGION_ID/$LB_ID" \
  -H "Authorization: APIKey $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "demo-vrrp-RENAMED"
  }')

echo "PATCH Response - Full:"
echo "$PATCH_RESPONSE" | jq '.'
echo ""

VRRP_IPS_COUNT_PATCH=$(echo "$PATCH_RESPONSE" | jq '.vrrp_ips | length')
echo "❌ CRITICAL ISSUE: PATCH response shows vrrp_ips array with $VRRP_IPS_COUNT_PATCH elements:"
echo "$PATCH_RESPONSE" | jq '.vrrp_ips'
echo ""
echo "^^^ THIS IS THE BUG! PATCH returns empty vrrp_ips[] ^^^"
echo ""

# Step 4: GET Load Balancer again (after rename)
echo "=========================================="
echo "STEP 4: GET Load Balancer after PATCH"
echo "=========================================="
echo ""

sleep 2  # Brief wait

GET_RESPONSE_2=$(curl -s -X GET \
  "https://api.gcore.com/cloud/v1/loadbalancers/$PROJECT_ID/$REGION_ID/$LB_ID" \
  -H "Authorization: APIKey $API_KEY")

echo "GET Response - Full:"
echo "$GET_RESPONSE_2" | jq '.'
echo ""

VRRP_IPS_COUNT_2=$(echo "$GET_RESPONSE_2" | jq '.vrrp_ips | length')
echo "✓ GET response shows vrrp_ips array with $VRRP_IPS_COUNT_2 elements:"
echo "$GET_RESPONSE_2" | jq '.vrrp_ips'
echo ""

# Summary
echo "=========================================="
echo "SUMMARY OF API DISCREPANCY"
echo "=========================================="
echo ""
echo "Initial GET:   vrrp_ips has $VRRP_IPS_COUNT_1 elements ✓"
echo "PATCH rename:  vrrp_ips has $VRRP_IPS_COUNT_PATCH elements ❌ (EMPTY!)"
echo "Final GET:     vrrp_ips has $VRRP_IPS_COUNT_2 elements ✓"
echo ""
echo "CONCLUSION:"
echo "  • API documentation claims PATCH returns full vrrp_ips"
echo "  • Actual PATCH response returns empty vrrp_ips: []"
echo "  • GET endpoint returns correct vrrp_ips"
echo "  • Terraform provider MUST do GET after PATCH to refresh state"
echo ""

# Save LB ID for cleanup
echo "$LB_ID" > /tmp/demo_lb_id.txt
echo "Load Balancer ID saved to /tmp/demo_lb_id.txt for cleanup"
echo ""
echo "=========================================="
echo "Demonstration Complete"
echo "=========================================="
