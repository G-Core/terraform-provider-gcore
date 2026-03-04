# Old Provider Comparison Testing

## Purpose

Compare Terraform plan outputs and state files between the old provider (`gcore_instancev2`, `gcore_volume`, etc.) and new provider (`gcore_cloud_instance`, `gcore_cloud_volume`, etc.) to verify:

1. **API calls match** - Both providers call the same endpoints
2. **State structure is compatible** - Migration path is smooth
3. **Computed fields behave consistently** - No unexpected drift
4. **Update operations work identically** - Flavor changes, volume resizes, etc.

## When to Use Old Provider Comparison

Use this approach when:
- ✅ Testing critical update operations (flavor change, volume resize)
- ✅ Verifying migration from old to new provider
- ✅ Investigating drift or unexpected behavior
- ✅ Documenting behavioral differences for users
- ✅ Validating that new provider doesn't break existing workflows

## Setup

### 1. Old Provider Location

The old provider is located at:
```bash
old_terraform_provider/
├── terraform-provider-gcore  # Pre-built binary
├── gcore/                    # Source code
│   ├── resource_gcore_instance.go
│   ├── resource_gcore_instancev2.go
│   ├── resource_gcore_volume.go
│   └── ...
└── README.md
```

### 2. Configure dev_overrides for Old Provider

Create a `.terraformrc` file pointing to the old provider:

```hcl
provider_installation {
  dev_overrides {
    "gcore/gcore" = "/Users/user/repos/gcore-terraform/old_terraform_provider"
  }
  direct {}
}
```

### 3. Test Directory Structure

Create a comparison test directory:

```bash
test-old-provider-comparison/
├── .terraformrc              # Points to old provider
├── main.tf                   # Test configuration
├── terraform.tfvars          # Credentials
├── old_provider.tf           # Copy of main.tf for reference
├── old_provider.tfstate      # State file saved for comparison
└── terraform.tfstate         # Current state
```

## Testing Workflow

### Step 1: Create Infrastructure with Old Provider

```bash
cd test-old-provider-comparison

# Set environment variables
export TF_CLI_CONFIG_FILE="$(pwd)/.terraformrc"

# Create terraform.tfvars with credentials
cat > terraform.tfvars <<EOF
gcore_api_key = "$GCORE_API_KEY"
EOF

# Create test configuration
cat > main.tf <<'EOF'
terraform {
  required_providers {
    gcore = {
      source  = "gcore/gcore"
    }
  }
}

provider "gcore" {
  permanent_api_token = var.gcore_api_key
}

variable "gcore_api_key" {
  type      = string
  sensitive = true
}

variable "project_id" {
  type    = number
  default = 379987
}

variable "region_id" {
  type    = number
  default = 76
}

variable "flavor_id" {
  description = "Instance flavor ID"
  type        = string
  default     = "g1-standard-1-2"
}

variable "volume_size" {
  description = "Boot volume size in GiB"
  type        = number
  default     = 15
}

# Old provider uses separate volume resource
resource "gcore_volume" "test_boot_volume" {
  name       = "test-old-provider-boot"
  type_name  = "standard"
  size       = var.volume_size
  image_id   = "8a59956d-ada3-42c6-aaf3-e2c38ba70f99" # ubuntu-24.04-x64
  project_id = var.project_id
  region_id  = var.region_id
}

# Old provider uses gcore_instancev2
resource "gcore_instancev2" "test_instance" {
  name       = "test-old-provider"
  flavor_id  = var.flavor_id
  project_id = var.project_id
  region_id  = var.region_id

  volume {
    volume_id = gcore_volume.test_boot_volume.id
  }

  interface {
    type            = "external"
    name            = "eth0"
    security_groups = []
  }

  metadata_map = {
    purpose = "old-provider-comparison"
  }
}

output "instance_id" {
  value = gcore_instancev2.test_instance.id
}

output "instance_flavor_id" {
  value = gcore_instancev2.test_instance.flavor_id
}

output "volume_id" {
  value = gcore_volume.test_boot_volume.id
}

output "volume_size" {
  value = gcore_volume.test_boot_volume.size
}
EOF

# Apply configuration
terraform apply -auto-approve
```

### Step 2: Capture Plan Outputs for Changes

**Test Flavor Change:**
```bash
echo "=== OLD PROVIDER - FLAVOR CHANGE ===" > comparison_output.txt
terraform plan -no-color -var="flavor_id=g1-standard-2-4" >> comparison_output.txt 2>&1
```

**Test Volume Resize:**
```bash
echo "=== OLD PROVIDER - VOLUME RESIZE ===" >> comparison_output.txt
terraform plan -no-color -var="volume_size=25" >> comparison_output.txt 2>&1
```

