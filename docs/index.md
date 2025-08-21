---
page_title: "Provider: Gcore"
description: |-
  Gcore Terraform Provider is used to manage Gcore resources.
---

# Gcore Provider

Gcore Terraform Provider allows you to automate the provisioning, management, and testing of your Gcore resources programatically.

## Authentication and Configuration

To start using the Gcore Terraform Provider you need to configure the provider with the proper credentials.

Configuration for the provider can be derived from multiple sources, which are applied in the following order:

1. Parameters in the provider configuration
2. Environment variables

### Provider Configuration

!> Warning: Hard-coded credentials are not recommended in any Terraform configuration and risk secret leakage should it ever be committed to a public version control system.

The [permanent API token](https://gcore.com/docs/account-settings/create-use-or-delete-a-permanent-api-token) can be provided by adding a `permanent_api_token` argument to the `gcore` provider block.

Example:

```terraform
provider "gcore" {
  permanent_api_token = "251$d3361.............1b35f26d8"
}
```

If needed, the **deprecated** `username` and `password` arguments can be added to the `gcore` provider block instead of a permanent API token.

Other settings that can be configured include:

- `api_endpoint`
- `gcore_cdn_api`
- `gcore_client_id`
- `gcore_cloud_api`
- `gcore_dns_api`
- `gcore_platform_api`
- `gcore_storage_api`

### Environment Variables

The [permanent API token](https://gcore.com/docs/account-settings/create-use-or-delete-a-permanent-api-token) can be provided by setting the `GCORE_PERMANENT_TOKEN` environment variable.

For example:

```terraform
provider "gcore" {}
```

```shell
export GCORE_PERMANENT_TOKEN='251$d3361.............1b35f26d8'
terraform plan
```

If needed, the **deprecated** username / password authentication can be used by setting the `GCORE_USERNAME` and `GCORE_PASSWORD` environment variables.

Other supported environment variables include:

- `GCORE_API_ENDPOINT`
- `GCORE_CDN_API`
- `GCORE_CLIENT_ID`
- `GCORE_CLOUD_API`
- `GCORE_DNS_API`
- `GCORE_PLATFORM_API`
- `GCORE_STORAGE_API`

## Example Usage

```terraform
terraform {
  required_version = ">= 0.13.0"
  required_providers {
    gcore = {
      source  = "G-Core/gcore"
      version = ">= 0.3.70"
    }
  }
}

provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

variable "project_id" {
  type    = number
  default = 1
}

variable "region_id" {
  type    = number
  default = 76
}

resource "gcore_keypair" "kp" {
  project_id  = var.project_id
  public_key  = "ssh-ed25519 AAAA...CjZ user@example.com"
  sshkey_name = "test_key"
}

resource "gcore_network" "network" {
  name       = "test_network"
  type       = "vxlan"
  region_id  = var.region_id
  project_id = var.project_id
}

resource "gcore_subnet" "subnet" {
  name            = "test_subnet"
  cidr            = "192.168.10.0/24"
  network_id      = gcore_network.network.id
  dns_nameservers = ["8.8.4.4", "1.1.1.1"]

  region_id  = var.region_id
  project_id = var.project_id
}

resource "gcore_subnet" "subnet2" {
  name            = "test_subnet_2"
  cidr            = "192.168.20.0/24"
  network_id      = gcore_network.network.id
  dns_nameservers = ["8.8.4.4", "1.1.1.1"]
  region_id       = var.region_id
  project_id      = var.project_id
}

resource "gcore_volume" "first_volume" {
  name       = "test_boot_volume_1"
  type_name  = "ssd_hiiops"
  image_id   = "8f0900ba-2002-4f79-b866-390444caa19e"
  size       = 10
  region_id  = var.region_id
  project_id = var.project_id
}

resource "gcore_volume" "second_volume" {
  name       = "test_boot_volume_2"
  type_name  = "ssd_hiiops"
  image_id   = "8f0900ba-2002-4f79-b866-390444caa19e"
  size       = 10
  region_id  = var.region_id
  project_id = var.project_id
}

resource "gcore_volume" "third_volume" {
  name = "test_data_volume"
  type_name = "ssd_hiiops"
  size = 6
  region_id = var.region_id
  project_id = var.project_id
}

resource "gcore_instancev2" "instance" {
  flavor_id    = "g1-standard-2-4"
  name         = "test_instance_1"
  keypair_name = gcore_keypair.kp.sshkey_name

  volume {
    source     = "existing-volume"
    volume_id  = gcore_volume.first_volume.id
    boot_index = 0
  }

  interface {
    type       = "subnet"
    network_id = gcore_network.network.id
    subnet_id  = gcore_subnet.subnet.id
    security_groups = ["11384ae2-2677-439c-8618-f350da006163"]
  }

  interface {
    type            = "subnet"
    network_id      = gcore_network.network.id
    subnet_id       = gcore_subnet.subnet2.id
    security_groups = ["11384ae2-2677-439c-8618-f350da006163"]
  }

  metadata_map = {
    owner = "username"
  }

  region_id  = var.region_id
  project_id = var.project_id
}

resource "gcore_loadbalancerv2" "lb" {
  project_id = var.project_id
  region_id  = var.region_id
  name       = "test_loadbalancer"
  flavor     = "lb1-1-2"
}

resource "gcore_lblistener" "listener" {
  project_id      = var.project_id
  region_id       = var.region_id
  name            = "test_listener"
  protocol        = "HTTP"
  protocol_port   = 80
  loadbalancer_id = gcore_loadbalancerv2.lb.id
}

resource "gcore_lbpool" "pl" {
  project_id      = var.project_id
  region_id       = var.region_id
  name            = "test_pool"
  protocol        = "HTTP"
  lb_algorithm    = "LEAST_CONNECTIONS"
  loadbalancer_id = gcore_loadbalancerv2.lb.id
  listener_id     = gcore_lblistener.listener.id
  health_monitor {
    type        = "PING"
    delay       = 60
    max_retries = 5
    timeout     = 10
  }
}

resource "gcore_lbmember" "lbm" {
  project_id    = var.project_id
  region_id     = var.region_id
  pool_id       = gcore_lbpool.pl.id
  instance_id   = gcore_instancev2.instance.id
  address       = tolist(gcore_instancev2.instance.interface).0.ip_address
  protocol_port = 8081
}

resource "gcore_instancev2" "instance2" {
  flavor_id    = "g1-standard-2-4"
  name         = "test_instance_2"
  keypair_name = gcore_keypair.kp.sshkey_name

  volume {
    source     = "existing-volume"
    volume_id  = gcore_volume.second_volume.id
    boot_index = 0
  }

  volume {
    source = "existing-volume"
    volume_id = gcore_volume.third_volume.id
    boot_index = 1
  }

  interface {
    type            = "subnet"
    network_id      = gcore_network.network.id
    subnet_id       = gcore_subnet.subnet.id
    security_groups = ["11384ae2-2677-439c-8618-f350da006163"]
  }

  metadata_map = {
    owner = "username"
  }

  region_id  = var.region_id
  project_id = var.project_id
}

resource "gcore_lbmember" "lbm2" {
  project_id    = var.project_id
  region_id     = var.region_id
  pool_id       = gcore_lbpool.pl.id
  instance_id   = gcore_instancev2.instance2.id
  address       = tolist(gcore_instancev2.instance2.interface).0.ip_address
  protocol_port = 8081
  weight        = 5
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `api_endpoint` (String) A single API endpoint for all products. Will be used when specific product API url is not defined. Can also be set with the GCORE_API_ENDPOINT environment variable.
- `gcore_api` (String, Deprecated) Region API.
- `gcore_cdn_api` (String) CDN API (define only if you want to override CDN API endpoint). Can also be set with the GCORE_CDN_API environment variable.
- `gcore_client_id` (String) Client ID. Can also be set with the GCORE_CLIENT_ID environment variable.
- `gcore_cloud_api` (String) Region API (define only if you want to override Region API endpoint). Can also be set with the GCORE_CLOUD_API environment variable.
- `gcore_dns_api` (String) DNS API (define only if you want to override DNS API endpoint). Can also be set with the GCORE_DNS_API environment variable.
- `gcore_fastedge_api` (String) FastEdge API (define only if you want to override FastEdge API endpoint). Can also be set with the GCORE_FASTEDGE_API environment variable.
- `gcore_iam_api` (String) IAM API (define only if you want to override IAM API endpoint). Can also be set with the GCORE_IAM_API environment variable.
- `gcore_platform` (String, Deprecated) Platform URL is used for generate JWT.
- `gcore_platform_api` (String) Platform URL is used for generate JWT (define only if you want to override Platform API endpoint). Can also be set with the GCORE_PLATFORM_API environment variable.
- `gcore_storage_api` (String) Storage API (define only if you want to override Storage API endpoint). Can also be set with the GCORE_STORAGE_API environment variable.
- `gcore_waap_api` (String) WAAP API (define only if you want to override WAAP API endpoint). Can also be set with the GCORE_WAAP_API environment variable.
- `ignore_creds_auth_error` (Boolean, Deprecated) Should be set to true when you are gonna to use storage resource with permanent API-token only.
- `password` (String, Deprecated) Gcore account password. Can also be set with the GCORE_PASSWORD environment variable.
- `permanent_api_token` (String, Sensitive) A permanent [API-token](https://gcore.com/docs/account-settings/create-use-or-delete-a-permanent-api-token). Can also be set with the GCORE_PERMANENT_TOKEN environment variable.
- `user_name` (String, Deprecated) Gcore account username. Can also be set with the GCORE_USERNAME environment variable.
