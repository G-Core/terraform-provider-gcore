#!/bin/bash
set -e

echo "=== Testing WithJSONSet Necessity ==="
echo "This test verifies if routes=[] is serialized automatically or if WithJSONSet is required"
echo ""

# Load credentials
if [ -f .env ]; then
    source .env
elif [ -f ../.env ]; then
    source ../.env
else
    echo "ERROR: .env not found"
    exit 1
fi

export TF_CLI_CONFIG_FILE=/Users/user/repos/gcore-terraform/.terraformrc

cd test-withjsonset

# Clean start
rm -rf .terraform terraform.tfstate* .terraform.lock.hcl 2>/dev/null || true

echo "=== Phase 1: Create router with 2 routes ==="
# Skip init - we have dev overrides in .terraformrc
terraform apply -auto-approve

ROUTER_ID=$(terraform output -raw router_id)
echo "Router ID: $ROUTER_ID"

ROUTES_COUNT=$(terraform output -json routes | jq '. | length')
echo "Initial routes count: $ROUTES_COUNT"

if [ "$ROUTES_COUNT" != "2" ]; then
    echo "❌ Failed to create router with 2 routes"
    exit 1
fi

echo ""
echo "=== Phase 2: Test deleting routes by setting routes = [] (empty block) ==="
echo "Modifying main.tf to have empty routes block..."

cat > main.tf << 'EOF'
terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {
}

resource "gcore_cloud_network_router" "test" {
  project_id = 379987
  region_id  = 76
  name       = "withjsonset-test"

  # Empty routes - explicitly setting routes to []
  routes = []
}

output "router_id" {
  value = gcore_cloud_network_router.test.id
}

output "routes" {
  value = gcore_cloud_network_router.test.routes
}
EOF

echo "Running terraform plan to see what will happen..."
TF_LOG=DEBUG terraform plan -out=phase2.tfplan 2>&1 | tee phase2_plan.log

echo ""
echo "Checking plan output for PATCH body..."
if grep -q '"routes":\[\]' phase2_plan.log; then
    echo "✅ Found routes:[] in PATCH body - empty array IS being serialized"
elif grep -q '"routes":null' phase2_plan.log; then
    echo "⚠️  Found routes:null in PATCH body - this might cause issues"
elif grep -q 'routes' phase2_plan.log | grep -v 'routes_changed'; then
    echo "⚠️  Found routes field in PATCH body"
else
    echo "❌ routes field NOT found in PATCH body - might be omitted!"
fi

echo ""
echo "Applying changes..."
terraform apply -auto-approve phase2.tfplan

ROUTES_COUNT=$(terraform output -json routes | jq '. | length')
echo "Routes count after empty block: $ROUTES_COUNT"

if [ "$ROUTES_COUNT" = "0" ]; then
    echo "✅ Empty routes block successfully deleted all routes"
else
    echo "❌ Empty routes block did not delete routes (count: $ROUTES_COUNT)"
fi

echo ""
echo "=== Phase 3: Add routes back for next test ==="
cat > main.tf << 'EOF'
terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {
}

resource "gcore_cloud_network_router" "test" {
  project_id = 379987
  region_id  = 76
  name       = "withjsonset-test"

  routes = [
    {
      destination = "10.0.1.0/24"
      nexthop     = "192.168.1.1"
    },
    {
      destination = "10.0.2.0/24"
      nexthop     = "192.168.1.2"
    }
  ]
}

output "router_id" {
  value = gcore_cloud_network_router.test.id
}

output "routes" {
  value = gcore_cloud_network_router.test.routes
}
EOF

terraform apply -auto-approve > /dev/null
ROUTES_COUNT=$(terraform output -json routes | jq '. | length')
echo "Routes restored: $ROUTES_COUNT"

echo ""
echo "=== Phase 4: Test removing routes attribute entirely ==="
echo "Modifying main.tf to have NO routes block..."

cat > main.tf << 'EOF'
terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {
}

resource "gcore_cloud_network_router" "test" {
  project_id = 379987
  region_id  = 76
  name       = "withjsonset-test"

  # No routes block at all - attribute removed from config
}

output "router_id" {
  value = gcore_cloud_network_router.test.id
}

output "routes" {
  value = gcore_cloud_network_router.test.routes
}
EOF

echo "Running terraform plan..."
TF_LOG=DEBUG terraform plan -out=phase4.tfplan 2>&1 | tee phase4_plan.log

echo ""
echo "Checking plan output for PATCH body..."
if grep -q '"routes":\[\]' phase4_plan.log; then
    echo "✅ Found routes:[] in PATCH body - ModifyPlan detected removal and set to []"
elif grep -q '"routes":null' phase4_plan.log; then
    echo "⚠️  Found routes:null in PATCH body"
elif grep -q 'routes' phase4_plan.log | grep -v 'routes_changed'; then
    echo "⚠️  Found routes field in PATCH body"
else
    echo "❌ routes field NOT found in PATCH body - might be omitted!"
fi

echo ""
echo "Applying changes..."
terraform apply -auto-approve phase4.tfplan

ROUTES_COUNT=$(terraform output -json routes | jq 'if . == null then 0 else . | length end')
echo "Routes count after removing attribute: $ROUTES_COUNT"

if [ "$ROUTES_COUNT" = "0" ]; then
    echo "✅ Removing routes attribute successfully deleted all routes"
else
    echo "❌ Removing routes attribute did not delete routes (count: $ROUTES_COUNT)"
fi

echo ""
echo "=== Phase 5: Check actual PATCH body sent to API ==="
echo "Looking for PATCH requests in debug logs..."

# Extract actual PATCH body from TF_LOG
PATCH_BODY=$(grep -A 20 'PATCH.*routers' phase4_plan.log | grep -m1 '"routes"' || echo "not found")
echo "PATCH body with routes field: $PATCH_BODY"

echo ""
echo "=== Cleanup ==="
terraform destroy -auto-approve

cd ..

echo ""
echo "=== FINDINGS SUMMARY ==="
echo ""
echo "1. Empty routes block (routes {}) behavior:"
if [ -f test-withjsonset/phase2_plan.log ]; then
    if grep -q '"routes":\[\]' test-withjsonset/phase2_plan.log; then
        echo "   ✅ routes=[] IS serialized (Pedro is CORRECT)"
    else
        echo "   ❌ routes=[] is NOT serialized (WithJSONSet might be needed)"
    fi
fi

echo ""
echo "2. Removed routes attribute behavior:"
if [ -f test-withjsonset/phase4_plan.log ]; then
    if grep -q '"routes":\[\]' test-withjsonset/phase4_plan.log; then
        echo "   ✅ ModifyPlan sets routes=[] and it IS serialized"
    else
        echo "   ❌ routes field might be omitted (WithJSONSet might be needed)"
    fi
fi

echo ""
echo "3. Recommendation:"
echo "   - If routes=[] IS serialized in both cases, WithJSONSet can be removed"
echo "   - If routes=[] is NOT serialized, WithJSONSet is REQUIRED"
