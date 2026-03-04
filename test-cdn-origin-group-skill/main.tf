terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# =============================================================================
# TC-1: Basic origin group with sources (public origins)
# =============================================================================
resource "gcore_cdn_origin_group" "tc1_public" {
  count = var.run_tc1 ? 1 : 0

  name                = "test-cdn-og-tc1-public-${var.test_suffix}"
  use_next            = true
  proxy_next_upstream = ["error", "timeout", "invalid_header", "http_500", "http_502"]
  auth_type           = "none"

  sources = [
    {
      source  = "93.184.216.34"
      enabled = true
      backup  = false
    },
    {
      source  = "93.184.216.35:8080"
      enabled = true
      backup  = true
    }
  ]
}

# =============================================================================
# TC-2: S3 auth with s3_type=other (BUG FIX VERIFICATION)
# This tests that s3_credentials_version persists in state
# =============================================================================
resource "gcore_cdn_origin_group" "tc2_s3_other" {
  count = var.run_tc2 ? 1 : 0

  name      = "test-cdn-og-tc2-s3other-${var.test_suffix}"
  auth_type = "awsSignatureV4"
  use_next  = false

  auth = {
    s3_credentials_version = var.tc2_credentials_version
    s3_storage_hostname    = "s3.example.com"
    s3_access_key_id       = "AKIAIOSFODNN7EXAMPLE"
    s3_bucket_name         = "my-test-bucket"
    s3_type                = "other"
    s3_secret_access_key   = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
  }
}

# =============================================================================
# TC-3: S3 auth with s3_type=amazon
# =============================================================================
resource "gcore_cdn_origin_group" "tc3_s3_amazon" {
  count = var.run_tc3 ? 1 : 0

  name      = "test-cdn-og-tc3-s3amazon-${var.test_suffix}"
  auth_type = "awsSignatureV4"
  use_next  = true
  proxy_next_upstream = ["error", "timeout"]

  auth = {
    s3_credentials_version = 1
    s3_access_key_id       = "AKIAIOSFODNN7EXAMPLE"
    s3_bucket_name         = "my-amazon-bucket"
    s3_type                = "amazon"
    s3_secret_access_key   = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYzzzzzzzzzz"
    s3_region              = "us-east-1"
  }
}

# =============================================================================
# TC-4: Update test - will modify sources after initial create
# =============================================================================
resource "gcore_cdn_origin_group" "tc4_update" {
  count = var.run_tc4 ? 1 : 0

  name                = var.tc4_name
  use_next            = var.tc4_use_next
  proxy_next_upstream = var.tc4_proxy_next_upstream
  auth_type           = "none"

  sources = var.tc4_sources
}

# =============================================================================
# Variables
# =============================================================================
variable "test_suffix" {
  description = "Suffix for unique resource names"
  default     = "skill"
}

variable "run_tc1" {
  default = false
}

variable "run_tc2" {
  default = false
}

variable "run_tc3" {
  default = false
}

variable "run_tc4" {
  default = false
}

variable "tc2_credentials_version" {
  default = 1
}

variable "tc4_name" {
  default = "test-cdn-og-tc4-update-skill"
}

variable "tc4_use_next" {
  default = true
}

variable "tc4_proxy_next_upstream" {
  default = ["error", "timeout"]
}

variable "tc4_sources" {
  default = [
    {
      source  = "93.184.216.40"
      enabled = true
      backup  = false
    }
  ]
}

# =============================================================================
# Outputs
# =============================================================================
output "tc1_id" {
  value = var.run_tc1 ? gcore_cdn_origin_group.tc1_public[0].id : null
}

output "tc1_state" {
  value = var.run_tc1 ? gcore_cdn_origin_group.tc1_public[0] : null
}

output "tc2_id" {
  value = var.run_tc2 ? gcore_cdn_origin_group.tc2_s3_other[0].id : null
}

output "tc2_auth" {
  value = var.run_tc2 ? gcore_cdn_origin_group.tc2_s3_other[0].auth : null
  sensitive = true
}

output "tc3_id" {
  value = var.run_tc3 ? gcore_cdn_origin_group.tc3_s3_amazon[0].id : null
}

output "tc3_auth" {
  value = var.run_tc3 ? gcore_cdn_origin_group.tc3_s3_amazon[0].auth : null
  sensitive = true
}

output "tc4_id" {
  value = var.run_tc4 ? gcore_cdn_origin_group.tc4_update[0].id : null
}
