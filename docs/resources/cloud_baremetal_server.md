---
page_title: "gcore_cloud_baremetal_server Resource - Gcore"
subcategory: ""
description: |-
  Bare metal servers are dedicated physical machines with direct hardware access, supporting provisioning, rebuilding, and network configuration within a cloud region.
---

# gcore_cloud_baremetal_server (Resource)

Bare metal servers are dedicated physical machines with direct hardware access, supporting provisioning, rebuilding, and network configuration within a cloud region.

## Example Usage

### Baremetal server with one public interface

Create a basic baremetal server with a single external IPv4 interface.

```terraform
# Create an SSH key for baremetal server access
resource "gcore_cloud_ssh_key" "my_key" {
  project_id = 1
  name       = "my-keypair"
  public_key = "ssh-ed25519 ...your public key... user@example.com"
}

# Create a baremetal server with a single external interface
resource "gcore_cloud_baremetal_server" "server" {
  project_id   = 1
  region_id    = 1
  flavor       = "bm1-infrastructure-small"
  name         = "my-bare-metal"
  image_id     = "0f25a566-91a4-4507-aa42-bdd732fb998d"
  ssh_key_name = gcore_cloud_ssh_key.my_key.name

  interfaces = [{
    type      = "external"
    ip_family = "ipv4"
  }]
}
```

### Baremetal server with two interfaces

Create a baremetal server with two network interfaces: one public and one private.

```terraform
# Create a private network and subnet (baremetal requires vlan, not vxlan)
resource "gcore_cloud_network" "network" {
  project_id = 1
  region_id  = 1
  name       = "my-network"
  type       = "vlan"
}

resource "gcore_cloud_network_subnet" "subnet" {
  project_id      = 1
  region_id       = 1
  name            = "my-subnet"
  cidr            = "192.168.10.0/24"
  network_id      = gcore_cloud_network.network.id
  dns_nameservers = ["8.8.4.4", "1.1.1.1"]
}

# Create an SSH key for baremetal server access
resource "gcore_cloud_ssh_key" "my_key" {
  project_id = 1
  name       = "my-keypair"
  public_key = "ssh-ed25519 ...your public key... user@example.com"
}

# Create a baremetal server with two interfaces: one public, one private
resource "gcore_cloud_baremetal_server" "server_with_two_interfaces" {
  project_id   = 1
  region_id    = 1
  flavor       = "bm1-infrastructure-small"
  name         = "my-bare-metal"
  image_id     = "0f25a566-91a4-4507-aa42-bdd732fb998d"
  ssh_key_name = gcore_cloud_ssh_key.my_key.name

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

### Windows baremetal server

Create a Windows baremetal server with a public interface.

```terraform
# Create a Windows baremetal server with a public interface
resource "gcore_cloud_baremetal_server" "windows_server" {
  project_id          = 1
  region_id           = 1
  flavor              = "bm1-infrastructure-small"
  name                = "my-windows-bare-metal"
  image_id            = "408a0e4d-6a28-4bae-93fa-f738d964f555"
  password_wo         = "my-s3cR3tP@ssw0rd"
  password_wo_version = 1

  interfaces = [{
    type      = "external"
    ip_family = "ipv4"
  }]
}
```

### Baremetal server with reserved public IP

Create a baremetal server using a pre-allocated reserved fixed IP address.

```terraform
# Create an SSH key for baremetal server access
resource "gcore_cloud_ssh_key" "my_key" {
  project_id = 1
  name       = "my-keypair"
  public_key = "ssh-ed25519 ...your public key... user@example.com"
}

# Reserve a public IP address
resource "gcore_cloud_reserved_fixed_ip" "external_fixed_ip" {
  project_id = 1
  region_id  = 1
  type       = "external"
}

