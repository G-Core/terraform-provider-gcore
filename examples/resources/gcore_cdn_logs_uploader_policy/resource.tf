resource "gcore_cdn_logs_uploader_policy" "example_cdn_logs_uploader_policy" {
  date_format = "[02/Jan/2006:15:04:05 -0700]"
  description = "New policy"
  escape_special_characters = true
  field_conversions = {
    request_time = {
      conversions = [{
        config = {
          factor = 1000
          precision = 0
          rounding = "nearest"
        }
        type = "scale"
      }]
    }
    upstream_cache_status = {
      conversions = [{
        config = {
          values = {
            HIT = "cached"
            MISS = "uncached"
          }
          default = "other"
        }
        type = "replace"
      }]
    }
  }
  field_delimiter = ","
  field_remap = {

  }
  field_separator = ";"
  fields = ["remote_addr", "request_time", "upstream_cache_status"]
  file_name_template = "{{YYYY}}_{{MM}}_{{DD}}_{{HH}}_{{mm}}_{{ss}}_access.log.gz"
  format_type = "json"
  include_empty_logs = true
  include_shield_logs = true
  log_sample_rate = 1
  name = "Policy"
  retry_interval_minutes = 32
  rotate_interval_minutes = 32
  rotate_threshold_lines = 5000
  rotate_threshold_mb = 252
  tags = {

  }
}
