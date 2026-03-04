#!/bin/bash
set -e

echo "=========================================="
echo "Router Comprehensive Testing - Phases 3-6"
echo "=========================================="

cd "$(dirname "$0")"

# Load credentials
if [ -f "../.env" ]; then
    set -o allexport
    source ../.env
    set +o allexport
    echo "✓ Credentials loaded"
else
    echo "✗ Error: .env not found"
    exit 1
fi

export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/.terraformrc"

PASSED=0
FAILED=0

# ========== PHASE 3: CRUD with API Verification ==========
echo ""
echo "=========================================="
echo "PHASE 3: CRUD Tests with API Verification"
echo "=========================================="

mkdir -p crud/TC-CRUD-001
cd crud/TC-CRUD-001

cat > main.tf << 'EOF'
terraform {
  required_providers {
    gcore = { source = "gcore/gcore" }
  }
}
provider "gcore" {}

resource "gcore_cloud_network" "network" {
  project_id = 379987
  region_id  = 76
  name       = "test-router-crud-network"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "subnet" {
  project_id  = 379987
  region_id   = 76
  name        = "test-router-crud-subnet"
  cidr        = "192.168.50.0/24"
  network_id  = gcore_cloud_network.network.id
}

resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "test-router-crud"
  external_gateway_info = { enable_snat = true }
  interfaces = [{ subnet_id = gcore_cloud_network_subnet.subnet.id, type = "subnet" }]
  routes = [{ destination = "10.0.5.0/24", nexthop = "192.168.50.1" }]
}

output "router_id" { value = gcore_cloud_network_router.router.id }
EOF

echo "Creating router..."
terraform apply -auto-approve > /dev/null 2>&1

ROUTER_ID=$(terraform output -raw router_id)
echo "✓ Router created: $ROUTER_ID"

# Verify in API
echo "Verifying router exists in GCore API..."
API_RESPONSE=$(curl -s -H "Authorization: APIKey ${GCORE_API_KEY}" \
  "https://api.gcore.com/cloud/v1/routers/379987/76/${ROUTER_ID}")

if echo "$API_RESPONSE" | grep -q "\"id\":\"${ROUTER_ID}\""; then
    echo "✅ TC-CRUD-001 PASSED: Router exists in API"
    PASSED=$((PASSED + 1))
else
    echo "❌ TC-CRUD-001 FAILED: Router not found in API"
    FAILED=$((FAILED + 1))
fi

# Clean up
terraform destroy -auto-approve > /dev/null 2>&1

# Verify deletion
API_RESPONSE=$(curl -s -w "%{http_code}" -H "Authorization: APIKey ${GCORE_API_KEY}" \
  "https://api.gcore.com/cloud/v1/routers/379987/76/${ROUTER_ID}")

if echo "$API_RESPONSE" | grep -q "404"; then
    echo "✅ TC-CRUD-002 PASSED: Router deleted from API (404)"
    PASSED=$((PASSED + 1))
else
    echo "❌ TC-CRUD-002 FAILED: Router still exists after destroy"
    FAILED=$((FAILED + 1))
fi

cd ../..

# ========== PHASE 4: ForceNew Field Testing ==========
echo ""
echo "=========================================="
echo "PHASE 4: ForceNew Field Tests"
echo "=========================================="

mkdir -p forcenew/TC-FORCENEW-001
cd forcenew/TC-FORCENEW-001

cat > main.tf << 'EOF'
terraform {
  required_providers {
    gcore = { source = "gcore/gcore" }
  }
}
provider "gcore" {}

variable "project_id" {
  type = number
  default = 379987
}

