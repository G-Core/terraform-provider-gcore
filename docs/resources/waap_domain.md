---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "gcore_waap_domain Resource - terraform-provider-gcore"
subcategory: ""
description: |-
  Represent WAAP domain
---

# gcore_waap_domain (Resource)

Represent WAAP domain

## Example Usage

```terraform
provider gcore {
  permanent_api_token = "768660$.............a43f91f"
}

resource "gcore_cdn_resource" "example" {
  cname  = "api.example.com"
  origin = "origin.example.com"
  options {
    waap { value = true }
  }
}

resource "gcore_waap_domain" "domain" {
  name   = gcore_cdn_resource.example.cname
  status = "monitor"

  settings {
    ddos {
      global_threshold     = 2000
      burst_threshold      = 1000
    }
    
    api {
      api_urls = [
        "https://api.example.com/v1",
        "https://api.example.com/v2"
      ]
    }
  }

  api_discovery_settings {
    description_file_location = "https://api.example.com/v1/openapi.json"
    description_file_scan_enabled = true
    description_file_scan_interval_hours = 24
    traffic_scan_enabled = true
    traffic_scan_interval_hours = 6
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the domain.

### Optional

- `api_discovery_settings` (Block List, Max: 1) (see [below for nested schema](#nestedblock--api_discovery_settings))
- `settings` (Block List, Max: 1) (see [below for nested schema](#nestedblock--settings))
- `status` (String) Status of the domain. It must be one of these values {active, monitor}.

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--api_discovery_settings"></a>
### Nested Schema for `api_discovery_settings`

Required:

- `description_file_location` (String) The URL of the API description file. This will be periodically scanned if `description_file_scan_enabled` is enabled. Supported formats are YAML and JSON, and it must adhere to OpenAPI versions 2, 3, or 3.1.

Optional:

- `description_file_scan_enabled` (Boolean) Indicates if periodic scan of the description file is enabled.
- `description_file_scan_interval_hours` (Number) The interval in hours for scanning the description file.
- `traffic_scan_enabled` (Boolean) Indicates if traffic scan is enabled.
- `traffic_scan_interval_hours` (Number) The interval in hours for scanning the traffic.


<a id="nestedblock--settings"></a>
### Nested Schema for `settings`

Optional:

- `api` (Block List, Max: 1) (see [below for nested schema](#nestedblock--settings--api))
- `ddos` (Block List, Max: 1) (see [below for nested schema](#nestedblock--settings--ddos))

<a id="nestedblock--settings--api"></a>
### Nested Schema for `settings.api`

Optional:

- `api_urls` (Set of String) List of API URL patterns.
- `is_api` (Boolean) Indicates if the domain is an API domain. All requests to an API domain are treated as API requests. If this is set to true then the api_urls field is ignored.


<a id="nestedblock--settings--ddos"></a>
### Nested Schema for `settings.ddos`

Optional:

- `burst_threshold` (Number) Burst threshold for DDoS protection.
- `global_threshold` (Number) Global threshold for DDoS protection.
