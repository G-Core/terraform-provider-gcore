---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "gcore_cdn_sslcert Resource - terraform-provider-gcorelabs"
subcategory: ""
description: |-
  
---

# gcore_cdn_sslcert (Resource)



## Example Usage

```terraform
provider gcore {
  user_name = "test"
  password = "test"
  gcore_platform = "https://api.gcdn.co"
  gcore_cdn_api = "https://api.gcdn.co"
}

variable "cert" {
  type = string
  sensitive = true
}

variable "private_key" {
  type = string
  sensitive = true
}

resource "gcore_cdn_sslcert" "cdnopt_cert" {
  name = "Test cert for cdnopt_bookatest_by"
  cert = var.cert
  private_key = var.private_key
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **cert** (String, Sensitive) The public part of the SSL certificate. All chain of the SSL certificate should be added.
- **name** (String) Name of the SSL certificate. Must be unique.
- **private_key** (String, Sensitive) The private key of the SSL certificate.

### Optional

- **id** (String) The ID of this resource.

### Read-Only

- **automated** (Boolean) The way SSL certificate was issued.
- **has_related_resources** (Boolean) It shows if the SSL certificate is used by a CDN resource.

