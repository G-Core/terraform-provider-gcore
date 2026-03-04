#!/bin/bash
# Comprehensive testing script for gcore_cdn_trusted_ca_certificate
set -euo pipefail

export GCORE_API_KEY='21788$1e278ce67b6aa33f178122658b1dd0210d0edff453d348acb9b68bffea6a635b7791925ddda198d5678a4dc20269fe04a263ca92c7e5aa41ea79075f89b66bf6'
export GCORE_CLOUD_PROJECT_ID=379987
export GCORE_CLOUD_REGION_ID=76

TESTDIR="$(cd "$(dirname "$0")" && pwd)"
RESULTS_DIR="$TESTDIR/results"
mkdir -p "$RESULTS_DIR"

CA_CERT="$(cat "$TESTDIR/ca-cert.pem")"

# Utility: run terraform command, capture output, return exit code
run_tf() {
    local label="$1"
    shift
    local outfile="$RESULTS_DIR/${label}.txt"
    echo "=== Running: terraform $* ==="
    terraform "$@" 2>&1 | tee "$outfile"
    local rc=${PIPESTATUS[0]}
    echo "=== Exit code: $rc ==="
    return $rc
}

########################################################################
# TEST 1: Basic Create
########################################################################
echo "############################################################"
echo "# TEST 1: Basic Create"
echo "############################################################"
TEST1_DIR="$TESTDIR/test1-basic-create"
rm -rf "$TEST1_DIR" && mkdir -p "$TEST1_DIR"

cat > "$TEST1_DIR/main.tf" <<TFEOF
terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

resource "gcore_cdn_trusted_ca_certificate" "test" {
  name            = "tf-test-cacert-basic"
  ssl_certificate = <<-EOT
${CA_CERT}
EOT
}

output "cert_id" {
  value = gcore_cdn_trusted_ca_certificate.test.id
}
output "cert_name" {
  value = gcore_cdn_trusted_ca_certificate.test.name
}
output "cert_issuer" {
  value = gcore_cdn_trusted_ca_certificate.test.cert_issuer
}
output "cert_subject_cn" {
  value = gcore_cdn_trusted_ca_certificate.test.cert_subject_cn
}
output "has_related_resources" {
  value = gcore_cdn_trusted_ca_certificate.test.has_related_resources
}
output "validity_not_after" {
  value = gcore_cdn_trusted_ca_certificate.test.validity_not_after
}
output "validity_not_before" {
  value = gcore_cdn_trusted_ca_certificate.test.validity_not_before
}
TFEOF

cd "$TEST1_DIR"

echo "--- Test 1a: First Apply (Create) ---"
run_tf "test1a_apply" apply -auto-approve -no-color || true

echo "--- Test 1b: Second Apply (Drift Check) ---"
run_tf "test1b_drift" apply -auto-approve -no-color || true

echo "--- Test 1c: Plan (should be empty) ---"
run_tf "test1c_plan" plan -no-color -detailed-exitcode || true

echo "--- Test 1d: Destroy ---"
run_tf "test1d_destroy" destroy -auto-approve -no-color || true

########################################################################
# TEST 2: Name Update (In-Place)
########################################################################
echo ""
echo "############################################################"
echo "# TEST 2: Name Update (In-Place)"
echo "############################################################"
TEST2_DIR="$TESTDIR/test2-name-update"
rm -rf "$TEST2_DIR" && mkdir -p "$TEST2_DIR"

cat > "$TEST2_DIR/main.tf" <<TFEOF
terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

variable "cert_name" {
  type    = string
  default = "tf-test-cacert-update-v1"
}

resource "gcore_cdn_trusted_ca_certificate" "test" {
  name            = var.cert_name
  ssl_certificate = <<-EOT
${CA_CERT}
EOT
}

output "cert_id" {
  value = gcore_cdn_trusted_ca_certificate.test.id
}
output "cert_name" {
  value = gcore_cdn_trusted_ca_certificate.test.name
}
TFEOF

cd "$TEST2_DIR"

echo "--- Test 2a: Create with original name ---"
run_tf "test2a_create" apply -auto-approve -no-color -var='cert_name=tf-test-cacert-update-v1' || true

echo "--- Test 2b: Update name (in-place, no recreation) ---"
run_tf "test2b_update" apply -auto-approve -no-color -var='cert_name=tf-test-cacert-update-v2' || true

echo "--- Test 2c: Drift check after update ---"
run_tf "test2c_drift" apply -auto-approve -no-color -var='cert_name=tf-test-cacert-update-v2' || true

echo "--- Test 2d: Plan after update (should be empty) ---"
run_tf "test2d_plan" plan -no-color -detailed-exitcode -var='cert_name=tf-test-cacert-update-v2' || true

