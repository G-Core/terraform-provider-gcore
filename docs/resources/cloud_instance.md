---
page_title: "gcore_cloud_instance Resource - Gcore"
subcategory: ""
description: |-
  Instances are cloud virtual machines with configurable CPU, memory, storage, and networking, supporting various operating systems and workloads.
---

# gcore_cloud_instance (Resource)

Instances are cloud virtual machines with configurable CPU, memory, storage, and networking, supporting various operating systems and workloads.

## Example Usage

### Instance with one public interface

Create a basic instance with a single external IPv4 interface.

```terraform
# Create an SSH key for instance access
resource "gcore_cloud_ssh_key" "my_key" {
  project_id = 1
  name       = "my-keypair"
  public_key = "ssh-ed25519 ...your public key... user@example.com"
}

# Create a boot volume from an image
resource "gcore_cloud_volume" "boot_volume" {
  project_id = 1
  region_id  = 1
  name       = "my-boot-volume"
  source     = "image"
  image_id   = "6dc4e521-0c72-462f-b2d4-306bcf15e227"
  size       = 20
  type_name  = "ssd_hiiops"
}

# Create an instance with a single external interface
resource "gcore_cloud_instance" "instance_with_one_interface" {
  project_id   = 1
  region_id    = 1
  flavor       = "g1-standard-2-4"
  name         = "my-instance"
  ssh_key_name = gcore_cloud_ssh_key.my_key.name

  volumes = [{ volume_id = gcore_cloud_volume.boot_volume.id }]

  interfaces = [{
    type      = "external"
    ip_family = "ipv4"
  }]
}
```

### Instance with two interfaces

Create an instance with two network interfaces: one public and one private.

```terraform
# Create a private network and subnet
resource "gcore_cloud_network" "network" {
  project_id = 1
  region_id  = 1
  name       = "my-network"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "subnet" {
  project_id      = 1
  region_id       = 1
  name            = "my-subnet"
  cidr            = "192.168.10.0/24"
  network_id      = gcore_cloud_network.network.id
  dns_nameservers = ["8.8.4.4", "1.1.1.1"]
}

# Create an SSH key for instance access
resource "gcore_cloud_ssh_key" "my_key" {
  project_id = 1
  name       = "my-keypair"
  public_key = "ssh-ed25519 ...your public key... user@example.com"
}

# Create a boot volume from an image
resource "gcore_cloud_volume" "boot_volume" {
  project_id = 1
  region_id  = 1
  name       = "my-boot-volume"
  source     = "image"
  image_id   = "6dc4e521-0c72-462f-b2d4-306bcf15e227"
  size       = 20
  type_name  = "ssd_hiiops"
}

# Create an instance with two interfaces: one public, one private
resource "gcore_cloud_instance" "instance_with_two_interfaces" {
  project_id   = 1
  region_id    = 1
  flavor       = "g1-standard-2-4"
  name         = "my-instance"
  ssh_key_name = gcore_cloud_ssh_key.my_key.name

  volumes = [{ volume_id = gcore_cloud_volume.boot_volume.id }]

  interfaces = [
    {
      type      = "external"
      ip_family = "ipv4"
    },
    {
      type       = "subnet"
      network_id = gcore_cloud_network.network.id
      subnet_id  = gcore_cloud_network_subnet.subnet.id
    },
  ]
}
```

### Windows instance

Create a Windows instance with a public interface.

```terraform
resource "gcore_cloud_volume" "boot_volume_windows" {
  project_id = 1
  region_id  = 1
  name       = "my-windows-boot-volume"
  source     = "image"
  image_id   = "a2c1681c-94e0-4aab-8fa3-09a8e662d4c0"
  size       = 50
  type_name  = "ssd_hiiops"
}

resource "gcore_cloud_instance" "windows_instance" {
  project_id = 1
  region_id  = 1
  flavor     = "g1w-standard-4-8"
  name       = "my-windows-instance"
  password   = "my-s3cR3tP@ssw0rd"

  volumes = [{ volume_id = gcore_cloud_volume.boot_volume_windows.id }]

  interfaces = [{
    type      = "external"
    ip_family = "ipv4"
  }]
}
```

### Dual-stack public interface

Create an instance with both IPv4 and IPv6 addresses on a single interface.

