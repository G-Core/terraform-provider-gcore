# Router Configuration Examples

Copy-paste examples for different testing scenarios.

## Example 1: Router with No Interface

```hcl
resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "qa-terr-router"

  external_gateway_info = {
    enable_snat = true
    type        = "default"
  }

  # No interfaces
  interfaces = []
}
```

## Example 2: Router with One Interface

```hcl
resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "qa-terr-router"

  external_gateway_info = {
    enable_snat = true
    type        = "default"
  }

  # One interface attached
  interfaces = [
    {
      subnet_id = gcore_cloud_network_subnet.sb.id
      type      = "subnet"
    }
  ]
}
```

## Example 3: Router with Multiple Subnets

First, add a second subnet:

```hcl
resource "gcore_cloud_network_subnet" "sb2" {
  project_id = 379987
  region_id  = 76
  name       = "sys2"
  cidr       = "192.168.1.0/24"
  network_id = gcore_cloud_network.nw.id
}
```

Then attach both:

```hcl
resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "qa-terr-router"

  external_gateway_info = {
    enable_snat = true
    type        = "default"
  }

  # Two interfaces attached
  interfaces = [
    {
      subnet_id = gcore_cloud_network_subnet.sb.id
      type      = "subnet"
    },
    {
      subnet_id = gcore_cloud_network_subnet.sb2.id
      type      = "subnet"
    }
  ]
}
```

## Example 4: Using Existing Subnet (Data Source)

If you have an existing subnet:

```hcl
# Read existing subnet
data "gcore_cloud_network_subnet" "existing" {
  project_id = 379987
  region_id  = 76
  subnet_id  = "59a5f550-7fac-4b02-a834-3385b48cc79b"
}

# Attach to router
resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "qa-terr-router"

  external_gateway_info = {
    enable_snat = true
    type        = "default"
  }

  interfaces = [
    {
      subnet_id = data.gcore_cloud_network_subnet.existing.id
      type      = "subnet"
    }
  ]
}
```

## Example 5: Router with External Gateway

```hcl
resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "qa-terr-router"

  external_gateway_info = {
    enable_snat = true
    type        = "manual"
    network_id  = "external-network-id"  # Replace with actual external network ID
  }

  interfaces = [
    {
      subnet_id = gcore_cloud_network_subnet.sb.id
      type      = "subnet"
    }
  ]
}
```

## Example 6: Router with Routes

```hcl
resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "qa-terr-router"

  external_gateway_info = {
    enable_snat = true
    type        = "default"
  }

  interfaces = [
    {
      subnet_id = gcore_cloud_network_subnet.sb.id
      type      = "subnet"
    }
  ]

  # Static routes
  routes = [
    {
      destination = "10.0.0.0/24"
      nexthop     = "192.168.0.1"
    }
  ]
}
```

## Testing Sequence

### Test A: Attach → Detach → Attach

**Step 1: Create with interface**
```hcl
interfaces = [
  { subnet_id = gcore_cloud_network_subnet.sb.id, type = "subnet" }
]
```
```bash
terraform apply
# Check: POST /attach
```

**Step 2: Remove interface**
```hcl
interfaces = []
```
```bash
terraform apply
# Check: POST /detach, NO PATCH
```

**Step 3: Add back**
```hcl
interfaces = [
  { subnet_id = gcore_cloud_network_subnet.sb.id, type = "subnet" }
]
```
```bash
terraform apply
# Check: POST /attach, NO PATCH
```

### Test B: Only Interface Changes (No Other Attributes)

**Initial state:**
```hcl
name = "qa-terr-router"
interfaces = [
  { subnet_id = gcore_cloud_network_subnet.sb.id, type = "subnet" }
]
```

**Change only interfaces:**
```hcl
name = "qa-terr-router"  # Unchanged
interfaces = []           # Changed
```

**Expected:**
- ✅ POST /detach
- ❌ NO PATCH (because only interfaces changed)

### Test C: Name + Interface Change

**Initial state:**
```hcl
name = "qa-terr-router"
interfaces = [
  { subnet_id = gcore_cloud_network_subnet.sb.id, type = "subnet" }
]
```

**Change both:**
```hcl
name = "qa-terr-router-v2"  # Changed
interfaces = []              # Changed
```

**Expected:**
- ✅ POST /detach
- ✅ PATCH /routers (for name change, but WITHOUT interfaces in payload)

### Test D: Only Name Change

**Initial state:**
```hcl
name = "qa-terr-router"
interfaces = [
  { subnet_id = gcore_cloud_network_subnet.sb.id, type = "subnet" }
]
```

**Change only name:**
```hcl
name = "qa-terr-router-renamed"  # Changed
interfaces = [                    # Unchanged
  { subnet_id = gcore_cloud_network_subnet.sb.id, type = "subnet" }
]
```

**Expected:**
- ✅ PATCH /routers (for name change only)
- ❌ NO attach/detach (interfaces unchanged)

## Quick Copy-Paste Test Configs

### Minimal Router

```hcl
terraform {
  required_providers {
    gcore = { source = "gcore/gcore" }
  }
}

provider "gcore" {}

resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "test-router"
  external_gateway_info = {
    enable_snat = true
    type        = "default"
  }
  interfaces = []
}
```

### Complete Example

```hcl
terraform {
  required_providers {
    gcore = { source = "gcore/gcore" }
  }
}

provider "gcore" {}

resource "gcore_cloud_network" "nw" {
  project_id    = 379987
  region_id     = 76
  create_router = true
  name          = "test-nw"
}

resource "gcore_cloud_network_subnet" "sb" {
  project_id = 379987
  region_id  = 76
  name       = "test-subnet"
  cidr       = "192.168.10.0/24"
  network_id = gcore_cloud_network.nw.id
}

resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "test-router"
  external_gateway_info = {
    enable_snat = true
    type        = "default"
  }
  interfaces = [
    {
      subnet_id = gcore_cloud_network_subnet.sb.id
      type      = "subnet"
    }
  ]
}
```
