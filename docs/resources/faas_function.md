---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "gcore_faas_function Resource - terraform-provider-gcore"
subcategory: ""
description: |-
  Represent FaaS function
---

# gcore_faas_function (Resource)

Represent FaaS function

## Example Usage

```terraform
provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_faas_function" "func" {
        project_id = 1
        region_id = 1
        name = "testf"
        namespace = "ns4test"
        description = "function description"
        envs = {
                BIG = "EXAMPLE2"
        }
        runtime = "go1.16.6"
        code_text = <<EOF
package kubeless

import (
        "github.com/kubeless/kubeless/pkg/functions"
)

func Run(evt functions.Event, ctx functions.Context) (string, error) {
        return "Hello World!!", nil
}
EOF
        timeout = 5
        flavor = "80mCPU-128MB"
        main_method = "Run"
        min_instances = 1
        max_instances = 2
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `code_text` (String)
- `flavor` (String)
- `main_method` (String) Main startup method name
- `max_instances` (Number) Autoscaling max number of instances
- `min_instances` (Number) Autoscaling min number of instances
- `name` (String)
- `namespace` (String) Namespace of the function
- `runtime` (String)
- `timeout` (Number)

### Optional

- `dependencies` (String) Function dependencies to install
- `description` (String)
- `disabled` (Boolean) Set to true if function is disabled
- `enable_api_key` (Boolean) Enable/Disable api key authorization
- `envs` (Map of String)
- `keys` (List of String) List of used api keys
- `project_id` (Number)
- `project_name` (String)
- `region_id` (Number)
- `region_name` (String)

### Read-Only

- `build_message` (String)
- `build_status` (String)
- `created_at` (String)
- `deploy_status` (Map of Number)
- `endpoint` (String)
- `id` (String) The ID of this resource.
- `status` (String)

## Import

Import is supported using the following syntax:

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can be used, for example:

```shell
# import using <project_id>:<region_id>:<namespace_name><function_name> format
terraform import gcore_faas_function.test 1:6:ns:test_func
```