```terraform
# Create an SSH key for instance access
resource "gcore_cloud_ssh_key" "my_key" {
  project_id = 1
  name       = "my-keypair"
  public_key = "ssh-ed25519 ...your public key... user@example.com"
}

# Create a boot volume from an image
resource "gcore_cloud_volume" "boot_volume" {
  project_id = 1
  region_id  = 1
  name       = "my-boot-volume"
  source     = "image"
  image_id   = "6dc4e521-0c72-462f-b2d4-306bcf15e227"
  size       = 20
  type_name  = "ssd_hiiops"
}

# Create an instance with dual-stack (IPv4 + IPv6) public interface
resource "gcore_cloud_instance" "instance_with_dualstack" {
  project_id   = 1
  region_id    = 1
  flavor       = "g1-standard-2-4"
  name         = "my-instance"
  ssh_key_name = gcore_cloud_ssh_key.my_key.name

  volumes = [{ volume_id = gcore_cloud_volume.boot_volume.id }]

  interfaces = [{
    type      = "external"
    ip_family = "dual"
  }]
}

output "addresses" {
  value = gcore_cloud_instance.instance_with_dualstack.addresses
}
```

### Instance with floating IP

Create an instance and attach a floating IP address for external access.

```terraform
# Create a private network and subnet
resource "gcore_cloud_network" "network" {
  project_id = 1
  region_id  = 1
  name       = "my-network"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "subnet" {
  project_id      = 1
  region_id       = 1
  name            = "my-subnet"
  cidr            = "192.168.10.0/24"
  network_id      = gcore_cloud_network.network.id
  dns_nameservers = ["8.8.4.4", "1.1.1.1"]
}

# Create an SSH key for instance access
resource "gcore_cloud_ssh_key" "my_key" {
  project_id = 1
  name       = "my-keypair"
  public_key = "ssh-ed25519 ...your public key... user@example.com"
}

# Create a boot volume from an image
resource "gcore_cloud_volume" "boot_volume" {
  project_id = 1
  region_id  = 1
  name       = "my-boot-volume"
  source     = "image"
  image_id   = "6dc4e521-0c72-462f-b2d4-306bcf15e227"
  size       = 20
  type_name  = "ssd_hiiops"
}

# Reserve a fixed IP on the private subnet
resource "gcore_cloud_reserved_fixed_ip" "fixed_ip" {
  project_id = 1
  region_id  = 1
  type       = "subnet"
  network_id = gcore_cloud_network.network.id
  subnet_id  = gcore_cloud_network_subnet.subnet.id
}

# Create a floating IP and associate it with the fixed IP
resource "gcore_cloud_floating_ip" "floating_ip" {
  project_id       = 1
  region_id        = 1
  fixed_ip_address = gcore_cloud_reserved_fixed_ip.fixed_ip.fixed_ip_address
  port_id          = gcore_cloud_reserved_fixed_ip.fixed_ip.port_id
}

# Create an instance with floating IP for external access
resource "gcore_cloud_instance" "instance_with_floating_ip" {
  project_id   = 1
  region_id    = 1
  flavor       = "g1-standard-2-4"
  name         = "my-instance"
  ssh_key_name = gcore_cloud_ssh_key.my_key.name

  volumes = [{ volume_id = gcore_cloud_volume.boot_volume.id }]

  interfaces = [{
    type    = "reserved_fixed_ip"
    port_id = gcore_cloud_reserved_fixed_ip.fixed_ip.port_id

    floating_ip = {
      source               = "existing"
      existing_floating_id = gcore_cloud_floating_ip.floating_ip.id
    }
  }]
}
```

### Instance with reserved public IP

Create an instance using a pre-allocated reserved fixed IP address.

```terraform
# Create an SSH key for instance access
resource "gcore_cloud_ssh_key" "my_key" {
  project_id = 1
  name       = "my-keypair"
  public_key = "ssh-ed25519 ...your public key... user@example.com"
}

# Create a boot volume from an image
resource "gcore_cloud_volume" "boot_volume" {
  project_id = 1
  region_id  = 1
  name       = "my-boot-volume"
  source     = "image"
  image_id   = "6dc4e521-0c72-462f-b2d4-306bcf15e227"
  size       = 20
  type_name  = "ssd_hiiops"
}

# Reserve a public IP address
resource "gcore_cloud_reserved_fixed_ip" "external_fixed_ip" {
  project_id = 1
  region_id  = 1
  type       = "external"
}

# Create an instance using the reserved public IP
resource "gcore_cloud_instance" "instance_with_reserved_address" {
  project_id   = 1
  region_id    = 1
  flavor       = "g1-standard-2-4"
  name         = "my-instance"
  ssh_key_name = gcore_cloud_ssh_key.my_key.name

  volumes = [{ volume_id = gcore_cloud_volume.boot_volume.id }]

  interfaces = [{
    type    = "reserved_fixed_ip"
    port_id = gcore_cloud_reserved_fixed_ip.external_fixed_ip.port_id
  }]
}
```

