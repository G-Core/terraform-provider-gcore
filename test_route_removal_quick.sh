#!/bin/bash
set -e

# Test just the route removal on existing router
cd test-route-fix-*

echo "Testing route removal with debug output..."
echo "Current router state:"
terraform show | grep -A5 "router.router"

echo ""
echo "Applying route removal..."
terraform apply -auto-approve -var="include_routes=false" 

echo ""
echo "Checking result..."
terraform show | grep -A5 "routes ="
