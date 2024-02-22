provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_baremetal" "bm" {
  name       = "test bm instance"
  region_id  = 1
  project_id = 1
  flavor_id  = "bm1-infrastructure-small"
  image_id   = "1ee7ccee-5003-48c9-8ae0-d96063af75b2" // your image id

  //additional interface, available type is 'subnet' or 'external'
  //  interface {
  //	type = "subnet"
  //	network_id = "9c7867fb-f404-4a2d-8bb5-24acf2fccaf1" //your network_id
  //	subnet_id = "b68ea6e2-c2b6-4a8d-95eb-7194d12a2156" // your subnet_id
  //  }

  //  interface {
  //	type = "external"
  //    is_parent = "true" // if is_parent = true interface cant be detached, and always connected first
  //  }

  keypair_name = "test" // your keypair name
}
