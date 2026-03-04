#!/bin/bash
set -e

cd /Users/user/repos/gcore-terraform/test-cloud-instance-skill

# Load environment
export GCORE_API_KEY='21788$1e278ce67b6aa33f178122658b1dd0210d0edff453d348acb9b68bffea6a635b7791925ddda198d5678a4dc20269fe04a263ca92c7e5aa41ea79075f89b66bf6'
export GCORE_CLIENT=3621
export GCORE_CLOUD_PROJECT_ID=379987
export GCORE_CLOUD_REGION_ID=76
export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/.terraformrc"
export HTTP_PROXY="http://127.0.0.1:9092"
export HTTPS_PROXY="http://127.0.0.1:9092"
export NO_PROXY="registry.terraform.io,releases.hashicorp.com"
export SSL_CERT_FILE="/Users/user/repos/gcore-terraform/ca-bundle.pem"

mkdir -p evidence

echo "============================================"
echo "TEST 1: Create instance with external interface + existing volume"
echo "Testing: volume_id serialization (JIRA bug), interface creation, drift"
echo "============================================"

echo ""
echo "=== Step 1: Terraform Plan ==="
terraform plan -out=tfplan 2>&1 | tee evidence/test1_plan.log

echo ""
echo "=== Step 2: Terraform Apply ==="
terraform apply -auto-approve tfplan 2>&1 | tee evidence/test1_apply.log

echo ""
echo "=== Step 3: Drift Check (2nd Plan - MUST show no changes) ==="
set +e
terraform plan -detailed-exitcode 2>&1 | tee evidence/test1_drift.log
DRIFT_EXIT=$?
set -e

if [ $DRIFT_EXIT -eq 0 ]; then
    echo ""
    echo "✅ PASS: No drift detected - infrastructure matches configuration"
else
    echo ""
    echo "❌ FAIL: Drift detected! Exit code: $DRIFT_EXIT"
fi

echo ""
echo "=== Step 4: Output values ==="
terraform output 2>&1 | tee evidence/test1_outputs.log

echo ""
echo "============================================"
echo "TEST 1 COMPLETE"
echo "============================================"
