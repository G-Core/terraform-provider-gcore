#!/bin/bash
set -e

echo "=========================================="
echo "Testing: Tags Inconsistency Error"
echo "GCLOUD2-20778 - tags_v2 issue"
echo "=========================================="
echo ""

# Setup
if [ -f ../.env ]; then
    source ../.env
    echo "✓ Loaded credentials"
else
    echo "❌ .env not found"
    exit 1
fi

export TF_CLI_CONFIG_FILE=/Users/user/repos/gcore-terraform/.terraformrc

echo ""
echo "Step 1: Creating LB WITHOUT tags..."
terraform apply -auto-approve

echo ""
echo "Step 2: Verifying LB created successfully..."
terraform show

echo ""
echo "Step 3: Adding tags to configuration..."
# Uncomment the tags block
sed -i.bak 's/  # tags = {/  tags = {/' main.tf
sed -i.bak 's/  #   "qa" = "load-balancer"/    "qa" = "load-balancer"/' main.tf
sed -i.bak 's/  # }/  }/' main.tf

echo "Updated configuration:"
cat main.tf | grep -A 3 "tags = {"

echo ""
echo "=========================================="
echo "Step 4: CRITICAL TEST - Applying tags"
echo "=========================================="
if terraform apply -auto-approve 2>&1 | tee apply_with_tags.log; then
    echo ""
    echo "✅ SUCCESS: Tags applied without error"
    echo "Bug NOT reproduced - tags work correctly"
    exit 0
else
    echo ""
    echo "❌ FAILED: Error occurred when applying tags"
    echo ""
    if grep -q "tags_v2.*appeared" apply_with_tags.log; then
        echo "🐛 BUG REPRODUCED: tags_v2 inconsistency error detected!"
        echo ""
        echo "Error details:"
        grep -A 5 "Error:" apply_with_tags.log || true
    else
        echo "Different error occurred - see apply_with_tags.log"
    fi
    exit 1
fi
