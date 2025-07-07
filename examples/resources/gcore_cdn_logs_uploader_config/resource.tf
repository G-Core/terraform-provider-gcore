provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_cdn_logs_uploader_policy" "policy_1" {
  name = "Main policy"
  fields = ["remote_addr", "remote_user", "time_local"]
}

resource "gcore_cdn_logs_uploader_target" "target_1" {
  name = "Main target"
  config {
    ftp {
      hostname = "ftp.example.com"
      user = "user"
      password = "password"
    }
  }
}


resource "gcore_cdn_logs_uploader_config" "config_1" {
  name = "Logs uploader config"
  policy = gcore_cdn_logs_uploader_policy.policy_1.id
  target = gcore_cdn_logs_uploader_target.target_1.id
  for_all_resources = true
}
