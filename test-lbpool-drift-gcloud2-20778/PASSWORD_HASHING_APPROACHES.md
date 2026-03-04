# Password Hashing Approaches in Terraform

## Problem: bcrypt() Function Causes Drift

The Terraform `bcrypt()` function generates a **new random salt** on every execution, causing the hash to change even with the same input. This results in spurious diffs and unnecessary resource updates.

**Official Documentation:**
- https://developer.hashicorp.com/terraform/language/functions/bcrypt

> "Since a bcrypt hash value includes a randomly selected salt, each call to this function will return a different value, even if the given string and cost are the same. Using this function directly with resource arguments will therefore cause spurious diffs."

---

## Three Approaches Compared

### ❌ Approach 1: bcrypt() Function (NOT RECOMMENDED)

```hcl
locals {
  hash = bcrypt("Test123!")
}

resource "gcore_cloud_load_balancer_listener" "ls" {
  # ... other config ...
  user_list = [{
    username           = "qauser"
    encrypted_password = local.hash  # ❌ Changes on every run!
  }]
}
```

**Result:** ❌ **DRIFT ON EVERY PLAN**
```
~ user_list = [
    ~ {
        ~ encrypted_password = "$2a$10$abc..." -> (known after apply)
      }
  ]
```

---

### ✅ Approach 2: random_password.bcrypt_hash (RECOMMENDED)

```hcl
terraform {
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "~> 3.6"
    }
  }
}

# Generate random password with stable bcrypt hash
resource "random_password" "qa_password" {
  length  = 16
  special = false

  lifecycle {
    ignore_changes = [length, special]  # Optional: prevent changes
  }
}

resource "gcore_cloud_load_balancer_listener" "ls" {
  # ... other config ...
  user_list = [{
    username           = "qauser"
    encrypted_password = random_password.qa_password.bcrypt_hash  # ✅ Stable!
  }]
}

output "password" {
  value     = random_password.qa_password.result
  sensitive = true
}
```

**How it works:**
- `random_password` generates a random password once
- The `bcrypt_hash` attribute is computed and **stored in state**
- Salt is persisted, so hash remains consistent across runs
- No drift on subsequent plans

**Result:** ✅ **NO DRIFT**
```
No changes. Your infrastructure matches the configuration.
```

**References:**
- https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/password
- https://github.com/hashicorp/terraform-provider-random/issues/102

---

### ✅ Approach 3: Hardcoded Hash (SIMPLE BUT MANUAL)

```hcl
resource "gcore_cloud_load_balancer_listener" "ls" {
  # ... other config ...
  user_list = [{
    username           = "qauser"
    encrypted_password = "$5$isRr.HJ1IrQP38.m$oViu3DJOpUG2ZsjCBtbITV3mqpxxbZfyWJojLPNSPO5"  # ✅ Stable!
  }]
}
```

**How to generate hash manually:**
```bash
# Python
python3 -c 'import crypt; print(crypt.crypt("Test123!", crypt.mksalt(crypt.METHOD_SHA256)))'

# Or use OpenSSL/htpasswd
htpasswd -nbB username password | cut -d: -f2
```

**Result:** ✅ **NO DRIFT**

---

## Recommendation Matrix

| Use Case | Recommended Approach | Why |
|----------|---------------------|-----|
| **Production deployments** | `random_password.bcrypt_hash` | Automatic, secure, no drift |
| **Drift testing** | Hardcoded hash | Simple, predictable |
| **One-time provisioning** | `bcrypt()` in provisioner | Acceptable for non-tracked operations |
| **Dynamic passwords** | `random_password.bcrypt_hash` | State-managed, consistent |

---

## Test Results (GCLOUD2-20778)

**Configuration:**
- Load Balancer + Listener with `connection_limit = 5000`
- `user_list` with bcrypt password
- Pool added after listener created

**Test 1: Hardcoded Hash**
```
✅ terraform plan #1: No changes
✅ terraform plan #2: No changes
✅ Pool added: No drift detected
```

**Test 2: bcrypt() Function**
```
❌ terraform plan #1: Shows change in encrypted_password
❌ terraform plan #2: Shows change in encrypted_password
⚠️  Drift unrelated to pool, caused by bcrypt() regeneration
```

**Test 3: random_password.bcrypt_hash** (see `main-random-password-example.tf`)
```
✅ terraform plan #1: No changes (after initial apply)
✅ terraform plan #2: No changes
✅ Pool added: No drift detected
```

---

## Conclusion

**For GCLOUD2-20778 drift testing:**
- ✅ Hardcoded hash approach is working correctly
- ✅ No drift from pool addition with `connection_limit = 5000`
- ✅ Issue appears to be resolved in current provider code

**For production use:**
- Use `random_password.bcrypt_hash` for automated, drift-free password management
- **Never use `bcrypt()` function directly in resource arguments**
