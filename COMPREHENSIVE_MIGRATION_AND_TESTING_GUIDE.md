# Comprehensive Terraform Provider Migration & Testing Guide

**For LLM-Assisted Migration of GCore Terraform Provider Resources**

---

## Table of Contents

1. [Introduction & Core Principles](#introduction--core-principles)
2. [Quick Start Decision Tree](#quick-start-decision-tree)
3. [Understanding API Operation Types](#understanding-api-operation-types)
4. [Pre-Migration Analysis](#pre-migration-analysis)
5. [Environment Setup & Testing Infrastructure](#environment-setup--testing-infrastructure)
6. [Migration Implementation Patterns](#migration-implementation-patterns)
7. [Real Infrastructure Testing Methodology](#real-infrastructure-testing-methodology)
8. [Common Issues & Solutions](#common-issues--solutions)
9. [QA Checklist](#qa-checklist)
10. [Troubleshooting & Debugging](#troubleshooting--debugging)
11. [Case Study: Load Balancer Pool Drift Fix](#case-study-load-balancer-pool-drift-fix)
12. [Reference Materials](#reference-materials)

---

## Introduction & Core Principles

### Purpose

This guide provides comprehensive instructions for migrating GCore Terraform provider resources from the old SDK-based implementation to the new Stainless-generated Framework-based implementation, with **proven real infrastructure testing methodology**.

### Core Principle

**⚠️ CRITICAL: Never make conclusions based on source code analysis alone. All hypotheses MUST be validated through real infrastructure testing.**

### Why This Matters

Our experience shows that:
- ✅ Code review alone missed a critical Pool drift bug
- ✅ Real infrastructure testing found it immediately
- ✅ State can lie; actual API calls and cloud resources are the truth
- ✅ Drift may not be apparent until resources exist in production

### What This Guide Covers

1. **Migration**: From manual HTTP handling → *AndPoll methods
2. **Testing**: Comprehensive real infrastructure validation
3. **Debugging**: Finding and fixing common issues
4. **Validation**: Ensuring actual correctness, not just apparent correctness

---

## Quick Start Decision Tree

```
┌─────────────────────────────────────────┐
│ Starting Resource Migration?            │
└─────────────┬───────────────────────────┘
              │
              ▼
   ┌──────────────────────────┐
   │ Does old provider use    │
   │ tasks.WaitTaskAndReturn? │
   └──────┬──────────┬────────┘
          │          │
      YES │          │ NO
          │          │
          ▼          ▼
  ┌───────────┐  ┌──────────┐
  │ Use       │  │ Use      │
  │ *AndPoll  │  │ Direct   │
  │ Methods   │  │ SDK      │
  │ (Async)   │  │ (Sync)   │
  └─────┬─────┘  └────┬─────┘
        │             │
        ▼             ▼
  ┌─────────────────────────┐
  │ Implement Migration     │
  │ (Section 6)             │
  └───────────┬─────────────┘
              │
              ▼
  ┌─────────────────────────┐
  │ Test with Real Infra    │
  │ (Section 7)             │
  └───────────┬─────────────┘
              │
              ▼
     ┌────────────────┐
     │ Issues Found?  │
     └───┬────────┬───┘
         │        │
      YES│        │NO
         │        │
         ▼        ▼
    ┌────────┐ ┌──────┐
    │Fix     │ │Done! │
    │(Sec 8) │ └──────┘
    └────┬───┘
         │
         └──────────┐
                    │
         ▼          │
    Test Again ─────┘
```

---

## Understanding API Operation Types

### Critical First Step

Before any migration, determine if the API is **synchronous** or **asynchronous**. This determines your entire implementation approach.

### Synchronous Operations (No Tasks)

**Characteristics:**
- API returns resource data immediately in response
- No task IDs returned
- No polling required
- Completes in single HTTP request/response

**Examples:**
- SSH Keys (`gcore_cloud_ssh_key`)
- FaaS Keys
- Storage SFTP Keys

**Code Pattern:**
```go
// Create - returns resource directly
sshKey, err := r.client.Cloud.SSHKeys.New(ctx, params,
    option.WithRequestBody("application/json", dataBytes),
    option.WithMiddleware(logging.Middleware(ctx)),
)
if err != nil {
    resp.Diagnostics.AddError("failed to make http request", err.Error())
    return
}
err = apijson.UnmarshalComputed([]byte(sshKey.RawJSON()), &data)

// Update - direct update, no polling
sshKey, err := r.client.Cloud.SSHKeys.Update(ctx, id, params, ...)

// Delete - direct delete, no polling
err := r.client.Cloud.SSHKeys.Delete(ctx, id, params, ...)
```

### Asynchronous Operations (Returns Tasks)

**Characteristics:**
- API returns task IDs that must be polled
- Resource may not be immediately available
- Requires waiting for task completion
- Multiple HTTP requests (initial + polling)

**Examples:**
- Load Balancers (`gcore_cloud_load_balancer`)
- Load Balancer Listeners (`gcore_cloud_load_balancer_listener`)
- Load Balancer Pools (`gcore_cloud_load_balancer_pool`)
- Volumes, Instances, Networks, Subnets

**Code Pattern:**
```go
// Create - automatically polls until complete
listener, err := r.client.Cloud.LoadBalancers.Listeners.NewAndPoll(ctx, params,
    option.WithRequestBody("application/json", dataBytes),
    option.WithMiddleware(logging.Middleware(ctx)),
)
if err != nil {
    resp.Diagnostics.AddError("failed to make http request", err.Error())
    return
}
err = apijson.UnmarshalComputed([]byte(listener.RawJSON()), &data)

// Update - automatically polls until complete
listener, err := r.client.Cloud.LoadBalancers.Listeners.UpdateAndPoll(ctx, id, params, ...)

// Delete - automatically polls until complete
err := r.client.Cloud.LoadBalancers.Listeners.DeleteAndPoll(ctx, id, params, ...)
```

### How to Determine Operation Type

#### Method 1: Check Old Provider Implementation

**Asynchronous (uses task waiting):**
```go
results, err := listeners.Create(client, opts).Extract()
taskID := results.Tasks[0]
listenerID, err := tasks.WaitTaskAndReturnResult(client, taskID, ...)
```

**Synchronous (direct extract):**
```go
kp, err := keypairs.Create(client, opts).Extract()
d.SetId(kp.ID)  // No task waiting
```

#### Method 2: Check SDK Method Signatures

```bash
# Check the SDK
go doc github.com/G-Core/gcore-go/cloud.SSHKeyService

# Synchronous - returns resource directly
func (r *SSHKeyService) New(...) (res *SSHKeyCreated, err error)

# Asynchronous - has *AndPoll variant
func (r *ListenerService) NewAndPoll(...) (res *Listener, err error)
```

#### Method 3: Test with Real Infrastructure

The most reliable method:

```bash
# Enable detailed logging
export TF_LOG=DEBUG
export TF_LOG_PATH=terraform.log

# Create resource
terraform apply

# Check logs for task IDs
grep -i "task" terraform.log

# If you see task polling or task IDs, it's asynchronous
```

---

## Pre-Migration Analysis

### Step 1: Locate Old Provider Implementation

**File Location:**
```
old_terraform_provider/gcore/resource_gcore_<resource_name>.go
```

**Example:**
- `resource_gcore_lblistener.go` → `gcore_cloud_load_balancer_listener`
- `resource_gcore_lbpool.go` → `gcore_cloud_load_balancer_pool`

### Step 2: Identify Custom Business Logic

**Critical areas to examine:**

#### Update Operations (MOST IMPORTANT)
Look for complex update logic:

```go
// Example: Old provider with Unset logic
if changed {
    results, err := listeners.Update(clientV2, d.Id(), updateOpts, ...)

    if toUnset {
        // Additional unset operation for clearing fields
        results, err := listeners.Unset(clientV2, d.Id(), unsetOpts, ...)
    }
}
```

**Questions to answer:**
- Does update call multiple API endpoints?
- Are there conditional updates based on field changes?
- Is there field clearing/unsetting logic?
- Are there dependencies between field updates?

#### Create Operations
Look for:
- Post-creation configuration
- Additional API calls after resource creation
- Field transformations
- Conditional logic based on input

#### Read Operations
Look for:
- Custom field mapping
- Computed field logic
- Data transformation

#### Delete Operations
Look for:
- Pre-deletion cleanup
- Cascade deletion handling
- Resource unlocking

### Step 3: Document Field Requirements

Create a manifest documenting:

**For each field:**
- `required` - User must provide
- `optional` - User may provide
- `computed` - API/provider computes, user cannot set
- `computed_optional` - User may provide OR API computes
- `force_new` - Changing this field requires resource recreation

**Example Manifest:**
```markdown
Resource: gcore_cloud_load_balancer_pool

Fields:
- name: required
- lb_algorithm: required
- protocol: required
- listener_id: optional, force_new
- load_balancer_id: optional, force_new
- healthmonitor.http_method: computed_optional (API defaults to "GET")
- healthmonitor.max_retries_down: computed_optional (API defaults to 3)
- timeout_client_data: optional
- creator_task_id: computed
- task_id: computed
- tasks: computed
```

### Step 4: Identify Drift-Prone Fields

**Red flags for potential drift issues:**

1. **API Returns Defaults Not in User Config**
   - Field is optional in Terraform
   - API returns a default value
   - Will cause drift if not marked `computed_optional`

2. **API Modifies User-Provided Values**
   - API rounds or normalizes values
   - API adds computed portions
   - Needs `computed_optional` or custom handling

3. **Computed Relationships**
   - Lists of related resources (listeners, loadbalancers, members)
   - Should be marked `computed` with proper tags

**How to spot these:**
```bash
# Create resource with minimal config
terraform apply

# Check state for fields you didn't provide
terraform show

# If you see values you didn't set, they may need computed_optional
```

---

## Environment Setup & Testing Infrastructure

### Step 1: Environment Variables

Create `.env` file in project root:

```bash
# /Users/user/repos/gcore-terraform/.env
GCORE_API_KEY=your_api_key_here
GCORE_CLOUD_PROJECT_ID=379987
GCORE_CLOUD_REGION_ID=76
```

**Load variables before testing:**
```bash
# Correct pattern - exports then unexports
set -o allexport; source .env; set +o allexport

# DO NOT USE these patterns (can fail with complex values):
# export $(grep -v '^#' .env | xargs)  # ❌
# Manual export commands                 # ❌
```

### Step 2: Provider Override

Create `.terraformrc` for development:

```hcl
# /Users/user/repos/gcore-terraform/.terraformrc
provider_installation {
  dev_overrides {
    "gcore/gcore" = "/Users/user/repos/gcore-terraform"
  }
  direct {}
}
```

**Set override:**
```bash
export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/.terraformrc"
```

### Step 3: Enable Provider Logging

**For debugging and validation:**

```bash
# Enable detailed Terraform logging
export TF_LOG=DEBUG
export TF_LOG_PATH=terraform-debug.log

# Provider will log:
# - API requests and responses
# - State changes
# - Plan calculations
# - Error details
```

### Step 4: Build Provider

```bash
cd /Users/user/repos/gcore-terraform

# Clean and rebuild
go clean -cache
go mod tidy
go build -o terraform-provider-gcore
```

### Step 5: Create Test Directory Structure

**Organize tests by purpose:**

```bash
test-<resource>-comprehensive/
├── README.md                  # Test documentation
├── run_test.sh                # Test runner script
├── drift/                     # Drift detection tests
│   ├── TC-DRIFT-001-baseline/
│   ├── TC-DRIFT-002-with-optional/
│   └── TC-DRIFT-003-with-nested/
├── update/                    # Update operation tests
│   ├── TC-UPDATE-001-name-change/
│   ├── TC-UPDATE-002-config-change/
│   └── TC-UPDATE-003-patch-vs-replace/
└── combined/                  # Full stack tests
    └── TC-COMBINED-001-full-stack/
```

**Example test runner (`run_test.sh`):**
```bash
#!/bin/bash
set -e

TEST_DIR=$1

if [ -z "$TEST_DIR" ]; then
    echo "Usage: $0 <test-directory>"
    exit 1
fi

echo "=========================================="
echo "Running test: $TEST_DIR"
echo "=========================================="

cd "$TEST_DIR"

# Source credentials
if [ -f "../../.env" ]; then
    source ../../.env
    export GCORE_API_KEY
    export GCORE_CLIENT
    export GCORE_CLOUD_PROJECT_ID
    export GCORE_CLOUD_REGION_ID
fi

# Copy .terraformrc if exists
if [ -f "../../.terraformrc" ]; then
    cp ../../.terraformrc .
fi

# Skip init with dev override
if [ -f ".terraformrc" ]; then
    echo "Step 1: Skipping terraform init (using provider dev override)"
else
    echo "Step 1: Terraform init"
    terraform init -upgrade
fi

# First apply
echo ""
echo "Step 2: First terraform apply"
terraform apply -auto-approve

# Drift check
echo ""
echo "Step 3: Check for configuration drift"
echo ""

if terraform plan -detailed-exitcode; then
    echo ""
    echo "✅ PASS: No drift detected"
    exit 0
else
    EXIT_CODE=$?
    if [ $EXIT_CODE -eq 2 ]; then
        echo ""
        echo "❌ FAIL: Drift detected"
        terraform plan
        exit 1
    else
        echo ""
        echo "❌ ERROR: terraform plan failed"
        exit 1
    fi
fi
```

---

## Migration Implementation Patterns

### Pattern 1: Asynchronous Operations (*AndPoll)

**Old Pattern (Manual HTTP + Task Waiting):**
```go
res := new(http.Response)
_, err = r.client.Cloud.LoadBalancers.Listeners.New(
    ctx,
    params,
    option.WithRequestBody("application/json", dataBytes),
    option.WithResponseBodyInto(&res),  // ❌ Manual response handling
    option.WithMiddleware(logging.Middleware(ctx)),
)
if err != nil {
    resp.Diagnostics.AddError("failed to make http request", err.Error())
    return
}
bytes, _ := io.ReadAll(res.Body)  // ❌ Manual body reading
err = apijson.Unmarshal(bytes, &data)  // ❌ Wrong unmarshaler
```

**New Pattern (*AndPoll):**
```go
listener, err := r.client.Cloud.LoadBalancers.Listeners.NewAndPoll(
    ctx,
    params,
    option.WithRequestBody("application/json", dataBytes),
    // ❌ REMOVE: option.WithResponseBodyInto(&res)
    option.WithMiddleware(logging.Middleware(ctx)),
)
if err != nil {
    resp.Diagnostics.AddError("failed to make http request", err.Error())
    return
}
// ✅ Use RawJSON() from SDK response
err = apijson.UnmarshalComputed([]byte(listener.RawJSON()), &data)
```

**Key Changes:**
1. ✅ `New()` → `NewAndPoll()`
2. ✅ Remove `option.WithResponseBodyInto(&res)`
3. ✅ Remove `io.ReadAll(res.Body)`
4. ✅ Use `listener.RawJSON()` instead
5. ✅ Use `UnmarshalComputed` (not `Unmarshal`)

**Apply to all CRUD operations:**
- Create: `NewAndPoll()`
- Update: `UpdateAndPoll()`
- Delete: `DeleteAndPoll()`

### Pattern 2: Synchronous Operations (Direct SDK)

**Old Pattern:**
```go
res := new(http.Response)
_, err = r.client.Cloud.SSHKeys.New(
    ctx,
    params,
    option.WithRequestBody("application/json", dataBytes),
    option.WithResponseBodyInto(&res),  // ❌ Manual response handling
    option.WithMiddleware(logging.Middleware(ctx)),
)
if err != nil {
    resp.Diagnostics.AddError("failed to make http request", err.Error())
    return
}
bytes, _ := io.ReadAll(res.Body)  // ❌ Manual body reading
err = apijson.Unmarshal(bytes, &data)  // ❌ Wrong unmarshaler
```

**New Pattern:**
```go
sshKey, err := r.client.Cloud.SSHKeys.New(
    ctx,
    params,
    option.WithRequestBody("application/json", dataBytes),
    // ❌ REMOVE: option.WithResponseBodyInto(&res)
    option.WithMiddleware(logging.Middleware(ctx)),
)
if err != nil {
    resp.Diagnostics.AddError("failed to make http request", err.Error())
    return
}
// ✅ Use RawJSON() from SDK response
err = apijson.UnmarshalComputed([]byte(sshKey.RawJSON()), &data)
```

**Key Differences from *AndPoll:**
- ✅ Keep method names: `New()`, `Update()`, `Delete()` (no AndPoll suffix)
- ✅ SDK returns resource directly, no task polling
- ✅ Still remove manual HTTP response handling
- ✅ Still use `RawJSON()` and `UnmarshalComputed`

### Critical: UnmarshalComputed vs Unmarshal

**When to use `UnmarshalComputed`:**

✅ **ALWAYS in these locations:**
- `Create()` method - After NewAndPoll or New
- `Read()` method - Reading resource state
- `ImportState()` method - Importing existing resources
- `Update()` method - After UpdateAndPoll or Update (if reading response)

**Why UnmarshalComputed:**
```go
// Model with computed_optional field
type Model struct {
    HTTPMethod types.String `json:"http_method,computed_optional"`
}

// apijson.Unmarshal - WRONG
// - Ignores ,computed_optional tag
// - Overwrites state with null if field not in response
// - Causes drift when API returns computed values

// apijson.UnmarshalComputed - CORRECT
// - Respects ,computed_optional tag
// - Preserves existing state if field not in response
// - Handles API-computed defaults correctly
```

**Example from real bug:**
```go
// BEFORE (caused drift)
err = apijson.Unmarshal(bytes, &data)

// AFTER (fixed drift)
err = apijson.UnmarshalComputed(bytes, &data)
```

This single line change fixed Pool resource drift detected through real infrastructure testing.

### Pattern 3: Update Operations

**Simple updates (most cases):**
```go
func (r *ResourceType) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data *Model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    var state *Model
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

    // Serialize for PATCH
    dataBytes, err := data.MarshalJSONForUpdate(*state)
    if err != nil {
        resp.Diagnostics.AddError("failed to serialize http request", err.Error())
        return
    }

    // Check if any changes needed
    dataStr := strings.TrimSpace(string(dataBytes))
    if dataStr == "{}" || dataStr == "null" || len(dataBytes) == 0 {
        // No changes - just refresh from API
        // Use Read logic here
        return
    }

    // Send PATCH request
    resource, err := r.client.Cloud.Resource.UpdateAndPoll(
        ctx,
        data.ID.ValueString(),
        params,
        option.WithRequestBody("application/json", dataBytes),
        option.WithMiddleware(logging.Middleware(ctx)),
    )
    if err != nil {
        resp.Diagnostics.AddError("failed to make http request", err.Error())
        return
    }

    err = apijson.UnmarshalComputed([]byte(resource.RawJSON()), &data)
    if err != nil {
        resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
        return
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
```

**Complex updates (special handling needed):**

If old provider had `Unset()` operations or conditional logic:
1. Document the business logic
2. Implement equivalent in new provider
3. Test thoroughly with real infrastructure
4. Verify PATCH operations are sent correctly

---

## Real Infrastructure Testing Methodology

### Core Principle

**Testing is not just running `terraform apply` and checking for errors.**

Real infrastructure testing means:
1. ✅ Verifying resources actually exist in cloud platform
2. ✅ Inspecting actual API calls (not just Terraform logs)
3. ✅ Validating operation types (POST/PATCH/PUT/DELETE)
4. ✅ Confirming behavior matches expectations
5. ✅ Testing failure scenarios and edge cases

### Phase 1: Drift Detection Testing

**Purpose:** Verify no false positives on second `terraform plan`

**Test Workflow:**

```bash
# 1. Create resource
cd test-drift/TC-DRIFT-001-baseline
terraform apply -auto-approve

# 2. Immediately run plan again
terraform plan -detailed-exitcode

# Expected result: Exit code 0 (no changes)
# If exit code 2: DRIFT DETECTED - BUG FOUND
```

**What this catches:**
- Fields marked `optional` that should be `computed_optional`
- Missing `UnmarshalComputed` in Read method
- Computed fields not properly handled
- API defaults causing state mismatch

**Example test configuration:**
```hcl
# Minimal config to test drift
resource "gcore_cloud_load_balancer_pool" "test" {
  name             = "test-pool-drift"
  lb_algorithm     = "ROUND_ROBIN"
  protocol         = "HTTP"
  listener_id      = gcore_cloud_load_balancer_listener.test.id
  load_balancer_id = gcore_cloud_load_balancer.test.id
  project_id       = var.project_id
  region_id        = var.region_id

  # Include optional nested structures that have API defaults
  healthmonitor = {
    delay          = 10
    max_retries    = 3
    timeout        = 5
    type           = "HTTP"
    url_path       = "/health"
    expected_codes = "200"
    # http_method not specified - API will default to "GET"
    # max_retries_down not specified - API will default to 3
  }

  timeout_client_data    = 50000
  timeout_member_connect = 5000
  timeout_member_data    = 50000
}
```

**Success criteria:**
```
terraform plan

No changes. Your infrastructure matches the configuration.

✅ PASS
```

**Failure example (actual bug found):**
```
terraform plan

Terraform will perform the following actions:

  # gcore_cloud_load_balancer_pool.test will be updated in-place
  ~ resource "gcore_cloud_load_balancer_pool" "test" {
      ~ healthmonitor = {
          - http_method      = "GET" -> null
          - max_retries_down = 3 -> null
            # (6 unchanged attributes hidden)
        }
      ~ listeners = [...] -> (known after apply)
      ~ members = [] -> (known after apply)
        ...
    }

❌ FAIL - Bug found: API defaults cause drift
```

### Phase 2: Update Operation Testing

**Purpose:** Verify updates use PATCH (not replacement) and work correctly

**Test Workflow:**

```bash
# 1. Create resource
cd test-update/TC-UPDATE-001-name-change
terraform apply -auto-approve -var="name=original-name"

# 2. Capture resource ID
RESOURCE_ID=$(terraform output -raw resource_id)

# 3. Update resource
terraform apply -auto-approve -var="name=updated-name"

# 4. Verify ID unchanged
NEW_ID=$(terraform output -raw resource_id)

if [ "$RESOURCE_ID" = "$NEW_ID" ]; then
    echo "✅ PASS: Resource updated in-place"
else
    echo "❌ FAIL: Resource was replaced"
fi

# 5. Check no drift after update
terraform plan -detailed-exitcode
```

**What this catches:**
- Unnecessary resource recreation (ForceNew incorrectly set)
- PATCH not being used (DELETE+POST instead)
- Post-update drift
- State inconsistencies after update

**Verification with logs:**
```bash
# Enable logging
export TF_LOG=DEBUG
export TF_LOG_PATH=update-test.log

# Run update
terraform apply

# Check API operation
grep -i "PATCH\|PUT\|POST.*pool" update-test.log

# Should see PATCH to /pools/{id}, not DELETE then POST
✅ PATCH /cloud/v1/lbpools/379987/76/{id}
❌ DELETE /cloud/v1/lbpools/379987/76/{id}
❌ POST /cloud/v1/lbpools/379987/76
```

**Example test configuration:**
```hcl
variable "pool_name" {
  default = "test-pool-update"
}

resource "gcore_cloud_load_balancer_pool" "test" {
  name             = var.pool_name  # Updatable field
  lb_algorithm     = "ROUND_ROBIN"
  protocol         = "HTTP"
  listener_id      = gcore_cloud_load_balancer_listener.test.id
  load_balancer_id = gcore_cloud_load_balancer.test.id
  project_id       = var.project_id
  region_id        = var.region_id
}

output "pool_id" {
  value = gcore_cloud_load_balancer_pool.test.id
}

output "pool_name" {
  value = gcore_cloud_load_balancer_pool.test.name
}
```

**Expected plan output:**
```
terraform plan -var="pool_name=updated-name"

Terraform will perform the following actions:

  # gcore_cloud_load_balancer_pool.test will be updated in-place
  ~ resource "gcore_cloud_load_balancer_pool" "test" {
      ~ name = "test-pool-update" -> "updated-name"
      # (other fields unchanged)
    }

Plan: 0 to add, 1 to change, 0 to destroy.

✅ PASS: Shows "update in-place", not replacement
```

### Phase 3: CRUD Testing with Real Verification

**Purpose:** Verify resources actually exist in cloud, not just in Terraform state

**Test Workflow:**

```bash
# 1. Create
terraform apply -auto-approve

# 2. Verify in cloud platform
POOL_ID=$(terraform output -raw pool_id)

# Option A: Use API directly
curl -H "Authorization: APIKey ${GCORE_API_KEY}" \
  "https://api.gcore.com/cloud/v1/lbpools/${PROJECT_ID}/${REGION_ID}/${POOL_ID}"

# Option B: Use gcloud CLI if available
gcloud lb pools describe $POOL_ID

# Expected: Resource exists with matching attributes

# 3. Update
terraform apply -auto-approve -var="name=updated"

# 4. Verify update in cloud
curl -H "Authorization: APIKey ${GCORE_API_KEY}" \
  "https://api.gcore.com/cloud/v1/lbpools/${PROJECT_ID}/${REGION_ID}/${POOL_ID}" \
  | jq '.name'

# Expected: "updated"

# 5. Delete
terraform destroy -auto-approve

# 6. Verify deletion
curl -H "Authorization: APIKey ${GCORE_API_KEY}" \
  "https://api.gcore.com/cloud/v1/lbpools/${PROJECT_ID}/${REGION_ID}/${POOL_ID}"

# Expected: 404 Not Found
```

**What this catches:**
- Terraform state says resource exists but it doesn't
- Updates don't actually apply to cloud resource
- Deletes don't actually remove cloud resource
- State drift from external changes

### Phase 4: ForceNew Field Testing

**Purpose:** Verify fields requiring replacement actually trigger replacement

**Test Workflow:**

```bash
# 1. Create with specific listener
terraform apply -auto-approve -var="listener_id=listener-1"
OLD_ID=$(terraform output -raw pool_id)

# 2. Change ForceNew field
terraform apply -auto-approve -var="listener_id=listener-2"
NEW_ID=$(terraform output -raw pool_id)

# 3. Verify resource was replaced
if [ "$OLD_ID" != "$NEW_ID" ]; then
    echo "✅ PASS: Resource correctly replaced"
else
    echo "❌ FAIL: Resource should have been replaced"
fi
```

**Expected plan output:**
```
terraform plan -var="listener_id=listener-2"

Terraform will perform the following actions:

  # gcore_cloud_load_balancer_pool.test must be replaced
-/+ resource "gcore_cloud_load_balancer_pool" "test" {
      ~ listener_id = "listener-1" -> "listener-2" # forces replacement
      ~ id          = "old-id" -> (known after apply)
        ...
    }

Plan: 1 to add, 0 to change, 1 to destroy.

✅ PASS: Shows "forces replacement"
```

### Phase 5: Import Testing

**Purpose:** Verify existing resources can be imported

**Test Workflow:**

```bash
# 1. Create resource manually or with Terraform
terraform apply -auto-approve
POOL_ID=$(terraform output -raw pool_id)

# 2. Remove from state (but not cloud)
terraform state rm gcore_cloud_load_balancer_pool.test

# 3. Import
terraform import gcore_cloud_load_balancer_pool.test \
  "${PROJECT_ID}/${REGION_ID}/${POOL_ID}"

# 4. Verify no changes needed
terraform plan -detailed-exitcode

# Expected: Exit code 0 (no changes)
```

**What this catches:**
- Import path format issues
- State incompatibility
- Missing computed fields after import
- Drift after import

### Phase 6: Edge Case Testing

**Test scenarios:**

1. **Empty optional fields:**
   ```hcl
   healthmonitor = {}  # All optional fields omitted
   ```

2. **Null vs empty string:**
   ```hcl
   description = ""  # Empty string
   # vs
   # description not specified (null)
   ```

3. **API computed defaults:**
   ```hcl
   # Don't specify fields with API defaults
   # Verify they don't cause drift
   ```

4. **Maximum values:**
   ```hcl
   connection_limit = 1000000  # Maximum allowed
   ```

5. **Concurrent operations:**
   ```bash
   # Create multiple resources in parallel
   terraform apply -parallelism=10
   ```

---

## Common Issues & Solutions

### Issue 1: Configuration Drift (False Positives)

**Symptom:**
```
terraform plan

  # resource shows changes but nothing actually changed
  ~ resource "gcore_cloud_load_balancer_pool" "test" {
      ~ healthmonitor = {
          - http_method = "GET" -> null
        }
    }
```

**Detection:**
- Second `terraform plan` after apply shows changes
- Fields you didn't configure appear then disappear
- Computed fields show as changing

**Root Cause:**
1. API returns default values not in user config
2. Field marked `optional` instead of `computed_optional`
3. `apijson.Unmarshal` used instead of `UnmarshalComputed` in Read method

**Fix:**

**Step 1: Fix Read method**
```go
// File: internal/services/cloud_<resource>/resource.go

// BEFORE (line ~249 in Read method)
err = apijson.Unmarshal(bytes, &data)

// AFTER
err = apijson.UnmarshalComputed(bytes, &data)
```

**Also check ImportState method:**
```go
// BEFORE (line ~329 in ImportState method)
err = apijson.Unmarshal(bytes, &data)

// AFTER
err = apijson.UnmarshalComputed(bytes, &data)
```

**Step 2: Fix model tags**
```go
// File: internal/services/cloud_<resource>/model.go

// BEFORE
HTTPMethod types.String `tfsdk:"http_method" json:"http_method,optional"`

// AFTER
HTTPMethod types.String `tfsdk:"http_method" json:"http_method,computed_optional"`
```

**Step 3: Fix schema**
```go
// File: internal/services/cloud_<resource>/schema.go

// BEFORE
"http_method": schema.StringAttribute{
    Description: "...",
    Optional:    true,
    ...
},

// AFTER
"http_method": schema.StringAttribute{
    Description: "...",
    Computed:    true,
    Optional:    true,
    ...
},
```

**Verification:**
```bash
# Rebuild
go build -o terraform-provider-gcore

# Test
cd test-drift
terraform destroy -auto-approve
terraform apply -auto-approve
terraform plan -detailed-exitcode

# Should exit with code 0 (no changes)
```

**Prevention:**
- Always use `UnmarshalComputed` in Read and ImportState
- Test drift detection for every resource
- Check API documentation for default values

### Issue 2: Forces Replacement (Incorrect ForceNew)

**Symptom:**
```
terraform plan

  # resource will be replaced
-/+ resource "gcore_cloud_load_balancer_pool" "test" {
      ~ name = "old-name" -> "new-name" # forces replacement
    }
```

**Detection:**
- Update plan shows replacement instead of in-place update
- Fields that should be updatable trigger recreation
- Resource ID changes after update

**Root Cause:**
1. Field has `RequiresReplace()` plan modifier but shouldn't
2. Field marked `optional` when should be `computed_optional`
3. API doesn't support PATCH for this field (verify in API docs)

**Fix:**

**If field SHOULD be updatable:**

See [MAKE_FIELD_OPTIONAL_COMPUTED.md](./MAKE_FIELD_OPTIONAL_COMPUTED.md) for complete process.

Quick summary:
1. Edit `api-schemas/scripts/config.yaml`
2. Add field as `computed_optional`
3. CI regenerates OpenAPI spec
4. Stainless regenerates provider code
5. Rebase your branch

**If field CANNOT be updated (API limitation):**
- Keep `RequiresReplace()` plan modifier
- Document in field description
- Add to changelog as known behavior

**Verification:**
```bash
# Test update
terraform apply -var="name=original"
OLD_ID=$(terraform output -raw id)

terraform apply -var="name=updated"
NEW_ID=$(terraform output -raw id)

# IDs should match
[ "$OLD_ID" = "$NEW_ID" ] && echo "✅ PASS" || echo "❌ FAIL"
```

### Issue 3: Empty Update API Calls

**Symptom:**
```
# Provider logs show
PATCH /cloud/v1/lbpools/{id}
Request body: {}

# But terraform plan showed changes
```

**Detection:**
- Check `TF_LOG=DEBUG` logs
- See PATCH requests with empty body `{}`
- Updates appear to work but nothing changes

**Root Cause:**
Field marked as `Computed` in both `CreateParams` and `Schema`, causing serialization to skip it.

**Fix:**
```go
// Check: internal/services/cloud_<resource>/schema.go
// Remove Computed from parameter structs, keep only in schema

// BEFORE
type CreateParams struct {
    Name types.String `json:"name,computed"`  // ❌ Remove computed here
}

// AFTER
type CreateParams struct {
    Name types.String `json:"name"`  // ✅ Not computed in params
}

// Schema should still have Computed:
"name": schema.StringAttribute{
    Computed: true,  // ✅ Keep computed in schema
    Optional: true,
}
```

**Verification:**
```bash
export TF_LOG=DEBUG
export TF_LOG_PATH=update.log

terraform apply -var="name=updated"

# Check request body
grep -A 10 "PATCH.*lbpool" update.log | grep -A 5 "Request body"

# Should show: {"name": "updated"}, not {}
```

### Issue 4: Missing *AndPoll Implementation

**Symptom:**
```
Error: context deadline exceeded
```

**Detection:**
- Long-running operations timeout
- Resources appear in incomplete state
- Task polling errors

**Root Cause:**
Using `New()` / `Update()` / `Delete()` instead of `NewAndPoll()` / `UpdateAndPoll()` / `DeleteAndPoll()` for asynchronous operations.

**Fix:**
```go
// BEFORE
listener, err := r.client.Cloud.LoadBalancers.Listeners.New(ctx, params, ...)

// AFTER
listener, err := r.client.Cloud.LoadBalancers.Listeners.NewAndPoll(ctx, params, ...)
```

Apply to all CRUD operations for asynchronous resources.

**Verification:**
- Operations complete successfully
- No timeout errors
- Resources reach ACTIVE state

---

## QA Checklist

### Pre-Migration (Old Provider Analysis)

- [ ] Located old provider implementation file
- [ ] Identified async vs sync operations (task waiting present?)
- [ ] Documented custom business logic in:
  - [ ] Create method
  - [ ] Read method
  - [ ] Update method (CRITICAL - check thoroughly)
  - [ ] Delete method
- [ ] Created field manifest (required/optional/computed/force_new)
- [ ] Identified potential drift-prone fields
- [ ] Tested old provider behavior with real infrastructure
- [ ] Captured example terraform plan/apply output
- [ ] Documented in Jira ticket

### Migration Implementation

- [ ] Determined correct pattern (sync vs async)
- [ ] Implemented *AndPoll methods (if async) or direct SDK (if sync)
- [ ] Replaced manual HTTP response handling:
  - [ ] Removed `option.WithResponseBodyInto(&res)`
  - [ ] Removed `io.ReadAll(res.Body)`
  - [ ] Using `resource.RawJSON()`
- [ ] Using `UnmarshalComputed` in:
  - [ ] Create method
  - [ ] Read method
  - [ ] ImportState method
  - [ ] Update method (if reading response)
- [ ] Ported custom business logic from old provider
- [ ] Configured computed_optional fields (if needed)
- [ ] Built provider successfully

### Real Infrastructure Testing

#### Drift Detection
- [ ] Created resource with terraform
- [ ] Second plan shows zero changes (no drift)
- [ ] Tested with minimal configuration
- [ ] Tested with full configuration including optionals
- [ ] Tested nested structures with API defaults

#### Update Operations
- [ ] Updates use PATCH (verified in logs)
- [ ] Resource ID remains unchanged after update
- [ ] No drift after update
- [ ] State correctly reflects changes
- [ ] Verified actual cloud resource updated

#### CRUD Operations
- [ ] Create: Resource exists in cloud platform
- [ ] Read: terraform refresh syncs state correctly
- [ ] Update: Changes applied to actual resource
- [ ] Delete: Resource removed from cloud platform
- [ ] Import: Existing resources can be imported

#### ForceNew Fields
- [ ] Fields with RequiresReplace trigger replacement
- [ ] Fields without RequiresReplace allow in-place update
- [ ] Resource ID changes for ForceNew fields
- [ ] Resource ID stable for updateable fields

#### Edge Cases
- [ ] Empty optional fields handled
- [ ] Null vs empty string differentiated
- [ ] Maximum field values accepted
- [ ] Concurrent operations succeed
- [ ] API computed defaults don't cause drift

### Post-Migration Validation

- [ ] All tests passed against real infrastructure
- [ ] No false drift detected
- [ ] Updates use PATCH appropriately
- [ ] ForceNew behavior correct
- [ ] Import functionality works
- [ ] State file structure matches old provider (if upgrading)
- [ ] Documented any breaking changes
- [ ] Updated Jira ticket with results
- [ ] Created pull request

### Pass/Fail Criteria

**✅ PASS if:**
- All drift tests show "No changes" on second plan
- Updates use PATCH (not DELETE+POST) where appropriate
- Resource IDs stable for non-ForceNew updates
- Real cloud resources match Terraform state
- Import works without drift
- No regressions from old provider

**❌ FAIL if:**
- Drift detected on second plan after clean apply
- Updates unnecessarily recreate resources
- State diverges from actual infrastructure
- Import causes drift
- Custom logic not migrated
- API calls show wrong operation types (POST instead of PATCH)

---

## Troubleshooting & Debugging

### Enable Detailed Logging

```bash
# Terraform core logging
export TF_LOG=DEBUG
export TF_LOG_PATH=terraform-debug.log

# Provider-specific logging (if implemented)
export TF_LOG_PROVIDER=TRACE

# Run operation
terraform apply

# Analyze logs
grep -i "error\|panic\|fail" terraform-debug.log
grep -i "api.*request\|api.*response" terraform-debug.log
grep -i "drift\|computed\|unmarshal" terraform-debug.log
```

### Analyzing API Calls

**Check operation types:**
```bash
# Look for HTTP methods
grep -E "(GET|POST|PUT|PATCH|DELETE).*http" terraform-debug.log

# Find what was sent in requests
grep -A 20 "Request.*lbpool" terraform-debug.log | head -40

# Find what was received in responses
grep -A 20 "Response.*lbpool" terraform-debug.log | head -40
```

**Verify PATCH vs replacement:**
```bash
# Update should show PATCH
grep "PATCH.*lbpools" terraform-debug.log

# Recreation shows DELETE then POST
grep -E "(DELETE|POST).*lbpools" terraform-debug.log
```

### Debugging Drift Issues

```bash
# 1. Create resource
terraform apply

# 2. Capture initial state
terraform show -json > state-initial.json

# 3. Run plan (should be no changes)
terraform plan -detailed-exitcode

# 4. If drift detected, capture plan
terraform plan -out=drift.tfplan
terraform show -json drift.tfplan > drift.json

# 5. Compare to find what's drifting
diff <(jq -S '.values.root_module.resources[] | select(.address=="gcore_cloud_load_balancer_pool.test")' state-initial.json) \
     <(jq -S '.planned_values.root_module.resources[] | select(.address=="gcore_cloud_load_balancer_pool.test")' drift.json)
```

### Debugging Update Issues

```bash
# Enable logging
export TF_LOG=DEBUG
export TF_LOG_PATH=update-debug.log

# Run update
terraform apply -var="name=updated"

# Check if PATCH was sent
grep "PATCH" update-debug.log

# Check request body
grep -A 10 "Request body" update-debug.log

# If body is empty {}, you have the duplicate Computed field issue
```

### Common Log Patterns

**Good patterns to see:**
```
✅ Using UnmarshalComputed for response data
✅ PATCH request to /lbpools/{id}
✅ Request body contains changed fields: {"name": "updated"}
✅ Response contains all expected fields
```

**Bad patterns (indicate bugs):**
```
❌ Using Unmarshal (should be UnmarshalComputed)
❌ POST and DELETE for update (should be PATCH)
❌ Request body empty: {}
❌ Response missing fields that were in request
❌ Drift detected on computed_optional fields
```

### When to Escalate

**Self-fix these:**
- UnmarshalComputed missing in Read/ImportState
- Schema missing Computed flag
- Model missing computed_optional tag
- Test configuration errors

**Escalate these:**
- API behavior doesn't match documentation
- SDK method signatures don't match expectations
- Stainless codegen produces incorrect output
- Fundamental API limitations (can't update field)

---

## Case Study: Load Balancer Pool Drift Fix

### Background

During comprehensive testing of Load Balancer resources, we encountered a critical bug in the Pool resource that was **not detected by code review** but was **immediately found through real infrastructure testing**.

### Problem Discovery

**Test Case:** TC-DRIFT-003 - Pool with healthmonitor

**Configuration:**
```hcl
resource "gcore_cloud_load_balancer_pool" "test" {
  name             = "test-pool-drift-03"
  lb_algorithm     = "ROUND_ROBIN"
  protocol         = "HTTP"
  listener_id      = gcore_cloud_load_balancer_listener.test.id
  load_balancer_id = gcore_cloud_load_balancer.test.id
  project_id       = var.project_id
  region_id        = var.region_id

  healthmonitor = {
    delay          = 10
    max_retries    = 3
    timeout        = 5
    type           = "HTTP"
    url_path       = "/health"
    expected_codes = "200"
    # Note: http_method and max_retries_down NOT specified
  }

  timeout_client_data    = 50000
  timeout_member_connect = 5000
  timeout_member_data    = 50000
}
```

**Test Steps:**
```bash
# 1. Apply configuration
terraform apply -auto-approve
# ✅ Success - pool created

# 2. Check for drift
terraform plan
```

**Result: ❌ FAIL - Drift detected**
```
Terraform will perform the following actions:

  # gcore_cloud_load_balancer_pool.test will be updated in-place
  ~ resource "gcore_cloud_load_balancer_pool" "test" {
      ~ creator_task_id = "..." -> (known after apply)
      ~ healthmonitor = {
          - http_method      = "GET" -> null
          - max_retries_down = 3 -> null
            # (6 unchanged attributes hidden)
        }
      ~ listeners = [...] -> (known after apply)
      ~ loadbalancers = [...] -> (known after apply)
      ~ members = [] -> (known after apply)
      ~ operating_status = "ONLINE" -> (known after apply)
      ~ provisioning_status = "ACTIVE" -> (known after apply)
      + task_id = (known after apply)
      + tasks = (known after apply)
    }

Plan: 0 to add, 1 to change, 0 to destroy.
```

### Root Cause Analysis

**Investigation revealed three related issues:**

1. **Read method using wrong unmarshaler:**
   ```go
   // Line 249: internal/services/cloud_load_balancer_pool/resource.go
   err = apijson.Unmarshal(bytes, &data)  // ❌ WRONG
   ```

2. **ImportState method using wrong unmarshaler:**
   ```go
   // Line 329: internal/services/cloud_load_balancer_pool/resource.go
   err = apijson.Unmarshal(bytes, &data)  // ❌ WRONG
   ```

3. **Model fields not marked computed_optional:**
   ```go
   // Lines 52-53: internal/services/cloud_load_balancer_pool/model.go
   HTTPMethod     types.String `json:"http_method,optional"`  // ❌ WRONG
   MaxRetriesDown types.Int64  `json:"max_retries_down,optional"`  // ❌ WRONG
   ```

4. **Schema fields not marked Computed:**
   ```go
   // Lines 156, 173: internal/services/cloud_load_balancer_pool/schema.go
   "http_method": schema.StringAttribute{
       Optional: true,  // ❌ Missing Computed: true
   }
   "max_retries_down": schema.Int64Attribute{
       Optional: true,  // ❌ Missing Computed: true
   }
   ```

**Why this caused drift:**
- User doesn't specify `http_method` or `max_retries_down`
- API computes defaults: `GET` and `3`
- On Read, `apijson.Unmarshal` ignores the computed values
- Terraform sees API values as drift from null

### Fix Implementation

**File 1: resource.go (2 changes)**
```go
// Line 249 - Read method
- err = apijson.Unmarshal(bytes, &data)
+ err = apijson.UnmarshalComputed(bytes, &data)

// Line 329 - ImportState method
- err = apijson.Unmarshal(bytes, &data)
+ err = apijson.UnmarshalComputed(bytes, &data)
```

**File 2: model.go (2 changes)**
```go
// Lines 52-53
- HTTPMethod     types.String `json:"http_method,optional"`
- MaxRetriesDown types.Int64  `json:"max_retries_down,optional"`
+ HTTPMethod     types.String `json:"http_method,computed_optional"`
+ MaxRetriesDown types.Int64  `json:"max_retries_down,computed_optional"`
```

**File 3: schema.go (2 changes)**
```go
// Line 158
"http_method": schema.StringAttribute{
    Description: "...",
+   Computed:    true,
    Optional:    true,
    Validators:  [...],
},

// Line 176
"max_retries_down": schema.Int64Attribute{
    Description: "...",
+   Computed:    true,
    Optional:    true,
    Validators:  [...],
},
```

### Verification

**Rebuild and retest:**
```bash
# Rebuild provider
go build -o terraform-provider-gcore

# Clean up old test resources
terraform destroy -auto-approve

# Re-run drift test
terraform apply -auto-approve
terraform plan -detailed-exitcode
```

**Result: ✅ PASS**
```
No changes. Your infrastructure matches the configuration.

✅ PASS: No drift detected
```

### Additional Testing

**Also tested update operations:**
```bash
# Test pool name update
terraform apply -var="pool_name=test-pool-renamed"
```

**Result: ✅ PASS**
```
  # gcore_cloud_load_balancer_pool.test will be updated in-place
  ~ resource "gcore_cloud_load_balancer_pool" "test" {
      ~ name = "test-pool-drift-03" -> "test-pool-renamed"
    }

Plan: 0 to add, 1 to change, 0 to destroy.

# Verified:
✅ PATCH operation sent (not DELETE+POST)
✅ Resource ID unchanged
✅ No drift after update
```

### Key Lessons Learned

1. **Real infrastructure testing is essential**
   - Code review didn't catch this bug
   - Only actual API interaction revealed the issue

2. **Systematic drift testing catches bugs early**
   - Simple test: create then plan
   - Immediately shows state handling issues

3. **UnmarshalComputed is not optional**
   - Required for proper computed field handling
   - Must be used in Read and ImportState

4. **Computed_optional pattern is critical**
   - For fields with API-computed defaults
   - Requires coordination between model, schema, and unmarshaler

5. **Pattern applies across resources**
   - Same fix pattern for Load Balancer and Listener
   - Systematic issue suggests checking all resources

### Statistics

**Testing Duration:** ~2 hours
**Tests Executed:** 4 (3 drift, 1 update)
**Bugs Found:** 1 (critical)
**Bugs Fixed:** 1 (100%)
**Code Changes:** 6 insertions, 4 deletions across 3 files
**Verification:** ✅ Passed all tests against real infrastructure

---

## Reference Materials

### Related Documentation

- **[MAKE_FIELD_OPTIONAL_COMPUTED.md](./MAKE_FIELD_OPTIONAL_COMPUTED.md)** - Detailed guide for configuring computed_optional fields in api-schemas
- **[terraform-resource-qa-checklist.md](./terraform-resource-qa-checklist.md)** - Original QA checklist (superseded by this guide)
- **[LLM_MIGRATION_INSTRUCTIONS.md](./LLM_MIGRATION_INSTRUCTIONS.md)** - Original migration instructions (superseded by this guide)

### Test Examples

**Location:** `/Users/user/repos/gcore-terraform/test-lb-comprehensive/`

**Structure:**
```
test-lb-comprehensive/
├── README.md
├── TESTING_PLAN.md
├── TEST_REPORT.md
├── FIXES_REQUIRED.md
├── run_test.sh
├── drift/
│   ├── TC-DRIFT-001-lb-no-changes/
│   ├── TC-DRIFT-002-listener-no-changes/
│   └── TC-DRIFT-003-pool-healthmonitor/
├── update/
│   └── TC-UPDATE-001-pool-name/
└── combined/
    └── TC-COMBINED-001-full-stack/
```

### Tool Setup

**Environment Variables:**
```bash
# .env file
GCORE_API_KEY=your_key
GCORE_CLOUD_PROJECT_ID=379987
GCORE_CLOUD_REGION_ID=76
```

**Provider Override:**
```hcl
# .terraformrc file
provider_installation {
  dev_overrides {
    "gcore/gcore" = "/Users/user/repos/gcore-terraform"
  }
  direct {}
}
```

**Logging:**
```bash
export TF_LOG=DEBUG
export TF_LOG_PATH=terraform-debug.log
export TF_CLI_CONFIG_FILE=".terraformrc"
```

### Quick Reference Commands

**Build provider:**
```bash
cd /Users/user/repos/gcore-terraform
go clean -cache && go mod tidy && go build -o terraform-provider-gcore
```

**Run drift test:**
```bash
cd test-resource/drift/TC-DRIFT-001
set -o allexport; source ../../.env; set +o allexport
export TF_CLI_CONFIG_FILE="../../.terraformrc"
terraform apply -auto-approve && terraform plan -detailed-exitcode
```

**Run update test:**
```bash
cd test-resource/update/TC-UPDATE-001
set -o allexport; source ../../.env; set +o allexport
export TF_CLI_CONFIG_FILE="../../.terraformrc"
./test.sh
```

**Check API calls:**
```bash
export TF_LOG=DEBUG TF_LOG_PATH=debug.log
terraform apply
grep -E "(GET|POST|PATCH|PUT|DELETE).*http" debug.log
```

### Decision Matrix

| Scenario | Use *AndPoll? | Use UnmarshalComputed? | Mark computed_optional? |
|----------|---------------|------------------------|-------------------------|
| Async API (returns tasks) | ✅ Yes | ✅ Yes | If API computes defaults |
| Sync API (direct response) | ❌ No | ✅ Yes | If API computes defaults |
| Field user provides | - | ✅ Yes | ❌ No (use optional) |
| Field API always computes | - | ✅ Yes | ❌ No (use computed) |
| Field user OR API provides | - | ✅ Yes | ✅ Yes |
| Update uses PATCH | - | ✅ Yes | - |
| Update recreates resource | - | ✅ Yes | Mark RequiresReplace |

---

## Summary

This comprehensive guide combines proven real infrastructure testing methodology with migration patterns to ensure reliable, high-quality Terraform provider resources.

**Key Takeaways:**

1. ✅ **Test with real infrastructure** - State can lie, actual API calls don't
2. ✅ **Systematic drift testing** - Second plan must show zero changes
3. ✅ **Update verification** - PATCH vs replacement matters
4. ✅ **UnmarshalComputed everywhere** - In Read, ImportState, and after updates
5. ✅ **Computed_optional for API defaults** - Prevent drift from computed values
6. ✅ **Validate, don't assume** - Every hypothesis must be proven

**Success Criteria:**
- Zero drift on second plan after clean apply
- Updates use PATCH where appropriate
- Resource IDs stable for non-ForceNew updates
- Real cloud resources match Terraform state
- All tests pass against actual infrastructure

By following this guide, you can confidently migrate and validate Terraform provider resources with the assurance that they work correctly in production environments.

---

**Document Version:** 1.0
**Last Updated:** 2025-11-04
**Based on:** Successful Load Balancer Pool drift fix and comprehensive testing
