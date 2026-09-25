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
