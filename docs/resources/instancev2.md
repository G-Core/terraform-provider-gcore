---
page_title: "gcore_instancev2 Resource - terraform-provider-gcore"
subcategory: ""
description: |-
  
Gcore Instance offer a flexible, powerful, and scalable solution for hosting applications and services.
Designed to meet a wide range of computing needs, our instances ensure optimal performance, reliability, and security for
your applications.
---

# gcore_instancev2 (Resource)


Gcore Instance offer a flexible, powerful, and scalable solution for hosting applications and services.
Designed to meet a wide range of computing needs, our instances ensure optimal performance, reliability, and security for
your applications.

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
  type       = "vxlan"
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
  name       = "ubuntu-22.04-x64"
  region_id  = data.gcore_region.region.id
  project_id = data.gcore_project.project.id
}

resource "gcore_volume" "boot_volume" {
  name       = "my-boot-volume"
  type_name  = "ssd_hiiops"
  size       = 5
  image_id   = data.gcore_image.ubuntu.id
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}

resource "gcore_keypair" "my_keypair" {
  project_id  = data.gcore_project.project.id
  sshkey_name = "my-keypair"
  public_key  = "ssh-ed25519 ...your public key... gcore@gcore.com"
}

data "gcore_securitygroup" "default" {
  name       = "default"
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}

data "gcore_image" "windows" {
  name       = "windows-server-2022"
  region_id  = data.gcore_region.region.id
  project_id = data.gcore_project.project.id
}
```

### Basic example

#### Creating instance with one public interface

```terraform
resource "gcore_instancev2" "instance_with_one_interface" {
  flavor_id     = "g1-standard-2-4"
  name          = "my-instance"
  keypair_name  = "my-keypair"

  volume {
    volume_id  = gcore_volume.boot_volume.id
    boot_index = 0
  }

  interface {
    type = "external"
    name = "my-external-interface"
    security_groups = [data.gcore_securitygroup.default.id]
  }

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}
```

#### Creating instance with two interfaces

This example demonstrates how to create an instance with two network interfaces: one public and one private.

```terraform
resource "gcore_instancev2" "instance_with_two_interface" {
  flavor_id     = "g1-standard-2-4"
  name          = "my-instance"
  keypair_name  = "my-keypair"

  volume {
    volume_id  = gcore_volume.boot_volume.id
    boot_index = 0
  }

  interface {
    type = "external"
    name = "my-external-interface"
    security_groups = [data.gcore_securitygroup.default.id]
  }

  interface {
    type = "subnet"
    name = "my-private-interface"
    security_groups = [data.gcore_securitygroup.default.id]

    network_id = gcore_network.network.id
    subnet_id = gcore_subnet.subnet.id
  }
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}
```

#### Creating Windows instance with one public interface

```terraform
resource "gcore_volume" "boot_volume_windows" {
  name       = "my-windows-boot-volume"
  type_name  = "ssd_hiiops"
  size       = 50
  image_id   = data.gcore_image.windows.id
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}

resource "gcore_instancev2" "windows_instance" {
  flavor_id     = "g1w-standard-4-8"
  name          = "my-windows-instance"
  password      = "my-s3cR3tP@ssw0rd"

  volume {
    volume_id  = gcore_volume.boot_volume_windows.id
    boot_index = 0
  }

  interface {
    type = "external"
    name = "my-external-interface"
    security_groups = [data.gcore_securitygroup.default.id]
  }

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}
```

### Advanced examples


#### Creating instance with a dual-stack public interface

This example demonstrates how to create an instance with a dual-stack public interface.
The instance has both an IPv4 and an IPv6 address.

```terraform
resource "gcore_instancev2" "instance_with_dualstack" {
  flavor_id     = "g1-standard-2-4"
  name          = "my-instance"
  keypair_name  = "my-keypair"

  volume {
    volume_id  = gcore_volume.boot_volume.id
    boot_index = 0
  }

  interface {
    type      = "external"
    ip_family = "dual"
    name      = "my-external-interface"
    security_groups = [data.gcore_securitygroup.default.id]
  }

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}

