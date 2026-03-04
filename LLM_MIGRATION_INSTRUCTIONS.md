# LLM Instructions for Gcore Terraform Provider Migration to *AndPoll Methods

## Overview

This document provides comprehensive instructions for migrating Gcore Terraform provider resources from manual HTTP response handling to *AndPoll methods, along with analysis of old vs new provider implementations.

## Identifying API Operation Types

### Synchronous vs Asynchronous Operations

**IMPORTANT:** Before migrating a resource, determine if the API operations are synchronous or asynchronous.

#### Synchronous Operations (No Tasks)

- API returns the resource data **immediately** in the response
- No task IDs are returned
- No polling required
- Examples: SSH keys, some simple CRUD operations

**Pattern for Synchronous Operations:**

```go
// Create
sshKey, err := r.client.Cloud.SSHKeys.New(ctx, params, ...)
if err != nil {
    resp.Diagnostics.AddError("failed to make http request", err.Error())
    return
}
err = apijson.UnmarshalComputed([]byte(sshKey.RawJSON()), &data)

// Update
sshKey, err := r.client.Cloud.SSHKeys.Update(ctx, id, params, ...)

// Delete
err := r.client.Cloud.SSHKeys.Delete(ctx, id, params, ...)
```

#### Asynchronous Operations (Returns Tasks)

- API returns task IDs that must be polled
- Resource may not be immediately available
- Requires waiting for task completion
- Examples: Load balancers, listeners, pools, volumes, instances

**Pattern for Asynchronous Operations:**

```go
// Create
listener, err := r.client.Cloud.LoadBalancers.Listeners.NewAndPoll(ctx, params, ...)

// Update
listener, err := r.client.Cloud.LoadBalancers.Listeners.UpdateAndPoll(ctx, id, params, ...)

// Delete
err := r.client.Cloud.LoadBalancers.Listeners.DeleteAndPoll(ctx, id, params, ...)
```

#### How to Determine Operation Type

**1. Check the old provider implementation:**

```go
// Asynchronous - Uses tasks.WaitTaskAndReturnResult()
results, err := listeners.Create(client, opts).Extract()
taskID := results.Tasks[0]
listenerID, err := tasks.WaitTaskAndReturnResult(client, taskID, ...)

// Synchronous - Direct extract without task waiting
kp, err := keypairs.Create(client, opts).Extract()
d.SetId(kp.ID)  // No task waiting
```

**2. Check SDK method signatures:**

```bash
# Check what the SDK method returns
go doc github.com/G-Core/gcore-go/cloud.SSHKeyService

# If it returns the resource directly, it's synchronous:
func (r *SSHKeyService) New(...) (res *SSHKeyCreated, err error)

# If it has an *AndPoll variant, it's asynchronous:
func (r *ListenerService) NewAndPoll(...) (res *Listener, err error)
```

**3. Check API documentation:**

- Look for "returns task" or "asynchronous operation" in API docs
- Synchronous operations return the resource in the response body immediately
- Asynchronous operations return task IDs to poll

#### Quick Reference: Known Resource Types

**Synchronous (No *AndPoll needed):**

- SSH Keys (`gcore_cloud_ssh_key`) - [Test example](test-ssh-key/main.tf)
- FaaS Keys
- Storage SFTP Keys

**Asynchronous (Needs *AndPoll):**

- Load Balancers (`gcore_cloud_load_balancer`)
- Load Balancer Listeners (`gcore_cloud_load_balancer_listener`)
- Load Balancer Pools (`gcore_cloud_load_balancer_pool`)
- Volumes
- Instances
- Networks
- Subnets

### Migration Decision Tree

```ini
Is this resource using manual HTTP response handling?
├── Yes
│   ├── Does the old provider use tasks.WaitTaskAndReturnResult()?
│   │   ├── Yes → Use *AndPoll methods (see section 1 below)
│   │   └── No → Use direct SDK methods without *AndPoll (see section 1.1 below)
│   └── Already using SDK methods → No changes needed
└── No → No migration needed
```

## Core Migration Requirements

### 1. Replace Manual HTTP Response Handling with *AndPoll Methods (Asynchronous APIs)

**Old Pattern (Manual HTTP handling):**

