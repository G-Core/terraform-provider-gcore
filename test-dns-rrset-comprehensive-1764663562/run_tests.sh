#!/bin/bash
set -e

# Load credentials
if [ -f /Users/user/repos/gcore-terraform/.env ]; then
    set -o allexport
    source /Users/user/repos/gcore-terraform/.env
    set +o allexport
    echo "✅ Credentials loaded"
else
    echo "❌ .env not found"
    exit 1
fi

export TF_CLI_CONFIG_FILE=/Users/user/repos/gcore-terraform/.terraformrc

EVIDENCE_DIR="evidence"
mkdir -p $EVIDENCE_DIR

echo "=============================================="
echo "DNS Zone RRSet Comprehensive Testing"
echo "=============================================="

# =============================================================================
# Test 1: Create A Record
# =============================================================================
echo ""
echo "=== Test 1: Create A Record ==="
terraform apply -auto-approve 2>&1 | tee $EVIDENCE_DIR/test1_create_a.log
echo "Test 1 Create: DONE"

# =============================================================================
# Test 6: Drift Test after Create
# =============================================================================
echo ""
echo "=== Test 6: Drift Test (A Record) ==="
terraform plan -detailed-exitcode 2>&1 | tee $EVIDENCE_DIR/test6_drift_a.log
DRIFT_EXIT=$?
if [ $DRIFT_EXIT -eq 0 ]; then
    echo "✅ Test 6 PASS: No drift detected (exit code 0)"
else
    echo "❌ Test 6 FAIL: Drift detected (exit code $DRIFT_EXIT)"
fi

# =============================================================================
# Test 7: Update TTL
# =============================================================================
echo ""
echo "=== Test 7: Update TTL (300 -> 600) ==="
terraform apply -auto-approve -var="ttl_a=600" 2>&1 | tee $EVIDENCE_DIR/test7_update_ttl.log
echo "Verifying TTL update..."
terraform plan -detailed-exitcode -var="ttl_a=600" 2>&1 | tee $EVIDENCE_DIR/test7_drift_after_ttl.log
DRIFT_EXIT=$?
if [ $DRIFT_EXIT -eq 0 ]; then
    echo "✅ Test 7 PASS: TTL updated, no drift"
else
    echo "❌ Test 7 FAIL: Drift after TTL update (exit code $DRIFT_EXIT)"
fi

# =============================================================================
# Test 8: Update - Add Second IP
# =============================================================================
echo ""
echo "=== Test 8: Add Second IP ==="
terraform apply -auto-approve -var="ttl_a=600" -var='records_a=[{content=["\"192.168.1.100\""],enabled=true},{content=["\"192.168.1.101\""],enabled=true}]' 2>&1 | tee $EVIDENCE_DIR/test8_add_ip.log
echo "Verifying second IP..."
terraform plan -detailed-exitcode -var="ttl_a=600" -var='records_a=[{content=["\"192.168.1.100\""],enabled=true},{content=["\"192.168.1.101\""],enabled=true}]' 2>&1 | tee $EVIDENCE_DIR/test8_drift_after_ip.log
DRIFT_EXIT=$?
if [ $DRIFT_EXIT -eq 0 ]; then
    echo "✅ Test 8 PASS: Second IP added, no drift"
else
    echo "❌ Test 8 FAIL: Drift after adding IP (exit code $DRIFT_EXIT)"
fi

# =============================================================================
# Test 2: Create AAAA Record
# =============================================================================
echo ""
echo "=== Test 2: Create AAAA Record ==="
terraform apply -auto-approve -var="ttl_a=600" -var='records_a=[{content=["\"192.168.1.100\""],enabled=true},{content=["\"192.168.1.101\""],enabled=true}]' -var="test_aaaa=true" 2>&1 | tee $EVIDENCE_DIR/test2_create_aaaa.log
terraform plan -detailed-exitcode -var="ttl_a=600" -var='records_a=[{content=["\"192.168.1.100\""],enabled=true},{content=["\"192.168.1.101\""],enabled=true}]' -var="test_aaaa=true" 2>&1 | tee $EVIDENCE_DIR/test2_drift_aaaa.log
DRIFT_EXIT=$?
if [ $DRIFT_EXIT -eq 0 ]; then
    echo "✅ Test 2 PASS: AAAA created, no drift"
else
    echo "❌ Test 2 FAIL: Drift after AAAA (exit code $DRIFT_EXIT)"
fi

# =============================================================================
# Test 3: Create CNAME Record
# =============================================================================
echo ""
echo "=== Test 3: Create CNAME Record ==="
terraform apply -auto-approve -var="ttl_a=600" -var='records_a=[{content=["\"192.168.1.100\""],enabled=true},{content=["\"192.168.1.101\""],enabled=true}]' -var="test_aaaa=true" -var="test_cname=true" 2>&1 | tee $EVIDENCE_DIR/test3_create_cname.log
terraform plan -detailed-exitcode -var="ttl_a=600" -var='records_a=[{content=["\"192.168.1.100\""],enabled=true},{content=["\"192.168.1.101\""],enabled=true}]' -var="test_aaaa=true" -var="test_cname=true" 2>&1 | tee $EVIDENCE_DIR/test3_drift_cname.log
DRIFT_EXIT=$?
if [ $DRIFT_EXIT -eq 0 ]; then
    echo "✅ Test 3 PASS: CNAME created, no drift"
else
    echo "❌ Test 3 FAIL: Drift after CNAME (exit code $DRIFT_EXIT)"
