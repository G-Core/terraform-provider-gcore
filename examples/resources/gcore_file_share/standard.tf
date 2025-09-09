provider "gcore" {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

data "gcore_project" "project" {
  name = "Default"
}

data "gcore_region" "region" {
  name = "Luxembourg-2"
}

resource "gcore_file_share" "file_share_standard" {
  name       = "tf-file-share-standard"
  size       = 20
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
  type_name  = "standard"
  protocol   = "NFS"

  network {
    network_id = "378ba73d-16c5-4a4e-a755-d9406dd73e63"
  }

  access {
    ip_address  = "10.95.129.0/24"
    access_mode = "rw"
  }
}
