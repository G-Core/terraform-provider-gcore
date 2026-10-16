---
page_title: "gcore_cloud_gpu_baremetal_cluster Resource - gcore"
subcategory: ""
description: |-
  
---

# gcore_cloud_gpu_baremetal_cluster (Resource)



## Example Usage

### Prerequisite

```terraform
# Prerequisite resources for GPU bare metal cluster examples

resource "gcore_cloud_network" "network" {
  project_id = 1
  region_id  = 1
  name       = "my-network"
  type       = "vlan"
}

resource "gcore_cloud_network_subnet" "subnet" {
  project_id = 1
  region_id  = 1
  name       = "my-subnet"
  cidr       = "192.168.10.0/24"
  network_id = gcore_cloud_network.network.id
}

resource "gcore_cloud_ssh_key" "keypair" {
  project_id = 1
  name       = "my-keypair"
  public_key = "ssh-ed25519 ...your public key... gcore@gcore.com"
}

data "gcore_cloud_file_share" "vast" {
  project_id = 1
  region_id  = 1
  find_one_by = {
    name = "my-files-share"
  }
}
```

### Creating GPU cluster with one public interface and tags

```terraform
# GPU bare metal cluster with one public interface and tags
resource "gcore_cloud_gpu_baremetal_cluster" "gpu_cluster" {
  project_id    = 1
  region_id     = 1
  flavor        = "bm3-ai-ndp2-1xlarge-h100-80-8"
  image_id      = "234c133c-b37e-4744-8a26-dc32fe407066"
  name          = "my-gpu-cluster"
  servers_count = 1

  servers_settings = {
    interfaces = [{
      type = "external"
    }]
    credentials = {
      ssh_key_name = "my-keypair"
    }
  }

  tags = {
    my-tag-key = "my-tag-value"
  }
}
```

### Creating GPU cluster with two interfaces

This example demonstrates how to create a baremetal GPU cluster with two network interfaces: one public and one private.

```terraform
# GPU bare metal cluster with two interfaces: one public and one private
resource "gcore_cloud_gpu_baremetal_cluster" "gpu_cluster" {
  project_id    = 1
  region_id     = 1
  flavor        = "bm3-ai-ndp2-1xlarge-h100-80-8"
  image_id      = "234c133c-b37e-4744-8a26-dc32fe407066"
  name          = "my-gpu-cluster"
  servers_count = 1

  servers_settings = {
    interfaces = [
      {
        type = "external"
      },
      {
        type       = "subnet"
        network_id = gcore_cloud_network.network.id
        subnet_id  = gcore_cloud_network_subnet.subnet.id
      },
    ]
    credentials = {
      ssh_key_name = "my-keypair"
    }
  }
}
```

### Creating GPU cluster with VAST file share

This example demonstrates how to create a baremetal GPU cluster with VAST file share manual mount.

```terraform
# GPU bare metal cluster with VAST file share manual mount
resource "gcore_cloud_file_share" "vast" {
  project_id = 1
  region_id  = 1
  name       = "tf-file-share-vast"
  size       = 10
  type_name  = "vast"
  protocol   = "NFS"
  share_settings = {
    allowed_characters = "LCD"
    path_length        = "LCD"
    root_squash        = true
  }
}

resource "gcore_cloud_gpu_baremetal_cluster" "gpu_cluster" {
  project_id    = 1
  region_id     = 1
  flavor        = "bm3-ai-ndp2-1xlarge-h100-80-8"
  image_id      = "234c133c-b37e-4744-8a26-dc32fe407066"
  name          = "cluster-for-vast"
  servers_count = 1

  servers_settings = {
    interfaces = [
      {
        type = "external"
      },
      {
        # This interface is required to ensure that the AI cluster
        # is connected to the same network as the file share.
        # Without it, mounting the NFS share will fail.
        type       = "subnet"
        network_id = gcore_cloud_file_share.vast.network_id
        subnet_id  = gcore_cloud_file_share.vast.subnet_id
      },
    ]
    credentials = {
      ssh_key_name = "my-keypair"
    }
    file_shares = [{
      id         = gcore_cloud_file_share.vast.id
      mount_path = "/mnt/vast"
    }]
    user_data = base64encode(<<-EOT
      #cloud-config
      runcmd:
        - mkdir -p /mnt/vast
        - apt-get update -y
        - apt-get install -y nfs-common
        - mount -o vers=3,nconnect=56,remoteports=dns,spread_reads,spread_writes,noextend ${data.gcore_cloud_file_share.vast.connection_point} /mnt/vast
    EOT
    )
  }

  tags = {
    my-tag-key = "my-tag-value"
  }
}
```