### Instance with custom security group

Create an instance with a custom security group allowing SSH, HTTP, and HTTPS inbound traffic.

```terraform
# Create an SSH key for instance access
resource "gcore_cloud_ssh_key" "my_key" {
  project_id = 1
  name       = "my-keypair"
  public_key = "ssh-ed25519 ...your public key... user@example.com"
}

# Create a boot volume from an image
resource "gcore_cloud_volume" "boot_volume" {
  project_id = 1
  region_id  = 1
  name       = "my-boot-volume"
  source     = "image"
  image_id   = "6dc4e521-0c72-462f-b2d4-306bcf15e227"
  size       = 20
  type_name  = "ssd_hiiops"
}

# Create a security group, then add rules as separate resources
resource "gcore_cloud_security_group" "web_server" {
  project_id = 1
  region_id  = 1
  name       = "web-server-only"
}

resource "gcore_cloud_security_group_rule" "egress_low" {
  project_id     = 1
  region_id      = 1
  group_id       = gcore_cloud_security_group.web_server.id
  direction      = "egress"
  ethertype      = "IPv4"
  protocol       = "tcp"
  port_range_min = 1
  port_range_max = 24
  description    = "Allow outgoing TCP except SMTP"
}

resource "gcore_cloud_security_group_rule" "egress_high" {
  project_id     = 1
  region_id      = 1
  group_id       = gcore_cloud_security_group.web_server.id
  direction      = "egress"
  ethertype      = "IPv4"
  protocol       = "tcp"
  port_range_min = 26
  port_range_max = 65535
  description    = "Allow outgoing TCP except SMTP"
}

resource "gcore_cloud_security_group_rule" "ssh" {
  project_id     = 1
  region_id      = 1
  group_id       = gcore_cloud_security_group.web_server.id
  direction      = "ingress"
  ethertype      = "IPv4"
  protocol       = "tcp"
  port_range_min = 22
  port_range_max = 22
  description    = "Allow SSH"
}

resource "gcore_cloud_security_group_rule" "http" {
  project_id     = 1
  region_id      = 1
  group_id       = gcore_cloud_security_group.web_server.id
  direction      = "ingress"
  ethertype      = "IPv4"
  protocol       = "tcp"
  port_range_min = 80
  port_range_max = 80
  description    = "Allow HTTP"
}

resource "gcore_cloud_security_group_rule" "https" {
  project_id     = 1
  region_id      = 1
  group_id       = gcore_cloud_security_group.web_server.id
  direction      = "ingress"
  ethertype      = "IPv4"
  protocol       = "tcp"
  port_range_min = 443
  port_range_max = 443
  description    = "Allow HTTPS"
}

# Create an instance with the custom security group
resource "gcore_cloud_instance" "instance_with_custom_sg" {
  project_id   = 1
  region_id    = 1
  flavor       = "g1-standard-2-4"
  name         = "my-instance"
  ssh_key_name = gcore_cloud_ssh_key.my_key.name

  volumes = [{ volume_id = gcore_cloud_volume.boot_volume.id }]

  interfaces = [{
    type      = "external"
    ip_family = "ipv4"
    security_groups = [{
      id = gcore_cloud_security_group.web_server.id
    }]
  }]

  security_groups = [{
    id = gcore_cloud_security_group.web_server.id
  }]
}
```

### Windows instance with two users

Create a Windows instance and use userdata to add a second user during provisioning.

