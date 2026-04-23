resource "gcore_storage_ssh_key" "example_storage_ssh_key" {
  name = "my-production-key"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAA... user@example.com"
}
