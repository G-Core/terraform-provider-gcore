resource "gcore_cdn_logs_uploader_policy" "example_cdn_logs_uploader_policy" {
  date_format = "[02/Jan/2006:15:04:05 -0700]"
  description = "New policy"
  escape_special_characters = true
  field_delimiter = ","
  field_separator = ";"
  fields = ["remote_addr", "status"]
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
