# Release Notes Examples

Real examples from this repository showing the target format. Use these as
style references when generating human-readable release notes.

Every release starts with the **disclaimer** block, followed by a
`## Release notes` heading, then the **Part 1 human-readable summary**,
followed by the **Part 2 auto-generated changelog** (which comes from the PR
body and is not shown in these examples).

## Example 1: Multi-product release with breaking changes

```markdown
> [!warning]
> v2 is a ground-up rewrite of the provider, featuring OpenAPI-spec-driven
> code generation and a move to terraform-plugin-framework under the hood.
>
> This is an **alpha** release, and **breaking changes** are expected.

If you'd like to try it out, pin the provider version **exactly** to
v2.0.0-alpha.2 in your Terraform configuration.

```hcl
terraform {
  required_providers {
    gcore = {
      source  = "G-Core/gcore"
      version = "2.0.0-alpha.2"
    }
  }
}
```

## Release notes

We're excited to announce version 2.0.0-alpha.2!

### **CDN**

* **`gcore_cdn_resource`**
  * Fixed drift detection for computed fields on `secondary_hostnames` and `options`
  * Fixed update serialization — resolved issues where unchanged nested blocks were sent in PATCH requests

* **`gcore_cdn_resource_rule`**
  * ⚠ BREAKING CHANGE: Renamed resource from `gcore_cdn_rule` to `gcore_cdn_resource_rule` to avoid naming conflict with the CDN resource `rule_overrides` attribute

### **Cloud**

* **`gcore_cloud_load_balancer`**
  * ⚠ BREAKING CHANGE: Removed inline `listeners` attribute — use the dedicated `gcore_cloud_load_balancer_listener` resource instead
  * Removed deprecated timeout fields (`timeout_client_data`, `timeout_member_connect`, `timeout_member_data`)

* **`gcore_cloud_network`**
  * ⚠ BREAKING CHANGE: Removed `create_router` attribute — use `gcore_cloud_network_router` resource instead
  * Updated acceptance tests and sweeper

* **`gcore_cloud_security_group`**
  * Migrated to v2 API endpoints for create and update operations

* **`gcore_cloud_security_group_rule`**
  * Added full CRUD support with v2 async endpoints
  * Added acceptance tests

### **DNS**

* **`gcore_dns_zone`**
  * Added `meta` attribute support with custom `MetaStringType` for flexible metadata handling
  * Added DNSSEC configuration support
  * Added import support

* **`gcore_dns_zone_rrset`**
  * Added full CRUD and import support
  * Added `MetaStringType` for meta fields and update support

### **FastEdge**

* **`gcore_fastedge_app`**
  * Fixed state drift — resolved inconsistencies between plan and apply for `template` and `binary` attributes
  * Added validation: exactly one of `template` or `binary` must be specified
  * Switched to PUT instead of PATCH for update operations

* **`gcore_fastedge_binary`**
  * Added file upload support for binary resources

* **`gcore_fastedge_secret`**
  * Added sensitive field handling and force delete support

### **Other**

* Fixed spurious update plans for float attributes after import
* Improved data source timeout schema handling
```

## Example 2: Smaller release with features and fixes

```markdown
> [!warning]
> v2 is a ground-up rewrite of the provider, featuring OpenAPI-spec-driven
> code generation and a move to terraform-plugin-framework under the hood.
>
> This is an **alpha** release, and **breaking changes** are expected.

If you'd like to try it out, pin the provider version **exactly** to
v2.0.0-alpha.3 in your Terraform configuration.

```hcl
terraform {
  required_providers {
    gcore = {
      source  = "G-Core/gcore"
      version = "2.0.0-alpha.3"
    }
  }
}
```

## Release notes

We're excited to announce version 2.0.0-alpha.3!

### **Cloud**

* **`gcore_cloud_floating_ip`**
  * Migrated to v2 API endpoint for updates
  * Added `UpdateAndPoll` support for async floating IP operations

* **`gcore_cloud_instance`**
  * Made `password` a write-only field — no longer stored in state
  * Fixed floating IP assign/unassign to use `UpdateAndPoll`

* **`gcore_cloud_secret`**
  * Converted `payload` fields to write-only attributes — sensitive data no longer stored in Terraform state

### **FastEdge**

* **`gcore_fastedge_app`**
  * Removed read-only `name` from `app_store` required fields — fixes creation errors when store name is server-assigned
```

## Example 3: Release with new resources and data sources

```markdown
> [!warning]
> v2 is a ground-up rewrite of the provider, featuring OpenAPI-spec-driven
> code generation and a move to terraform-plugin-framework under the hood.
>
> This is an **alpha** release, and **breaking changes** are expected.

If you'd like to try it out, pin the provider version **exactly** to
v2.0.0-alpha.4 in your Terraform configuration.

```hcl
terraform {
  required_providers {
    gcore = {
      source  = "G-Core/gcore"
      version = "2.0.0-alpha.4"
    }
  }
}
```

## Release notes

We're excited to announce version 2.0.0-alpha.4!

### **CDN**

* **`gcore_cdn_origin_group`**
  * New resource — manage CDN origin groups with full CRUD support

* **`gcore_cdn_trusted_ca_certificate`**
  * New resource — upload and manage trusted CA certificates for CDN

### **Cloud**

* **`gcore_cloud_file_share`**
  * New resource — manage NFS/CIFS file shares with full lifecycle support
  * Made `name` required in schema

* **`gcore_cloud_file_share_access_rule`**
  * New resource — manage access rules for file shares

* **`gcore_cloud_gpu_virtual_cluster`**
  * New resource — manage GPU virtual clusters with drift prevention and import support

* **`gcore_cloud_k8s_cluster`**
  * ⚠ BREAKING CHANGE: Renamed all `k8` references to `k8s` — attribute names and resource internals updated for consistency
  * Removed `ddos_profile` from resource and data source schemas
  * Added full lifecycle management with pool operations

### **WAAP**

* **`gcore_waap_domain`**
  * New resource — manage WAAP domain protection with full CRUD support
```

## Style Rules (inferred from examples)

1. **Disclaimer**: Always starts the release notes. Uses GitHub `> [!warning]`
   callout syntax. Includes the version number in both prose and HCL block.
2. **Separator**: A `## Release notes` heading separates the disclaimer from
   the human-readable summary.
3. **Opening line**: Always `We're excited to announce version {VERSION}!`
   (immediately after the `## Release notes` heading).
4. **Product area headers**: `### **{Area}**` — bold inside h3.
5. **Sub-area items**: `* **\`gcore_{resource}\`**` — bold + backtick, followed
   by indented child bullets.
6. **Separate sub-areas** for resources and data sources (e.g.,
   `gcore_cloud_network` and `gcore_cloud_networks` are distinct sub-areas).
7. **Breaking changes**: Inline with `⚠ BREAKING CHANGE:` prefix, describe
   what was removed/changed and why.
8. **Deprecations**: ``Deprecated `{attribute}` attribute — use
   `{alternative}` instead``
9. **New resources**: `New resource — {description}`
10. **New attributes**: ``Added `{attribute}` attribute — {description}``
11. **Fixes**: `Fixed {description} — {detail}`
12. **Attribute names**: Always snake_case Terraform names in backticks (e.g.,
    `create_router`, `admin_state_up`) — what users see in `.tf` files.
13. **Descriptions**: Short, specific, user-actionable. No commit hashes in
    Part 1.
14. **Product area order**: Alphabetical. Place "Other" last.
15. **Sub-area order**: Alphabetical by resource name within each product area.
