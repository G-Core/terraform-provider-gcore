resource "gcore_reservedfixedip" "fixed_ip_external" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  type       = "external"

  is_vip     = false
}