### Step 3: Save State and Configuration

```bash
# Save state file for later comparison
cp terraform.tfstate old_provider.tfstate

# Save configuration for reference
cp main.tf old_provider.tf

# Important: List what's in the state
terraform show -no-color > old_provider_show.txt
```

### Step 4: Compare with New Provider

Create equivalent test with new provider:

```bash
cd ../test-new-provider-comparison

cat > main.tf <<'EOF'
terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {
  api_key = var.gcore_api_key
}

variable "gcore_api_key" {
  type      = string
  sensitive = true
}

variable "project_id" {
  type    = number
  default = 379987
}

variable "region_id" {
  type    = number
  default = 76
}

variable "flavor" {
  description = "Instance flavor"
  type        = string
  default     = "g1-standard-1-2"
}

variable "volume_size" {
  description = "Boot volume size in GiB"
  type        = number
  default     = 15
}

# New provider integrates volumes into instance
resource "gcore_cloud_instance" "test_instance" {
  project_id = var.project_id
  region_id  = var.region_id
  name       = "test-new-provider"
  flavor     = var.flavor

  boot_volume {
    size      = var.volume_size
    source    = "image"
    image_id  = "8a59956d-ada3-42c6-aaf3-e2c38ba70f99"
    type_name = "standard"
  }

  network_interface {
    type = "external"
  }
}

output "instance_id" {
  value = gcore_cloud_instance.test_instance.id
}

output "instance_flavor" {
  value = gcore_cloud_instance.test_instance.flavor
}

output "volume_size" {
  value = gcore_cloud_instance.test_instance.volumes[0].size
}
EOF

# Use new provider config
export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/.terraformrc"

terraform apply -auto-approve

# Capture same tests
echo "=== NEW PROVIDER - FLAVOR CHANGE ===" > comparison_output.txt
terraform plan -no-color -var="flavor=g1-standard-2-4" >> comparison_output.txt 2>&1

echo "=== NEW PROVIDER - VOLUME RESIZE ===" >> comparison_output.txt
terraform plan -no-color -var="volume_size=25" >> comparison_output.txt 2>&1

cp terraform.tfstate new_provider.tfstate
terraform show -no-color > new_provider_show.txt
```

## Comparison Points

### 1. Terraform Plan Output

**What to Compare:**
```bash
# Side-by-side comparison
diff -u old_provider_flavor_plan.txt new_provider_flavor_plan.txt

# Key indicators:
# ✅ Both show "update in-place" (not replacement)
# ✅ Both show computed attributes as "(known after apply)"
# ✅ Resource IDs remain stable
# ✅ No "forces replacement" annotations
```

**Expected Similarities:**
- Both providers show `~ update in-place`
- Computed fields show `(known after apply)` in both
- Interface/network fields refresh in both
- No instance replacement in either case

**Acceptable Differences:**
- Old provider: separate `gcore_volume` resource
- New provider: integrated `volumes` block
- Old provider: `flavor_id` field
- New provider: `flavor` field
- Field names may differ but behavior should match

### 2. State File Structure

**Compare State Contents:**
```bash
# Extract resource data
jq '.resources[] | select(.type | contains("instance"))' old_provider.tfstate > old_instance_state.json
jq '.resources[] | select(.type | contains("instance"))' new_provider.tfstate > new_instance_state.json

# Compare computed vs configured fields
jq '.resources[].instances[].attributes' old_provider.tfstate
jq '.resources[].instances[].attributes' new_provider.tfstate
```

**Key Fields to Compare:**
- Instance ID (should be different but stable across updates)
- Flavor/flavor_id values
- Volume size and configuration
- Computed fields (addresses, status, vm_state)
- Network interfaces

### 3. API Calls (with mitmproxy)

If using mitmproxy capture:

```bash
# Old provider API calls
mitmdump -r old_provider_flow.mitm | grep -E "(POST|PATCH|PUT|DELETE)" > old_api_calls.txt

# New provider API calls
mitmdump -r new_provider_flow.mitm | grep -E "(POST|PATCH|PUT|DELETE)" > new_api_calls.txt

# Compare endpoints used
diff -u old_api_calls.txt new_api_calls.txt
```

**Expected Matches:**
- ✅ Same endpoints for create (POST /instances)
- ✅ Same endpoints for flavor change (POST /changeflavor)
- ✅ Same endpoints for volume resize (POST /extend)
- ✅ Same endpoint for delete (DELETE /instances)

## Common Findings

### Normal Behavior (Both Providers)

