resource "gcore_cloud_file_share" "file_share_standard" {
  project_id = 1
  region_id  = 1

  name      = "tf-file-share-standard"
  size      = 20
  type_name = "standard"
  protocol  = "NFS"

  network = {
    network_id = "378ba73d-16c5-4a4e-a755-d9406dd73e63"
  }
}
