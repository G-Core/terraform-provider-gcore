terraform {
  required_version = ">= 1.11" # write-only arguments
  required_providers {
    gcore = { source = "G-Core/gcore" }
  }
}

variable "sftp_password" {
  type      = string
  sensitive = true
  ephemeral = true
}

resource "gcore_storage_sftp" "example_storage_sftp" {
  location_name          = "s-region-1"
  name                   = "my-sftp-storage"
  password_wo            = var.sftp_password
  password_wo_version    = 1
  expires                = "2 years 6 months"
  has_custom_config_file = false
  is_http_disabled       = false
  server_alias           = "my-storage.example.com"
  ssh_key_ids            = [1, 2, 3]
}
