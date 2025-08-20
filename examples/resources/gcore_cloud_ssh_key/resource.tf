resource "gcore_cloud_ssh_key" "example_cloud_ssh_key" {
  project_id = 1
  name = "my-ssh-key"
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIIjxL6g1II8NsO8odvBwGKvq2Dx/h/xrvsV9b9LVIYKm my-username@my-hostname"
  shared_in_project = true
}
