terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# Test 1: Public Origins with multiple sources
resource "gcore_cdn_origin_group" "public_origins" {
  count = var.test_public_origins ? 1 : 0

  name      = "test-public-origins-new"
  auth_type = "none"
  use_next  = true

  sources = [
    {
      source  = "google.com"
      enabled = true
      backup  = false
    },
    {
      source  = "cloudflare.com"
      enabled = true
      backup  = true
    }
  ]
}

# Test 2: Gcore S3 Storage Auth (s3_type = "other")
resource "gcore_cdn_origin_group" "gcore_s3" {
  count = var.test_gcore_s3 ? 1 : 0

  name      = "test-gcore-s3-new"
  auth_type = "awsSignatureV4"

  auth = {
    s3_type              = "other"
    s3_storage_hostname  = "s-ed1.cloud.gcore.lu"
    s3_access_key_id     = var.s3_access_key_id
    s3_secret_access_key = var.s3_secret_access_key
    s3_bucket_name       = "cdn-test-bucket"
  }
}

# Test 3: Update test - modify name
resource "gcore_cdn_origin_group" "update_test" {
  count = var.test_update ? 1 : 0

  name      = var.update_test_name
  auth_type = "none"
  use_next  = var.update_test_use_next

  sources = [
    {
      source  = "github.com"
      enabled = true
      backup  = false
    }
  ]
}

# Variables
variable "test_public_origins" {
  type    = bool
  default = false
}

variable "test_gcore_s3" {
  type    = bool
  default = false
}

variable "test_update" {
  type    = bool
  default = false
}

variable "s3_access_key_id" {
  type      = string
  default   = "8JFFS2NK771RDN7D5QB6"
  sensitive = true
}

variable "s3_secret_access_key" {
  type      = string
  default   = "PVzuij5qBpmavBJInB4eCtwMtkBEvo1k37ubj7z1"
  sensitive = true
}

variable "update_test_name" {
  type    = string
  default = "test-update-origin-group"
}

variable "update_test_use_next" {
  type    = bool
  default = true
}

# Outputs
output "public_origins_id" {
  value = var.test_public_origins ? gcore_cdn_origin_group.public_origins[0].id : null
}

output "gcore_s3_id" {
  value = var.test_gcore_s3 ? gcore_cdn_origin_group.gcore_s3[0].id : null
}

output "update_test_id" {
  value = var.test_update ? gcore_cdn_origin_group.update_test[0].id : null
}