Both old and new providers exhibit these behaviors:

1. **Computed Attributes Refresh**
   - Fields like `addresses`, `status`, `vm_state` show `(known after apply)`
   - This is normal - providers re-read state after update operations
   - Interface blocks may be removed/added during updates

2. **Update In-Place**
   - Flavor changes: `~ flavor_id` or `~ flavor`
   - Volume resizes: `~ size`
   - No "forces replacement" or "destroy and create"

3. **Stable Resource IDs**
   - Instance ID doesn't change during updates
   - Only changes during full recreation

### Differences to Document

| Aspect | Old Provider | New Provider |
|--------|--------------|--------------|
| **Volume Management** | Separate `gcore_volume` resource | Nested `volumes` block |
| **Flavor Field** | `flavor_id` (string) | `flavor` (string) |
| **Resource Count** | 2 resources (volume + instance) | 1 resource (instance with volumes) |
| **Interface Config** | May have drift issues | Better drift handling |

## Example Comparison Report

```markdown
# Old vs New Provider Comparison - Instance Resize Operations

## Test Scenario
- **Operation**: Flavor change (g1-standard-1-2 → g1-standard-2-4)
- **Old Provider**: gcore_instancev2
- **New Provider**: gcore_cloud_instance

## Infrastructure Created
- **Old Provider Instance**: `1ff550f0-c515-43d0-889d-f98c1ee90df4`
- **New Provider Instance**: `b634def3-4c2e-42e7-acf8-d0b9341b28cc`
- **Old Provider Volume**: `06d10d4a-ec92-4efc-b9f9-a683e5abe241` (separate resource)
- **New Provider Volume**: Integrated in instance

## Terraform Plan Comparison

### Flavor Change

**Old Provider Output:**
```
  ~ resource "gcore_instancev2" "test_instance" {
      ~ flavor_id = "g1-standard-1-2" -> "g1-standard-2-4"

      - interface {
          - ip_address = "85.204.246.136" -> null
          ...
        }
      + interface {
          + ip_address = (known after apply)
          ...
        }
    }

Plan: 0 to add, 1 to change, 0 to destroy.
```

**New Provider Output:**
```
  ~ resource "gcore_cloud_instance" "test_instance" {
      ~ flavor = "g1-standard-1-2" -> "g1-standard-2-4"

      ~ addresses = (known after apply)
      ~ status = (known after apply)
      ~ vm_state = (known after apply)
    }

Plan: 0 to add, 1 to change, 0 to destroy.
```

## Findings

### ✅ Similarities (Correct Behavior)
1. Both show `update in-place`
2. Both preserve instance ID
3. Both show computed attributes as `(known after apply)`
4. Both complete successfully without errors

### ⚠️ Differences (Expected)
1. Field name: `flavor_id` vs `flavor`
2. Old provider shows interface drift
3. New provider shows more computed fields refreshing
4. Old provider manages volumes separately

### API Verification
- Both use `POST /v1/instances/{id}/resize` endpoint
- Both poll for task completion
- Both re-read instance state after operation

## Conclusion

**New provider behavior matches old provider** - both perform in-place updates correctly. Differences in terraform output are cosmetic (computed field display) and don't affect functionality.

User concern about `(known after apply)` is **NORMAL** - old provider shows the same pattern.
```

## Cleanup

```bash
# Destroy old provider resources
cd test-old-provider-comparison
export TF_CLI_CONFIG_FILE="$(pwd)/.terraformrc"
terraform destroy -auto-approve

# Destroy new provider resources
cd ../test-new-provider-comparison
export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/.terraformrc"
terraform destroy -auto-approve
```

## Troubleshooting

### Old Provider Not Found

```bash
# Check old provider binary exists
ls -la old_terraform_provider/terraform-provider-gcore

# Verify .terraformrc path
cat .terraformrc

# The path should be absolute and point to old_terraform_provider directory
```

### Quota Exceeded

```bash
# Clean up orphaned volumes first
curl -X GET -H "Authorization: APIKey $GCORE_API_KEY" \
  "https://api.gcore.com/cloud/v1/volumes/379987/76?has_attachments=false" | jq

# Delete unused volumes
curl -X DELETE -H "Authorization: APIKey $GCORE_API_KEY" \
  "https://api.gcore.com/cloud/v1/volumes/379987/76/{volume_id}"
```

### Provider Version Conflicts

If terraform complains about duplicate providers:
```bash
# Make sure only one .tf file exists in test directory
ls -la *.tf

# Remove old_provider.tf if running tests (keep for reference only)
mv old_provider.tf old_provider.tf.bak
```
