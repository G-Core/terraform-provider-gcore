resource "gcore_reservedfixedip" "external_fixed_ip" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
  type       = "external"
}

resource "gcore_baremetal" "baremetal_with_reserved_address" {
  flavor_id     = "bm1-infrastructure-small"
  name          = "my-baremetal"
  keypair_name  = "my-keypair"
  image_id      = data.gcore_image.ubuntu.id

  interface {
    type    = "reserved_fixed_ip"
    port_id = gcore_reservedfixedip.external_fixed_ip.port_id
  }

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}