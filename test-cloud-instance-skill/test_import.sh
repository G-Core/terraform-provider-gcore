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

echo "============================================"
echo "IMPORT TEST: Create -> Remove from state -> Import -> Drift check"
echo "============================================"

echo ""
echo "=== Step 1: Create fresh instance ==="
terraform apply -auto-approve 2>&1 | tee evidence/import_create.log

echo ""
echo "=== Step 2: Get instance ID ==="
INSTANCE_ID=$(terraform output -raw instance_id)
VOLUME_ID=$(terraform output -raw boot_volume_id)
echo "Instance ID: $INSTANCE_ID"
echo "Volume ID: $VOLUME_ID"

echo ""
echo "=== Step 3: Remove instance from state (simulate import scenario) ==="
terraform state rm gcore_cloud_instance.test_external 2>&1 | tee evidence/import_state_rm.log

echo ""
echo "=== Step 4: Import instance back ==="
# Import format: project_id/region_id/instance_id
terraform import gcore_cloud_instance.test_external "379987/76/$INSTANCE_ID" 2>&1 | tee evidence/import_import.log

echo ""
echo "=== Step 5: CRITICAL - Drift check after import (MUST show no changes) ==="
set +e
terraform plan -detailed-exitcode 2>&1 | tee evidence/import_drift.log
DRIFT_EXIT=$?
set -e

if [ $DRIFT_EXIT -eq 0 ]; then
    echo ""
    echo "✅ PASS: No drift detected after import"
else
    echo ""
    echo "❌ FAIL: Drift detected after import! Exit code: $DRIFT_EXIT"
fi

echo ""
echo "=== Step 6: Cleanup ==="
terraform destroy -auto-approve 2>&1 | tee evidence/import_cleanup.log

echo ""
echo "============================================"
echo "IMPORT TEST COMPLETE"
echo "============================================"
