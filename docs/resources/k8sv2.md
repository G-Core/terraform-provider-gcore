---
page_title: "gcore_k8sv2 Resource - terraform-provider-gcore"
subcategory: ""
description: |-
  Represent k8s cluster with one default pool.
---

# gcore_k8sv2 (Resource)

Represent k8s cluster with one default pool.

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

resource "gcore_keypair" "my_keypair" {
  project_id  = data.gcore_project.project.id
  sshkey_name = "my-keypair"
  public_key  = "ssh-ed25519 ...your public key... gcore@gcore.com"
}
```

### Creating a managed Kubernetes cluster in a private network

```terraform
resource "gcore_network" "network" {
  name       = "my-network"
  type       = "vxlan"
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}

resource "gcore_subnet" "subnet" {
  name       = "my-subnet"
  cidr       = "192.168.10.0/24"
  network_id = gcore_network.network.id

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}

resource "gcore_k8sv2" "cluster" {
  project_id    = data.gcore_project.project.id
  region_id     = data.gcore_region.region.id
  name          = "my-k8s-cluster"
  fixed_network = gcore_network.network.id
  fixed_subnet  = gcore_subnet.subnet.id
  keypair       = gcore_keypair.my_keypair.sshkey_name
  version       = "v1.31.9"
  pool {
    name               = "my-k8s-pool"
    flavor_id          = "g1-standard-2-4"
    servergroup_policy = "soft-anti-affinity"
    min_node_count     = 1
    max_node_count     = 1
    boot_volume_size   = 10
    boot_volume_type   = "standard"
  }
}

data "gcore_k8sv2_kubeconfig" "config" {
  cluster_name = gcore_k8sv2.cluster.name
  region_id    = data.gcore_region.region.id
  project_id   = data.gcore_project.project.id
}

// to store kubeconfig in a file pls use
// terraform output -raw kubeconfig > config.yaml
output "kubeconfig" {
  value = data.gcore_k8sv2_kubeconfig.config.kubeconfig
}
```

### Creating a managed Kubernetes cluster in a public network with advanced DDoS protection

```terraform
resource "gcore_k8sv2" "cluster" {
  project_id    = data.gcore_project.project.id
  region_id     = data.gcore_region.region.id
  name          = "my-k8s-cluster"
  keypair       = gcore_keypair.my_keypair.sshkey_name
  version       = "v1.31.9"
  pool {
    name             = "my-k8s-pool"
    flavor_id        = "g1-standard-2-4"
    servergroup_policy = "soft-anti-affinity"
    min_node_count   = 1
    max_node_count   = 1
    boot_volume_size = 10
    boot_volume_type = "standard"
    is_public_ipv4 = true
  }
  ddos_profile {
    enabled = true
    fields {
      base_field = 1353
      field_value = jsonencode(["AF"])
    }
    fields {
      base_field = 1354
      field_value = jsonencode(50)
    }
    fields {
      base_field = 1355
      field_value = jsonencode(150)
    }
    fields {
      base_field = 1356
      field_value = jsonencode(300)
    }
    fields {
      base_field = 1357
      field_value = jsonencode(300)
    }

    fields {
      base_field = 1352
      field_value = jsonencode([
        {
          "sip_list":["192.168.0.1","10.10.0.1"],
          "dport_list": ["27015","27025"],
          "proto_list": ["udp"],
          "sport_list": ["27025"],
          "policy": "DROP"
        }
      ])
    }
    profile_template = 1128
  }
}

data "gcore_k8sv2_kubeconfig" "config" {
  cluster_name       = gcore_k8sv2.cluster.name
  region_id          = data.gcore_region.region.id
  project_id         = data.gcore_project.project.id
}

