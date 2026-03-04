# GPU Virtual Image Test - GCLOUD2-21134

This directory contains test files for `gcore_gpu_virtual_image` resource migration from old provider to new generated provider.

## Analysis Summary

### Old Provider (G-Core/gcore)
- **Resource**: `gcore_gpu_virtual_image`
- **Operations**: Create, Read, Delete (NO Update - all attributes are ForceNew)
- **Create behavior**: Calls API, waits for task completion via `tasks.WaitTaskAndReturnResult`, extracts image ID from completed task
- **Delete behavior**: Calls API, does NOT wait for task completion

### New SDK (G-Core/gcore-go)
- **Service**: `GPUVirtualClusterImageService`
- **Methods**:
  - `Upload` - returns `TaskIDList` (async)
  - `UploadAndPoll` - uploads and polls task until completion, returns `GPUImage`
  - `Delete` - returns `TaskIDList` (async)
  - `DeleteAndPoll` - deletes and polls task until completion
  - `Get` - returns `GPUImage`
  - `List` - returns `GPUImageList`

### Current Generated Provider Issues
The generated provider (`gcore_cloud_gpu_baremetal_cluster_image`) does NOT use Poll methods:
- Create calls `Upload` but doesn't poll - resource ID is set to task ID, not image ID
- Delete calls `Delete` but doesn't poll - deletion may not complete

### Required Fix
Implement Poll methods support in the generated provider to:
1. Use `UploadAndPoll` in Create to wait for image creation and get real image ID
2. Use `DeleteAndPoll` in Delete to wait for deletion completion

## Running the Test

### Prerequisites
```bash
export GCORE_API_KEY="your-api-token"
export GCORE_CLOUD_PROJECT_ID=379987
export GCORE_CLOUD_REGION_ID=76
```

### Run
```bash
./run_test.sh
```

### Manual Steps
```bash
# Use old provider terraformrc (no dev overrides)
export TF_CLI_CONFIG_FILE="$(pwd)/.terraformrc.old"

# Create tfvars
cat > terraform.tfvars <<EOF
gcore_api_token = "$GCORE_API_KEY"
project_id      = $GCORE_CLOUD_PROJECT_ID
region_id       = $GCORE_CLOUD_REGION_ID
EOF

# Run terraform
terraform init
terraform plan
terraform apply
```

## Files for JIRA Ticket
After running, save these files:
- `main.tf` - Terraform configuration
- `terraform.tfstate` - State file showing resource was created
- `terraform.tfvars.example` - Example variables
