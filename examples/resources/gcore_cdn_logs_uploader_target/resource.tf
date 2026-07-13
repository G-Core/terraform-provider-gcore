resource "gcore_cdn_logs_uploader_target" "example_cdn_logs_uploader_target" {
  config = {
    access_key_id = "access_key_id"
    bucket_name = "bucket_name"
    directory = "directory"
    endpoint = "endpoint"
    region = "region"
    secret_access_key = "secret_access_key"
    use_path_style = true
  }
  storage_type = "s3_gcore"
  description = "description"
  name = "name"
}
