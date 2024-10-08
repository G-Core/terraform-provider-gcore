---
page_title: "gcore_lbmember Resource - terraform-provider-gcore"
subcategory: ""
description: |-
  Represent load balancer member
---

# gcore_lbmember (Resource)

Represent load balancer member

## Example Usage

### Prerequisite

```terraform
provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

data "gcore_project" "project" {
  name = "Default"
}

data "gcore_region" "region" {
  name = "Luxembourg-2"
}
```

```terraform
resource "gcore_loadbalancerv2" "lb" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name       = "My first complex load balancer"
  flavor     = "lb1-1-2"
}

resource "gcore_lblistener" "http_80" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  loadbalancer_id = gcore_loadbalancerv2.lb.id

  name          = "http-80"
  protocol      = "HTTP"
  protocol_port = 80
}

resource "gcore_lbpool" "http" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  loadbalancer_id = gcore_loadbalancerv2.lb.id
  listener_id     = gcore_lblistener.http_80.id

  name            = "My HTTP pool"
  protocol        = "HTTP"
  lb_algorithm    = "ROUND_ROBIN"

  health_monitor {
    type        = "TCP"
    delay       = 10
    max_retries = 3
    timeout     = 5
  }
}
```

### Public member

```terraform
resource "gcore_lbmember" "public_member" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  pool_id       = gcore_lbpool.http.id

  address       = "8.8.8.8"
  protocol_port = 80
  weight        = 1
}
```

### Private member

```terraform
resource "gcore_network" "private_network" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name = "my-private-network"
}

resource "gcore_subnet" "private_subnet" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  cidr       = "10.0.0.0/24"
  name       = "my-private-network-subnet"
  network_id = gcore_network.private_network.id
}

resource "gcore_reservedfixedip" "fixed_ip" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  type             = "ip_address"
  network_id       = gcore_network.private_network.id
  subnet_id        = gcore_subnet.private_subnet.id
  fixed_ip_address = "10.0.0.10"
  is_vip           = false
}

resource "gcore_lbmember" "private_member" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  pool_id = gcore_lbpool.http.id

  address   = gcore_reservedfixedip.fixed_ip.fixed_ip_address
  subnet_id = gcore_reservedfixedip.fixed_ip.subnet_id

  protocol_port = 80
  weight        = 1
}
```

### Private Instance member

```terraform
resource "gcore_network" "instance_member_private_network" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name = "my-private-network"
}

resource "gcore_subnet" "instance_member_private_subnet" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  cidr       = "10.0.0.0/24"
  name       = "my-private-network-subnet"
  network_id = gcore_network.instance_member_private_network.id
}

resource "gcore_reservedfixedip" "instance_member_fixed_ip" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  type             = "ip_address"
  network_id       = gcore_network.instance_member_private_network.id
  subnet_id        = gcore_subnet.instance_member_private_subnet.id
  fixed_ip_address = "10.0.0.11"
  is_vip           = false
}

data "gcore_image" "ubuntu" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name       = "ubuntu-22.04"
}

data "gcore_securitygroup" "default" {
  name       = "default"
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}

resource "gcore_volume" "instance_member_volume" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name       = "boot volume"
  type_name  = "ssd_hiiops"
  size       = 10
  image_id   = data.gcore_image.ubuntu.id
}


resource "gcore_instancev2" "instance_member" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name_template = "ed-c16-{ip_octets}"
  flavor_id  = "g1-standard-1-2"

  volume {
    volume_id  = gcore_volume.instance_member_volume.id
    boot_index = 0
  }

  interface {
    type            = "reserved_fixed_ip"
    name            = "my-private-network-interface"
    port_id         = gcore_reservedfixedip.instance_member_fixed_ip.port_id
    security_groups = [data.gcore_securitygroup.default.id]
  }
}

resource "gcore_lbmember" "instance_member" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  pool_id       = gcore_lbpool.http.id

  instance_id = gcore_instancev2.instance_member.id
  address       = gcore_reservedfixedip.instance_member_fixed_ip.fixed_ip_address
  protocol_port = 80
  weight        = 1
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `address` (String) IP address to communicate with real server.
- `pool_id` (String) ID of the target load balancer pool to attach newly created member.
- `protocol_port` (Number) Port to communicate with real server.

### Optional

- `instance_id` (String) ID of the gcore_instance.
- `project_id` (Number) ID of the desired project to create load balancer member in. Alternative for `project_name`. One of them should be specified.
- `project_name` (String) Name of the desired project to create load balancer member in. Alternative for `project_id`. One of them should be specified.
- `region_id` (Number) ID of the desired region to create load balancer member in. Alternative for `region_name`. One of them should be specified.
- `region_name` (String) Name of the desired region to create load balancer member in. Alternative for `region_id`. One of them should be specified.
- `subnet_id` (String) ID of the subnet in which real server placed.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `weight` (Number) Value between 0 and 256, default 1.

### Read-Only

- `id` (String) The ID of this resource.
- `last_updated` (String) Datetime when load balancer member was updated at the last time.
- `operating_status` (String) Operating status of this member.

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `update` (String)





## Import

Import is supported using the following syntax:

```shell
# import using <project_id>:<region_id>:<lbmember>:<pool_id> format
terraform import gcore_lbmember.lbmember1 1:6:a775dd94-4e9c-4da7-9f0e-ffc9ae34446b:447d2959-8ae0-4ca0-8d47-9f050a3637d7
```

