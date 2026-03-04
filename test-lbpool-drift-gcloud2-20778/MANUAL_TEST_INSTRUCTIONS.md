# Manual Test Instructions for GCLOUD2-20778 Drift Issue

## Bug Description (from Jira Comment - 2025-11-20)

After adding an LB pool to a listener with `user_list` configured, Terraform detects resource drift showing:
- `stats` attribute as "(known after apply)"
- `encrypted_password` in `user_list` as changing to "(known after apply)"

## Prerequisites

1. Ensure `.env` file exists with credentials:
   ```bash
   export GCORE_API_KEY="your-api-key"
   ```

2. Terraform provider override configured in `.terraformrc`

## Step-by-Step Test Instructions

### Step 1: Initial Setup - Create LB and Listener WITHOUT Pool

Create `step1.tf` with just LB and Listener:

```hcl
terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

data "gcore_cloud_projects" "my_projects" {
  name = "default"
}

locals {
  project_id = [for p in data.gcore_cloud_projects.my_projects.items : p.id]
}

data "gcore_cloud_region" "rg" {
  region_id = 76
}

resource "gcore_cloud_load_balancer" "lb" {
  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id
  flavor     = "lb1-1-2"
  name       = "qa-drift-test-lb"
}

resource "gcore_cloud_load_balancer_listener" "ls" {
  project_id       = local.project_id[0]
  region_id        = data.gcore_cloud_region.rg.id
  load_balancer_id = gcore_cloud_load_balancer.lb.id
  name             = "qa-ls-name"
  protocol         = "HTTP"
  protocol_port    = 80

  # This user_list is key to reproducing the bug
  user_list = [
    {
      username           = "testuser"
      encrypted_password = "$5$isRr.HJ1IrQP38.m$oViu3DJOpUG2ZsjCBtbITV3mqpxxbZfyWJojLPNSPO5"
    }
  ]
}
```

**Run:**
```bash
export TF_CLI_CONFIG_FILE=/Users/user/repos/gcore-terraform/.terraformrc
source ../.env
terraform apply -auto-approve
```

**Expected:** LB and Listener created successfully.

**Verify no drift:**
```bash
terraform plan
```

**Expected:** "No changes. Your infrastructure matches the configuration."

---

### Step 2: Add the Pool - This is where drift might appear

Now add the pool resource to the **same file** (`step1.tf`):

```hcl
# Add this to the end of step1.tf

resource "gcore_cloud_load_balancer_pool" "lb_pool" {
  project_id   = local.project_id[0]
  region_id    = data.gcore_cloud_region.rg.id
  lb_algorithm = "LEAST_CONNECTIONS"
  name         = "pool-drift-test"
  protocol     = "HTTP"
  listener_id  = gcore_cloud_load_balancer_listener.ls.id
}
```

**Run:**
```bash
terraform apply -auto-approve
```

**Expected:** Pool created and attached to listener successfully.

---

### Step 3: Check for Drift - This is the critical test

**Run:**
```bash
terraform plan
```

**What to look for:**

#### ✅ If the bug is FIXED (current state):
```
No changes. Your infrastructure matches the configuration.
```

#### ❌ If the bug EXISTS (as reported in Jira):
```
Terraform will perform the following actions:

  # gcore_cloud_load_balancer_listener.ls will be updated in-place
  ~ resource "gcore_cloud_load_balancer_listener" "ls" {
        id                     = "a6ea2cc2-5cd5-459d-a4a1-b2cffdbed0c3"
        name                   = "qa-ls-name"
      + stats                  = (known after apply)
      ~ user_list              = [
          ~ {
              ~ encrypted_password = "$2a$10$vp25Soo.i6aYcyWtrfV7SeBsXMa1GjMRHZRJyMCTAiY/T6j8kXv7u" -> (known after apply)
                # (1 unchanged attribute hidden)
            },
        ]
        # (15 unchanged attributes hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.
```

---

### Step 4: Additional Verification (Optional)

If you want to see what happens after applying the drift:

```bash
# If drift was detected
terraform apply -auto-approve

# Then check again
terraform plan
```

**Question:** Does the drift keep appearing on every plan, or does it stabilize?

---

## Quick One-Line Test

If you want to test quickly using the existing test directory:

```bash
cd /Users/user/repos/gcore-terraform/test-lbpool-drift-gcloud2-20778
source ../.env
export TF_CLI_CONFIG_FILE=/Users/user/repos/gcore-terraform/.terraformrc
terraform plan
```

**Current Result:** No changes detected (bug is fixed)

---

## Cleanup

```bash
terraform destroy -auto-approve
```

---

## Understanding the Root Cause

The bug occurred because:

1. **Without `UseStateForUnknown()` plan modifiers:**
   - When pool is added, Terraform refreshes the listener state
   - API returns `stats` and `user_list` with potentially different representations
   - Terraform sees them as "unknown" and marks them as "(known after apply)"
   - This creates false drift detection

2. **With `UseStateForUnknown()` plan modifiers (the fix):**
   - When field value isn't changing in config
   - And Terraform doesn't know the future value from API
   - It uses the current state value instead of showing "(known after apply)"
   - No false drift

---

## Files in This Test Directory

- `main.tf` - Full test configuration (already set up)
- `FINDINGS.md` - Detailed investigation report
- `MANUAL_TEST_INSTRUCTIONS.md` - This file

---

## Notes

- The bug was reported on **2025-11-20**
- The fix (commit f999c1c) was applied on **2025-11-18** (2 days earlier)
- This suggests the QA tester was using an older build
- Current provider code already has the fix applied