# Create a baremetal server using the reserved public IP
resource "gcore_cloud_baremetal_server" "server_with_reserved_address" {
  project_id   = 1
  region_id    = 1
  flavor       = "bm1-infrastructure-small"
  name         = "my-bare-metal"
  image_id     = "0f25a566-91a4-4507-aa42-bdd732fb998d"
  ssh_key_name = gcore_cloud_ssh_key.my_key.name

  interfaces = [{
    type    = "reserved_fixed_ip"
    port_id = gcore_cloud_reserved_fixed_ip.external_fixed_ip.port_id
  }]
}
```

### Windows baremetal server with two users

Create a Windows baremetal server and use userdata to add a second user during provisioning.

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

# Create a Windows baremetal server with userdata to add a second user
resource "gcore_cloud_baremetal_server" "windows_with_userdata" {
  project_id          = 1
  region_id           = 1
  flavor              = "bm1-infrastructure-small"
  name                = "my-windows-bare-metal"
  image_id            = "408a0e4d-6a28-4bae-93fa-f738d964f555"
  password_wo         = "my-s3cR3tP@ssw0rd"
  password_wo_version = 1
  user_data           = base64encode(var.second_user_userdata)

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
- `interfaces` (Attributes List) A list of network interfaces for the server. You can create one or more interfaces - private, public, or both. (see [below for nested schema](#nestedatt--interfaces))

### Optional

> **NOTE**: [Write-only arguments](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments) are supported in Terraform 1.11 and later.

- `app_config` (Map of String) Parameters for the application template if creating the instance from an `apptemplate`.
- `apptemplate_id` (String) Apptemplate ID. Either `image_id` or `apptemplate_id` is required.
- `image_id` (String) Image ID. Either `image_id` or `apptemplate_id` is required.
- `name` (String) Server name.
- `name_template` (String) If you want server names to be automatically generated based on IP addresses, you can provide a name template instead of specifying the name manually. The template should include a placeholder that will be replaced during provisioning. Supported placeholders are: `{ip_octets}` (last 3 octets of the IP), `{two_ip_octets}`, and `{one_ip_octet}`.
- `password_wo` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) For Linux instances, 'username' and 'password' are used to create a new user. When only 'password' is provided, it is set as the password for the default user of the image. For Windows instances, 'username' cannot be specified. Use the 'password' field to set the password for the 'Admin' user on Windows. Use the 'user_data' field to provide a script to create new users on Windows. The password of the Admin user cannot be updated via 'user_data'.
- `password_wo_version` (Number) Version of the password write-only field. Increment this value to trigger a replacement when changing the password.
- `project_id` (Number) Project ID
- `region_id` (Number) Region ID
- `ssh_key_name` (String) Specifies the name of the SSH keypair, created via the
[/v1/`ssh_keys` endpoint](/docs/api-reference/cloud/ssh-keys/add-or-generate-ssh-key).
- `tags` (Map of String) Key-value tags to associate with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Both tag keys and values have a maximum length of 255 characters. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.
- `user_data` (String) String in base64 format. For Linux instances, 'user_data' is ignored when 'password' field is provided. For Windows instances, Admin user password is set by 'password' field and cannot be updated via 'user_data'. Examples of the `user_data`: https://cloudinit.readthedocs.io/en/latest/topics/examples.html
- `username` (String) For Linux instances, 'username' and 'password' are used to create a new user. For Windows instances, 'username' cannot be specified. Use 'password' field to set the password for the 'Admin' user on Windows.

### Read-Only

- `addresses` (Map of List of Object) Map of `network_name` to list of addresses in that network
- `blackhole_ports` (Attributes List) IP addresses of the instances that are blackholed by DDoS mitigation system (see [below for nested schema](#nestedatt--blackhole_ports))
- `created_at` (String) Datetime when bare metal server was created
- `fixed_ip_assignments` (Attributes List) Fixed IP assigned to instance (see [below for nested schema](#nestedatt--fixed_ip_assignments))
- `id` (String) The ID of this resource.
- `instance_isolation` (Attributes) Instance isolation information (see [below for nested schema](#nestedatt--instance_isolation))
- `region` (String) Region name
- `status` (String) Bare metal server status
Available values: "ACTIVE", "BUILD", "DELETED", "ERROR", "HARD_REBOOT", "MIGRATING", "PASSWORD", "PAUSED", "REBOOT", "REBUILD", "RESCUE", "RESIZE", "REVERT_RESIZE", "SHELVED", "SHELVED_OFFLOADED", "SHUTOFF", "SOFT_DELETED", "SUSPENDED", "UNKNOWN", "VERIFY_RESIZE".
- `vm_state` (String) Bare metal server state
Available values: "active", "building", "deleted", "error", "paused", "rescued", "resized", "shelved", "shelved_offloaded", "soft-deleted", "stopped", "suspended".

<a id="nestedatt--interfaces"></a>
### Nested Schema for `interfaces`

Required:

- `type` (String) A public IP address will be assigned to the instance.
Available values: "external", "subnet", "any_subnet", "reserved_fixed_ip".

Optional:

- `floating_ip` (Attributes) Allows the instance to have a public IP that can be reached from the internet. (see [below for nested schema](#nestedatt--interfaces--floating_ip))
- `interface_name` (String) Interface name. Defaults to `null` and is returned as `null` in the API response if not set.
- `ip_address` (String) You can specify a specific IP address from your subnet.
- `ip_family` (String) Specify `ipv4`, `ipv6`, or `dual` to enable both.
Available values: "dual", "ipv4", "ipv6".
- `network_id` (String) The network where the instance will be connected.
- `port_group` (Number) Specifies the trunk group to which this interface belongs. Applicable only for bare metal servers. Each unique port group is mapped to a separate trunk port. Use this to control how interfaces are grouped across trunks.
- `port_id` (String) Network ID the subnet belongs to. Port will be plugged in this network.
- `subnet_id` (String) The instance will get an IP address from this subnet.

<a id="nestedatt--interfaces--floating_ip"></a>
### Nested Schema for `interfaces.floating_ip`

Required:

- `source` (String) A new floating IP will be created and attached to the instance. A floating IP is a public IP that makes the instance accessible from the internet, even if it only has a private IP. It works like SNAT, allowing outgoing and incoming traffic.
Available values: "new", "existing".

Optional:

- `existing_floating_id` (String) An existing available floating IP id must be specified if the source is set to `existing`



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
$ terraform import gcore_cloud_baremetal_server.example '<project_id>/<region_id>/<server_id>'
```

