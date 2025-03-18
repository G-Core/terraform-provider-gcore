terraform {
  required_providers {
    gcore = {
      source = "local.gcore.com/repo/gcore"
    }
  }
}

provider gcore {
  permanent_api_token = "asd123"
  gcore_cdn_api       = "http://localhost:8000"
}

resource "gcore_cdn_rule_template" "cdn_example_com_rule_template_1" {
  name        = "All PNG images template"
  rule        = "/folder/images/*.png"
  rule_type   = 0
  weight      = 1
  override_origin_protocol = "HTTPS"

  options {
    edge_cache_settings {
      default = "14d"
    }
    gzip_on {
      value = true
    }
    ignore_query_string {
      value = true
    }
  }
}

resource "gcore_cdn_logs_uploader_policy" "example_policy_1" {
  name = "Logs uploader policy"
  description = "Policy for logs uploader"
  # rotate_threshold_mb = 50
  field_separator = ";"
  fields = ["remote_addr", "remote_user"]
}

resource "gcore_cdn_logs_uploader_policy" "example_policy_2" {
  name = "Logs uploader policy"
  description = "Policy for logs uploader 2"
  # rotate_threshold_mb = 10
}

resource "gcore_cdn_logs_uploader_target" "example_target_1" {
  name = "Logs uploader target"
  description = "Target for logs uploader"
  # policy_id = gcore_cdn_logs_uploader_policy.example_policy_1.id
  storage_type = "sftp"
  config {
    sftp {
      hostname = "ftp.example.com"
      user = "user3"
      password = "password"
    }
    # s3_amazon {
    #   access_key_id = "access_key_id"
    #   secret_access_key = "secret_access"
    #   bucket_name = "bucket_name"
    #   region = "region"
    # }
  }
}
