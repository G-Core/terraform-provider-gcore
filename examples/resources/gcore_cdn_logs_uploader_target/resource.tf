provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}


resource "gcore_cdn_logs_uploader_target" "target_1" {
  name = "Logs uploader target"
  description = "Target for logs uploader"
  config {
    sftp {
      hostname = "ftp.example.com"
      user = "user"
      password = "password"
    }
  }
}

resource "gcore_cdn_logs_uploader_target" "target_2" {
  name = "Target 2"
  config {
    s3_oss {
      access_key_id = "access_key_id"
      secret_access_key = "secret_access123456789"
      bucket_name = "bucket_name"
      region = "region"
      directory = "directory"
    }
  }
}
