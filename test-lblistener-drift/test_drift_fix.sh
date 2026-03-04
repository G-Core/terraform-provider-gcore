#!/bin/bash
set -e

echo "=== Testing Drift Fix ==="
echo ""

# Source credentials
source ../.env

echo "Step 1: First terraform apply (creating resources)..."
terraform apply -auto-approve

echo ""
echo "=== First apply complete! ==="
echo ""
echo "Step 2: Second terraform apply (checking for drift)..."
echo ""

terraform apply

echo ""
echo "=== Test complete ==="
