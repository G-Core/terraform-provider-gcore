variable "s3_access_key" {
  type      = string
  sensitive = true
}

variable "s3_secret_key" {
  type      = string
  sensitive = true
}

# Create an origin group with S3 authentication
resource "gcore_cdn_origin_group" "s3_origin" {
  name      = "s3-origin-group"
  auth_type = "awsSignatureV4"

  auth = {
    s3_type                 = "amazon"
    s3_access_key_id        = var.s3_access_key
    s3_secret_access_key    = var.s3_secret_key
    s3_bucket_name          = "my-bucket"
    s3_region               = "eu-west-1"
    s3_credentials_version  = 1
  }
}
