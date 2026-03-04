#!/bin/bash
set -e

echo "=========================================="
echo "Router Configuration Drift Test"
echo "=========================================="
echo

# Load environment variables
echo "Loading environment variables..."
cd /Users/user/repos/gcore-terraform/test-router-drift
set -o allexport
source ../.env
set +o allexport
export TF_CLI_CONFIG_FILE="../.terraformrc"

echo "=========================================="
echo "STEP 1: First Apply (Create Infrastructure)"
echo "=========================================="
echo

terraform apply -auto-approve 2>&1 | tee apply1.log

echo
echo "First apply completed. Saving state checkpoint..."
cp terraform.tfstate terraform.tfstate.apply1
echo

echo "=========================================="
echo "STEP 2: Second Apply (Drift Detection)"
echo "=========================================="
echo "This should show 'No changes' if there is no drift."
echo

terraform apply -auto-approve 2>&1 | tee apply2.log

echo
echo "Second apply completed. Saving state checkpoint..."
cp terraform.tfstate terraform.tfstate.apply2
echo

echo "=========================================="
echo "STEP 3: Analysis"
echo "=========================================="
echo

# Check if second apply had changes
if grep -q "No changes" apply2.log; then
    echo "✅ SUCCESS: No configuration drift detected"
    echo "   Second apply showed 'No changes'"
    DRIFT_DETECTED=0
elif grep -q "Apply complete! Resources: 0 added, 0 changed, 0 destroyed" apply2.log; then
    echo "✅ SUCCESS: No configuration drift detected"
    echo "   Second apply made no changes"
    DRIFT_DETECTED=0
else
    echo "⚠️  WARNING: Configuration drift detected!"
    echo "   Second apply made changes"
    DRIFT_DETECTED=1

    # Show what changed
    echo
    echo "Changes detected:"
    grep -A 5 "Terraform will perform the following actions:" apply2.log || true
fi

echo
echo "=========================================="
echo "STEP 4: Detailed State Comparison"
echo "=========================================="
echo

# Compare states
echo "Comparing router states between first and second apply..."
echo

echo "Router: router_basic"
echo "-------------------"
diff <(terraform show -json terraform.tfstate.apply1 2>/dev/null | python3 -c "import sys, json; d=json.load(sys.stdin); r=[x for x in d['values']['root_module']['resources'] if x['address']=='gcore_cloud_network_router.router_basic']; print(json.dumps(r[0]['values'], indent=2, sort_keys=True) if r else 'Not found')") \
     <(terraform show -json terraform.tfstate.apply2 2>/dev/null | python3 -c "import sys, json; d=json.load(sys.stdin); r=[x for x in d['values']['root_module']['resources'] if x['address']=='gcore_cloud_network_router.router_basic']; print(json.dumps(r[0]['values'], indent=2, sort_keys=True) if r else 'Not found')") \
     || echo "  No differences"
echo

echo "Router: router_interfaces"
echo "-------------------------"
diff <(terraform show -json terraform.tfstate.apply1 2>/dev/null | python3 -c "import sys, json; d=json.load(sys.stdin); r=[x for x in d['values']['root_module']['resources'] if x['address']=='gcore_cloud_network_router.router_interfaces']; print(json.dumps(r[0]['values'], indent=2, sort_keys=True) if r else 'Not found')") \
     <(terraform show -json terraform.tfstate.apply2 2>/dev/null | python3 -c "import sys, json; d=json.load(sys.stdin); r=[x for x in d['values']['root_module']['resources'] if x['address']=='gcore_cloud_network_router.router_interfaces']; print(json.dumps(r[0]['values'], indent=2, sort_keys=True) if r else 'Not found')") \
     || echo "  No differences"
echo

echo "Router: router_routes"
echo "---------------------"
diff <(terraform show -json terraform.tfstate.apply1 2>/dev/null | python3 -c "import sys, json; d=json.load(sys.stdin); r=[x for x in d['values']['root_module']['resources'] if x['address']=='gcore_cloud_network_router.router_routes']; print(json.dumps(r[0]['values'], indent=2, sort_keys=True) if r else 'Not found')") \
     <(terraform show -json terraform.tfstate.apply2 2>/dev/null | python3 -c "import sys, json; d=json.load(sys.stdin); r=[x for x in d['values']['root_module']['resources'] if x['address']=='gcore_cloud_network_router.router_routes']; print(json.dumps(r[0]['values'], indent=2, sort_keys=True) if r else 'Not found')") \
     || echo "  No differences"
echo

echo "Router: router_complete"
echo "-----------------------"
diff <(terraform show -json terraform.tfstate.apply1 2>/dev/null | python3 -c "import sys, json; d=json.load(sys.stdin); r=[x for x in d['values']['root_module']['resources'] if x['address']=='gcore_cloud_network_router.router_complete']; print(json.dumps(r[0]['values'], indent=2, sort_keys=True) if r else 'Not found')") \
     <(terraform show -json terraform.tfstate.apply2 2>/dev/null | python3 -c "import sys, json; d=json.load(sys.stdin); r=[x for x in d['values']['root_module']['resources'] if x['address']=='gcore_cloud_network_router.router_complete']; print(json.dumps(r[0]['values'], indent=2, sort_keys=True) if r else 'Not found')") \
     || echo "  No differences"
echo

echo "=========================================="
echo "STEP 5: Summary"
echo "=========================================="
echo

if [ $DRIFT_DETECTED -eq 0 ]; then
    echo "✅ DRIFT TEST PASSED"
    echo ""
    echo "Summary:"
    echo "  - First apply: Created 4 routers"
    echo "  - Second apply: No changes detected"
    echo "  - Configuration is stable"
else
    echo "❌ DRIFT TEST FAILED"
    echo ""
    echo "Summary:"
    echo "  - First apply: Created 4 routers"
    echo "  - Second apply: Detected changes"
    echo "  - Configuration drift exists"
    echo ""
    echo "Review apply2.log for details"
fi

echo
echo "Test artifacts saved:"
echo "  - apply1.log: First apply output"
echo "  - apply2.log: Second apply output"
echo "  - terraform.tfstate.apply1: State after first apply"
echo "  - terraform.tfstate.apply2: State after second apply"
echo

exit $DRIFT_DETECTED
