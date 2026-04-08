provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_cdn_origingroup" "origin_group_1" {
  name     = "origin_group_1"
  use_next = true
  origin {
    source  = "example.com"
    enabled = true
  }
  origin {
    source  = "mirror.example.com"
    enabled = true
    backup  = true
  }
}

resource "gcore_cdn_origingroup" "amazon_s3_origin_group" {
  name = "amazon_s3_origin_group"
  auth {
    s3_type              = "amazon"
    s3_access_key_id     = "123*******************"
    s3_secret_access_key = "123*******************"
    s3_bucket_name       = "bucket-name"
    s3_region            = "eu-south-2"
  }
}

resource "gcore_cdn_origingroup" "other_s3_origin_group" {
  name = "other_s3_origin_group"
  auth {
    s3_type              = "other"
    s3_storage_hostname  = "s3.example.com"
    s3_access_key_id     = "123*******************"
    s3_secret_access_key = "123*******************"
    s3_bucket_name       = "bucket-name"
  }
}

# Mixed origin group with host and S3 origins
resource "gcore_cdn_origingroup" "mixed_origin_group" {
  name     = "mixed_origin_group"
  use_next = true

  origin {
    source               = "cdn.example.com"
    enabled              = true
    host_header_override = "origin.example.com"
  }

  origin {
    origin_type = "s3"
    enabled     = true
    backup      = true
    config {
      s3_type              = "amazon"
      s3_bucket_name       = "my-bucket"
      s3_region            = "eu-west-1"
      s3_access_key_id     = "123*******************"
      s3_secret_access_key = "123*******************"
    }
  }
}

# S3-only origin group using the new origin block syntax
resource "gcore_cdn_origingroup" "s3_origin_group_new" {
  name     = "s3_origin_group_new"
  use_next = true

  origin {
    origin_type          = "s3"
    enabled              = true
    host_header_override = "storage.example.com"
    config {
      s3_type              = "other"
      s3_storage_hostname  = "s3.example.com"
      s3_bucket_name       = "my-bucket"
      s3_access_key_id     = "123*******************"
      s3_secret_access_key = "123*******************"
    }
  }
}
