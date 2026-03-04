#!/bin/bash
set -euo pipefail

export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/.terraformrc"
export HTTP_PROXY="http://127.0.0.1:9092"
export HTTPS_PROXY="http://127.0.0.1:9092"
export NO_PROXY="registry.terraform.io,releases.hashicorp.com"
export SSL_CERT_FILE="/Users/user/repos/gcore-terraform/ca-bundle.pem"
export REQUESTS_CA_BUNDLE="/Users/user/repos/gcore-terraform/ca-bundle.pem"
eval "$(grep -v '^#' /Users/user/repos/gcore-terraform/.env | sed 's/^/export /')"

cd /Users/user/repos/gcore-terraform/test-dns-zone-skill

TF="terraform"

echo "=== TEST 1: Create minimal zone (name only) ==="
$TF apply -auto-approve -no-color 2>&1 | tee evidence/test1_create.log
echo ""
echo "=== TEST 1: Drift check ==="
$TF apply -auto-approve -no-color 2>&1 | tail -5
$TF plan -detailed-exitcode -no-color 2>&1 | tee evidence/test1_drift.log
T1_EXIT=$?
echo "TEST 1 DRIFT EXIT CODE: $T1_EXIT"

echo ""
echo "=== TEST 2: Update with SOA fields ==="
$TF apply -auto-approve -no-color \
  -var='contact=admin@test.com' \
  -var='refresh=7200' \
  -var='retry=1800' \
  -var='expiry=604800' \
  -var='nx_ttl=600' \
  2>&1 | tee evidence/test2_soa_update.log

echo ""
echo "=== TEST 2: Drift check ==="
$TF plan -detailed-exitcode -no-color \
  -var='contact=admin@test.com' \
  -var='refresh=7200' \
  -var='retry=1800' \
  -var='expiry=604800' \
  -var='nx_ttl=600' \
  2>&1 | tee evidence/test2_drift.log
T2_EXIT=$?
echo "TEST 2 DRIFT EXIT CODE: $T2_EXIT"

echo ""
echo "=== TEST 3: Update SOA fields (change contact and refresh) ==="
$TF apply -auto-approve -no-color \
  -var='contact=changed@test.com' \
  -var='refresh=3600' \
  -var='retry=1800' \
  -var='expiry=604800' \
  -var='nx_ttl=600' \
  2>&1 | tee evidence/test3_soa_change.log

echo ""
echo "=== TEST 3: Drift check ==="
$TF plan -detailed-exitcode -no-color \
  -var='contact=changed@test.com' \
  -var='refresh=3600' \
  -var='retry=1800' \
  -var='expiry=604800' \
  -var='nx_ttl=600' \
  2>&1 | tee evidence/test3_drift.log
T3_EXIT=$?
echo "TEST 3 DRIFT EXIT CODE: $T3_EXIT"

echo ""
echo "=== TEST 4: Add meta field ==="
$TF apply -auto-approve -no-color \
  -var='contact=changed@test.com' \
  -var='refresh=3600' \
  -var='retry=1800' \
  -var='expiry=604800' \
  -var='nx_ttl=600' \
  -var='meta={"webhook":"https://example.com/hook"}' \
  2>&1 | tee evidence/test4_meta_add.log

echo ""
echo "=== TEST 4: Drift check (meta should not have server-injected keys) ==="
$TF plan -detailed-exitcode -no-color \
  -var='contact=changed@test.com' \
  -var='refresh=3600' \
  -var='retry=1800' \
  -var='expiry=604800' \
  -var='nx_ttl=600' \
  -var='meta={"webhook":"https://example.com/hook"}' \
  2>&1 | tee evidence/test4_drift.log
T4_EXIT=$?
echo "TEST 4 DRIFT EXIT CODE: $T4_EXIT"

echo ""
echo "=== TEST 5: Enable DNSSEC ==="
$TF apply -auto-approve -no-color \
  -var='contact=changed@test.com' \
  -var='refresh=3600' \
  -var='retry=1800' \
  -var='expiry=604800' \
  -var='nx_ttl=600' \
  -var='meta={"webhook":"https://example.com/hook"}' \
  -var='dnssec_enabled=true' \
  2>&1 | tee evidence/test5_dnssec_enable.log

echo ""
echo "=== TEST 5: Drift check ==="
$TF plan -detailed-exitcode -no-color \
  -var='contact=changed@test.com' \
  -var='refresh=3600' \
  -var='retry=1800' \
  -var='expiry=604800' \
  -var='nx_ttl=600' \
  -var='meta={"webhook":"https://example.com/hook"}' \
  -var='dnssec_enabled=true' \
  2>&1 | tee evidence/test5_drift.log