// to store kubeconfig in a file pls use
// terraform output -raw kubeconfig > config.yaml
output "kubeconfig" {
  value = data.gcore_k8sv2_kubeconfig.config.kubeconfig
}
```

### Creating a managed Kubernetes cluster in a public network with custom security group rules

```terraform
resource "gcore_k8sv2" "cluster" {
  project_id    = data.gcore_project.project.id
  region_id     = data.gcore_region.region.id
  name          = "my-k8s-cluster"
  keypair       = gcore_keypair.my_keypair.sshkey_name
  version       = "v1.31.9"
  pool {
    name             = "my-k8s-pool"
    flavor_id        = "g1-standard-2-4"
    servergroup_policy = "soft-anti-affinity"
    min_node_count   = 1
    max_node_count   = 1
    boot_volume_size = 10
    boot_volume_type = "standard"
    is_public_ipv4 = true
  }
  security_group_rules {
    direction      = "ingress"
    ethertype      = "IPv4"
    protocol       = "tcp"
    port_range_min = 80
    port_range_max = 80
  }
}

data "gcore_k8sv2_kubeconfig" "config" {
  cluster_name       = gcore_k8sv2.cluster.name
  region_id          = data.gcore_region.region.id
  project_id         = data.gcore_project.project.id
}

// to store kubeconfig in a file pls use
// terraform output -raw kubeconfig > config.yaml
output "kubeconfig" {
  value = data.gcore_k8sv2_kubeconfig.config.kubeconfig
}
```

### Creating a managed Kubernetes cluster with vast share integration

```terraform
resource "gcore_file_share" "vast" {
  name       = "tf-file-share-vast"
  size       = 10
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
  type_name  = "vast"
  protocol   = "NFS"
  share_settings {
    allowed_characters = "LCD"
    path_length = "LCD"
    root_squash = true
  }
}

resource "gcore_network" "network" {
  name       = "my-network"
  type       = "vxlan"
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}

resource "gcore_subnet" "subnet" {
  name       = "my-subnet"
  cidr       = "192.168.10.0/24"
  network_id = gcore_network.network.id

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}

