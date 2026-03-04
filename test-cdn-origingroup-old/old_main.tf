terraform {
  required_providers {
    gcore = {
      source  = "local/gcore/gcore"
      version = "1.0.0"
    }
  }
}

variable "gcore_api_key" {
  type      = string
  sensitive = true
}

provider "gcore" {
  permanent_api_token = var.gcore_api_key
}

# Test 1: Public Origins
resource "gcore_cdn_origingroup" "public_origins" {
  name     = "test-public-origins-old"
  use_next = true

  origin {
    source  = "google.com"
    enabled = true
    backup  = false
  }

  origin {
    source  = "cloudflare.com"
    enabled = true
    backup  = true
  }
}

# Test 2: Gcore S3 Storage Auth
resource "gcore_cdn_origingroup" "gcore_s3" {
  name = "test-gcore-s3-old"

  auth {
    s3_type              = "other"
    s3_storage_hostname  = "s-ed1.cloud.gcore.lu"
    s3_access_key_id     = "8JFFS2NK771RDN7D5QB6"
    s3_secret_access_key = "PVzuij5qBpmavBJInB4eCtwMtkBEvo1k37ubj7z1"
    s3_bucket_name       = "cdn-test-bucket"
  }
}

# Outputs
output "public_origins_id" {
  value = gcore_cdn_origingroup.public_origins.id
}

output "gcore_s3_id" {
  value = gcore_cdn_origingroup.gcore_s3.id
}
