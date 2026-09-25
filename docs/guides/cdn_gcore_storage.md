---
page_title: "Serve CDN content from and upload CDN logs to Gcore Object Storage"
subcategory: ""
description: |-
  Bind a CDN origin group and a CDN logs uploader target to a Gcore Object Storage bucket.
---

# Serve CDN content from and upload CDN logs to Gcore Object Storage

A CDN origin group and a CDN logs uploader target can be bound to a Gcore Object Storage and one of its buckets.
You only reference the storage and the bucket: the CDN fills in the endpoint, region and credentials itself.

* In `gcore_cdn_origingroup`, use an `origin` with `origin_type = "s3"` and a `config` with `s3_type = "gcore"`, `storage_id` and `s3_bucket_name`.
  The CDN gets read-only access to the bucket.
* In `gcore_cdn_logs_uploader_target`, use a `config.s3_gcore` block with `storage_id` and `bucket_name`.
  The CDN gets write access to the bucket.

With a storage binding, the provider rejects the fields that the CDN fills in:

* origin group: `s3_access_key_id`, `s3_secret_access_key`, `s3_region`, `s3_storage_hostname` and the origin's `host_header_override`;
* logs uploader target: `access_key_id`, `secret_access_key`, `region`, `endpoint` and `use_path_style`.

## Requirements

* The Gcore Object Storage integration must be enabled for your account. Contact support if the API reports that you cannot use it.
* Only Standard storages are supported.
* The storage must belong to the same account as the CDN resources.
* The bucket must already exist. The CDN never creates buckets, so reference `gcore_storage_s3_bucket` resources to create them first.

If any enabled origin in an origin group is bound to a Gcore Object Storage, the CDN pulls from all origins of that group over HTTPS.

## Example

```terraform
provider "gcore" {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_storage_s3" "cdn" {
  name     = "cdn-storage"
  location = "s-region-1"
}

resource "gcore_storage_s3_bucket" "content" {
  storage_id = gcore_storage_s3.cdn.storage_id
  name       = "cdn-content"
}

resource "gcore_storage_s3_bucket" "logs" {
  storage_id = gcore_storage_s3.cdn.storage_id
  name       = "cdn-logs"
}

resource "gcore_cdn_origingroup" "storage" {
  name     = "storage-origin"
  use_next = true

  origin {
    origin_type = "s3"
    enabled     = true
    config {
      s3_type        = "gcore"
      storage_id     = gcore_storage_s3.cdn.storage_id
      s3_bucket_name = gcore_storage_s3_bucket.content.name
    }
  }
}

resource "gcore_cdn_resource" "cdn" {
  cname        = "cdn.example.com"
  origin_group = gcore_cdn_origingroup.storage.id
}

resource "gcore_cdn_logs_uploader_policy" "cdn" {
  name   = "storage-logs-policy"
  fields = ["remote_addr", "time_local", "request", "status"]
}

resource "gcore_cdn_logs_uploader_target" "storage" {
  name = "storage-logs-target"
  config {
    s3_gcore {
      storage_id  = gcore_storage_s3.cdn.storage_id
      bucket_name = gcore_storage_s3_bucket.logs.name
      directory   = "cdn"
    }
  }
}

resource "gcore_cdn_logs_uploader_config" "cdn" {
  name      = "storage-logs"
  policy    = gcore_cdn_logs_uploader_policy.cdn.id
  target    = gcore_cdn_logs_uploader_target.storage.id
  resources = [gcore_cdn_resource.cdn.id]
}
```
