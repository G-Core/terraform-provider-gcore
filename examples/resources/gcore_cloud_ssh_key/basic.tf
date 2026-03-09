# Upload an existing SSH public key
resource "gcore_cloud_ssh_key" "deployer" {
  project_id = 1
  name       = "deployer-key"
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIJe8rDJP... deployer@workstation"
}
