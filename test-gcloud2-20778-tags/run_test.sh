#!/bin/bash
set -e

echo "=========================================="
echo "GCLOUD2-20778 Drift Reproduction Test"
echo "Testing: Drift after adding LB pool"
echo "=========================================="
echo ""

# Setup environment
if [ -f ../.env ]; then
    source ../.env
    echo "✓ Loaded credentials from .env"
else
    echo "❌ .env file not found"
    exit 1
fi

export TF_CLI_CONFIG_FILE=/Users/user/repos/gcore-terraform/.terraformrc

# Initialize
echo ""
echo "Step 1: Initialize Terraform..."
terraform init

# Step 2: Create LB + Listener with user_list (NO POOL YET)
echo ""
echo "Step 2: Creating LB + Listener with user_list..."
terraform apply -auto-approve

# Verify no drift after initial creation
echo ""
echo "Step 3: Checking for drift after initial creation..."
if terraform plan -detailed-exitcode; then
    echo "✓ No drift detected after initial creation"
else
    echo "⚠ Unexpected drift detected after initial creation"
fi

# Step 4: Add the pool by uncommenting it
echo ""
echo "Step 4: Adding LB pool to configuration..."
sed -i.bak 's/^# resource "gcore_cloud_load_balancer_pool"/resource "gcore_cloud_load_balancer_pool"/' main.tf
sed -i.bak 's/^#   /  /' main.tf
sed -i.bak 's/^# }/}/' main.tf

# Apply the change (add pool)
echo ""
echo "Step 5: Applying change (adding pool)..."
terraform apply -auto-approve

# Step 6: THIS IS THE CRITICAL TEST - Check for drift after adding pool
echo ""
echo "=========================================="
echo "CRITICAL TEST: Checking for drift after adding pool"
echo "=========================================="
if terraform plan -detailed-exitcode; then
    echo ""
    echo "✅ PASS: No drift detected"
    echo "The bug is FIXED or cannot be reproduced"
    exit 0
else
    EXIT_CODE=$?
    echo ""
    echo "❌ FAIL: Drift detected!"
    echo "Exit code: $EXIT_CODE (2 = changes detected)"
    echo ""
    echo "Saving drift output..."
    terraform plan -no-color > drift_output.txt
    echo "Drift details saved to drift_output.txt"
    echo ""
    echo "BUG REPRODUCED: After adding pool, Terraform wants to update listener"
    exit 1
fi