## Changing the image, SSH key or user data re-images the servers

Changing `image_id`, `servers_settings.credentials.ssh_key_name` or `servers_settings.user_data`
is applied in place: the cluster keeps its ID and the servers keep their addresses, but the API
reinstalls the operating system on every server whose settings are out of date, and anything on
their local disks is lost. Terraform shows this as an ordinary in-place update, so the provider
emits a warning at plan time. Read it before applying, and revert the attribute if the re-image
is not intended.

The API refuses this operation unless the caller confirms it (`max_disruption = "rebuild"` on
`apply_settings`). The provider sends that confirmation on your behalf; the plan-time warning is
the only prompt, and a non-interactive apply (`-auto-approve`, or applying a saved plan) proceeds
past it.

The warning is also emitted when one of these attributes is not yet known at plan time (for
example, interpolated from a resource that does not exist yet). The apply then re-images the
servers only if the resolved value actually differs.

~> **On a plan that replaces the cluster** (a `flavor` change, or `terraform plan -replace`) no
in-place re-image takes place: every server is destroyed and created again. The provider
suppresses the warning and the credential checks below when the replacement is caused by an
attribute change. A `-replace` given on the command line, or a `replace_triggered_by`, is
resolved by Terraform after the provider has answered and is not visible to it: on such a plan
the warning may still appear and the credential checks still apply.

