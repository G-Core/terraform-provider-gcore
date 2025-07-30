resource "gcore_cloud_file_share" "example_cloud_file_share" {
  project_id = 1
  region_id = 1
  name = "test-share-file-system"
  network = {
    network_id = "024a29e9-b4b7-4c91-9a46-505be123d9f8"
    subnet_id = "91200a6c-07e0-42aa-98da-32d1f6545ae7"
  }
  protocol = "NFS"
  size = 5
  access = [{
    access_mode = "ro"
    ip_address = "10.0.0.1"
  }]
  tags = {
    my-tag = "my-tag-value"
  }
  volume_type = "default_share_type"
}
