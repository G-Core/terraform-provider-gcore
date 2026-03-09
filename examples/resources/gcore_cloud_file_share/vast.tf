resource "gcore_cloud_file_share" "file_share_vast" {
  project_id = 1
  region_id  = 1

  name      = "tf-file-share-vast"
  size      = 10
  type_name = "vast"
  protocol  = "NFS"

  share_settings = {
    allowed_characters = "LCD"
    path_length        = "LCD"
    root_squash        = true
  }
}
