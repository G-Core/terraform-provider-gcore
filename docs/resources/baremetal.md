---
page_title: "gcore_baremetal Resource - terraform-provider-gcore"
subcategory: ""
description: |-
  Represent baremetal instance
---

# gcore_baremetal (Resource)

Represent baremetal instance

## Example Usage

##### Prerequisite

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

resource "gcore_network" "network" {
  name       = "my-network"
  type       = "vlan"
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}

resource "gcore_subnet" "subnet" {
  name            = "my-subnet"
  cidr            = "192.168.10.0/24"
  network_id      = gcore_network.network.id

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}

data "gcore_image" "ubuntu" {
  name         = "ubuntu-24.04-x64"
  is_baremetal = true
  region_id    = data.gcore_region.region.id
  project_id   = data.gcore_project.project.id
}

resource "gcore_keypair" "keypair" {
  project_id  = data.gcore_project.project.id
  sshkey_name = "my-keypair"
  public_key  = "ssh-ed25519 ...your public key... gcore@gcore.com"
}

data "gcore_image" "windows" {
  name         = "windows-server-standard-2022-ironic"
  is_baremetal = true
  region_id    = data.gcore_region.region.id
  project_id   = data.gcore_project.project.id
}
```

### Basic example

#### Creating baremetal instance with one public interface

```terraform
resource "gcore_baremetal" "baremetal_with_one_interface" {
  flavor_id     = "bm1-infrastructure-small"
  name          = "my-baremetal"
  keypair_name  = "my-keypair"
  image_id      = data.gcore_image.ubuntu.id

  interface {
    type = "external"
  }

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}
```

#### Creating instance with two interfaces

This example demonstrates how to create a baremetal instance with two network interfaces: one public and one private.

```terraform
resource "gcore_baremetal" "baremetal_with_two_interface" {
  flavor_id     = "bm1-infrastructure-small"
  name          = "my-baremetal"
  keypair_name  = "my-keypair"
  image_id      = data.gcore_image.ubuntu.id

  interface {
    type = "external"
  }

  interface {
    type = "subnet"

    network_id = gcore_network.network.id
    subnet_id = gcore_subnet.subnet.id
  }
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}
```

#### Creating Windows baremetal instance with one public interface

```terraform
resource "gcore_baremetal" "windows_baremetal" {
  flavor_id     = "bm1-infrastructure-small"
  name          = "my-windows-baremetal"
  password      = "my-s3cR3tP@ssw0rd"
  image_id      = data.gcore_image.windows.id

  interface {
    type = "external"
  }

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}
```

### Advanced examples

#### Creating baremetal instance with floating ip

```terraform
resource "gcore_reservedfixedip" "fixed_ip" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
  type       = "subnet"
  network_id = gcore_network.network.id
  subnet_id  = gcore_subnet.subnet.id
}

resource "gcore_floatingip" "floating_ip" {
  project_id       = data.gcore_project.project.id
  region_id        = data.gcore_region.region.id
  fixed_ip_address = gcore_reservedfixedip.fixed_ip.fixed_ip_address
  port_id          = gcore_reservedfixedip.fixed_ip.port_id
}

resource "gcore_baremetal" "baremetal_with_floating_ip" {
  flavor_id     = "bm1-infrastructure-small"
  name          = "my-baremetal"
  keypair_name  = "my-keypair"
  image_id      = data.gcore_image.ubuntu.id

  interface {
    type    = "reserved_fixed_ip"
    port_id = gcore_reservedfixedip.fixed_ip.port_id

    existing_fip_id = gcore_floatingip.floating_ip.id
  }

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}
```

#### Creating instance with a reserved public interface

```terraform
resource "gcore_reservedfixedip" "external_fixed_ip" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
  type       = "external"
}

resource "gcore_baremetal" "baremetal_with_reserved_address" {
  flavor_id     = "bm1-infrastructure-small"
  name          = "my-baremetal"
  keypair_name  = "my-keypair"
  image_id      = data.gcore_image.ubuntu.id

  interface {
    type    = "reserved_fixed_ip"
    port_id = gcore_reservedfixedip.external_fixed_ip.port_id
  }

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}
```


#### Creating Windows baremetal instance with two users

This example shows how to create a Windows baremetal instance with two users. The second user is added by using
the userdata feature to automate the creation process.


```terraform
variable "second_user_userdata" {
 description = "This is a variable of type string"
 type        = string
 default     = <<EOF
<powershell>
# Be sure to set the username and password on these two lines. Of course this is not a good
# security practice to include a password at command line.
$User = "SecondUser"
$Password = ConvertTo-SecureString "s3cR3tP@ssw0rd" -AsPlainText -Force
New-LocalUser $User -Password $Password
Add-LocalGroupMember -Group "Remote Desktop Users" -Member $User
Add-LocalGroupMember -Group "Administrators" -Member $User
</powershell>
EOF
}

resource "gcore_baremetal" "baremetal_windows_with_userdata" {
  flavor_id     = "bm1-infrastructure-small"
  name          = "my-windows-baremetal"
  password      = "my-s3cR3tP@ssw0rd"
  user_data     = base64encode(var.second_user_userdata)
  image_id      = data.gcore_image.windows.id

  interface {
    type = "external"
  }

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `flavor_id` (String)
- `interface` (Block List, Min: 1) (see [below for nested schema](#nestedblock--interface))

### Optional

- `app_config` (Map of String)
- `apptemplate_id` (String)
- `image_id` (String)
- `keypair_name` (String)
- `last_updated` (String)
- `metadata` (Block List, Deprecated) (see [below for nested schema](#nestedblock--metadata))
- `metadata_map` (Map of String)
- `name` (String)
- `name_template` (String)
- `name_templates` (List of String, Deprecated)
- `password` (String)
- `project_id` (Number)
- `project_name` (String)
- `region_id` (Number)
- `region_name` (String)
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `user_data` (String)
- `username` (String)

### Read-Only

- `addresses` (List of Object) (see [below for nested schema](#nestedatt--addresses))
- `flavor` (Map of String)
- `id` (String) The ID of this resource.
- `status` (String)
- `vm_state` (String)

<a id="nestedblock--interface"></a>
### Nested Schema for `interface`

Required:

- `type` (String) Available value is 'subnet', 'any_subnet', 'external', 'reserved_fixed_ip'

Optional:

- `existing_fip_id` (String)
- `fip_source` (String)
- `ip_address` (String)
- `is_parent` (Boolean) If not set will be calculated after creation. Trunk interface always attached first. Can't detach interface if is_parent true. Fields affect only on creation
- `network_id` (String) required if type is 'subnet' or 'any_subnet'
- `order` (Number) Order of attaching interface. Trunk interface always attached first, fields affect only on creation
- `port_id` (String) required if type is  'reserved_fixed_ip'
- `subnet_id` (String) required if type is 'subnet'


<a id="nestedblock--metadata"></a>
### Nested Schema for `metadata`

Required:

- `key` (String)
- `value` (String)


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)


<a id="nestedatt--addresses"></a>
### Nested Schema for `addresses`

Read-Only:

- `net` (List of Object) (see [below for nested schema](#nestedobjatt--addresses--net))

<a id="nestedobjatt--addresses--net"></a>
### Nested Schema for `addresses.net`

Read-Only:

- `addr` (String)
- `type` (String)






## Import

Import is supported using the following syntax:

```shell
# import using <project_id>:<region_id>:<instance_id> format
terraform import gcore_baremetal.instance1 1:6:447d2959-8ae0-4ca0-8d47-9f050a3637d7
```