echo "--- Test 2e: Destroy ---"
run_tf "test2e_destroy" destroy -auto-approve -no-color -var='cert_name=tf-test-cacert-update-v2' || true

########################################################################
# TEST 3: Import + Drift Check
########################################################################
echo ""
echo "############################################################"
echo "# TEST 3: Import + Drift Check"
echo "############################################################"
TEST3_DIR="$TESTDIR/test3-import"
rm -rf "$TEST3_DIR" && mkdir -p "$TEST3_DIR"

cat > "$TEST3_DIR/main.tf" <<TFEOF
terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

resource "gcore_cdn_trusted_ca_certificate" "test" {
  name            = "tf-test-cacert-import"
  ssl_certificate = <<-EOT
${CA_CERT}
EOT
}
TFEOF

cd "$TEST3_DIR"

echo "--- Test 3a: Create resource ---"
run_tf "test3a_create" apply -auto-approve -no-color || true

# Get the ID
CERT_ID=$(terraform output -raw -state="$TEST3_DIR/terraform.tfstate" cert_id 2>/dev/null || terraform show -json | python3 -c "import sys,json; s=json.load(sys.stdin); print([r for r in s.get('values',{}).get('root_module',{}).get('resources',[]) if r['type']=='gcore_cdn_trusted_ca_certificate'][0]['values']['id'])" 2>/dev/null || echo "")

if [ -z "$CERT_ID" ]; then
  # Fallback: extract from state
  CERT_ID=$(terraform show -no-color | grep '^\s*id\s*=' | head -1 | sed 's/.*=\s*//' | tr -d ' "')
fi

echo "Certificate ID for import: $CERT_ID"
echo "$CERT_ID" > "$RESULTS_DIR/test3_cert_id.txt"

echo "--- Test 3b: Remove from state ---"
run_tf "test3b_remove" state rm gcore_cdn_trusted_ca_certificate.test || true

echo "--- Test 3c: Import by ID ---"
run_tf "test3c_import" import -no-color gcore_cdn_trusted_ca_certificate.test "$CERT_ID" || true

echo "--- Test 3d: Plan after import (drift check) ---"
run_tf "test3d_plan" plan -no-color -detailed-exitcode || true

echo "--- Test 3e: Apply after import (should be no-op or minimal) ---"
run_tf "test3e_apply" apply -auto-approve -no-color || true

echo "--- Test 3f: Second apply after import (drift check) ---"
run_tf "test3f_drift" plan -no-color -detailed-exitcode || true

echo "--- Test 3g: Destroy ---"
run_tf "test3g_destroy" destroy -auto-approve -no-color || true

########################################################################
# TEST 4: Data Source
########################################################################
echo ""
echo "############################################################"
echo "# TEST 4: Data Source"
echo "############################################################"
TEST4_DIR="$TESTDIR/test4-datasource"
rm -rf "$TEST4_DIR" && mkdir -p "$TEST4_DIR"

cat > "$TEST4_DIR/main.tf" <<TFEOF
terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

resource "gcore_cdn_trusted_ca_certificate" "test" {
  name            = "tf-test-cacert-datasource"
  ssl_certificate = <<-EOT
${CA_CERT}
EOT
}

data "gcore_cdn_trusted_ca_certificate" "lookup" {
  id = gcore_cdn_trusted_ca_certificate.test.id
}

output "resource_name" {
  value = gcore_cdn_trusted_ca_certificate.test.name
}
output "datasource_name" {
  value = data.gcore_cdn_trusted_ca_certificate.lookup.name
}
output "datasource_cert_issuer" {
  value = data.gcore_cdn_trusted_ca_certificate.lookup.cert_issuer
}
output "datasource_cert_subject_cn" {
  value = data.gcore_cdn_trusted_ca_certificate.lookup.cert_subject_cn
}
output "datasource_validity_not_after" {
  value = data.gcore_cdn_trusted_ca_certificate.lookup.validity_not_after
}
output "names_match" {
  value = gcore_cdn_trusted_ca_certificate.test.name == data.gcore_cdn_trusted_ca_certificate.lookup.name
}
TFEOF

cd "$TEST4_DIR"

echo "--- Test 4a: Apply (create resource + read data source) ---"
run_tf "test4a_apply" apply -auto-approve -no-color || true

echo "--- Test 4b: Drift check ---"
run_tf "test4b_drift" plan -no-color -detailed-exitcode || true

echo "--- Test 4c: Destroy ---"
run_tf "test4c_destroy" destroy -auto-approve -no-color || true

echo ""
echo "############################################################"
echo "# ALL TESTS COMPLETE"
echo "############################################################"
echo "Results saved to: $RESULTS_DIR"
ls -la "$RESULTS_DIR"
