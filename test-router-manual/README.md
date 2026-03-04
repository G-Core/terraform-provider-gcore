# Manual Router Testing Example

Simple Terraform configuration for manual testing of router interface operations.

## Setup

```bash
# From the gcore-terraform root directory

# Terminal 1: Start mitmproxy (optional)
./run_mitm.sh

# Terminal 2: Set environment and run terraform
source ./set_env.sh
cd test-router-manual
```

## Basic Usage

### 1. Initialize

```bash
terraform init
```

### 2. Create resources

```bash
terraform plan
terraform apply
```

**Creates:**
- Network: `qa-terr-nw`
- Subnet: `sys` (192.168.0.0/24)
- Router: `qa-terr-router` with 1 interface attached

### 3. View outputs

```bash
terraform output
```

### 4. Clean up

```bash
terraform destroy
```

## Testing Scenarios

### Scenario 1: Verify Initial Creation

**main.tf:**
```hcl
# Router with one interface (current config)
interfaces = [
  {
    subnet_id = gcore_cloud_network_subnet.sb.id
    type      = "subnet"
  }
]
```

**Run:**
```bash
terraform apply
```

**Expected HTTP calls:**
- POST /networks (create network)
- POST /subnets (create subnet)
- POST /routers (create router)
- POST /routers/.../attach (attach subnet)

### Scenario 2: Remove Interface (Detach)

**Edit main.tf:**
```hcl
# Remove interface
interfaces = []
```

**Run:**
```bash
terraform apply
```

**Expected HTTP calls:**
- ✅ POST /routers/.../detach (detach subnet)
- ❌ NO PATCH /routers

**Verify:**
```bash
# If using mitmproxy
mitmdump -nr flow.mitm "~u /detach"
# Should show: POST .../detach

mitmdump -nr flow.mitm "~m PATCH & ~u /routers/"
# Should be empty (no PATCH requests)
```

### Scenario 3: Add Interface Back (Attach)

**Edit main.tf:**
```hcl
# Add interface back
interfaces = [
  {
    subnet_id = gcore_cloud_network_subnet.sb.id
    type      = "subnet"
  }
]
```

**Run:**
```bash
terraform apply
```

**Expected HTTP calls:**
- ✅ POST /routers/.../attach (attach subnet)
- ❌ NO PATCH /routers

**Verify:**
```bash
mitmdump -nr flow.mitm "~u /attach"
# Should show: POST .../attach
```

### Scenario 4: Rename Router (PATCH allowed)

**Edit main.tf:**
```hcl
resource "gcore_cloud_network_router" "router" {
  # ...
  name = "qa-terr-router-renamed"  # Changed name
  # interfaces unchanged
}
```

**Run:**
```bash
terraform apply
```

**Expected HTTP calls:**
- ✅ PATCH /routers (update name)
- ❌ NO attach/detach (interfaces unchanged)

### Scenario 5: Rename + Remove Interface

**Edit main.tf:**
```hcl
resource "gcore_cloud_network_router" "router" {
  # ...
  name = "qa-terr-router-v2"  # Changed name
  interfaces = []              # Removed interface
}
```

**Run:**
```bash
terraform apply
```

**Expected HTTP calls:**
- ✅ POST /routers/.../detach (remove interface)
- ✅ PATCH /routers (update name, but WITHOUT interfaces in payload)

**Verify PATCH doesn't include interfaces:**
```bash
mitmdump -nr flow.mitm "~m PATCH & ~u /routers/" | grep -i "interfaces"
# Should be empty (interfaces not in PATCH payload)
```

## Quick Test Script

Create `test.sh` in this directory:

```bash
#!/bin/bash
set -e

echo "=== Test 1: Create with interface ==="
terraform apply -auto-approve
sleep 2

echo "=== Test 2: Remove interface ==="
# Edit main.tf: interfaces = []
sed -i.bak 's/subnet_id = gcore_cloud_network_subnet.sb.id/# REMOVED/' main.tf
sed -i '' 's/type      = "subnet"/# REMOVED/' main.tf
terraform apply -auto-approve
sleep 2

echo "=== Test 3: Add interface back ==="
# Restore main.tf
mv main.tf.bak main.tf
terraform apply -auto-approve

echo "=== Test complete ==="
terraform destroy -auto-approve
```

## Verification Commands

### Check HTTP requests (with mitmproxy)

```bash
# Count attach operations
mitmdump -nr flow.mitm "~u /attach" | grep -c POST

# Count detach operations
mitmdump -nr flow.mitm "~u /detach" | grep -c POST

# Count PATCH to routers
mitmdump -nr flow.mitm "~m PATCH & ~u /routers/" | wc -l

# View attach request details
mitmdump -nr flow.mitm "~u /attach" | less

# View detach request details
mitmdump -nr flow.mitm "~u /detach" | less
```

### Check Terraform state

```bash
# Show router interfaces in state
terraform state show gcore_cloud_network_router.router | grep interfaces

# Show all resources
terraform state list

# Show specific resource
terraform state show gcore_cloud_network_router.router
```

### Check provider logs

```bash
# View Terraform debug log
tail -f terraform_manual.log | grep -i router

# Search for attach/detach in logs
grep -i "attach\|detach" terraform_manual.log

# Search for PATCH in logs
grep -i "PATCH" terraform_manual.log
```

## Expected Behavior

### ✅ Correct Behavior

1. **Adding interface** → Uses POST /attach endpoint
2. **Removing interface** → Uses POST /detach endpoint
3. **Changing name only** → Uses PATCH /routers (without interfaces)
4. **Changing name + interface** → Uses POST attach/detach + PATCH (without interfaces in payload)

### ❌ Incorrect Behavior (Old Bug)

1. Adding interface → POST /attach **AND** PATCH /routers with interfaces ❌
2. Removing interface → PATCH /routers with empty interfaces ❌

## Troubleshooting

### Issue: Can't create resources

**Check credentials:**
```bash
echo $GCORE_CLOUD_API_KEY
# Should show your API key

# Or check .env
cat ../.env
```

### Issue: Dev overrides not working

**Check .terraformrc:**
```bash
cat ../.terraformrc
```

Should contain:
```hcl
provider_installation {
  dev_overrides {
    "gcore/gcore" = "/Users/user/repos/gcore-terraform"
  }
  direct {}
}
```

### Issue: Can't see HTTP requests

**Make sure mitmproxy is running:**
```bash
# Check if running
lsof -i :9092

# Check proxy env vars
echo $HTTP_PROXY
echo $HTTPS_PROXY
# Should output: http://127.0.0.1:9092
```

## Files

- `main.tf` - Main Terraform configuration
- `.terraform/` - Terraform working directory (created by init)
- `terraform.tfstate` - Terraform state file
- `.terraform.lock.hcl` - Provider lock file

## Clean Up

```bash
# Destroy resources
terraform destroy -auto-approve

# Clean Terraform files
rm -rf .terraform .terraform.lock.hcl terraform.tfstate*

# Stop mitmproxy
# Press Ctrl+C in terminal running ./run_mitm.sh
```
