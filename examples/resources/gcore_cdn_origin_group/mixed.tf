# Mixed origin group with host and S3 origins
resource "gcore_cdn_origin_group" "mixed" {
  name                   = "mixed-origin-group"
  use_next               = true
  s3_credentials_version = 1

  sources = [
    {
      source  = "cdn.example.com"
      enabled = true
    },
    {
      origin_type = "s3"
      enabled     = true
      backup      = true
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
