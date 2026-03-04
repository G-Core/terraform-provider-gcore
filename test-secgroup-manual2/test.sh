#!/bin/bash
set -e

echo "=== Multi-Security Group Test ==="
echo "This test creates 2 security groups with rules added via loops"
echo ""

# Check for .env file
if [ ! -f ../.env ]; then
    echo "Error: ../.env file not found"
    echo "Please create ../.env with:"
    echo "  export GCORE_API_KEY='your-api-key'"
    echo "  export GCORE_CLOUD_PROJECT_ID=your-project-id"
    echo "  export GCORE_CLOUD_REGION_ID=your-region-id"
    exit 1
fi

source ../.env
export GCORE_API_KEY GCORE_CLOUD_PROJECT_ID GCORE_CLOUD_REGION_ID

echo "Step 1: Initialize Terraform"
terraform init
echo ""

echo "Step 2: Create 2 security groups + 6 rules (3 web + 3 database)"
terraform apply -auto-approve
echo ""

echo "Step 3: Verify - check outputs"
terraform output
echo ""

echo "Step 4: Check for drift"
if terraform plan -detailed-exitcode; then
    echo "✓ No drift detected"
else
    exit_code=$?
    if [ $exit_code -eq 2 ]; then
        echo "✗ Drift detected!"
        exit 1
    fi
fi
echo ""

echo "Step 5: Update web security group name"
sed -i.bak 's/name        = "test-web-tier"/name        = "test-web-tier-updated"/' main.tf
terraform plan
echo ""
read -p "Apply the web tier name change? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    terraform apply -auto-approve
    echo ""
    echo "Step 6: Verify database rules were NOT replaced"
    terraform show | grep -A 5 "gcore_cloud_security_group_rule.database_rules"
fi

echo ""
echo "=== Test Complete ==="
echo "Summary:"
terraform output summary