```terraform
variable "second_user_userdata" {
  description = "PowerShell script to create a second Windows user"
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

resource "gcore_cloud_volume" "boot_volume_windows_userdata" {
  project_id = 1
  region_id  = 1
  name       = "my-windows-boot-volume"
  source     = "image"
  image_id   = "a2c1681c-94e0-4aab-8fa3-09a8e662d4c0"
  size       = 50
  type_name  = "ssd_hiiops"
}

resource "gcore_cloud_instance" "windows_with_userdata" {
  project_id = 1
  region_id  = 1
  flavor     = "g1w-standard-4-8"
  name       = "my-windows-instance"
  password   = "my-s3cR3tP@ssw0rd"
  user_data  = base64encode(var.second_user_userdata)

  volumes = [{ volume_id = gcore_cloud_volume.boot_volume_windows_userdata.id }]

  interfaces = [{
    type      = "external"
    ip_family = "ipv4"
  }]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `flavor` (String) The flavor of the instance.
- `interfaces` (Attributes List) A list of network interfaces for the instance. You can create one or more interfaces - private, public, or both. (see [below for nested schema](#nestedatt--interfaces))
- `volumes` (Attributes List) List of existing volumes to attach to the instance. Create volumes separately using gcore_cloud_volume resource. (see [below for nested schema](#nestedatt--volumes))

### Optional

> **NOTE**: [Write-only arguments](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments) are supported in Terraform 1.11 and later.

- `allow_app_ports` (Boolean) Set to `true` if creating the instance from an `apptemplate`. This allows application ports in the security group for instances created from a marketplace application template.
- `configuration` (Map of String) Parameters for the application template if creating the instance from an `apptemplate`.
- `name` (String) Instance name.
- `name_template` (String) If you want the instance name to be automatically generated based on IP addresses, you can provide a name template instead of specifying the name manually. The template should include a placeholder that will be replaced during provisioning. Supported placeholders are: `{ip_octets}` (last 3 octets of the IP), `{two_ip_octets}`, and `{one_ip_octet}`.
- `password_wo` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) For Linux instances, 'username' and 'password' are used to create a new user. When only 'password' is provided, it is set as the password for the default user of the image. For Windows instances, 'username' cannot be specified. Use the 'password' field to set the password for the 'Admin' user on Windows. Use the 'user_data' field to provide a script to create new users on Windows. The password of the Admin user cannot be updated via 'user_data'.
- `password_wo_version` (Number) Instance password write-only version. Used to trigger updates of the write-only password field.
- `project_id` (Number) Project ID. If not specified, uses GCORE_CLOUD_PROJECT_ID environment variable.
- `region_id` (Number) Region ID. If not specified, uses GCORE_CLOUD_REGION_ID environment variable.
- `security_groups` (Attributes List) Specifies security group UUIDs to be applied to all instance network interfaces. (see [below for nested schema](#nestedatt--security_groups))
- `servergroup_id` (String) Placement group ID for instance placement policy.

Supported group types:
- `anti-affinity`: Ensures instances are placed on different hosts for high availability.
- `affinity`: Places instances on the same host for low-latency communication.
- `soft-anti-affinity`: Tries to place instances on different hosts but allows sharing if needed.
- `ssh_key_name` (String) Specifies the name of the SSH keypair, created via the
[/v1/`ssh_keys` endpoint](/docs/api-reference/cloud/ssh-keys/add-or-generate-ssh-key).
- `tags` (Map of String) Key-value tags to associate with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Both tag keys and values have a maximum length of 255 characters. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.
- `user_data` (String) String in base64 format. For Linux instances, 'user_data' is ignored when 'password' field is provided. For Windows instances, Admin user password is set by 'password' field and cannot be updated via 'user_data'. Examples of the `user_data`: https://cloudinit.readthedocs.io/en/latest/topics/examples.html
- `username` (String) For Linux instances, 'username' and 'password' are used to create a new user. For Windows instances, 'username' cannot be specified. Use 'password' field to set the password for the 'Admin' user on Windows.
- `vm_state` (String) Virtual machine state. Set to 'active' to start the instance or 'stopped' to stop it.
Available values: "active", "stopped".

### Read-Only

- `addresses` (Map of List of Object) Map of `network_name` to list of addresses in that network
- `blackhole_ports` (Attributes List) IP addresses of the instances that are blackholed by DDoS mitigation system (see [below for nested schema](#nestedatt--blackhole_ports))
- `created_at` (String) Datetime when instance was created
- `creator_task_id` (String) Task that created this entity
- `ddos_profile` (Attributes) Advanced DDoS protection profile. It is always `null` if query parameter `with_ddos=true` is not set. (see [below for nested schema](#nestedatt--ddos_profile))
- `fixed_ip_assignments` (Attributes List) Fixed IP assigned to instance (see [below for nested schema](#nestedatt--fixed_ip_assignments))
- `id` (String) The ID of this resource.
- `instance_description` (String) Instance description
- `instance_isolation` (Attributes) Instance isolation information (see [below for nested schema](#nestedatt--instance_isolation))
- `region` (String) Region name
- `status` (String) Instance status
Available values: "ACTIVE", "BUILD", "DELETED", "ERROR", "HARD_REBOOT", "MIGRATING", "PASSWORD", "PAUSED", "REBOOT", "REBUILD", "RESCUE", "RESIZE", "REVERT_RESIZE", "SHELVED", "SHELVED_OFFLOADED", "SHUTOFF", "SOFT_DELETED", "SUSPENDED", "UNKNOWN", "VERIFY_RESIZE".
- `task_state` (String) Task state

<a id="nestedatt--interfaces"></a>
### Nested Schema for `interfaces`

Required:

- `type` (String) A public IP address will be assigned to the instance.
Available values: "external", "subnet", "any_subnet", "reserved_fixed_ip".

Optional:

- `floating_ip` (Attributes) Allows the instance to have a public IP that can be reached from the internet. (see [below for nested schema](#nestedatt--interfaces--floating_ip))
- `interface_name` (String) Interface name. Defaults to `null` and is returned as `null` in the API response if not set.
- `ip_address` (String) IP address assigned to this interface. Can be specified for subnet type, computed for other types.
- `ip_family` (String) Specify `ipv4`, `ipv6`, or `dual` to enable both.
Available values: "dual", "ipv4", "ipv6".
- `network_id` (String) The network where the instance will be connected.
- `port_id` (String) Port ID for the interface. Required for reserved_fixed_ip type, computed for other types.
- `security_groups` (Attributes List) Specifies security group UUIDs to be applied to the instance network interface. (see [below for nested schema](#nestedatt--interfaces--security_groups))
- `subnet_id` (String) The instance will get an IP address from this subnet.

<a id="nestedatt--interfaces--floating_ip"></a>
### Nested Schema for `interfaces.floating_ip`

Required:

- `source` (String) A new floating IP will be created and attached to the instance. A floating IP is a public IP that makes the instance accessible from the internet, even if it only has a private IP. It works like SNAT, allowing outgoing and incoming traffic.
Available values: "new", "existing".

Optional:

- `existing_floating_id` (String) An existing available floating IP id must be specified if the source is set to `existing`


<a id="nestedatt--interfaces--security_groups"></a>
### Nested Schema for `interfaces.security_groups`

Required:

- `id` (String) Resource ID



<a id="nestedatt--volumes"></a>
### Nested Schema for `volumes`

Required:

- `volume_id` (String) ID of an existing volume to attach to the instance.

Optional:

- `attachment_tag` (String) Block device attachment tag. Used to identify the device in the guest OS (e.g., 'vdb', 'data-disk'). Not exposed in user-visible tags.
- `boot_index` (Number) Boot device index (creation-only). 0 = primary boot, positive = secondary bootable, negative = not bootable. Cannot be changed after instance creation.


<a id="nestedatt--security_groups"></a>
### Nested Schema for `security_groups`

Required:

- `id` (String) Resource ID


<a id="nestedatt--blackhole_ports"></a>
### Nested Schema for `blackhole_ports`

Read-Only:

- `alarm_end` (String) A date-time string giving the time that the alarm ended. If not yet ended, time will be given as 0001-01-01T00:00:00Z
- `alarm_start` (String) A date-time string giving the time that the alarm started
- `alarm_state` (String) Current state of alarm
Available values: "ACK_REQ", "ALARM", "ARCHIVED", "CLEAR", "CLEARING", "CLEARING_FAIL", "END_GRACE", "END_WAIT", "MANUAL_CLEAR", "MANUAL_CLEARING", "MANUAL_CLEARING_FAIL", "MANUAL_MITIGATING", "MANUAL_STARTING", "MANUAL_STARTING_FAIL", "MITIGATING", "STARTING", "STARTING_FAIL", "START_WAIT", "ack_req", "alarm", "archived", "clear", "clearing", "clearing_fail", "end_grace", "end_wait", "manual_clear", "manual_clearing", "manual_clearing_fail", "manual_mitigating", "manual_starting", "manual_starting_fail", "mitigating", "start_wait", "starting", "starting_fail".
- `alert_duration` (String) Total alert duration
- `destination_ip` (String) Notification destination IP address
- `id` (Number)


<a id="nestedatt--ddos_profile"></a>
### Nested Schema for `ddos_profile`

Read-Only:

- `fields` (Attributes List) List of configured field values for the protection profile (see [below for nested schema](#nestedatt--ddos_profile--fields))
- `id` (Number) Unique identifier for the DDoS protection profile
- `options` (Attributes) Configuration options controlling profile activation and BGP routing (see [below for nested schema](#nestedatt--ddos_profile--options))
- `profile_template` (Attributes) Complete template configuration data used for this profile (see [below for nested schema](#nestedatt--ddos_profile--profile_template))
- `profile_template_description` (String) Detailed description of the protection template used for this profile
- `protocols` (Attributes List) List of network protocols and ports configured for protection (see [below for nested schema](#nestedatt--ddos_profile--protocols))
- `site` (String) Geographic site identifier where the protection is deployed
- `status` (Attributes) Current operational status and any error information for the profile (see [below for nested schema](#nestedatt--ddos_profile--status))

<a id="nestedatt--ddos_profile--fields"></a>
### Nested Schema for `ddos_profile.fields`

Read-Only:

- `base_field` (Number) ID of DDoS profile field
- `default` (String) Predefined default value for the field if not specified
- `description` (String) Detailed description explaining the field's purpose and usage guidelines
- `field_name` (String) Name of DDoS profile field
- `field_type` (String) Data type classification of the field (e.g., string, integer, array)
- `field_value` (String) Complex value. Only one of 'value' or 'field_value' must be specified.
- `id` (Number) Unique identifier for the DDoS protection field
- `name` (String) Human-readable name of the protection field
- `required` (Boolean) Indicates whether this field must be provided when creating a protection profile
- `validation_schema` (String) JSON schema defining validation rules and constraints for the field value
- `value` (String) Basic type value. Only one of 'value' or 'field_value' must be specified.


<a id="nestedatt--ddos_profile--options"></a>
### Nested Schema for `ddos_profile.options`

Read-Only:

- `active` (Boolean) Controls whether the DDoS protection profile is enabled and actively protecting the resource
- `bgp` (Boolean) Enables Border Gateway Protocol (BGP) routing for DDoS protection traffic


<a id="nestedatt--ddos_profile--profile_template"></a>
### Nested Schema for `ddos_profile.profile_template`

Read-Only:

- `description` (String) Detailed description explaining the template's purpose and use cases
- `fields` (Attributes List) List of configurable fields that define the template's protection parameters (see [below for nested schema](#nestedatt--ddos_profile--profile_template--fields))
- `id` (Number) Unique identifier for the DDoS protection template
- `name` (String) Human-readable name of the protection template

<a id="nestedatt--ddos_profile--profile_template--fields"></a>
### Nested Schema for `ddos_profile.profile_template.fields`

Read-Only:

- `default` (String) Predefined default value for the field if not specified
- `description` (String) Detailed description explaining the field's purpose and usage guidelines
- `field_type` (String) Data type classification of the field (e.g., string, integer, array)
- `id` (Number) Unique identifier for the DDoS protection field
- `name` (String) Human-readable name of the protection field
- `required` (Boolean) Indicates whether this field must be provided when creating a protection profile
- `validation_schema` (String) JSON schema defining validation rules and constraints for the field value



<a id="nestedatt--ddos_profile--protocols"></a>
### Nested Schema for `ddos_profile.protocols`

Read-Only:

- `port` (String) Network port number for which protocols are configured
- `protocols` (List of String) List of network protocols enabled on the specified port


<a id="nestedatt--ddos_profile--status"></a>
### Nested Schema for `ddos_profile.status`

Read-Only:

- `error_description` (String) Detailed error message describing any issues with the profile operation
- `status` (String) Current operational status of the DDoS protection profile



<a id="nestedatt--fixed_ip_assignments"></a>
### Nested Schema for `fixed_ip_assignments`

Read-Only:

- `external` (Boolean) Is network external
- `ip_address` (String) Ip address
- `subnet_id` (String) Interface subnet id


<a id="nestedatt--instance_isolation"></a>
### Nested Schema for `instance_isolation`

Read-Only:

- `reason` (String) The reason of instance isolation if it is isolated from external internet.


## Import

Import is supported using the following syntax:

```shell
$ terraform import gcore_cloud_instance.example '<project_id>/<region_id>/<instance_id>'
```

