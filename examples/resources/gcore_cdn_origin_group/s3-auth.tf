variable "s3_access_key" {
  type      = string
  sensitive = true
}

variable "s3_secret_key" {
  type      = string
  sensitive = true
}

# Create an origin group with an S3 origin (Amazon S3)
resource "gcore_cdn_origin_group" "s3_amazon" {
  name                   = "s3-amazon-origin-group"
  s3_credentials_version = 1

  sources = [
    {
      origin_type = "s3"
      enabled     = true
      config = {
        s3_type              = "amazon"
        s3_bucket_name       = "my-bucket"
        s3_access_key_id     = var.s3_access_key
        s3_secret_access_key = var.s3_secret_key
        s3_region            = "eu-west-1"
      }
    }
  ]
}

# Create an origin group with an S3 origin (other S3-compatible storage)
resource "gcore_cdn_origin_group" "s3_other" {
  name                   = "s3-other-origin-group"
  s3_credentials_version = 1

  sources = [
    {
      origin_type = "s3"
      enabled     = true
      config = {
        s3_type              = "other"
        s3_bucket_name       = "my-bucket"
        s3_access_key_id     = var.s3_access_key
        s3_secret_access_key = var.s3_secret_key
        s3_storage_hostname  = "s3.example.com"
      }
    }
  ]
}