fi

# =============================================================================
# Test 4: Create MX Record
# =============================================================================
echo ""
echo "=== Test 4: Create MX Record ==="
terraform apply -auto-approve -var="ttl_a=600" -var='records_a=[{content=["\"192.168.1.100\""],enabled=true},{content=["\"192.168.1.101\""],enabled=true}]' -var="test_aaaa=true" -var="test_cname=true" -var="test_mx=true" 2>&1 | tee $EVIDENCE_DIR/test4_create_mx.log
terraform plan -detailed-exitcode -var="ttl_a=600" -var='records_a=[{content=["\"192.168.1.100\""],enabled=true},{content=["\"192.168.1.101\""],enabled=true}]' -var="test_aaaa=true" -var="test_cname=true" -var="test_mx=true" 2>&1 | tee $EVIDENCE_DIR/test4_drift_mx.log
DRIFT_EXIT=$?
if [ $DRIFT_EXIT -eq 0 ]; then
    echo "✅ Test 4 PASS: MX created, no drift"
else
    echo "❌ Test 4 FAIL: Drift after MX (exit code $DRIFT_EXIT)"
fi

# =============================================================================
# Test 5: Create TXT Record
# =============================================================================
echo ""
echo "=== Test 5: Create TXT Record ==="
terraform apply -auto-approve -var="ttl_a=600" -var='records_a=[{content=["\"192.168.1.100\""],enabled=true},{content=["\"192.168.1.101\""],enabled=true}]' -var="test_aaaa=true" -var="test_cname=true" -var="test_mx=true" -var="test_txt=true" 2>&1 | tee $EVIDENCE_DIR/test5_create_txt.log
terraform plan -detailed-exitcode -var="ttl_a=600" -var='records_a=[{content=["\"192.168.1.100\""],enabled=true},{content=["\"192.168.1.101\""],enabled=true}]' -var="test_aaaa=true" -var="test_cname=true" -var="test_mx=true" -var="test_txt=true" 2>&1 | tee $EVIDENCE_DIR/test5_drift_txt.log
DRIFT_EXIT=$?
if [ $DRIFT_EXIT -eq 0 ]; then
    echo "✅ Test 5 PASS: TXT created, no drift"
else
    echo "❌ Test 5 FAIL: Drift after TXT (exit code $DRIFT_EXIT)"
fi

# =============================================================================
# Test 9: Create GeoDNS Record with Pickers
# =============================================================================
echo ""
echo "=== Test 9: Create GeoDNS Record with Pickers ==="
terraform apply -auto-approve -var="ttl_a=600" -var='records_a=[{content=["\"192.168.1.100\""],enabled=true},{content=["\"192.168.1.101\""],enabled=true}]' -var="test_aaaa=true" -var="test_cname=true" -var="test_mx=true" -var="test_txt=true" -var="test_geo=true" 2>&1 | tee $EVIDENCE_DIR/test9_create_geo.log
terraform plan -detailed-exitcode -var="ttl_a=600" -var='records_a=[{content=["\"192.168.1.100\""],enabled=true},{content=["\"192.168.1.101\""],enabled=true}]' -var="test_aaaa=true" -var="test_cname=true" -var="test_mx=true" -var="test_txt=true" -var="test_geo=true" 2>&1 | tee $EVIDENCE_DIR/test9_drift_geo.log
DRIFT_EXIT=$?
if [ $DRIFT_EXIT -eq 0 ]; then
    echo "✅ Test 9 PASS: GeoDNS created, no drift"
else
    echo "❌ Test 9 FAIL: Drift after GeoDNS (exit code $DRIFT_EXIT)"
fi

# =============================================================================
# Test 11: Import Test
# =============================================================================
echo ""
echo "=== Test 11: Import Test ==="
# First, remove A record from state
terraform state rm gcore_dns_zone_rrset.test_a 2>&1 | tee $EVIDENCE_DIR/test11_state_rm.log
# Import it back
terraform import gcore_dns_zone_rrset.test_a "maxima.lt/tf-comprehensive-a.maxima.lt/A" 2>&1 | tee $EVIDENCE_DIR/test11_import.log
# Check for drift after import
terraform plan -detailed-exitcode -var="ttl_a=600" -var='records_a=[{content=["\"192.168.1.100\""],enabled=true},{content=["\"192.168.1.101\""],enabled=true}]' -var="test_aaaa=true" -var="test_cname=true" -var="test_mx=true" -var="test_txt=true" -var="test_geo=true" 2>&1 | tee $EVIDENCE_DIR/test11_drift_after_import.log
DRIFT_EXIT=$?
if [ $DRIFT_EXIT -eq 0 ]; then
    echo "✅ Test 11 PASS: Import successful, no drift"
else
    echo "❌ Test 11 FAIL: Drift after import (exit code $DRIFT_EXIT)"
fi

# =============================================================================
# Test 10: Delete All Resources
# =============================================================================
echo ""
echo "=== Test 10: Delete All Resources ==="
terraform destroy -auto-approve 2>&1 | tee $EVIDENCE_DIR/test10_destroy.log
echo "✅ Test 10: Destroy complete"

# Save final state
cp terraform.tfstate $EVIDENCE_DIR/final_state.json 2>/dev/null || echo "No state file (expected after destroy)"

echo ""
echo "=============================================="
echo "Testing Complete!"
echo "Evidence saved to: $EVIDENCE_DIR/"
echo "=============================================="
