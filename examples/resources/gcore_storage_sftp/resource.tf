resource "gcore_storage_sftp" "example_storage_sftp" {
  location_name = "s-region-1"
  name = "my-sftp-storage"
  password_mode = "auto"
  expires = "2 years 6 months"
  has_custom_config_file = false
  is_http_disabled = false
  server_alias = "my-storage.example.com"
  sftp_password = "sftp_password"
  ssh_key_ids = [1, 2, 3]
}