```go
res := new(http.Response)
_, err = r.client.Cloud.LoadBalancers.Listeners.New(
    ctx,
    params,
    option.WithRequestBody("application/json", dataBytes),
    option.WithResponseBodyInto(&res),
    option.WithMiddleware(logging.Middleware(ctx)),
)
if err != nil {
    resp.Diagnostics.AddError("failed to make http request", err.Error())
    return
}
bytes, _ := io.ReadAll(res.Body)
err = apijson.UnmarshalComputed(bytes, &data)
```

**New Pattern (*AndPoll methods):**

```go
listener, err := r.client.Cloud.LoadBalancers.Listeners.NewAndPoll(
    ctx,
    params,
    option.WithRequestBody("application/json", dataBytes),
    option.WithMiddleware(logging.Middleware(ctx)),
)
if err != nil {
    resp.Diagnostics.AddError("failed to make http request", err.Error())
    return
}
err = apijson.UnmarshalComputed([]byte(listener.RawJSON()), &data)
```

### 1.1. Replace Manual HTTP Response Handling with Direct SDK Methods (Synchronous APIs)

**Old Pattern (Manual HTTP handling):**

```go
res := new(http.Response)
_, err = r.client.Cloud.SSHKeys.New(
    ctx,
    params,
    option.WithRequestBody("application/json", dataBytes),
    option.WithResponseBodyInto(&res),
    option.WithMiddleware(logging.Middleware(ctx)),
)
if err != nil {
    resp.Diagnostics.AddError("failed to make http request", err.Error())
    return
}
bytes, _ := io.ReadAll(res.Body)
err = apijson.UnmarshalComputed(bytes, &data)
```

**New Pattern (Direct SDK response):**

```go
sshKey, err := r.client.Cloud.SSHKeys.New(
    ctx,
    params,
    option.WithRequestBody("application/json", dataBytes),
    option.WithMiddleware(logging.Middleware(ctx)),
)
if err != nil {
    resp.Diagnostics.AddError("failed to make http request", err.Error())
    return
}
err = apijson.UnmarshalComputed([]byte(sshKey.RawJSON()), &data)
```

**Key differences from *AndPoll pattern:**

- Use `New()`, `Update()`, `Delete()` directly (no `AndPoll` suffix)
- SDK returns the resource object immediately (e.g., `*SSHKeyCreated`, `*SSHKey`)
- No task polling needed
- Still remove `option.WithResponseBodyInto(&res)`
- Still use `resource.RawJSON()` instead of `io.ReadAll(res.Body)`

### 2. Key Changes Required

#### For Asynchronous APIs:

**Create Method:**

- Replace `New()` with `NewAndPoll()`
- Remove `option.WithResponseBodyInto(&res)`
- Replace `io.ReadAll(res.Body)` with `[]byte(listener.RawJSON())`
- Remove manual HTTP response handling

**Update Method:**

- Replace `Update()` with `UpdateAndPoll()`
- Same response handling changes as Create

**Delete Method:**

- Replace `Delete()` with `DeleteAndPoll()`
- Remove manual response body reading
- *AndPoll methods handle async operations automatically

#### For Synchronous APIs:

**Create Method:**

- Keep `New()` (no *AndPoll suffix)
- Remove `option.WithResponseBodyInto(&res)`
- Remove `io.ReadAll(res.Body)`
- Use SDK response object directly: `[]byte(resource.RawJSON())`

**Update Method:**

- Keep `Update()` (no *AndPoll suffix)
- Same response handling changes as Create

**Delete Method:**

- Keep `Delete()` (no *AndPoll suffix)
- Remove manual response body reading
- Returns error directly (no resource object)

### 3. Response Processing Changes

**Avoid:**

- `apijson.MarshalRoot()` - deprecated serialization method
- Manual `io.ReadAll(res.Body)`
- Manual HTTP response object handling
- Direct task polling logic

**Use:**

- `listener.RawJSON()` for response data extraction
- `apijson.UnmarshalComputed()` for deserialization
- Built-in async operation polling in *AndPoll methods

## Analysis Task for LLMs

When analyzing resources for migration, LLMs must:

### 1. Explore Similar Resources

- Find corresponding resource in old Terraform provider (`old_terraform_provider/gcore/`)
- Compare implementation patterns between old and new providers
- Identify special logic that needs to be transferred

### 2. Special Logic Identification

Look for these patterns that require careful migration:

#### Update Strategies:

- **Old provider pattern:** Complex update logic with separate `Update()` and `Unset()` calls

