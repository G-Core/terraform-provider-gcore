#!/bin/bash
# Demonstration: bcrypt() vs random_password.bcrypt_hash drift behavior

set -e

echo "=================================================="
echo "Comparing bcrypt() Function vs random_password"
echo "=================================================="

# Create test directory
mkdir -p /tmp/bcrypt-test-{function,resource}

echo ""
echo "=== Test 1: bcrypt() Function (causes drift) ==="
cat > /tmp/bcrypt-test-function/main.tf <<'EOF'
locals {
  hash = bcrypt("Test123!")
}

output "hash" {
  value = local.hash
}
EOF

cd /tmp/bcrypt-test-function
terraform init -upgrade > /dev/null 2>&1
echo "Run 1:"
terraform apply -auto-approve 2>/dev/null | grep "hash ="
echo "Run 2:"
terraform apply -auto-approve 2>/dev/null | grep "hash ="
echo "Run 3:"
terraform apply -auto-approve 2>/dev/null | grep "hash ="
echo "❌ Notice: Hash changes on every run!"

echo ""
echo "=== Test 2: random_password.bcrypt_hash (stable) ==="
cat > /tmp/bcrypt-test-resource/main.tf <<'EOF'
terraform {
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "~> 3.6"
    }
  }
}

resource "random_password" "test" {
  length  = 12
  special = false
}

output "hash" {
  value = random_password.test.bcrypt_hash
  sensitive = false
}
EOF

cd /tmp/bcrypt-test-resource
terraform init -upgrade > /dev/null 2>&1
echo "Initial apply:"
terraform apply -auto-approve > /dev/null 2>&1
HASH1=$(terraform output -raw hash 2>/dev/null)
echo "hash = $HASH1"

echo "Run 2 (should be same):"
terraform plan > /dev/null 2>&1
HASH2=$(terraform output -raw hash 2>/dev/null)
echo "hash = $HASH2"

echo "Run 3 (should be same):"
terraform plan > /dev/null 2>&1
HASH3=$(terraform output -raw hash 2>/dev/null)
echo "hash = $HASH3"

if [ "$HASH1" = "$HASH2" ] && [ "$HASH2" = "$HASH3" ]; then
    echo "✅ Success: Hash is stable across runs!"
else
    echo "❌ Error: Hash changed unexpectedly"
fi

echo ""
echo "=== Cleanup ==="
cd /tmp
rm -rf bcrypt-test-{function,resource}
echo "Test directories removed"

echo ""
echo "=== Summary ==="
echo "❌ bcrypt() function: Generates new hash every time (NOT recommended for resources)"
echo "✅ random_password.bcrypt_hash: Generates once, stores in state (RECOMMENDED)"
echo ""
echo "Documentation:"
echo "- bcrypt(): https://developer.hashicorp.com/terraform/language/functions/bcrypt"
echo "- random_password: https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/password"
