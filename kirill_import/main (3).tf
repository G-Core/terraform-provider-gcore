# resource "random_password" "vm_password" {
#   length           = 16
#   special          = true
#   override_special = "!#$%&*()-_=+[]{}<>:?"
# }


######## MANUAL CHECKS ########

data "gcore_cloud_projects" "my_projects" {
  name = "default"
}

locals {
  project_id = [for p in data.gcore_cloud_projects.my_projects.items : p.id]
}

data "gcore_cloud_region" "rg" {
  region_id = 76
}

# resource "gcore_cloud_floating_ip" "fip" {
#   project_id = local.project_id[0]
#   region_id  = data.gcore_cloud_region.rg.id
# }

# Volume
# resource "gcore_cloud_volume" "boot_vol" {
#   project_id = local.project_id[0]
#   region_id  = data.gcore_cloud_region.rg.id
#   image_id = "f84ddba3-7a5a-4199-931a-250e981d16fb"
#   name = "bt_vol"
#   size = 5
#   source = "image"
#   attachment_tag = "device-tag"
#   tags = {
#     foo = "my-tag-value"
#   }
#   type_name = "standard"
# }

# resource "gcore_cloud_instance" "ext_instance" {
#   project_id = local.project_id[0]
#   region_id  = data.gcore_cloud_region.rg.id
#   flavor   = "g3a-standard-1-2"
#   name = "instance-luxembourg-2-01"
#   volumes = [{ 
#     boot_index = 0
#     volume_id = gcore_cloud_volume.boot_vol.id
#     }]
#   interfaces = [{
#     type      = "external"
#   },
#   {
#     type = "subnet"
#     network_id = gcore_cloud_network.prvt_nw.id
#     subnet_id = gcore_cloud_network_subnet.sb_sys.id
#     floating_ip = {
#       source = "new"
#       # source = "existing"
#       # existing_floating_id = gcore_cloud_floating_ip.fip.id
#     } 
#   }]
# }


# resource "gcore_cloud_network" "prvt_nw" {
#   project_id = local.project_id[0]
#   region_id  = data.gcore_cloud_region.rg.id
#   name = "my network"
#   create_router = true
#   tags = {
#     my-tag = "my-tag-value"
#   }
#   type = "vxlan"
# }

# resource "gcore_cloud_network_subnet" "sb_sys" {
#   project_id = local.project_id[0]
#   region_id  = data.gcore_cloud_region.rg.id
#   cidr = "192.168.10.0/24"
#   name = "sys"
#   network_id = gcore_cloud_network.prvt_nw.id
#   connect_to_network_router = true
#   enable_dhcp = true
#   gateway_ip = "192.168.10.1"
#   ip_version = 4
#   tags = {
#     my-tag = "my-tag-value"
#   }
# }

# # resource "gcore_cloud_floating_ip" "fip" {
# #   project_id = local.project_id[0]
# #   region_id  = data.gcore_cloud_region.rg.id
# # }

# resource "gcore_cloud_instance" "test" {
#   project_id = local.project_id[0]
#   region_id  = data.gcore_cloud_region.rg.id

#   flavor   = "g1-standard-2-4"
#   name     = "qa-vm-tf-renamed"
#   volumes = [{ volume_id = gcore_cloud_volume.boot_vol.id }]

#   interfaces = [{
#     type      = "external"
#   },
#   {
#     type = "any_subnet"
#     network_id = gcore_cloud_network.prvt_nw.id
#     security_groups = [{
#       id = "9c916e52-cf6f-4790-8a94-715b8a0c8fd3"
#     }]
#     floating_ip = {
#       source = "new"
#       # source = "existing"
#       # existing_floating_id = gcore_cloud_floating_ip.fip.id
#     } 
#   }]
#   security_groups = [{
#     id = "07409229-a210-4771-bf9a-42b20bc1886d"
#   }]

#   # interfaces = [{
#   #   type      = "external"
#   # }]
#   ssh_key_name = "qa-prod-tk-def"
#   servergroup_id = "c5366929-1ec4-40e7-b5d3-faf55f48d8ae"
#   user_data = base64encode(<<-EOT
#   #cloud-config
#   runcmd:
#     - mkdir -p /mnt/vast
#     - apt-get update -y
#     - apt-get install -y nginx
#   EOT
#     )
# }

########## IMPORTING CHECKS ############

resource "gcore_cloud_instance" "qa_imp" {
  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id
  name      = "qa-imp-vm"
  flavor    = "g3a-standard-1-2"
  volumes    =   [{ volume_id = "d94da554-4a7c-40ab-bc30-9677dd32aa61" },
                  { volume_id = "00969da0-50d7-427f-a00b-c537351b8d82" }]

  interfaces =  [{ type      = "external"}]

}
