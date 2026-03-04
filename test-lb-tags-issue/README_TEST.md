# Tags Inconsistency Error Reproduction Test

## Error to Reproduce

```
Error: Provider produced inconsistent result after apply

When applying changes to gcore_cloud_load_balancer.lb, provider
"provider["local.gcore.com/repo/gcore"]" produced an unexpected new value: .tags_v2: new
element 0 has appeared.
```

## Prerequisites

1. **Environment setup:**
   ```bash
   # Ensure .env file exists in parent directory
   cat ../.env
   # Should contain:
   # GCORE_API_KEY='...'
   # GCORE_CLIENT=...
   # GCORE_CLOUD_PROJECT_ID=...
   # GCORE_CLOUD_REGION_ID=76
   ```

2. **Provider override:**
   ```bash
   # Ensure .terraformrc exists
   cat ../.terraformrc
   # Should contain:
   # provider_installation {
   #   dev_overrides {
   #     "gcore/gcore" = "/Users/user/repos/gcore-terraform"
   #   }
   #   direct {}
   # }
   ```

3. **Provider built:**
   ```bash
   cd .. && go build -o terraform-provider-gcore_v99.0.0
   ```

## Test Files

- **reproduce_tags_error.tf** - Main test configuration with 4 scenarios
- **Makefile** - Automated test commands
- **README_TEST.md** - This file

## Quick Start

### Option 1: Using Makefile (Recommended)

```bash
# Run all automated tests
make test-all

# Or run individual scenarios
make test-scenario-1    # Create LB with tags from start
make test-scenario-3    # Exact Jira reproduction
```

### Option 2: Manual Testing

```bash
# Source environment
source ../.env
export TF_CLI_CONFIG_FILE=../.terraformrc

# Apply configuration
terraform apply -auto-approve

# Check for errors in output
```

## Test Scenarios

### Scenario 1: Create LB with Tags from Start

**Purpose:** Verify tags work when set during initial creation

**Resource:** `gcore_cloud_load_balancer.lb_with_tags`

**Steps:**
```bash
make test-scenario-1
```

**Expected:** ✅ No error, tags and tags_v2 both populated correctly

**What to Check:**
```bash
terraform output lb_with_tags_info
# Should show:
# tags = {"qa" = "load-balancer", "environment" = "test"}
# tags_v2 = [{"key"="qa", "value"="load-balancer", "read_only"=false}, ...]
```

---

### Scenario 2: Add Tags to Existing LB (Manual)

**Purpose:** Reproduce error when adding tags to LB that had none

**Resource:** `gcore_cloud_load_balancer.lb_add_tags_later`

**Steps:**

1. **Create LB without tags:**
   ```bash
   # Tags are commented out initially
   make test-scenario-2
   ```

2. **Verify LB exists without tags:**
   ```bash
   terraform show | grep -A 20 "lb_add_tags_later"
   # Should show: tags = null, tags_v2 = []
   ```

3. **Add tags:**
   - Edit `reproduce_tags_error.tf`
   - Uncomment the `tags` block for `lb_add_tags_later`
   - Save file

4. **Apply the change:**
   ```bash
   terraform apply -auto-approve 2>&1 | tee scenario2_result.log
   ```

5. **Check for error:**
   ```bash
   grep -i "inconsistent\|tags_v2.*appeared" scenario2_result.log
   ```

**Expected Result:**
- ❌ If bug exists: Error about `tags_v2: new element 0 has appeared`
- ✅ If bug fixed: Apply succeeds, tags added correctly

---

### Scenario 3: Exact Jira Reproduction

**Purpose:** Test exact configuration from Jira issue

**Resource:** `gcore_cloud_load_balancer.lb`

**Steps:**
```bash
make test-scenario-3
```

**Configuration:**
```hcl
resource "gcore_cloud_load_balancer" "lb" {
  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id
  flavor     = "lb1-2-4"
  name       = "qa-lb-name"

  tags = {
    "qa" = "load-balancer"
  }
}
```

**Expected:** Should match Jira error if bug exists

**What to Check:**
```bash
# Check apply output
terraform output lb_main_info

# Verify state
terraform state show gcore_cloud_load_balancer.lb
```

---

### Scenario 4: Modify Existing Tags (Manual)