resource "gcore_k8sv2" "cluster" {
  project_id    = data.gcore_project.project.id
  region_id     = data.gcore_region.region.id
  name          = "my-k8s-cluster"
  fixed_network = gcore_network.network.id
  fixed_subnet  = gcore_subnet.subnet.id
  keypair       = gcore_keypair.my_keypair.sshkey_name
  version       = "v1.33.3"

  cni {
    provider = "cilium"
  }

  csi {
    nfs {
      vast_enabled = true
    }
  }

  pool {
    name               = "gpu-1"
    flavor_id          = "bm3-ai-ndp-1xlarge-h100-80-8"
    is_public_ipv4     = false
    min_node_count     = 1
    max_node_count     = 1
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `keypair` (String) Name of the keypair used for SSH access to nodes.
- `name` (String) Cluster name.
- `pool` (Block List, Min: 1) (see [below for nested schema](#nestedblock--pool))
- `version` (String) Kubernetes version.

### Optional

- `authentication` (Block List, Max: 1) Cluster authentication configuration. (see [below for nested schema](#nestedblock--authentication))
- `autoscaler_config` (Map of String) Cluster autoscaler configuration params. Keys and values are expected to follow the cluster-autoscaler option format.
- `cni` (Block List, Max: 1) Cluster CNI configuration. (see [below for nested schema](#nestedblock--cni))
- `csi` (Block List, Max: 1) Container Storage Interface (CSI) driver settings. (see [below for nested schema](#nestedblock--csi))
- `ddos_profile` (Block List, Max: 1) DDoS profile configuration. (see [below for nested schema](#nestedblock--ddos_profile))
- `fixed_network` (String) Fixed network used to allocate network addresses for cluster nodes.
- `fixed_subnet` (String) Fixed subnet used to allocate network addresses for cluster nodes. Subnet should have a router.
- `is_ipv6` (Boolean) Enable public IPv6 address.
- `pods_ip_pool` (String) Pods IPv4 IP pool in CIDR notation.
- `pods_ipv6_pool` (String) Pods IPv6 IP pool in CIDR notation.
- `project_id` (Number)
- `project_name` (String)
- `region_id` (Number)
- `region_name` (String)
- `security_group_rules` (Block Set) Firewall rules control what inbound(ingress) and outbound(egress) traffic is allowed to enter or leave a Instance. At least one 'egress' rule should be set (see [below for nested schema](#nestedblock--security_group_rules))
- `services_ip_pool` (String) Services IPv4 IP pool in CIDR notation.
- `services_ipv6_pool` (String) Services IPv6 IP pool in CIDR notation.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `created_at` (String) Cluster creation date.
- `creator_task_id` (String)
- `id` (String) The ID of this resource.
- `is_public` (Boolean) True if the cluster is public.
- `security_group_id` (String) Security group ID.
- `status` (String) Cluster status.
- `task_id` (String)

<a id="nestedblock--pool"></a>
### Nested Schema for `pool`

Required:

- `flavor_id` (String) Cluster pool node flavor ID. Changing the value of this attribute will trigger recreation of the cluster pool.
- `min_node_count` (Number) Minimum number of nodes in the cluster pool.
- `name` (String) Cluster pool name. Changing the value of this attribute will trigger recreation of the cluster pool.

Optional:

- `auto_healing_enabled` (Boolean) Enable/disable auto healing of cluster pool nodes.
- `boot_volume_size` (Number) Cluster pool boot volume size. Must be set only for VM pools. Changing the value of this attribute will trigger recreation of the cluster pool.
- `boot_volume_type` (String) Cluster pool boot volume type. Must be set only for VM pools. Available values are 'standard', 'ssd_hiiops', 'cold', 'ultra'. Changing the value of this attribute will trigger recreation of the cluster pool.
- `crio_config` (Map of String) Crio configuration for pool nodes. Keys and values are expected to follow the crio option format.
- `is_public_ipv4` (Boolean) Assign public IPv4 address to nodes in this pool. Changing the value of this attribute will trigger recreation of the cluster pool.
- `kubelet_config` (Map of String) Kubelet configuration for pool nodes. Keys and values are expected to follow the kubelet configuration file format.
- `labels` (Map of String) Labels applied to the cluster pool nodes.
- `max_node_count` (Number) Maximum number of nodes in the cluster pool.
- `servergroup_policy` (String) Server group policy: anti-affinity, soft-anti-affinity or affinity
- `taints` (Map of String) Taints applied to the cluster pool nodes.

Read-Only:

- `created_at` (String) Cluster pool creation date.
- `node_count` (Number) Current node count in the cluster pool.
- `servergroup_id` (String) Server group id
- `servergroup_name` (String) Server group name
- `status` (String) Cluster pool status.


<a id="nestedblock--authentication"></a>
### Nested Schema for `authentication`

Optional:

- `oidc` (Block List, Max: 1) OpenID Connect configuration settings. (see [below for nested schema](#nestedblock--authentication--oidc))

<a id="nestedblock--authentication--oidc"></a>
### Nested Schema for `authentication.oidc`

Optional:

- `client_id` (String) A client id that all tokens must be issued for.
- `groups_claim` (String) JWT claim to use as the user's group.
- `groups_prefix` (String) Prefix prepended to group claims to prevent clashes with existing names.
- `issuer_url` (String) URL of the provider that allows the API server to discover public signing keys. Only URLs that use the https:// scheme are accepted.
- `required_claims` (Map of String) A map describing required claims in the ID Token. Each claim is verified to be present in the ID Token with a matching value.
- `signing_algs` (Set of String) Accepted signing algorithms. Supported values are: RS256, RS384, RS512, ES256, ES384, ES512, PS256, PS384, PS512.
- `username_claim` (String) JWT claim to use as the user name. When not specified, the `sub` claim will be used.
- `username_prefix` (String) Prefix prepended to username claims to prevent clashes with existing names.



<a id="nestedblock--cni"></a>
### Nested Schema for `cni`

Optional:

- `cilium` (Block List, Max: 1) Cilium CNI configuration. (see [below for nested schema](#nestedblock--cni--cilium))
- `provider` (String) CNI provider to use when creating the cluster. Supported values are: calico, cilium. The default value is calico.

<a id="nestedblock--cni--cilium"></a>
### Nested Schema for `cni.cilium`

Optional:

- `encryption` (Boolean) Enables transparent network encryption. The default value is false.
- `hubble_relay` (Boolean) Enables Hubble Relay. The default value is false.
- `hubble_ui` (Boolean) Enables Hubble UI. Requires `hubble_relay=true`. The default value is false.
- `lb_acceleration` (Boolean) Enables load balancer acceleration via XDP. The default value is false.
- `lb_mode` (String) The operation mode of load balancing for remote backends. Supported values are snat, dsr, hybrid. The default value is snat.
- `mask_size` (Number) Specifies the size allocated from pods_ip_pool CIDR to node.ipam.podCIDRs. The default value is 24.
- `mask_size_v6` (Number) Specifies the size allocated from pods_ipv6_pool CIDR to node.ipam.podCIDRs. The default value is 120.
- `routing_mode` (String) Enables native-routing mode or tunneling mode. The default value is tunnel.
- `tunnel` (String) Tunneling protocol to use in tunneling mode and for ad-hoc tunnels. The default value is geneve.



<a id="nestedblock--csi"></a>
### Nested Schema for `csi`

Optional:

- `nfs` (Block List, Max: 1) NFS CSI driver settings. (see [below for nested schema](#nestedblock--csi--nfs))

<a id="nestedblock--csi--nfs"></a>
### Nested Schema for `csi.nfs`

Optional:

- `vast_enabled` (Boolean) Enable or disable VAST NFS integration. The default value is false. When set to true, a dedicated StorageClass will be created in the cluster for each VAST NFS file share defined in the cloud. All file shares created prior to cluster creation will be available immediately, while those created afterward may take a few minutes for to appear.



<a id="nestedblock--ddos_profile"></a>
### Nested Schema for `ddos_profile`

Required:

- `enabled` (Boolean) Indicates if the DDoS profile is enabled.

Optional:

- `fields` (Block List) List of fields for the DDoS profile. (see [below for nested schema](#nestedblock--ddos_profile--fields))
- `profile_template` (Number) Profile template ID.
- `profile_template_name` (String) Profile template name.

<a id="nestedblock--ddos_profile--fields"></a>
### Nested Schema for `ddos_profile.fields`

Required:

- `base_field` (Number) Base field ID.

Optional:

- `field_value` (String) Complex value. Only one of 'value' or 'field_value' must be specified.
- `value` (String) Basic type value. Only one of 'value' or 'field_value' must be specified.



<a id="nestedblock--security_group_rules"></a>
### Nested Schema for `security_group_rules`

Required:

- `direction` (String) Available value is 'ingress', 'egress'
- `ethertype` (String) Available value is 'IPv4', 'IPv6'
- `protocol` (String) Available value is udp,tcp,any,ipv6-icmp,ipv6-route,ipv6-opts,ipv6-nonxt,ipv6-frag,ipv6-encap,icmp,ah,dccp,egp,esp,gre,igmp,ospf,pgm,rsvp,sctp,udplite,vrrp,51,50,112,0,4,ipip,ipencap

Optional:

- `description` (String)
- `port_range_max` (Number)
- `port_range_min` (Number)
- `remote_group_id` (String)
- `remote_ip_prefix` (String)

Read-Only:

- `created_at` (String)
- `id` (String)
- `updated_at` (String)


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `update` (String)





## Import

Import is supported using the following syntax:

```shell
# import using <project_id>:<region_id>:<cluster_name> format
terraform import gcore_k8sv2.cluster 1:6:my-k8s-cluster
```

