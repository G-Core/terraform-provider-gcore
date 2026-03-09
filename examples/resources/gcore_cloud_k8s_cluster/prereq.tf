terraform {
  required_providers {
    gcore = {
      source  = "gcore/gcore"
      version = ">= 0.1"
    }
  }
}

provider "gcore" {
  api_key = "your-api-key"
}

resource "gcore_cloud_ssh_key" "my_keypair" {
  project_id = 1
  name       = "my-keypair"
  public_key = "ssh-ed25519 ...your public key... gcore@gcore.com"
}
