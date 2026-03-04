#!/bin/bash
set -e

echo "=== Cleanup Multi-Security Group Test ==="

if [ ! -f ../.env ]; then
    echo "Error: ../.env file not found"
    exit 1
fi

source ../.env
export GCORE_API_KEY GCORE_CLOUD_PROJECT_ID GCORE_CLOUD_REGION_ID

echo "Destroying all resources..."
terraform destroy -auto-approve

echo ""
echo "✓ All resources destroyed"
