# SSH Key Resource Security Fix

## Summary

Modified the `gcore_cloud_ssh_key` resource to follow the old (more secure) Terraform provider behavior by making `public_key` required and removing `private_key` from state storage.

## Security Issue

**Previous Behavior (INSECURE):**
- `public_key` was **OPTIONAL** - users could omit it to have the platform generate keys
- `private_key` was **COMPUTED** and stored in Terraform state file
- This created a security risk because:
  - State files are often stored in version control or shared storage
  - Private keys should never be stored in state files
  - Anyone with access to the state file could access the private key

**New Behavior (SECURE):**
- `public_key` is **REQUIRED** - users must provide their own public key
- `private_key` is completely removed from the schema and state
- Users must generate SSH keypairs locally and keep the private key secure
- Only the public key is uploaded to the platform and stored in state

## Changes Made

### 1. Schema Changes (schema.go)

**public_key field:**
```go
// Before
"public_key": schema.StringAttribute{
    Description: "...If you want the platform to generate an Ed25519 key pair for you, leave this field empty...",
    Optional:    true,
    PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
},

// After
"public_key": schema.StringAttribute{
    Description: "...You must provide your own public key...Generate your SSH keypair locally using `ssh-keygen`...",
    Required:    true,
    PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
},
```

**private_key field:**
```go
// Before - Field existed and stored in state
"private_key": schema.StringAttribute{
    Description: "The private part of an SSH key...",
    Computed:    true,
},

// After - Field completely removed
// (No longer exists in schema)
```

### 2. Model Changes (model.go)

```go
// Before
type CloudSSHKeyModel struct {
    ID              types.String      `tfsdk:"id" json:"id,computed"`
    ProjectID       types.Int64       `tfsdk:"project_id" path:"project_id,optional"`
    Name            types.String      `tfsdk:"name" json:"name,required"`
    PublicKey       types.String      `tfsdk:"public_key" json:"public_key,optional"`  // Was optional
    SharedInProject types.Bool        `tfsdk:"shared_in_project" json:"shared_in_project,computed_optional"`
    CreatedAt       timetypes.RFC3339 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
    Fingerprint     types.String      `tfsdk:"fingerprint" json:"fingerprint,computed"`
    PrivateKey      types.String      `tfsdk:"private_key" json:"private_key,computed,no_refresh"`  // Removed
    State           types.String      `tfsdk:"state" json:"state,computed"`
}

// After
type CloudSSHKeyModel struct {
    ID              types.String      `tfsdk:"id" json:"id,computed"`
    ProjectID       types.Int64       `tfsdk:"project_id" path:"project_id,optional"`
    Name            types.String      `tfsdk:"name" json:"name,required"`
    PublicKey       types.String      `tfsdk:"public_key" json:"public_key,required"`  // Now required
    SharedInProject types.Bool        `tfsdk:"shared_in_project" json:"shared_in_project,computed_optional"`
    CreatedAt       timetypes.RFC3339 `tfsdk:"created_at" json:"created_at,computed" format:"date-time"`
    Fingerprint     types.String      `tfsdk:"fingerprint" json:"fingerprint,computed"`
    State           types.String      `tfsdk:"state" json:"state,computed"`
    // PrivateKey field removed entirely
}
```

### 3. Resource Implementation (resource.go)

No changes were needed to resource.go - the API behavior remains the same, but now users must always provide a public_key.

## Comparison with Old Provider

### Old Provider (resource_gcore_keypair.go) - SECURE
```go
"public_key": &schema.Schema{
    Type:     schema.TypeString,
    Required: true,      // Public key was required
    ForceNew: true,
},
// No private_key field at all
```

### New Provider - NOW MATCHES OLD SECURE BEHAVIOR
```go
"public_key": schema.StringAttribute{
    Required:      true,  // Now required like old provider
    PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
},
// No private_key field (removed)
```

## API Behavior

The API supports two modes, but we now only support Mode 1 (secure):

**Mode 1: Import existing public key (SUPPORTED - SECURE)**
- User provides `public_key` in request
- API returns keypair info WITHOUT `private_key`
- User already has the private key locally
- ✅ This is the secure approach we now enforce

**Mode 2: Generate new keypair (NO LONGER SUPPORTED)**
- User omits `public_key` in request
- API generates new Ed25519 keypair
- API returns both `public_key` and `private_key` once
- ❌ We no longer support this to avoid storing private keys in state

## Usage Example

### Correct Usage (Required)