```go
// Example from old lblistener
if changed {
    results, err := listeners.Update(clientV2, d.Id(), updateOpts, ...)
    // Wait for task completion
    if toUnset {
        // Additional unset operation for clearing fields
        results, err := listeners.Unset(clientV2, d.Id(), unsetOpts, ...)
    }
}
```

- **New provider pattern:** Simplified with *AndPoll methods

```go
listener, err := r.client.Cloud.LoadBalancers.Listeners.UpdateAndPoll(
    ctx,
    data.ID.ValueString(),
    params,
    option.WithRequestBody("application/json", dataBytes),
    option.WithMiddleware(logging.Middleware(ctx)),
)
```

#### Alternative Endpoints:

- Check if old provider uses different API versions (v1 vs v2)
- Identify cases where multiple API calls are needed
- Look for specialized update endpoints that avoid resource recreation

#### Conflict Retry Logic:

- Old provider: `ConflictRetryAmount` and `ConflictRetryInterval`
- New provider: Built into *AndPoll methods

### 3. Task Management Differences

**Old Provider (Manual Task Handling):**

```go
taskID := results.Tasks[0]
listenerID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, timeout, func(task tasks.TaskID) (interface{}, error) {
    taskInfo, err := tasks.Get(client, string(task)).Extract()
    if err != nil {
        return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
    }
    listenerID, err := listeners.ExtractListenerIDFromTask(taskInfo)
    if err != nil {
        return nil, fmt.Errorf("cannot retrieve LBListener ID from task info: %w", err)
    }
    return listenerID, nil
})
```

**New Provider (Automatic with *AndPoll):**

```go
// Task polling is handled automatically by *AndPoll methods
// No manual task management needed
```

## Schema Migration Analysis

### Framework Changes

- **Old:** Terraform Plugin SDK v2 (`github.com/hashicorp/terraform-plugin-sdk/v2`)
- **New:** Terraform Plugin Framework (`github.com/hashicorp/terraform-plugin-framework`)

### SDK Changes

- **Old:** `gcorelabscloud-go` SDK
- **New:** `gcore-go` SDK

## Migration Diff Table

*Note: Output format is designed for easy copying to Confluence - use list format instead of tables.*

### Load Balancer Listener

__Old resource name:__ `gcore_lblistener`
__New resource name:__ `gcore_cloud_load_balancer_listener`
__Fields added in new:__ `creator_task_id`, `task_id`, `tasks`, `stats`, `insert_headers`
__Fields removed in new:__ `project_name`, `region_name`, `last_updated`
__Notes:__ SDK v2→Framework migration, *AndPoll methods, Enhanced validation, Async task tracking
__Jira ticket:__ N/A

__Old data source name:__ `data.gcore_lblistener`
__New data source name:__ `data.gcore_cloud_load_balancer_listener`
__Fields added in new:__ `show_stats`, `stats`, `insert_headers`
__Fields removed in new:__ `project_name`, `region_name`
__Notes:__ Framework migration, Enhanced data retrieval, Statistics support
__Jira ticket:__ N/A

### Load Balancer Pool

__Old resource name:__ `gcore_lbpool`
__New resource name:__ `gcore_cloud_load_balancer_pool`
__Fields added in new:__ `creator_task_id`, `task_id`, `tasks`
__Fields removed in new:__ `project_name`, `region_name`, `last_updated`
__Notes:__ SDK v2→Framework migration, *AndPoll methods implemented
__Jira ticket:__ N/A

### Load Balancer

__Old resource name:__ `gcore_loadbalancer`
__New resource name:__ `gcore_cloud_load_balancer`
__Fields added in new:__ `creator_task_id`, `task_id`, `tasks`
__Fields removed in new:__ `project_name`, `region_name`, `last_updated`
__Notes:__ SDK v2→Framework migration, *AndPoll methods
__Jira ticket:__ N/A

### Detailed Field Analysis

#### Fields Added in New Provider

1. __`creator_task_id`__ (Computed)

   - Description: Task that created this entity
   - Type: String
   - Purpose: Audit tracking

2. __`task_id`__ (Computed)

   - Description: UUID of active task holding resource lock
   - Type: String
   - Purpose: Prevents concurrent modifications

3. **`tasks`** (Computed)

   - Description: List of task IDs for async operations
   - Type: List of Strings
   - Purpose: Monitor operation progress

