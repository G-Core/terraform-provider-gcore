provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

data "gcore_project" "pr" {
  name = "test"
}

data "gcore_region" "rg" {
  name = "ED-10 Preprod"
}

data "gcore_floatingip" "ip" {
  floating_ip_address = "10.100.179.172"
  region_id           = data.gcore_region.rg.id
  project_id          = data.gcore_project.pr.id
}

data "gcore_floatingip" "ip_by_fixed_ip" {
  fixed_ip_address    = "192.168.42.12"
  region_id           = data.gcore_region.rg.id
  project_id          = data.gcore_project.pr.id
}

output "view" {
  value = data.gcore_floatingip.ip
}
