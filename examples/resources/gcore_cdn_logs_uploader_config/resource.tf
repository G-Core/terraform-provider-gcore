provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_cdn_logs_uploader_policy" "policy_1" {
  name = "Main policy"
  fields = ["remote_addr", "remote_user", "time_local"]
}

resource "gcore_cdn_logs_uploader_target" "target_1" {
  name = "Main target"
  config {
    ftp {
      hostname = "ftp.example.com"
      user = "user"
      password = "password"
    }
  }
}


resource "gcore_cdn_logs_uploader_config" "config_1" {
  name = "Logs uploader config"
  policy = gcore_cdn_logs_uploader_policy.policy_1.id
  target = gcore_cdn_logs_uploader_target.target_1.id
  for_all_resources = true
}

# Serve content from and upload logs to a Gcore Object Storage: the CDN fills in the endpoint, region and credentials
resource "gcore_storage_s3" "storage" {
  name     = "cdn-storage"
  location = "s-region-1"
}

resource "gcore_storage_s3_bucket" "content" {
  storage_id = gcore_storage_s3.storage.storage_id
  name       = "cdn-content"
}

resource "gcore_storage_s3_bucket" "logs" {
  storage_id = gcore_storage_s3.storage.storage_id
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
      storage_id     = gcore_storage_s3.storage.storage_id
      s3_bucket_name = gcore_storage_s3_bucket.content.name
    }
  }
}

resource "gcore_cdn_resource" "storage" {
  cname        = "cdn.example.com"
  origin_group = gcore_cdn_origingroup.storage.id
}

resource "gcore_cdn_logs_uploader_target" "storage" {
  name = "Storage target"
  config {
    s3_gcore {
      storage_id  = gcore_storage_s3.storage.storage_id
      bucket_name = gcore_storage_s3_bucket.logs.name
      directory   = "cdn"
    }
  }
}

resource "gcore_cdn_logs_uploader_config" "storage" {
  name      = "Storage logs uploader config"
  policy    = gcore_cdn_logs_uploader_policy.policy_1.id
  target    = gcore_cdn_logs_uploader_target.storage.id
  resources = [gcore_cdn_resource.storage.id]
}