4. **`stats`** (Computed, Data Source only)

   - Description: Load balancer statistics
   - Type: Nested Object
   - Fields: `active_connections`, `bytes_in`, `bytes_out`, `request_errors`, `total_connections`
   - Purpose: Performance monitoring

5. __`insert_headers`__ (Computed)

   - Description: Dictionary of additional HTTP header insertions
   - Type: JSON String
   - Purpose: Advanced HTTP header manipulation

6. __`show_stats`__ (Data Source only)

   - Description: Flag to include statistics in response
   - Type: Boolean
   - Purpose: Control data retrieval scope

#### Fields Removed from Old Provider

1. __`project_name`__

   - Old Type: String
   - Old Purpose: Alternative to `project_id` for project identification
   - Removal Reason: Simplified schema, ID-based approach preferred

2. __`region_name`__

   - Old Type: String
   - Old Purpose: Alternative to `region_id` for region identification
   - Removal Reason: Simplified schema, ID-based approach preferred

3. __`last_updated`__

   - Old Type: String (RFC850 timestamp)
   - Old Purpose: Track last modification time
   - Removal Reason: Replaced by system-managed task tracking

#### Fields with Changed Validation

1. __`connection_limit`__

   - Old: Optional, 1-1,000,000 range (no default)
   - New: Optional, -1-1,000,000 range (default: 100000)
   - Change: Added support for -1 (unlimited), explicit default

2. **`protocol`**

   - Old: Custom validation function with enum check
   - New: `stringvalidator.OneOfCaseInsensitive()` with case-insensitive matching
   - Change: More flexible validation, better error messages

#### Implementation Changes

1. **SDK Migration**

   - Old: `gcorelabscloud-go` SDK
   - New: `gcore-go` SDK
   - Impact: Different API client structure

2. **Framework Migration**

   - Old: Terraform Plugin SDK v2
   - New: Terraform Plugin Framework
   - Impact: Type system changes, validation improvements

3. **Task Management**

   - Old: Manual task polling with `tasks.WaitTaskAndReturnResult()`
   - New: Automatic with `*AndPoll()` methods
   - Impact: Simplified async operation handling

4. **Update Logic**

   - Old: Separate `Update()` and `Unset()` operations with manual coordination
   - New: Single `UpdateAndPoll()` operation
   - Impact: Reduced complexity, better error handling

5. **Response Handling**

   - Old: Manual HTTP response parsing with `io.ReadAll()`
   - New: `listener.RawJSON()` with built-in parsing
   - Impact: More reliable data extraction

## Key Implementation Notes

1. **Async Operations:** *AndPoll methods handle polling automatically
2. **Error Handling:** Simplified with built-in retry logic
3. **Response Processing:** Use `RawJSON()` instead of manual HTTP body reading
4. **Task Management:** Automatic with *AndPoll, no manual task polling needed
5. **Validation:** Enhanced schema validation in new Framework
6. __Computed Fields:__ New provider adds system fields like `task_id`, `tasks`

## Testing Requirements

### Environment Setup

#### 1. Environment Variables (.env file)

Create a `.env` file in the project root with required credentials:

```bash
# .env file - Create this at /Users/user/repos/gcore-terraform/.env
GCORE_API_KEY=your_api_key_here
GCORE_CLOUD_PROJECT_ID=379987
GCORE_CLOUD_REGION_ID=76
```

**IMPORTANT:** Load environment variables before running Terraform commands. Use this exact pattern:

```bash
# This pattern ensures variables are exported correctly and then unexported
set -o allexport; source .env; set +o allexport
```

**Why this pattern?**

- `set -o allexport` - Makes all variables automatically exported
- `source .env` - Loads variables from the .env file
- `set +o allexport` - Disables auto-export to avoid polluting environment

**Do NOT use:**

- `export $(grep -v '^#' .env | xargs)` - Can fail with complex values
- Manual `export` commands - Easy to miss variables

#### 2. Provider Override (.terraformrc)

Create `.terraformrc` file in project root for development override:

```hcl
# .terraformrc - Create this at /Users/user/repos/gcore-terraform/.terraformrc
provider_installation {
  dev_overrides {
    "gcore/gcore" = "/Users/user/repos/gcore-terraform"
  }
  direct {}
}
```

**Override the Terraform CLI config:**

```bash
# Use relative path if running from project root
export TF_CLI_CONFIG_FILE=".terraformrc"

# Or use absolute path from anywhere
export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/.terraformrc"
```

