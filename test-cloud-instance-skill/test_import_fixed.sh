#!/bin/bash
set -e

cd /Users/user/repos/gcore-terraform/test-cloud-instance-skill

# Load environment
export GCORE_API_KEY='21788$1e278ce67b6aa33f178122658b1dd0210d0edff453d348acb9b68bffea6a635b7791925ddda198d5678a4dc20269fe04a263ca92c7e5aa41ea79075f89b66bf6'
export GCORE_CLIENT=3621
export GCORE_CLOUD_PROJECT_ID=379987
export GCORE_CLOUD_REGION_ID=76
export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/.terraformrc"
unset HTTP_PROXY
unset HTTPS_PROXY

mkdir -p evidence

# The instance from the previous test still exists
INSTANCE_ID="3b277ae2-5104-43d8-89cb-07fe692c6fbb"
VOLUME_ID="0d23dedf-8b9c-4001-9d36-efc0646e8338"

echo "============================================"
echo "IMPORT TEST (FIXED): Using correct import ID format"
echo "============================================"

echo ""
echo "Instance ID: $INSTANCE_ID"
echo "Import format: 379987/76/$INSTANCE_ID"

echo ""
echo "=== Step 1: Import instance with correct format ==="
terraform import gcore_cloud_instance.test_external "379987/76/$INSTANCE_ID" 2>&1 | tee evidence/import_fixed.log

echo ""
echo "=== Step 2: CRITICAL - Drift check after import ==="
set +e
terraform plan -detailed-exitcode 2>&1 | tee evidence/import_drift_fixed.log
DRIFT_EXIT=$?
set -e

if [ $DRIFT_EXIT -eq 0 ]; then
    echo ""
    echo "✅ PASS: No drift detected after import"
elif [ $DRIFT_EXIT -eq 2 ]; then
    echo ""
    echo "❌ FAIL: Drift detected after import! Exit code: $DRIFT_EXIT"
else
    echo ""
    echo "⚠️ Error running plan. Exit code: $DRIFT_EXIT"
fi

echo ""
echo "=== Step 3: Cleanup ==="
terraform destroy -auto-approve 2>&1 | tee evidence/import_cleanup_fixed.log

echo ""
echo "============================================"
echo "IMPORT TEST COMPLETE"
echo "============================================"
