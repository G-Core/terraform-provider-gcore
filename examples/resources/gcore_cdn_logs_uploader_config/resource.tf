resource "gcore_cdn_logs_uploader_config" "example_cdn_logs_uploader_config" {
  name = "name"
  policy = 0
  target = 0
  enabled = true
  for_all_resources = true
  resources = [0]
}