**Expected behavior:**

- Terraform will warn: "Provider development overrides are in effect"
- `terraform init` may error but will still use the override
- The built provider binary will be used instead of registry version

#### 3. Build and Test Updated Provider

**Complete build and test workflow:**

```bash
# 1. Navigate to project root
cd /Users/user/repos/gcore-terraform

# 2. Build the provider
go clean -cache
go mod tidy
go build -o terraform-provider-gcore

# 3. Create test directory
mkdir -p test-myresource
cd test-myresource

# 4. Create test configuration (see example below)
cat > main.tf << 'EOF'
# Your test configuration here
EOF

# 5. Run terraform with environment variables loaded
set -o allexport; source .env; set +o allexport && \
export TF_CLI_CONFIG_FILE=".terraformrc" && \
terraform plan
```

**Single-line command for all Terraform operations:**

```bash
# Plan
cd test-myresource && set -o allexport; source ../.env; set +o allexport && export TF_CLI_CONFIG_FILE="../.terraformrc" && terraform plan

# Apply
cd test-myresource && set -o allexport; source ../.env; set +o allexport && export TF_CLI_CONFIG_FILE="../.terraformrc" && terraform apply -auto-approve

# Destroy
cd test-myresource && set -o allexport; source ../.env; set +o allexport && export TF_CLI_CONFIG_FILE="../.terraformrc" && terraform destroy -auto-approve
```

**Test configuration example:**

```hcl
# main.tf
terraform {
  required_providers {
    gcore = {
      source  = "gcore/gcore"
      version = "~> 1.0"
    }
  }
}

provider "gcore" {
  # Reads from environment variables:
  # - GCORE_API_KEY
  # - GCORE_CLOUD_PROJECT_ID
  # - GCORE_CLOUD_REGION_ID
  # Do NOT hardcode credentials in the provider block
}

# Test load balancer
resource "gcore_cloud_load_balancer" "test" {
  name       = "test-lb-${formatdate("YYYY-MM-DD-hhmm", timestamp())}"
  flavor     = "lb1-1-2"
  project_id = 379987
  region_id  = 76
}

# Test listener with new resource name
resource "gcore_cloud_load_balancer_listener" "test" {
  name            = "test-listener-andpoll"
  loadbalancer_id = gcore_cloud_load_balancer.test.id
  protocol        = "HTTP"
  protocol_port   = 80
  project_id      = 379987
  region_id       = 76
}
```

**Testing commands workflow:**

```bash
# Note: Skip 'terraform init' when using dev_overrides

# 1. Validate configuration
terraform validate

# 2. Plan deployment
terraform plan

# 3. Apply changes
terraform apply -auto-approve

# 4. Test updates (modify main.tf, then apply again)
terraform apply -auto-approve

# 5. Clean up
terraform destroy -auto-approve
```

#### 4. Testing Provider from a Different Branch

**Workflow for testing changes from a specific branch:**

```bash
# 1. Clone or navigate to provider repository
cd /Users/user/repos/gcore-terraform

# 2. Checkout the branch with changes
git fetch origin
git checkout feature/your-branch-name

# 3. Build the provider from the branch
go clean -cache
go mod tidy
go build -o terraform-provider-gcore

# 4. The .terraformrc override will now use this branch's binary
# No additional configuration needed!

# 5. Run your tests
cd test-myresource
set -o allexport; source ../.env; set +o allexport && \
export TF_CLI_CONFIG_FILE="../.terraformrc" && \
terraform apply -auto-approve
```

**Testing multiple branches:**

```bash
# Switch branches and rebuild
git checkout main
go build -o terraform-provider-gcore
# Test with main branch...

git checkout feature/new-feature
go build -o terraform-provider-gcore
# Test with feature branch...
```

**Using git worktrees for parallel testing:**

```bash
# Create worktree for feature branch
git worktree add ../gcore-terraform-feature feature/your-branch

# Terminal 1: Test main branch
cd /Users/user/repos/gcore-terraform
go build -o terraform-provider-gcore
cd test-resource-main
set -o allexport; source ../.env; set +o allexport && \
export TF_CLI_CONFIG_FILE="../.terraformrc" && \
terraform apply

# Terminal 2: Test feature branch
cd /Users/user/repos/gcore-terraform-feature
go build -o terraform-provider-gcore
cd test-resource-feature
set -o allexport; source ../.env; set +o allexport && \
export TF_CLI_CONFIG_FILE="../.terraformrc" && \
terraform apply
```

