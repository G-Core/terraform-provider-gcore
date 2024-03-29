---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "gcore_secret Data Source - terraform-provider-gcore"
subcategory: ""
description: |-
  Represent secret
---

# gcore_secret (Data Source)

Represent secret

## Example Usage

```terraform
provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

data "gcore_project" "pr" {
  name = "test"
}

data "gcore_region" "rg" {
  name = "ED-10 Preprod"
}

data "gcore_secret" "lb_https" {
  name       = "lb_https"
  region_id  = data.gcore_region.rg.id
  project_id = data.gcore_project.pr.id
}

output "view" {
  value = data.gcore_secret.lb_https
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String)

### Optional

- `project_id` (Number)
- `project_name` (String)
- `region_id` (Number)
- `region_name` (String)

### Read-Only

- `algorithm` (String)
- `bit_length` (Number)
- `content_types` (Map of String)
- `created` (String) Datetime when the secret was created. The format is 2025-12-28T19:14:44.180394
- `expiration` (String) Datetime when the secret will expire. The format is 2025-12-28T19:14:44.180394
- `id` (String) The ID of this resource.
- `mode` (String)
- `status` (String)