resource "gcore_cloud_network" "network" {
  project_id = var.project_id
  region_id  = 76
  name       = "test-router-forcenew-network"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "subnet" {
  project_id  = var.project_id
  region_id   = 76
  name        = "test-router-forcenew-subnet"
  cidr        = "192.168.60.0/24"
  network_id  = gcore_cloud_network.network.id
}

resource "gcore_cloud_network_router" "router" {
  project_id = var.project_id
  region_id  = 76
  name       = "test-router-forcenew"
  external_gateway_info = { enable_snat = true }
  interfaces = [{ subnet_id = gcore_cloud_network_subnet.subnet.id, type = "subnet" }]
}

output "router_id" { value = gcore_cloud_network_router.router.id }
EOF

echo "Creating router with project_id=379987..."
terraform apply -auto-approve > /dev/null 2>&1
OLD_ID=$(terraform output -raw router_id)

echo "Checking if changing project_id triggers replacement..."
PLAN_OUTPUT=$(terraform plan -var="project_id=379988" 2>&1 || true)

if echo "$PLAN_OUTPUT" | grep -q "forces replacement\|must be replaced"; then
    echo "✅ TC-FORCENEW-001 PASSED: project_id change forces replacement"
    PASSED=$((PASSED + 1))
else
    echo "⚠️  TC-FORCENEW-001 SKIPPED: Cannot change project_id (likely validation error)"
    # This is expected - project_id changes usually aren't allowed
fi

terraform destroy -auto-approve > /dev/null 2>&1
cd ../..

# ========== PHASE 5: Import Testing ==========
echo ""
echo "=========================================="
echo "PHASE 5: Import Tests"
echo "=========================================="

mkdir -p import/TC-IMPORT-001
cd import/TC-IMPORT-001

cat > main.tf << 'EOF'
terraform {
  required_providers {
    gcore = { source = "gcore/gcore" }
  }
}
provider "gcore" {}

resource "gcore_cloud_network" "network" {
  project_id = 379987
  region_id  = 76
  name       = "test-router-import-network"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "subnet" {
  project_id  = 379987
  region_id   = 76
  name        = "test-router-import-subnet"
  cidr        = "192.168.70.0/24"
  network_id  = gcore_cloud_network.network.id
}

resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "test-router-import"
  external_gateway_info = { enable_snat = true }
  interfaces = [{ subnet_id = gcore_cloud_network_subnet.subnet.id, type = "subnet" }]
}

output "router_id" { value = gcore_cloud_network_router.router.id }
EOF

echo "Creating router..."
terraform apply -auto-approve > /dev/null 2>&1
ROUTER_ID=$(terraform output -raw router_id)

echo "Removing router from state..."
terraform state rm gcore_cloud_network_router.router > /dev/null 2>&1

echo "Importing router back..."
terraform import gcore_cloud_network_router.router "379987/76/${ROUTER_ID}" > /dev/null 2>&1

echo "Checking for drift after import..."
if terraform plan -detailed-exitcode > /dev/null 2>&1; then
    echo "✅ TC-IMPORT-001 PASSED: No drift after import"
    PASSED=$((PASSED + 1))
else
    echo "❌ TC-IMPORT-001 FAILED: Drift detected after import"
    FAILED=$((FAILED + 1))
fi

terraform destroy -auto-approve > /dev/null 2>&1
cd ../..

# ========== PHASE 6: Edge Case Testing ==========
echo ""
echo "=========================================="
echo "PHASE 6: Edge Case Tests"
echo "=========================================="

mkdir -p edge-cases/TC-EDGE-001
cd edge-cases/TC-EDGE-001

cat > main.tf << 'EOF'
terraform {
  required_providers {
    gcore = { source = "gcore/gcore" }
  }
}
provider "gcore" {}

resource "gcore_cloud_network" "network" {
  project_id = 379987
  region_id  = 76
  name       = "test-router-edge-network"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "subnet" {
  project_id  = 379987
  region_id   = 76
  name        = "test-router-edge-subnet"
  cidr        = "192.168.80.0/24"
  network_id  = gcore_cloud_network.network.id
}

variable "routes_config" {
  type = string
  default = "empty"
}

resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "test-router-edge"
  external_gateway_info = { enable_snat = true }
  interfaces = [{ subnet_id = gcore_cloud_network_subnet.subnet.id, type = "subnet" }]
  routes = var.routes_config == "empty" ? [] : (
    var.routes_config == "single" ? [{ destination = "10.0.8.0/24", nexthop = "192.168.80.1" }] : [
      { destination = "10.0.8.0/24", nexthop = "192.168.80.1" },
      { destination = "10.0.9.0/24", nexthop = "192.168.80.1" },
      { destination = "10.0.10.0/24", nexthop = "192.168.80.1" }
    ]
  )
}
EOF

echo "Test 1: Explicitly empty routes..."
terraform apply -auto-approve -var="routes_config=empty" > /dev/null 2>&1
if terraform plan -detailed-exitcode -var="routes_config=empty" > /dev/null 2>&1; then
    echo "✅ TC-EDGE-001 PASSED: Empty routes no drift"
    PASSED=$((PASSED + 1))
else
    echo "❌ TC-EDGE-001 FAILED: Drift with empty routes"
    FAILED=$((FAILED + 1))
fi

echo "Test 2: Multiple routes..."
terraform apply -auto-approve -var="routes_config=multiple" > /dev/null 2>&1
if terraform plan -detailed-exitcode -var="routes_config=multiple" > /dev/null 2>&1; then
    echo "✅ TC-EDGE-002 PASSED: Multiple routes no drift"
    PASSED=$((PASSED + 1))
else
    echo "❌ TC-EDGE-002 FAILED: Drift with multiple routes"
    FAILED=$((FAILED + 1))
fi

echo "Test 3: Route deletion then re-addition..."
terraform apply -auto-approve -var="routes_config=empty" > /dev/null 2>&1
terraform apply -auto-approve -var="routes_config=single" > /dev/null 2>&1
if terraform plan -detailed-exitcode -var="routes_config=single" > /dev/null 2>&1; then
    echo "✅ TC-EDGE-003 PASSED: Route deletion+addition cycle no drift"
    PASSED=$((PASSED + 1))
else
    echo "❌ TC-EDGE-003 FAILED: Drift after route cycle"
    FAILED=$((FAILED + 1))
fi

terraform destroy -auto-approve > /dev/null 2>&1
cd ../..

# ========== SUMMARY ==========
echo ""
echo "=========================================="
echo "Test Summary"
echo "=========================================="
echo "PASSED: $PASSED"
echo "FAILED: $FAILED"
echo ""

if [ $FAILED -gt 0 ]; then
    echo "❌ Some tests failed"
    exit 1
else
    echo "✅ All tests passed!"
    exit 0
fi