Removing `servers_settings.credentials.ssh_key_name` from the configuration is rejected: the API
has no operation that takes a key away from running servers. Keep the key; to build the cluster
again without it, change an attribute that forces replacement or destroy and recreate it.
Changing `username` or `password_wo_version` is rejected as well, since the settings update
accepts only `ssh_key_name`.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `flavor` (String) Cluster flavor ID
- `image_id` (String) System image ID
- `name` (String) Cluster name
- `servers_count` (Number) Number of servers in the cluster
- `servers_settings` (Attributes) Configuration settings for the servers in the cluster (see [below for nested schema](#nestedatt--servers_settings))

### Optional

- `project_id` (Number) Project ID
- `region_id` (Number) Region ID
- `tags` (Map of String) Key-value tags to associate with the resource. A tag is a key-value pair that can be associated with a resource, enabling efficient filtering and grouping for better organization and management. Both tag keys and values have a maximum length of 255 characters. Some tags are read-only and cannot be modified by the user. Tags are also integrated with cost reports, allowing cost data to be filtered based on tag keys or values.

### Read-Only

- `created_at` (String) Cluster creation date time
- `has_pending_changes` (Boolean) True if any server in the cluster has pending (not yet applied) settings changes
- `id` (String) The ID of this resource.
- `managed_by` (String) User type managing the resource
Available values: "k8s", "user".
- `servers_ids` (List of String) List of cluster nodes
- `status` (String) Cluster status
Available values: "active", "creating", "degraded", "deleting", "error", "rebooting", "rebuilding", "resizing", "shutoff".
- `updated_at` (String) Cluster update date time

<a id="nestedatt--servers_settings"></a>
### Nested Schema for `servers_settings`

Required:

- `interfaces` (Attributes List) Subnet IPs and floating IPs (see [below for nested schema](#nestedatt--servers_settings--interfaces))

Optional:

- `credentials` (Attributes) Optional server access credentials (see [below for nested schema](#nestedatt--servers_settings--credentials))
- `file_shares` (Attributes List) List of file shares to be mounted across the cluster. (see [below for nested schema](#nestedatt--servers_settings--file_shares))
- `security_groups` (Attributes List, Deprecated) Deprecated. Use per-interface `security_groups` inside `interfaces[]` instead. Cannot be combined with per-interface `security_groups`. If omitted everywhere, the project's default security group is applied. (see [below for nested schema](#nestedatt--servers_settings--security_groups))
- `user_data` (String) Optional custom user data (Base64-encoded)

<a id="nestedatt--servers_settings--interfaces"></a>
### Nested Schema for `servers_settings.interfaces`

Required:

- `type` (String) Available values: "external", "subnet", "any_subnet".

Optional:

- `floating_ip` (Attributes) Floating IP config for this subnet attachment. Only applies when `type` is "subnet" or "any_subnet". (see [below for nested schema](#nestedatt--servers_settings--interfaces--floating_ip))
- `ip_family` (String) Which subnets should be selected: IPv4, IPv6, or use dual stack.
Only applies when `type` is "external" or "any_subnet". Must not be set when `type` is "subnet", where the IP family is determined by the referenced subnet.
Available values: "dual", "ipv4", "ipv6".
- `network_id` (String) Network ID the subnet belongs to. Port will be plugged in this network. Required when `type` is "subnet" or "any_subnet"; must not be set when `type` is "external".
- `port_security_enabled` (Boolean) Controls port security for this interface. When omitted, the default applies (port security enabled, default security group attached). When false, the port is created with port security off and no security group attached; `security_groups` must not be set in that case. Not allowed for interfaces on a public network, nor for bare metal servers without a DPU (their ports cannot enforce port security).
- `security_groups` (Attributes List) Security group UUIDs applied to this interface. If omitted (or empty), the top-level `security_groups` value applies; if both are omitted, the project's default security group is applied. (see [below for nested schema](#nestedatt--servers_settings--interfaces--security_groups))
- `subnet_id` (String) Port is assigned an IP address from this subnet. Required when `type` is "subnet"; must not be set for any other interface type.

Read-Only:

- `name` (String) Interface name

<a id="nestedatt--servers_settings--interfaces--floating_ip"></a>
### Nested Schema for `servers_settings.interfaces.floating_ip`

Required:

- `source` (String) Available values: "new".


<a id="nestedatt--servers_settings--interfaces--security_groups"></a>
### Nested Schema for `servers_settings.interfaces.security_groups`

Required:

- `id` (String) Resource ID



<a id="nestedatt--servers_settings--credentials"></a>
### Nested Schema for `servers_settings.credentials`

Optional:

- `password_wo` (String, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) Used to set the password for the specified 'username' on Linux instances. If 'username' is not provided, the password is applied to the default user of the image. Mutually exclusive with '`user_data`' - only one can be specified.
- `password_wo_version` (Number) Version of the password write-only field. Increment this value to trigger an update when changing the password.
- `ssh_key_name` (String) Specifies the name of the SSH keypair, created via the
[/v1/`ssh_keys` endpoint](/docs/api-reference/cloud/ssh-keys/add-or-generate-ssh-key).
- `username` (String) The 'username' and 'password' fields create a new user on the system


<a id="nestedatt--servers_settings--file_shares"></a>
### Nested Schema for `servers_settings.file_shares`

Required:

- `id` (String) Unique identifier of the file share in UUID format.
- `mount_path` (String) Absolute mount path inside the system where the file share will be mounted.


<a id="nestedatt--servers_settings--security_groups"></a>
### Nested Schema for `servers_settings.security_groups`

Required:

- `id` (String) Resource ID


## Import

Import is supported using the following syntax:

```shell
$ terraform import gcore_cloud_gpu_baremetal_cluster.example '<project_id>/<region_id>/<cluster_id>'
```

~> **Clusters created with `username` and `password_wo`.** The API never returns the password or
`password_wo_version`, so after an import (with `terraform import` or an `import` block) the
configured credentials cannot be matched against the cluster, and plans fail with "Unsupported
credential update", including the refresh plan `terraform destroy` runs first. Before planning, write the
credentials the cluster was created with and add this `lifecycle` block, and keep it for as long as
the cluster exists (the credentials cannot change without recreating it anyway):

```terraform
resource "gcore_cloud_gpu_baremetal_cluster" "example" {
  # ...
  servers_settings = {
    # ...
    credentials = {
      username            = "admin"
      password_wo         = var.cluster_password
      password_wo_version = 1
    }
  }

  lifecycle {
    ignore_changes = [servers_settings.credentials]
  }
}
```

While the block is present, all changes to `credentials` are ignored: bumping `password_wo_version`
does **not** change the password, and nothing reports that. A replacement (for example a `flavor`
change) still builds the new cluster with the configured credentials, so they must be real, never
placeholders. A cluster that uses `ssh_key_name`
imports without this block and keeps its in-place key change, which the block would ignore too.
Use the block only for imported clusters: on a cluster Terraform created, it hides credentials
added or changed in the configuration, and the cluster keeps the ones it was created with.

