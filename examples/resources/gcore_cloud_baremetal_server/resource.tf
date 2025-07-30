resource "gcore_cloud_baremetal_server" "example_cloud_baremetal_server" {
  project_id = 1
  region_id = 1
  flavor = "bm2-hf-medium"
  interfaces = [{
    type = "external"
    interface_name = "eth0"
    ip_family = "ipv4"
    port_group = 0
  }]
  app_config = {

  }
  apptemplate_id = "apptemplate_id"
  ddos_profile = {
    profile_template = 123
    fields = [{
      base_field = 10
      field_name = "field_name"
      field_value = [45046, 45047]
      value = "value"
    }]
  }
  image_id = "image_id"
  name = "my-bare-metal"
  name_template = "name_template"
  password = "password"
  ssh_key_name = "my-ssh-key"
  tags = {
    my-tag = "my-tag-value"
  }
  user_data = "user_data"
  username = "username"
}
