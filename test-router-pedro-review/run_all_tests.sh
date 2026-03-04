#!/bin/bash
set -e

echo "=========================================="
echo "Router Testing - Pedro's Review Fixes"
echo "GCLOUD2-21144"
echo "=========================================="
echo ""

# Load credentials
source ../.env

# Set up Terraform config
export TF_CLI_CONFIG_FILE=/Users/user/repos/gcore-terraform/.terraformrc

# Clean start
rm -f terraform.tfstate* .terraform.lock.hcl
rm -f mitm_requests.log flow.mitm

echo "Test environment ready"
echo "  Project: $GCORE_CLOUD_PROJECT_ID"
echo "  Region: $GCORE_CLOUD_REGION_ID"
echo ""

# Function to run terraform and capture logs
run_test() {
  local test_num=$1
  local test_name=$2

  echo "=========================================="
  echo "Test #$test_num: $test_name"
  echo "=========================================="

  # Apply changes
  terraform apply -auto-approve 2>&1 | tee evidence/test_${test_num}_apply.log

  # Capture state
  terraform show -json > evidence/test_${test_num}_state.json

  # Critical: Drift check
  echo ""
  echo "Running drift check..."
  if terraform plan -detailed-exitcode > evidence/test_${test_num}_plan.log 2>&1; then
    echo "✅ PASS: No drift detected"
  else
    EXIT_CODE=$?
    if [ $EXIT_CODE -eq 2 ]; then
      echo "❌ FAIL: Drift detected!"
      cat evidence/test_${test_num}_plan.log
      exit 1
    fi
  fi

  # Save router ID for tracking
  ROUTER_ID=$(terraform output -raw router_id 2>/dev/null || echo "")
  echo "Router ID: $ROUTER_ID"

  # Copy MITM logs for this test
  if [ -f mitm_requests.log ]; then
    cp mitm_requests.log evidence/test_${test_num}_mitm.log
  fi

  echo ""
}

echo "Ready to run tests. Ensure mitmproxy is running in another terminal:"
echo "  Terminal 1: ./scripts/run_mitm.sh"
echo ""
read -p "Press Enter to continue..."