**Purpose:** Test changing tags on LB that already has them

**Resource:** `gcore_cloud_load_balancer.lb_modify_tags`

**Steps:**

1. **Create LB with initial tag:**
   ```bash
   make test-scenario-4
   ```

2. **Verify initial state:**
   ```bash
   terraform output lb_modify_tags_info
   # Should show: tags = {"initial" = "tag"}
   ```

3. **Modify tags:**
   - Edit `reproduce_tags_error.tf`
   - Comment out current tags block
   - Uncomment the alternative tags block
   - Save file

4. **Apply changes:**
   ```bash
   terraform apply -auto-approve 2>&1 | tee scenario4_result.log
   ```

5. **Check for error:**
   ```bash
   grep -i "inconsistent\|tags_v2" scenario4_result.log
   ```

**Expected Result:**
- ✅ Tags should update without error
- `tags_v2` should reflect new tags

---

## What to Look For

### Success Indicators (Bug Fixed)

```
Apply complete! Resources: X added, Y changed, 0 destroyed.

Outputs:

lb_with_tags_info = {
  "id" = "abc-123..."
  "tags" = tomap({
    "qa" = "load-balancer"
  })
  "tags_v2" = tolist([
    {
      "key" = "qa"
      "read_only" = false
      "value" = "load-balancer"
    },
  ])
}
```

### Error Indicators (Bug Exists)

```
╷
│ Error: Provider produced inconsistent result after apply
│
│ When applying changes to gcore_cloud_load_balancer.lb, provider
│ "provider["local.gcore.com/repo/gcore"]" produced an unexpected new value:
│ .tags_v2: new element 0 has appeared.
│
│ This is a bug in the provider, which should be reported in the provider's
│ own issue tracker.
╵
```

### Drift Indicators

After successful apply, run:
```bash
terraform plan -detailed-exitcode
```

**Expected:** Exit code 0 (no changes)

**If drift exists:**
```
Terraform will perform the following actions:

  # gcore_cloud_load_balancer.lb will be updated in-place
  ~ resource "gcore_cloud_load_balancer" "lb" {
      ~ tags_v2 = [...] -> (known after apply)
  }
```

---

## Debugging Steps

### 1. Check Provider Version
```bash
ls -lh ../terraform-provider-gcore_v99.0.0
git log --oneline -5
```

### 2. Enable Debug Logging
```bash
export TF_LOG=DEBUG
export TF_LOG_PATH=./terraform-debug.log
terraform apply -auto-approve
grep -i "tags" terraform-debug.log
```

### 3. Inspect State File
```bash
# Check tags and tags_v2 in state
jq '.resources[] | select(.type=="gcore_cloud_load_balancer") | {name: .name, tags: .instances[0].attributes.tags, tags_v2: .instances[0].attributes.tags_v2}' terraform.tfstate
```

### 4. Check API Responses
```bash
# If you have mitmproxy or similar
grep -A 50 "PATCH.*loadbalancers" mitm_requests.log
```

### 5. Compare with Old Provider
```bash
# Test same config with old provider
cd ../test-old-provider-comparison
# Copy reproduce_tags_error.tf and test
```

---

## Cleanup

```bash
# Destroy all resources
make destroy

# Or manually
terraform destroy -auto-approve

# Clean up artifacts
make clean
```

---

## Reporting Results

After testing, document:

1. **Which scenarios failed/passed:**
   - [ ] Scenario 1: Create with tags
   - [ ] Scenario 2: Add tags later
   - [ ] Scenario 3: Jira exact reproduction
   - [ ] Scenario 4: Modify tags

2. **Error messages:** (copy full error)

3. **Provider version/commit:**
   ```bash
   git rev-parse HEAD
   ```

4. **State file comparison:**
   ```bash
   # Before apply
   terraform show -json > before.json

   # After apply (if successful)
   terraform show -json > after.json

   # Compare
   diff before.json after.json
   ```

5. **Logs:** Attach `terraform-debug.log` if available

---

## Expected Outcome

With the current code (commit `216b73a` and later), all scenarios should **PASS** without errors.

If any scenario fails with the tags inconsistency error, this indicates:
- The bug still exists
- There's a specific condition triggering it
- Further investigation needed in the provider code
