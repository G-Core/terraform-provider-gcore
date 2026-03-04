terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# ============================================================================
# Test 1: Public Origins - Basic creation with multiple sources
# ============================================================================
resource "gcore_cdn_origin_group" "public_origins" {
  name      = var.public_origins_name
  auth_type = "none"
  use_next  = var.use_next

  sources = var.sources
}

# ============================================================================
# Test 2: S3 Auth - Gcore S3 storage with awsSignatureV4
# ============================================================================
resource "gcore_cdn_origin_group" "s3_auth" {
  count = var.test_s3_auth ? 1 : 0

  name      = "test-s3-auth-skill"
  auth_type = "awsSignatureV4"

  auth = {
    s3_type              = "other"
    s3_storage_hostname  = "s-ed1.cloud.gcore.lu"
    s3_access_key_id     = var.s3_access_key_id
    s3_secret_access_key = var.s3_secret_access_key
    s3_bucket_name       = "cdn-test-bucket"
  }
}

# ============================================================================
# Variables
# ============================================================================
variable "public_origins_name" {
  type    = string
  default = "test-public-origins-skill"
}

variable "use_next" {
  type    = bool
  default = true
}

variable "sources" {
  type = list(object({
    source  = string
    enabled = bool
    backup  = bool
  }))
  default = [
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

variable "test_s3_auth" {
  type    = bool
  default = false
}

variable "s3_access_key_id" {
  type      = string
  default   = "TEST_ACCESS_KEY_ID"
  sensitive = true
}

variable "s3_secret_access_key" {
  type      = string
  default   = "TEST_SECRET_ACCESS_KEY_12345678"
  sensitive = true
}

# ============================================================================
# Outputs
# ============================================================================
output "public_origins_id" {
  value = gcore_cdn_origin_group.public_origins.id
}

output "public_origins_name" {
  value = gcore_cdn_origin_group.public_origins.name
}

output "public_origins_path" {
  value = gcore_cdn_origin_group.public_origins.path
}

output "public_origins_has_related_resources" {
  value = gcore_cdn_origin_group.public_origins.has_related_resources
}

output "public_origins_proxy_next_upstream" {
  value = gcore_cdn_origin_group.public_origins.proxy_next_upstream
}

output "s3_auth_id" {
  value = var.test_s3_auth ? gcore_cdn_origin_group.s3_auth[0].id : null
}
