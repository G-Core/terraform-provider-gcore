terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# Test from Jira comment - reproduces "Provider produced inconsistent result after apply"
resource "gcore_cdn_origin_group" "example_cdn_origin_group" {
  name                = "test-cdn-og-gcloud2-22597-v7"
  use_next            = true
  proxy_next_upstream = ["error", "timeout", "invalid_header", "http_500", "http_502"]
  auth_type           = "awsSignatureV4"
  auth = {
    s3_credentials_version = 1
    s3_storage_hostname    = "s3.amazonaws.com"
    s3_access_key_id       = "AKIAIOSFODNN7EXAMPLE"
    s3_bucket_name         = "my-bucket"
    s3_type                = "other"
    s3_secret_access_key   = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
    s3_region              = "us-east-1"
  }
}
