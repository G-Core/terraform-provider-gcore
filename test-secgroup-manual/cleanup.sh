#!/bin/bash

echo "Cleaning up test resources..."

# Load environment
cd /Users/user/repos/gcore-terraform
set -o allexport
source .env
set +o allexport
cd test-secgroup-manual

# Destroy
terraform destroy -auto-approve

echo "✓ Cleanup complete!"