```hcl
# Generate SSH key locally first
# Run: ssh-keygen -t ed25519 -f ~/.ssh/gcore_key -N ""

resource "gcore_cloud_ssh_key" "example" {
  name              = "my-ssh-key"
  project_id        = 379987
  public_key        = file("~/.ssh/gcore_key.pub")  # Required!
  shared_in_project = true
}

output "ssh_key_id" {
  value = gcore_cloud_ssh_key.example.id
}

output "fingerprint" {
  value = gcore_cloud_ssh_key.example.fingerprint
}

# Note: private_key is NOT available in output or state
# Users must keep ~/.ssh/gcore_key secure on their local machine
```

### Incorrect Usage (Will Fail)

```hcl
# This will fail with: "The argument 'public_key' is required"
resource "gcore_cloud_ssh_key" "example" {
  name       = "my-ssh-key"
  project_id = 379987
  # public_key is omitted - ERROR!
}
```

## Testing Results

### Test 1: Create with public_key (SUCCESS)
```bash
$ terraform apply
gcore_cloud_ssh_key.test_user_provided: Creating...
gcore_cloud_ssh_key.test_user_provided: Creation complete after 1s [id=a4009035-30f9-44a0-a692-fc97bf3d2edb]

Outputs:
ssh_key_fingerprint = "b8:8a:11:01:b1:3d:1e:d2:a0:02:c9:66:37:a5:8e:9f"
ssh_key_id = "a4009035-30f9-44a0-a692-fc97bf3d2edb"
ssh_key_name = "test-ssh-key-user-provided"
```

### Test 2: Verify no private_key in state (SUCCESS)
```bash
$ cat terraform.tfstate | grep -i "private"
# No results - private_key is NOT stored in state ✅
```

### Test 3: Omit public_key (FAILS AS EXPECTED)
```bash
$ terraform validate
Error: Missing required argument
The argument "public_key" is required, but no definition was found.
```

## Migration Guide for Users

### For Users of Old Provider
No changes needed! The behavior matches the old provider exactly:
- `public_key` is required (same as before)
- No `private_key` in state (same as before)

### For Users of Current Provider (Before This Fix)

**If you were using generated keys (omitting public_key):**

1. **Before upgrading**, extract your private keys from state:
   ```bash
   terraform show -json | grep private_key
   # Save these private keys securely!
   ```

2. **After upgrading**, you must provide public keys:
   ```hcl
   resource "gcore_cloud_ssh_key" "example" {
     name       = "my-key"
     project_id = 379987
     public_key = file("~/.ssh/id_ed25519.pub")  # Now required
   }
   ```

**If you were already providing public_key:**
- No changes needed
- Your configuration will continue to work
- State will no longer store private_key (more secure)

## Security Best Practices

### How to Generate SSH Keys Locally

```bash
# Generate Ed25519 key (recommended)
ssh-keygen -t ed25519 -f ~/.ssh/gcore_key -C "my-email@example.com"

# Generate RSA key (alternative)
ssh-keygen -t rsa -b 4096 -f ~/.ssh/gcore_key -C "my-email@example.com"

# Set proper permissions
chmod 600 ~/.ssh/gcore_key
chmod 644 ~/.ssh/gcore_key.pub
```

### Using the Keys with Terraform

```hcl
resource "gcore_cloud_ssh_key" "my_key" {
  name       = "my-gcore-key"
  project_id = var.project_id
  public_key = file("~/.ssh/gcore_key.pub")  # Public key is safe to share
}

# The private key (~/.ssh/gcore_key) stays on your local machine
# Never commit private keys to version control!
```

### Why This is More Secure

1. **Private keys never leave your machine** - They're not sent to the API or stored in Terraform state
2. **State files can be safely shared** - Only public keys (which are safe) are in state
3. **No accidental exposure** - You can commit state files to git without exposing private keys
4. **Industry standard practice** - This matches how AWS, GCP, Azure handle SSH keys
5. **User control** - Users manage their own key lifecycle and security

## Files Modified

1. `internal/services/cloud_ssh_key/schema.go` - Made public_key required, removed private_key
2. `internal/services/cloud_ssh_key/model.go` - Updated model to match schema changes
3. `test-ssh-key/main.tf` - Updated test to use local public key

## Conclusion

This change aligns the new Terraform provider with the security best practices established in the old provider. Users are now required to generate and manage their own SSH keypairs locally, with only the public key being uploaded to the platform. Private keys never enter the Terraform state, significantly reducing the risk of accidental exposure.

**Migration Impact:**
- ✅ **Low risk** for users already providing public_key
- ⚠️ **Breaking change** for users relying on key generation (must now generate keys locally)
- 🔒 **Security improvement** - No more private keys in state files