### Validation Steps

After migration:

1. Test resource creation with real cloud resources
2. Verify update operations work correctly
3. Confirm delete operations complete successfully
4. Validate that old resource names are rejected
5. Ensure new resource names work as expected
6. Check that *AndPoll methods handle async operations properly
7. Verify environment variable loading works
8. Confirm provider override is functioning

## Files to Check During Migration

1. **Resource Implementation:** `internal/services/*/resource.go`
2. **Schema Definition:** `internal/services/*/schema.go`
3. __Data Source:__ `internal/services/*/data_source.go`
4. **Provider Registration:** `internal/provider.go`
5. __Old Implementation:__ `old_terraform_provider/gcore/resource_gcore_*.go`

This migration improves reliability, reduces complexity, and provides better async operation handling in the Gcore Terraform provider.

## Handling Computed Optional Fields in Terraform Provider

### Problem: Terraform "Forces Replacement" for Updatable Fields

When Terraform shows "forces replacement" for fields that should allow in-place updates, it's often because the field is marked as `Optional: true` but not `Computed: true`. This causes Terraform to recreate resources instead of updating them.

### Identifying Fields That Need computed_optional Configuration

Look for these symptoms:

1. **Terraform plan shows "forces replacement"** for fields that should be updatable
2. **Field has update logic** in the provider's Update method
3. **API supports in-place updates** for the field (e.g., volume resize, field modifications)
4. **Field schema is Optional but not Computed** in the generated Terraform provider

### Common Fields That Need computed_optional:

- **Volume size** - Can be updated via resize API
- **Network configuration** - Gateway IPs, DNS servers that may be computed by API
- **Resource limits** - CPU, memory that may be rounded by API
- **Any field where API might compute a value** if user doesn't provide one

### Solution: Configure in api-schemas config.yaml

Instead of modifying OpenAPI spec directly, use the `extra_terraform` configuration in `/Users/user/repos/api-schemas/scripts/config.yaml`.

#### Example Configuration:

```yaml
specs:
  - url: https://api.gcore.com/cloud/docs/openapi.yaml
    product: cloud
    title: Cloud API
    extra_terraform:
      components:
        schemas:
          CreateSubnetSerializer:
            properties:
              ip_version: required
              dns_nameservers: computed_optional
              host_routes: computed_optional
              gateway_ip: computed_optional
          VolumeSerializer:
            properties:
              size: computed_optional
          # Add more schemas as needed
```

#### Configuration Values:

- **`required`**: Must be provided by user, cannot be null
- **`optional`**: Can be provided by user, can be null
- **`computed`**: Cannot be provided by user, computed by API/provider
- __`computed_optional`__: Can be provided by user OR computed by API/provider

### Step-by-Step Process:

1. **Identify the issue:**

```bash
terraform plan
# Look for "forces replacement" on fields that should be updatable
```

2. **Find the schema name:**

```bash
# Check which schema is used in the OpenAPI spec
grep -n "YourResourceSerializer:" /Users/user/repos/api-schemas/openapi.yaml
```

3. **Add configuration:**

```yaml
# In /Users/user/repos/api-schemas/scripts/config.yaml
extra_terraform:
  components:
    schemas:
      YourResourceSerializer:
        properties:
          your_field: computed_optional
```

4. **Regenerate and test:**

```bash
# Regenerate OpenAPI spec with new configuration
cd /Users/user/repos/api-schemas
make generate  # or whatever command regenerates the spec

# Regenerate Terraform provider using Stainless
# Rebuild and test
go build -o terraform-provider-gcore
terraform plan  # Should no longer show "forces replacement"
```

### Verification:

After applying the configuration, the generated Terraform schema should have:

```go
"your_field": schema.StringAttribute{
    Optional: true,
    Computed: true,
    PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
},
```

This allows Terraform to:

- Accept user input for the field
- Handle cases where API computes/modifies the value
- Perform in-place updates without forcing replacement

### Common Patterns:

1. **Volume size** - User sets size, API ensures minimum values
2. **Network fields** - User may omit, API provides defaults
3. **Resource configurations** - User sets preferences, API may adjust
4. **Any field with server-side defaults or modifications**

This approach maintains clean OpenAPI specs while providing Terraform-specific configuration through Stainless.