T5_EXIT=$?
echo "TEST 5 DRIFT EXIT CODE: $T5_EXIT"

echo ""
echo "=== TEST 6: Disable DNSSEC ==="
$TF apply -auto-approve -no-color \
  -var='contact=changed@test.com' \
  -var='refresh=3600' \
  -var='retry=1800' \
  -var='expiry=604800' \
  -var='nx_ttl=600' \
  -var='meta={"webhook":"https://example.com/hook"}' \
  -var='dnssec_enabled=false' \
  2>&1 | tee evidence/test6_dnssec_disable.log

echo ""
echo "=== TEST 6: Drift check ==="
$TF plan -detailed-exitcode -no-color \
  -var='contact=changed@test.com' \
  -var='refresh=3600' \
  -var='retry=1800' \
  -var='expiry=604800' \
  -var='nx_ttl=600' \
  -var='meta={"webhook":"https://example.com/hook"}' \
  -var='dnssec_enabled=false' \
  2>&1 | tee evidence/test6_drift.log
T6_EXIT=$?
echo "TEST 6 DRIFT EXIT CODE: $T6_EXIT"

echo ""
echo "=== TEST 7: Import ==="
# Save current state for comparison
$TF show -json -no-color > evidence/test7_pre_import_state.json 2>/dev/null

# Remove from state
$TF state rm gcore_dns_zone.test -no-color 2>&1 | tee evidence/test7_import.log

# Import
$TF import -no-color \
  -var='contact=changed@test.com' \
  -var='refresh=3600' \
  -var='retry=1800' \
  -var='expiry=604800' \
  -var='nx_ttl=600' \
  -var='meta={"webhook":"https://example.com/hook"}' \
  -var='dnssec_enabled=false' \
  gcore_dns_zone.test "test-tf-dns-1770800090.com" 2>&1 | tee -a evidence/test7_import.log

echo ""
echo "=== TEST 7: Drift check after import ==="
$TF plan -detailed-exitcode -no-color \
  -var='contact=changed@test.com' \
  -var='refresh=3600' \
  -var='retry=1800' \
  -var='expiry=604800' \
  -var='nx_ttl=600' \
  -var='meta={"webhook":"https://example.com/hook"}' \
  -var='dnssec_enabled=false' \
  2>&1 | tee evidence/test7_drift.log
T7_EXIT=$?
echo "TEST 7 DRIFT EXIT CODE: $T7_EXIT"

echo ""
echo "=== TEST 8: Destroy ==="
$TF destroy -auto-approve -no-color \
  -var='contact=changed@test.com' \
  -var='refresh=3600' \
  -var='retry=1800' \
  -var='expiry=604800' \
  -var='nx_ttl=600' \
  -var='meta={"webhook":"https://example.com/hook"}' \
  -var='dnssec_enabled=false' \
  2>&1 | tee evidence/test8_destroy.log

echo ""
echo "========================"
echo "=== RESULTS SUMMARY ==="
echo "========================"
echo "Test 1 (Create minimal + drift): EXIT=$T1_EXIT $([ $T1_EXIT -eq 0 ] && echo 'PASS' || echo 'FAIL')"
echo "Test 2 (SOA fields + drift): EXIT=$T2_EXIT $([ $T2_EXIT -eq 0 ] && echo 'PASS' || echo 'FAIL')"
echo "Test 3 (SOA update + drift): EXIT=$T3_EXIT $([ $T3_EXIT -eq 0 ] && echo 'PASS' || echo 'FAIL')"
echo "Test 4 (Meta field + drift): EXIT=$T4_EXIT $([ $T4_EXIT -eq 0 ] && echo 'PASS' || echo 'FAIL')"
echo "Test 5 (DNSSEC enable + drift): EXIT=$T5_EXIT $([ $T5_EXIT -eq 0 ] && echo 'PASS' || echo 'FAIL')"
echo "Test 6 (DNSSEC disable + drift): EXIT=$T6_EXIT $([ $T6_EXIT -eq 0 ] && echo 'PASS' || echo 'FAIL')"
echo "Test 7 (Import + drift): EXIT=$T7_EXIT $([ $T7_EXIT -eq 0 ] && echo 'PASS' || echo 'FAIL')"
echo "Test 8 (Destroy): PASS (completed)"
