#!/bin/bash
set -e

# Load credentials
if [ -f /Users/user/repos/gcore-terraform/.env ]; then
    set -o allexport
    source /Users/user/repos/gcore-terraform/.env
    set +o allexport
    echo "Credentials loaded"
else
    echo "ERROR: .env file not found"
    exit 1
fi

# Set terraform config
export TF_CLI_CONFIG_FILE=/Users/user/repos/gcore-terraform/.terraformrc
export TF_LOG=DEBUG
export TF_LOG_PATH=terraform_test.log

# Run terraform
echo "=== Running terraform plan ==="
terraform plan -out=tfplan 2>&1 | tee plan_output.txt

echo ""
echo "=== Running terraform apply ==="
terraform apply -auto-approve tfplan 2>&1 | tee apply_output.txt

echo ""
echo "=== Running drift test (second plan) ==="
terraform plan -detailed-exitcode 2>&1 | tee drift_output.txt
DRIFT_EXIT=$?
if [ $DRIFT_EXIT -eq 0 ]; then
    echo "DRIFT TEST PASSED: No changes detected"
elif [ $DRIFT_EXIT -eq 2 ]; then
    echo "DRIFT TEST FAILED: Changes detected!"
else
    echo "DRIFT TEST ERROR: Exit code $DRIFT_EXIT"
fi
