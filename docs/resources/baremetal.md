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

- `flavor_id` (String) The ID of the flavor (type of server configuration). This field is required. Example: 'bm1-hf-medium-4x1nic'
- `interface` (Block List, Min: 1) (see [below for nested schema](#nestedblock--interface))

### Optional

- `app_config` (Map of String) Parameters for the application template from the marketplace. This could include parameters required for app setup. Example: {'shadowsocks_method': 'chacha20-ietf-poly1305', 'shadowsocks_password': '123'}
- `apptemplate_id` (String) The ID of the application template to use. Provide either 'apptemplate_id' or 'image_id', but not both
- `image_id` (String) The ID of the image to use. The image will be used to provision the bare metal server. Provide either 'image_id' or 'apptemplate_id', but not both
- `keypair_name` (String) The name of the SSH keypair to use for the baremetal
- `last_updated` (String) The date and time when the baremetal server was last updated
- `metadata` (Block List, Deprecated) (see [below for nested schema](#nestedblock--metadata))
- `metadata_map` (Map of String) A map of metadata items. Key-value pairs for instance metadata. Example: {'environment': 'production', 'owner': 'user'}
- `name` (String) The name of the baremetal server. If not provided, it will be generated automatically. Example: 'bm-server-01'
- `name_template` (String) The template used to generate server names. You can use forms 'ip_octets', 'two_ip_octets', 'one_ip_octet'. Example: 'server-${ip_octets}'
- `name_templates` (List of String, Deprecated) Deprecated. List of baremetal names which will be changed by template
- `password` (String) The password for accessing the baremetal server. This parameter is used to set a password for the 'Admin' user on a Windows instance, a default user or a new user on a Linux instance
- `project_id` (Number) Project ID, only one of project_id or project_name should be set
- `project_name` (String) Project name, only one of project_id or project_name should be set
- `region_id` (Number) Region ID, only one of region_id or region_name should be set
- `region_name` (String) Region name, only one of region_id or region_name should be set
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `user_data` (String) User data string in base64 format. This is passed to the instance at launch. For Linux instances, 'user_data' is ignored when 'password' field is provided. For Windows instances, Admin user password is set by 'password' field and cannot be updated via 'user_data'
- `username` (String) A name of a new user in the Linux instance. It may be passed with a 'password' parameter

### Read-Only

- `addresses` (List of Object) (see [below for nested schema](#nestedatt--addresses))
- `flavor` (Map of String) Details about the flavor (server configuration) including RAM, vCPU, etc.
- `id` (String) The ID of this resource.
- `status` (String) The current status of the baremetal server.
- `vm_state` (String) The state of the virtual machine

<a id="nestedblock--interface"></a>
### Nested Schema for `interface`

Required:

- `type` (String) The type of the network interface. Available value is 'subnet', 'any_subnet', 'external', 'reserved_fixed_ip'

Optional:

- `existing_fip_id` (String) The ID of the existing floating IP that will be attached to the interface
- `fip_source` (String) The source of floating IP. Can be 'new' or 'existing'
- `ip_address` (String) The IP address for the interface
- `is_parent` (Boolean) Indicates whether this interface is the parent. If not set will be calculated after creation. Trunk interface always attached first. Can't detach interface if is_parent true. Fields affect only on creation
- `network_id` (String) The network ID to attach the interface to. Required if type is 'subnet' or 'any_subnet'
- `order` (Number) Order of attaching interface. Trunk (parent) interface always attached first, fields affect only on creation
- `port_id` (String) The port ID for reserved fixed IP. Required if type is  'reserved_fixed_ip'
- `subnet_id` (String) The subnet ID to attach the interface to. Required if type is 'subnet'


<a id="nestedblock--metadata"></a>
### Nested Schema for `metadata`

Required:

- `key` (String) Metadata key
- `value` (String) Metadata value


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

