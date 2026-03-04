# Instance Resource Test

This directory contains tests for the `gcore_cloud_instance` resource after migrating to *AndPoll methods.

## Setup

1. Ensure environment variables are set in `.env` file at project root:
   ```
   GCORE_API_KEY=your_api_key
   GCORE_CLOUD_PROJECT_ID=379987
   GCORE_CLOUD_REGION_ID=76
   ```

2. Build the provider:
   ```bash
   cd /Users/user/repos/gcore-terraform
   go build -o terraform-provider-gcore
   ```

3. Ensure `.terraformrc` override is configured at project root

## Running Tests

### Test 1: Basic CRUD Operations

```bash
# Load environment and run plan
cd test-instance
set -o allexport; source ../.env; set +o allexport && \
export TF_CLI_CONFIG_FILE="../.terraformrc" && \
terraform plan

# Create instance (should wait for ACTIVE state)
terraform apply -auto-approve

# Verify instance is in ACTIVE state
terraform show | grep vm_state

# Delete instance (should wait for complete removal)
terraform destroy -auto-approve
```

### Test 2: Update Operations

After creating the instance, modify `main.tf` to test updates:

```bash
# Update name
# Change: name = "test-instance-andpoll-ORIGINAL"
# To:     name = "test-instance-andpoll-UPDATED"

terraform apply -auto-approve
# Should update in-place, not force replacement

# Update tags
# Add or modify tags in the tags block

terraform apply -auto-approve
# Should update in-place
```

### Test 3: Drift Detection

```bash
# After creating instance, manually change something via API/UI
# Then run:
terraform refresh
terraform plan
# Should detect the drift
```

## Expected Behavior

### After Migration to NewAndPoll:
- ✅ Instance should be in **ACTIVE** state after creation (not BUILD)
- ✅ Create operation should wait for instance to be fully provisioned
- ✅ Delete operation should wait for complete removal
- ✅ Name and tags updates should work without replacement
- ✅ Changes to flavor, interfaces, volumes should force replacement

### Changes Made:
1. Create: Uses `NewAndPoll()` instead of `New()` - waits for task completion
2. Delete: Uses `DeleteAndPoll()` instead of `Delete()` - waits for removal
3. Update: Uses SDK response directly (synchronous for name/tags)
4. Read: Uses SDK response directly with proper 404 handling

## Test Results

| Test | Status | Notes |
|------|--------|-------|
| Basic Create | ⏳ Pending | |
| Instance reaches ACTIVE | ⏳ Pending | |
| Read after create | ⏳ Pending | |
| Update name | ⏳ Pending | |
| Update tags | ⏳ Pending | |
| Drift detection | ⏳ Pending | |
| Delete with polling | ⏳ Pending | |
