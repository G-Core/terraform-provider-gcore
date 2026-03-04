terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# Test scenario - WITH boot_index explicitly specified
resource "gcore_cloud_instance" "qa_imp" {
  project_id = 379987
  region_id  = 76
  name       = "boot-index-drift-test"
  flavor     = "g1-standard-1-2"
  volumes    = [
    { volume_id = "86bc910b-da7b-4542-803c-1a5beed72e37", boot_index = 0 },
    { volume_id = "077d940d-1c93-4e91-be6e-4dd0e5274ced", boot_index = -1 }
  ]

  interfaces = [{ type = "external" }]
}

output "instance_id" {
  value = gcore_cloud_instance.qa_imp.id
}
