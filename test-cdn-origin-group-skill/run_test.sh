#!/bin/bash
set -e

export GCORE_API_KEY='21788$1e278ce67b6aa33f178122658b1dd0210d0edff453d348acb9b68bffea6a635b7791925ddda198d5678a4dc20269fe04a263ca92c7e5aa41ea79075f89b66bf6'
export HTTP_PROXY="http://127.0.0.1:9092"
export HTTPS_PROXY="http://127.0.0.1:9092"
export NO_PROXY="registry.terraform.io,releases.hashicorp.com"

cd /Users/user/repos/gcore-terraform/test-cdn-origin-group-skill

case "$1" in
  tc1)
    echo "=== TC-1: Basic origin group with sources ==="
    terraform apply -auto-approve -var="run_tc1=true" 2>&1 | tee evidence/tc1_apply.log
    echo ""
    echo "=== TC-1: Drift check ==="
    terraform plan -detailed-exitcode -var="run_tc1=true" 2>&1 | tee evidence/tc1_drift.log
    ;;
  tc2)
    echo "=== TC-2: S3 auth with s3_type=other ==="
    terraform apply -auto-approve -var="run_tc2=true" 2>&1 | tee evidence/tc2_apply.log
    echo ""
    echo "=== TC-2: Drift check ==="
    terraform plan -detailed-exitcode -var="run_tc2=true" 2>&1 | tee evidence/tc2_drift.log
    echo ""
    echo "=== TC-2: Verify s3_credentials_version in state ==="
    terraform show -json | jq '.values.root_module.resources[] | select(.type=="gcore_cdn_origin_group") | .values.auth'
    ;;
  tc3)
    echo "=== TC-3: S3 auth with s3_type=amazon ==="
    terraform apply -auto-approve -var="run_tc3=true" 2>&1 | tee evidence/tc3_apply.log
    echo ""
    echo "=== TC-3: Drift check ==="
    terraform plan -detailed-exitcode -var="run_tc3=true" 2>&1 | tee evidence/tc3_drift.log
    echo ""
    echo "=== TC-3: Verify s3_credentials_version in state ==="
    terraform show -json | jq '.values.root_module.resources[] | select(.type=="gcore_cdn_origin_group") | .values.auth'
    ;;
  tc4)
    echo "=== TC-4: Update test - initial create ==="
    terraform apply -auto-approve -var="run_tc4=true" 2>&1 | tee evidence/tc4_create.log
    echo ""
    BEFORE_ID=$(terraform show -json | jq -r '.values.root_module.resources[] | select(.type=="gcore_cdn_origin_group" and .name=="tc4_update") | .values.id')
    echo "ID before update: $BEFORE_ID"
    echo ""
    echo "=== TC-4: Update sources and name ==="
    terraform apply -auto-approve -var="run_tc4=true" \
      -var='tc4_name=test-cdn-og-tc4-updated-skill' \
      -var='tc4_proxy_next_upstream=["error", "timeout", "http_500"]' \
      -var='tc4_sources=[{source="93.184.216.41", enabled=true, backup=false}]' 2>&1 | tee evidence/tc4_update.log
    AFTER_ID=$(terraform show -json | jq -r '.values.root_module.resources[] | select(.type=="gcore_cdn_origin_group" and .name=="tc4_update") | .values.id')
    echo "ID after update: $AFTER_ID"
    if [ "$BEFORE_ID" = "$AFTER_ID" ]; then
      echo "PASS: Resource updated in-place (same ID)"
    else
      echo "FAIL: Resource was recreated (ID changed)"
    fi
    ;;
  tc5)
    echo "=== TC-5: Credential version bump ==="
    terraform apply -auto-approve -var="run_tc2=true" -var="tc2_credentials_version=2" 2>&1 | tee evidence/tc5_bump.log
    echo ""
    echo "=== TC-5: Verify new version in state ==="
    terraform show -json | jq '.values.root_module.resources[] | select(.type=="gcore_cdn_origin_group" and .name=="tc2_s3_other") | .values.auth.s3_credentials_version'
    ;;
  import)
    echo "=== TC-6: Import test ==="
    IMPORT_ID=$2
    if [ -z "$IMPORT_ID" ]; then
      echo "Usage: $0 import <origin_group_id>"
      exit 1
    fi
    # First destroy to clear state
    terraform state rm gcore_cdn_origin_group.tc1_public[0] 2>/dev/null || true
    # Then import
    terraform import -var="run_tc1=true" 'gcore_cdn_origin_group.tc1_public[0]' "$IMPORT_ID" 2>&1 | tee evidence/tc6_import.log
    terraform plan -var="run_tc1=true" 2>&1 | tee evidence/tc6_plan_after_import.log
    ;;
  state)
    echo "=== Current state ==="
    terraform show -json | jq '.values.root_module.resources[] | select(.type=="gcore_cdn_origin_group")'
    ;;
  destroy)
    echo "=== Destroying all resources ==="
    terraform destroy -auto-approve -var="run_tc1=true" -var="run_tc2=true" -var="run_tc3=true" -var="run_tc4=true" 2>&1 | tee evidence/destroy.log
    ;;
  *)
    echo "Usage: $0 {tc1|tc2|tc3|tc4|tc5|import <id>|state|destroy}"
    exit 1
    ;;
esac
