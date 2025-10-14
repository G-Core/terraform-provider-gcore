provider "gcore" {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

data "gcore_project" "project" {
  name = "Default"
}

data "gcore_region" "region" {
  name = "Sines-2"
}

resource "gcore_file_share" "file_share_vast" {
  name       = "tf-file-share-vast"
  size       = 10
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
  type_name  = "vast"
  protocol   = "NFS"
  share_settings {
    allowed_characters = "LCD"
    path_length = "LCD"
    root_squash = true
  }
}
