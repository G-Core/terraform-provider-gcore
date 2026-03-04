# Sensitive: true Analysis for encrypted_password

## What Does `Sensitive: true` Mean in Terraform?

### Behavior

When an attribute is marked as `Sensitive: true`, Terraform:

1. **Redacts values in CLI output**
   ```
   # Without Sensitive:
   + encrypted_password = "$5$isRr.HJ1IrQP38.m$oViu3DJOpUG2ZsjCBtbITV3mqpxxbZfyWJojLPNSPO5"

   # With Sensitive:
   + encrypted_password = (sensitive value)
   ```

2. **Hides values in plan/apply output**
   - Plan shows: `(sensitive value)` instead of actual value
   - Apply shows: `(sensitive value)` instead of actual value
   - State file: **STILL CONTAINS THE ACTUAL VALUE** (state is not encrypted by default)

3. **Prevents accidental exposure in logs**
   - CI/CD logs won't show the value
   - Terminal history won't capture the value
   - Screenshots won't reveal the value

4. **Taints dependent values**
   - Any expression using a sensitive value becomes sensitive
   - Prevents indirect leakage through string interpolation

### Important: What Sensitive Does NOT Do

❌ **Does NOT encrypt state file** - Value is still stored in plaintext in `terraform.tfstate`
❌ **Does NOT encrypt remote state** - Unless backend has encryption enabled
❌ **Does NOT validate the value** - Just hides it from output
❌ **Does NOT prevent API logging** - API requests still contain the value
❌ **Does NOT make the value more secure** - Just reduces exposure surface

---

## Current Implementation Comparison

### Old Provider (terraform-provider-gcore)
```go
"encrypted_password": &schema.Schema{
    Type:        schema.TypeString,
    Description: "Encrypted password (hash) to auth via Basic Authentication",
    Sensitive:   true,  // ✅ Marked as sensitive
    Required:    true,
},
```

### New Provider (gcore-terraform)
```go
"encrypted_password": schema.StringAttribute{
    Description: "Encrypted password to auth via Basic Authentication",
    Required:    true,
    // ❌ NOT marked as sensitive
},
```

---

## Should encrypted_password Be Marked Sensitive?

### Arguments FOR Marking It Sensitive ✅

1. **Defense in Depth**
   - Even though it's already encrypted/hashed, hiding it prevents:
     - Hash cracking attempts from logs
     - Rainbow table attacks on exposed hashes
     - Unnecessary exposure in CI/CD pipelines

2. **Security Best Practice**
   - Security principle: Don't expose credentials/hashes unnecessarily
   - Even bcrypt/sha256 hashes can be cracked with enough resources
   - Better safe than sorry

3. **Compliance Requirements**
   - Some security audits flag ANY password-related fields in logs
   - SOC2/ISO27001 compliance may require hiding password hashes

4. **Consistency with Old Provider**
   - Old provider had it as sensitive
   - Migration from old to new shouldn't reduce security posture
   - Users might expect this behavior

5. **Prevents Social Engineering**
   - Attackers can't see which hash algorithm is used
   - Harder to profile the system's security setup

### Arguments AGAINST Marking It Sensitive ❌

1. **It's Already Encrypted**
   - The field name itself says "encrypted_password"
   - It's a one-way hash (bcrypt, sha256, etc.)
   - Not a plaintext password
   - Example value: `$5$isRr.HJ1IrQP38.m$oViu3DJOpUG2ZsjCBtbITV3mqpxxbZfyWJojLPNSPO5`

2. **Debugging Difficulty**
   - Can't see in plan output if hash changed
   - Harder to troubleshoot configuration issues
   - Can't easily verify the hash format is correct

3. **False Sense of Security**
   - Marking as sensitive doesn't encrypt state file
   - State file still contains the hash in plaintext
   - If state is compromised, hash is exposed anyway

4. **User Responsibility**
   - Users should generate secure hashes
   - Users should secure their state files
   - Terraform provider shouldn't hide implementation details

---

## Recommendation

### **YES, it SHOULD be marked `Sensitive: true`** ✅

**Reasoning:**

1. **Better Safe Than Sorry**
   - Minimal cost (just hides from output)
   - Significant benefit (reduces exposure surface)
   - No functional downside

2. **Security Best Practice**
   - Even hashed passwords should be treated as sensitive
   - Modern password hashes (bcrypt) are designed to be slow, but not impossible to crack
   - Exposure in logs could enable offline attacks

3. **Backwards Compatibility**
   - Old provider had this as sensitive
   - Users migrating expect same security level
   - Removing sensitivity is a security regression

4. **Compliance**
   - Many security standards require hiding password hashes
   - Easier to pass security audits

### How to Fix

In Plugin Framework (new provider):

```go
"encrypted_password": schema.StringAttribute{
    Description: "Encrypted password to auth via Basic Authentication",
    Required:    true,
    Sensitive:   true,  // ← Add this
},
```

---

## Real-World Impact Example

### Without Sensitive (Current)
```bash
$ terraform plan

  # gcore_cloud_load_balancer_listener.ls will be updated in-place
  ~ resource "gcore_cloud_load_balancer_listener" "ls" {
      ~ user_list = [
          ~ {
              ~ encrypted_password = "$5$isRr.HJ1IrQP38.m$oViu3DJOpUG2ZsjCBtbITV3mqpxxbZfyWJojLPNSPO5" -> "$2a$10$abc..."
            },
        ]
    }
```
**Problem:** Hash is visible in CI/CD logs, screenshots, terminal history

### With Sensitive
```bash
$ terraform plan

  # gcore_cloud_load_balancer_listener.ls will be updated in-place
  ~ resource "gcore_cloud_load_balancer_listener" "ls" {
      ~ user_list = [
          ~ {
              ~ encrypted_password = (sensitive value)
            },
        ]
    }
```
**Benefit:** Hash is hidden, reducing attack surface

---

## Additional Context: Hash Security

Even "encrypted" passwords (hashes) can be attacked:

1. **Hash Algorithms Shown in Output:**
   - `$5$` = SHA-256 crypt
   - `$2a$10$` = bcrypt (cost factor 10)
   - This reveals your hashing strategy

2. **Offline Attacks:**
   - Attackers can use GPUs to crack hashes offline
   - bcrypt is slow, but not impossible
   - With hash visible, they can work on it indefinitely

3. **Rainbow Tables:**
   - Pre-computed hash tables exist
   - Common passwords can be cracked instantly
   - Even salted hashes aren't immune to targeted attacks

---

## Conclusion

**YES, mark `encrypted_password` as `Sensitive: true`**

This is a simple change with significant security benefits and no real downsides. It aligns with:
- Security best practices
- The old provider's behavior
- User expectations
- Compliance requirements

The fact that it's a hash (not plaintext) doesn't mean it should be exposed freely in logs and output.