output "addresses" {
  value = gcore_instancev2.instance_with_dualstack.addresses
}
```

#### Creating instance with floating ip

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

resource "gcore_instancev2" "instance_with_floating_ip" {
  flavor_id     = "g1-standard-2-4"
  name          = "my-instance"
  keypair_name  = "my-keypair"

  volume {
    volume_id  = gcore_volume.boot_volume.id
    boot_index = 0
  }

  interface {
    type    = "reserved_fixed_ip"
    name    = "my-floating-ip-interface"
    port_id = gcore_reservedfixedip.fixed_ip.port_id

    existing_fip_id = gcore_floatingip.floating_ip.id
    security_groups = [gcore_securitygroup.default.id]
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

resource "gcore_instancev2" "instance_with_reserved_address" {
  flavor_id     = "g1-standard-2-4"
  name          = "my-instance"
  keypair_name  = "my-keypair"

  volume {
    volume_id  = gcore_volume.boot_volume.id
    boot_index = 0
  }

  interface {
    type    = "reserved_fixed_ip"
    name    = "my-reserved-public-interface"
    port_id = gcore_reservedfixedip.external_fixed_ip.port_id
    security_groups = [data.gcore_securitygroup.default.id]
  }

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}
```


#### Creating instance with custom security group

This example demonstrates how to create an instance with a custom security group. The security group allows all
incoming traffic on ports 22, 80, and 443. Outgoing traffic is allowed on all ports, except port 25 for security reasons.


```terraform
resource "gcore_securitygroup" "web_server_security_group" {
  name       = "web server only"
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  security_group_rules {
    direction      = "egress"
    ethertype      = "IPv4"
    protocol       = "tcp"
    port_range_min = 1
    port_range_max = 24
  }

  security_group_rules {
    direction      = "egress"
    ethertype      = "IPv4"
    protocol       = "tcp"
    port_range_min = 26
    port_range_max = 65535
  }

  security_group_rules {
    direction      = "ingress"
    ethertype      = "IPv4"
    protocol       = "tcp"
    port_range_min = 22
    port_range_max = 22
  }

  security_group_rules {
    direction      = "ingress"
    ethertype      = "IPv4"
    protocol       = "tcp"
    port_range_min = 80
    port_range_max = 80
  }

  security_group_rules {
    direction      = "ingress"
    ethertype      = "IPv4"
    protocol       = "tcp"
    port_range_min = 443
    port_range_max = 443
  }

}

resource "gcore_instancev2" "instance_with_custom_security_group" {
  flavor_id     = "g1-standard-2-4"
  name          = "my-instance"
  keypair_name  = "my-keypair"

  volume {
    volume_id  = gcore_volume.boot_volume.id
    boot_index = 0
  }

  interface {
    type = "external"
    name = "my-external-interface"

    security_groups = [gcore_securitygroup.web_server_security_group.id]
  }

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}
```


#### Creating Windows instance with two users

This example shows how to create a Windows instance with two users. The second user is added by using
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

resource "gcore_volume" "boot_volume_windows_with_userdata" {
  name       = "my-windows-boot-volume"
  type_name  = "ssd_hiiops"
  size       = 50
  image_id   = data.gcore_image.windows.id
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}

