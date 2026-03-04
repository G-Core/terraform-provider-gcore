# Gcore Terraform Provider

The [Gcore Terraform provider](https://registry.terraform.io/providers/G-Core/gcore/latest/docs) provides convenient access to
the [Gcore REST API](https://api.gcore.com/docs) from Terraform.

It is generated with [Stainless](https://www.stainless.com/).

## Requirements

This provider requires Terraform CLI 1.0 or later. You can [install it for your system](https://developer.hashicorp.com/terraform/install)
on Hashicorp's website.

## Usage

Add the following to your `main.tf` file:

<!-- x-release-please-start-version -->

```hcl
# Declare the provider and version
terraform {
  required_providers {
    gcore = {
      source  = "G-Core/gcore"
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
  cloud_polling_timeout_seconds = 0
}

# Configure a resource
resource "gcore_cloud_secret" "example_cloud_secret" {
  project_id = 1
  region_id = 1
  name = "Load balancer certificate #1"
  payload = {
    certificate = "<certificate>"
    certificate_chain = "<certificate_chain>"
    private_key = "<private_key>"
  }
  expiration = "2019-12-27T18:11:19.117Z"
}
```

<!-- x-release-please-end -->

Initialize your project by running `terraform init` in the directory.

Additional examples can be found in the [./examples](./examples) folder within this repository, and you can
refer to the full documentation on [the Terraform Registry](https://registry.terraform.io/providers/G-Core/gcore/latest/docs).

### Provider Options

When you initialize the provider, the following options are supported. It is recommended to use environment variables for sensitive values like access tokens.
If an environment variable is provided, then the option does not need to be set in the terraform source.

| Property                       | Environment variable     | Required | Default value |
| ------------------------------ | ------------------------ | -------- | ------------- |
| api_key                        | `GCORE_API_KEY`          | true     | —             |
| cloud_region_id                | `GCORE_CLOUD_REGION_ID`  | false    | —             |
| cloud_project_id               | `GCORE_CLOUD_PROJECT_ID` | false    | —             |
| cloud_polling_timeout_seconds  | -                        | false    | `7200`        |
| cloud_polling_interval_seconds | -                        | false    | `3`           |

## Semantic versioning

This package generally follows [SemVer](https://semver.org/spec/v2.0.0.html) conventions, though certain backwards-incompatible changes may be released as minor versions:

1. Changes to library internals which are technically public but not intended or documented for external use. _(Please open a GitHub issue to let us know if you are relying on such internals.)_
2. Changes that we do not expect to impact the vast majority of users in practice.

We take backwards-compatibility seriously and work hard to ensure you can rely on a smooth upgrade experience.

We are keen for your feedback; please open an [issue](https://www.github.com/G-Core/terraform-provider-gcore/issues) with questions, bugs, or suggestions.

## Contributing

See [the contributing documentation](./CONTRIBUTING.md).
