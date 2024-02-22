resource "gcore_lbmember" "public_member" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  pool_id       = gcore_lbpool.http.id

  address       = "8.8.8.8"
  protocol_port = 80
  weight        = 1
}
