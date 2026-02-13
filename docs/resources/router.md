---
page_title: "gcore_router Resource - terraform-provider-gcore"
subcategory: ""
description: |-
  Represent router. Router enables you to dynamically exchange routes between networks
---

# gcore_router (Resource)

Represent router. Router enables you to dynamically exchange routes between networks

## Example Usage

#### Prerequisite

```terraform
provider "gcore" {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

data "gcore_project" "project" {
  name = "Default"
}

data "gcore_region" "region" {
  name = "Luxembourg-2"
}
```

### Basic Router

```terraform
resource "gcore_router" "basic" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name = "my-router"
}
```

### Router with Default External Gateway

Uses the default external network for internet access. This is the simplest way to give your router internet connectivity.

```terraform
resource "gcore_router" "with_default_gateway" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name = "router-with-internet"

  external_gateway_info {
    type        = "default"
    enable_snat = true
  }
}
```

### Router with Manual External Gateway

Specifies a particular external network for the gateway. Use this when you need control over which external network to use.

```terraform
# First, find the external network you want to use
data "gcore_network" "external" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name = "external-network"
}

resource "gcore_router" "with_manual_gateway" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name = "router-manual-gateway"

  external_gateway_info {
    type        = "manual"
    enable_snat = true
    network_id  = data.gcore_network.external.id
  }
}
```

### Router with Interfaces

Connects the router to multiple subnets, enabling routing between them.

```terraform
# Create networks and subnets first
resource "gcore_network" "network_a" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name          = "network-a"
  create_router = false
}

resource "gcore_subnet" "subnet_a" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name            = "subnet-a"
  cidr            = "192.168.1.0/24"
  network_id      = gcore_network.network_a.id
  gateway_ip      = "192.168.1.1"
  dns_nameservers = ["8.8.8.8", "8.8.4.4"]
}

resource "gcore_network" "network_b" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name          = "network-b"
  create_router = false
}

resource "gcore_subnet" "subnet_b" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name            = "subnet-b"
  cidr            = "192.168.2.0/24"
  network_id      = gcore_network.network_b.id
  gateway_ip      = "192.168.2.1"
  dns_nameservers = ["8.8.8.8", "8.8.4.4"]
}

# Create router connecting both subnets
resource "gcore_router" "multi_subnet" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name = "multi-subnet-router"

  external_gateway_info {
    type        = "default"
    enable_snat = true
  }

  interfaces {
    type      = "subnet"
    subnet_id = gcore_subnet.subnet_a.id
  }

  interfaces {
    type      = "subnet"
    subnet_id = gcore_subnet.subnet_b.id
  }
}
```

### Router with Static Routes

Configures custom static routes for directing traffic to specific destinations.

```terraform
resource "gcore_router" "with_routes" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name = "router-with-static-routes"

  external_gateway_info {
    type        = "default"
    enable_snat = true
  }

  # Static route to reach 10.0.0.0/8 via a specific next hop
  routes {
    destination = "10.0.0.0/8"
    nexthop     = "192.168.1.254"
  }

  # Static route to reach another network segment
  routes {
    destination = "172.16.0.0/16"
    nexthop     = "192.168.1.253"
  }
}
```

### Complete Setup

Full example creating a private network, subnet, and router with internet connectivity.

```terraform
# Complete example: Private network with internet access via router

# 1. Create a private network
resource "gcore_network" "private" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name          = "my-private-network"
  create_router = false
}

# 2. Create a subnet in the private network
resource "gcore_subnet" "private" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name            = "my-private-subnet"
  cidr            = "192.168.100.0/24"
  network_id      = gcore_network.private.id
  gateway_ip      = "192.168.100.1"
  dns_nameservers = ["8.8.8.8", "8.8.4.4"]
}

# 3. Create a router with external gateway for internet access
resource "gcore_router" "internet_gateway" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name = "internet-gateway-router"

  # Connect to external network for internet access
  external_gateway_info {
    type        = "default"
    enable_snat = true
  }

  # Connect the private subnet to the router
  interfaces {
    type      = "subnet"
    subnet_id = gcore_subnet.private.id
  }
}

# Output the router ID
output "router_id" {
  value = gcore_router.internet_gateway.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the router.

### Optional

- `external_gateway_info` (Block List, Max: 1) External gateway configuration for the router. (see [below for nested schema](#nestedblock--external_gateway_info))
- `interfaces` (Block Set) Set of interfaces attached to the router. Each interface connects the router to a subnet. (see [below for nested schema](#nestedblock--interfaces))
- `project_id` (Number) The id of the project. Either 'project_id' or 'project_name' must be specified.
- `project_name` (String) The name of the project. Either 'project_id' or 'project_name' must be specified.
- `region_id` (Number) The id of the region. Either 'region_id' or 'region_name' must be specified.
- `region_name` (String) The name of the region. Either 'region_id' or 'region_name' must be specified.
- `routes` (Block List) List of custom static routes to be advertised by the router. (see [below for nested schema](#nestedblock--routes))

### Read-Only

- `id` (String) The ID of this resource.
- `last_updated` (String) The timestamp of the last update.

<a id="nestedblock--external_gateway_info"></a>
### Nested Schema for `external_gateway_info`

Optional:

- `enable_snat` (Boolean) Whether SNAT (Source Network Address Translation) is enabled on the external gateway.
- `network_id` (String) Id of the external network
- `type` (String) Must be 'manual' or 'default'

Read-Only:

- `external_fixed_ips` (List of Object) List of external fixed IPs assigned to the router's gateway. (see [below for nested schema](#nestedatt--external_gateway_info--external_fixed_ips))

<a id="nestedatt--external_gateway_info--external_fixed_ips"></a>
### Nested Schema for `external_gateway_info.external_fixed_ips`

Read-Only:

- `ip_address` (String)
- `subnet_id` (String)



<a id="nestedblock--interfaces"></a>
### Nested Schema for `interfaces`

Required:

- `subnet_id` (String) Subnet for router interface must have a gateway IP
- `type` (String) must be 'subnet'

Read-Only:

- `ip_address` (String) The IP address assigned to the router interface.
- `mac_address` (String) The MAC address of the router interface.
- `network_id` (String) The network ID the interface is connected to.
- `port_id` (String) The port ID of the router interface.


<a id="nestedblock--routes"></a>
### Nested Schema for `routes`

Required:

- `destination` (String) The CIDR of the destination network.
- `nexthop` (String) IPv4 address to forward traffic to if it's destination IP matches 'destination' CIDR





## Import

Import is supported using the following syntax:

```shell
# import using <project_id>:<region_id>:<router_id> format
terraform import gcore_router.router1 1:6:447d2959-8ae0-4ca0-8d47-9f050a3637d7
```