resource "gcore_instancev2" "instance_windows_with_userdata" {
  flavor_id     = "g1w-standard-4-8"
  name          = "my-windows-instance"
  password      = "my-s3cR3tP@ssw0rd"
  user_data     = base64encode(var.second_user_userdata)

  volume {
    volume_id  = gcore_volume.boot_volume_windows_with_userdata.id
    boot_index = 0
  }

  interface {
    type = "external"
    name = "my-external-interface"
    security_groups = [data.gcore_securitygroup.default.id]
  }

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `flavor_id` (String) Flavor ID
- `interface` (Block Set, Min: 1) List of interfaces for the instance. You can detach the interface from the instance by removing the
interface from the instance resource and attach the interface by adding the interface resource
inside an instance resource. (see [below for nested schema](#nestedblock--interface))
- `volume` (Block Set, Min: 1) List of volumes for the instance. You can detach the volume from the instance by removing the
volume from the instance resource. You cannot detach the boot volume. You can attach a data volume
by adding the volume resource inside an instance resource. (see [below for nested schema](#nestedblock--volume))

### Optional

- `allow_app_ports` (Boolean) If true, application ports will be allowed in the security group for instances created
				from the marketplace application template
- `configuration` (Block List) Parameters for the application template from the marketplace (see [below for nested schema](#nestedblock--configuration))
- `keypair_name` (String) Name of the keypair to use for the instance
- `metadata_map` (Map of String) Create one or more metadata items for the instance
- `name` (String) Name of the instance.
- `name_template` (String) Instance name template. You can use forms 'ip_octets', 'two_ip_octets', 'one_ip_octet'
- `password` (String, Sensitive) For Linux instances, 'username' and 'password' are used to create a new user.
When only 'password' is provided, it is set as the password for the default user of the image. 'user_data' is ignored
when 'password' is specified. For Windows instances, 'username' cannot be specified. Use the 'password' field to set
the password for the 'Admin' user on Windows. Use the 'user_data' field to provide a script to create new users
on Windows. The password of the Admin user cannot be updated via 'user_data'
- `project_id` (Number) Project ID, only one of project_id or project_name should be set
- `project_name` (String) Project name, only one of project_id or project_name should be set
- `region_id` (Number) Region ID, only one of region_id or region_name should be set
- `region_name` (String) Region name, only one of region_id or region_name should be set
- `server_group` (String) ID of the server group to use for the instance
- `user_data` (String) String in base64 format. For Linux instances, 'user_data' is ignored when 'password' field is provided.
For Windows instances, Admin user password is set by 'password' field and cannot be updated via 'user_data'
- `username` (String) For Linux instances, 'username' and 'password' are used to create a new user. For Windows
instances, 'username' cannot be specified. Use 'password' field to set the password for the 'Admin' user on Windows.
- `vm_state` (String) Current vm state, use stopped to stop vm and active to start

### Read-Only

- `addresses` (List of Object) List of instance addresses (see [below for nested schema](#nestedatt--addresses))
- `flavor` (Map of String) Flavor details, RAM, vCPU, etc.
- `id` (String) The ID of this resource.
- `last_updated` (String)
- `status` (String) Status of the instance

<a id="nestedblock--interface"></a>
### Nested Schema for `interface`

Required:

- `name` (String) Name of interface, should be unique for the instance
- `security_groups` (List of String) list of security group IDs, they will be attached to exact interface

Optional:

- `existing_fip_id` (String) The id of the existing floating IP that will be attached to the interface
- `ip_address` (String) IP address for the interface.
- `ip_family` (String) IP family for the interface, available values are 'dual', 'ipv4' and 'ipv6'
- `network_id` (String) required if type is 'subnet' or 'any_subnet'
- `order` (Number) Order of attaching interface
- `port_id` (String) required if type is  'reserved_fixed_ip'
- `subnet_id` (String) required if type is 'subnet'
- `type` (String) Available value is 'subnet', 'any_subnet', 'external', 'reserved_fixed_ip'


<a id="nestedblock--volume"></a>
### Nested Schema for `volume`

Required:

- `volume_id` (String)

Optional:

- `boot_index` (Number) If boot_index==0 volumes can not detached

Read-Only:

- `attachment_tag` (String) Tag for the volume attachment
- `delete_on_termination` (Boolean) Delete volume on termination
- `id` (String)
- `image_id` (String) Image ID for the volume
- `name` (String) Name of the volume
- `size` (Number) Size of the volume in GiB
- `type_name` (String) Volume type name


<a id="nestedblock--configuration"></a>
### Nested Schema for `configuration`

Required:

- `key` (String)
- `value` (String)


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
terraform import gcore_instance.instance1 1:6:447d2959-8ae0-4ca0-8d47-9f050a3637d7
```

