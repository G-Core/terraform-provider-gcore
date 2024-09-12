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
