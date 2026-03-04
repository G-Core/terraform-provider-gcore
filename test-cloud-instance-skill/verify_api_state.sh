#!/bin/bash
cd /Users/user/repos/gcore-terraform/test-cloud-instance-skill

# Load environment without proxy
export GCORE_API_KEY='21788$1e278ce67b6aa33f178122658b1dd0210d0edff453d348acb9b68bffea6a635b7791925ddda198d5678a4dc20269fe04a263ca92c7e5aa41ea79075f89b66bf6'
export GCORE_CLIENT=3621
export GCORE_CLOUD_PROJECT_ID=379987
export GCORE_CLOUD_REGION_ID=76
export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/.terraformrc"
unset HTTP_PROXY
unset HTTPS_PROXY
export NO_PROXY="*"

echo "=== Verifying API state directly ==="
terraform refresh 2>&1 | tail -20

echo ""
echo "=== Current outputs after refresh ==="
terraform output
