---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "gcore_registry_credential Resource - terraform-provider-gcore"
subcategory: ""
description: |-
  Represent inference registry credential
---

# gcore_registry_credential (Resource)

Represent inference registry credential

## Example Usage

```terraform
provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

data "gcore_project" "project" {
  name = "Default"
}

resource "gcore_registry_credential" "creds" {
  name = "docker-io"
  username = "username"
  password = "passwd"
  registry_url = "docker.io"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String)
- `password` (String)
- `registry_url` (String)
- `username` (String)

### Optional

- `project_id` (Number)
- `project_name` (String)

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can be used, for example:

```shell
# import using <project_id>:<credentials_name> format
terraform import gcore_registry_credential.dockerio 1:docekrio
```
