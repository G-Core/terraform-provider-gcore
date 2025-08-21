# Gcore Terraform Provider

The [Gcore Terraform provider](https://registry.terraform.io/providers/stainless-sdks/gcore/latest/docs) provides convenient access to
the [Gcore REST API](https://api.gcore.com/docs) from Terraform.

It is generated with [Stainless](https://www.stainless.com/).

## Requirements

This provider requires Terraform CLI 1.0 or later. You can [install it for your system](https://developer.hashicorp.com/terraform/install)
on Hashicorp's website.

## Usage

Add the following to your `main.tf` file:

```hcl
# Declare the provider and version
terraform {
  required_providers {
    gcore = {
      source  = "stainless-sdks/gcore"
      version = "~> 0.0.1-alpha.0"
    }
  }
}

# Initialize the provider
provider "gcore" {
  api_key = "My API Key" # or set GCORE_API_KEY env variable
  cloud_project_id = 0 # or set GCORE_CLOUD_PROJECT_ID env variable
  cloud_region_id = 0 # or set GCORE_CLOUD_REGION_ID env variable
  cloud_polling_interval_seconds = 0
}

# Configure a resource
resource "gcore_cloud_load_balancer" "example_cloud_load_balancer" {
  project_id = 0
  region_id = 0
  flavor = "lb1-1-2"
  floating_ip = {
    existing_floating_id = "c64e5db1-5f1f-43ec-a8d9-5090df85b82d"
    source = "existing"
  }
  listeners = [{
    name = "my_listener"
    protocol = "HTTP"
    protocol_port = 80
    allowed_cidrs = ["10.0.0.0/8"]
    connection_limit = 100000
    insert_x_forwarded = false
    pools = [{
      lb_algorithm = "LEAST_CONNECTIONS"
      name = "pool_name"
      protocol = "HTTP"
      ca_secret_id = "ca_secret_id"
      crl_secret_id = "crl_secret_id"
      healthmonitor = {
        delay = 10
        max_retries = 3
        timeout = 5
        type = "HTTP"
        expected_codes = "200,301,302"
        http_method = "GET"
        max_retries_down = 3
        url_path = "/"
      }
      listener_id = "listener_id"
      loadbalancer_id = "bbb35f84-35cc-4b2f-84c2-a6a29bba68aa"
      members = [{
        address = "192.168.1.101"
        protocol_port = 8000
        admin_state_up = true
        backup = true
        instance_id = "a7e7e8d6-0bf7-4ac9-8170-831b47ee2ba9"
        monitor_address = "monitor_address"
        monitor_port = 0
        subnet_id = "32283b0b-b560-4690-810c-f672cbb2e28d"
        weight = 2
      }, {
        address = "192.168.1.102"
        protocol_port = 8000
        admin_state_up = true
        backup = true
        instance_id = "169942e0-9b53-42df-95ef-1a8b6525c2bd"
        monitor_address = "monitor_address"
        monitor_port = 0
        subnet_id = "32283b0b-b560-4690-810c-f672cbb2e28d"
        weight = 4
      }]
      secret_id = "secret_id"
      session_persistence = {
        type = "APP_COOKIE"
        cookie_name = "cookie_name"
        persistence_granularity = "persistence_granularity"
        persistence_timeout = 0
      }
      timeout_client_data = 50000
      timeout_member_connect = 50000
      timeout_member_data = 0
    }]
    secret_id = "f2e734d0-fa2b-42c2-ad33-4c6db5101e00"
    sni_secret_id = ["f2e734d0-fa2b-42c2-ad33-4c6db5101e00", "eb121225-7ded-4ff3-ae1f-599e145dd7cb"]
    timeout_client_data = 50000
    timeout_member_connect = 50000
    timeout_member_data = 0
    user_list = [{
      encrypted_password = "$5$isRr.HJ1IrQP38.m$oViu3DJOpUG2ZsjCBtbITV3mqpxxbZfyWJojLPNSPO5"
      username = "admin"
    }]
  }]
  logging = {
    destination_region_id = 1
    enabled = true
    retention_policy = {
      period = 45
    }
    topic_name = "my-log-name"
  }
  name = "new_load_balancer"
  name_template = "lb_name_template"
  preferred_connectivity = "L2"
  tags = {
    my-tag = "my-tag-value"
  }
  vip_ip_family = "dual"
  vip_network_id = "ac307687-31a4-4a11-a949-6bea1b2878f5"
  vip_port_id = "ff83e13a-b256-4be2-ba5d-028d3f0ab450"
  vip_subnet_id = "4e7802d3-5023-44b8-b298-7726558fddf4"
}
```

Initialize your project by running `terraform init` in the directory.

Additional examples can be found in the [./examples](./examples) folder within this repository, and you can
refer to the full documentation on [the Terraform Registry](https://registry.terraform.io/providers/stainless-sdks/gcore/latest/docs).

### Provider Options

When you initialize the provider, the following options are supported. It is recommended to use environment variables for sensitive values like access tokens.
If an environment variable is provided, then the option does not need to be set in the terraform source.

| Property                       | Environment variable     | Required | Default value |
| ------------------------------ | ------------------------ | -------- | ------------- |
| api_key                        | `GCORE_API_KEY`          | true     | —             |
| cloud_region_id                | `GCORE_CLOUD_REGION_ID`  | false    | —             |
| cloud_project_id               | `GCORE_CLOUD_PROJECT_ID` | false    | —             |
| cloud_polling_interval_seconds | -                        | false    | `3`           |

## Semantic versioning

This package generally follows [SemVer](https://semver.org/spec/v2.0.0.html) conventions, though certain backwards-incompatible changes may be released as minor versions:

1. Changes to library internals which are technically public but not intended or documented for external use. _(Please open a GitHub issue to let us know if you are relying on such internals.)_
2. Changes that we do not expect to impact the vast majority of users in practice.

We take backwards-compatibility seriously and work hard to ensure you can rely on a smooth upgrade experience.

We are keen for your feedback; please open an [issue](https://www.github.com/stainless-sdks/gcore-terraform/issues) with questions, bugs, or suggestions.

## Contributing

See [the contributing documentation](./CONTRIBUTING.md).
