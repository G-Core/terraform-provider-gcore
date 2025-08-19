provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_cdn_logs_uploader_policy" "policy_1" {
  name = "Main policy"
  rotate_threshold_lines = 10
  field_separator = ";"
  fields = ["remote_addr", "remote_user", "time_local", "custom_field_1", "custom_field_2"]
  include_empty_logs = true
  file_name_template = "{{CNAME}}.{{HOST}}.{{TIMESTAMP}}.log"
  tags = {
    custom_field_1 = "value1"
    custom_field_2 = "value2"
  }
}
