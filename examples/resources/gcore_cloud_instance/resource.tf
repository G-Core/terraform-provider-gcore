resource "gcore_cloud_instance" "example_cloud_instance" {
  project_id = 1
  region_id = 1
  flavor = "g2-standard-32-64"
  interfaces = [{
    type = "external"
    interface_name = "eth0"
    ip_family = "ipv4"
    security_groups = [{
      id = "ae74714c-c380-48b4-87f8-758d656cdad6"
    }]
  }]
  volumes = [{
    image_id = "e460e48c-6836-447e-bc9c-16fc4225d318"
    source = "image"
    attachment_tag = "boot"
    boot_index = 0
    delete_on_termination = false
    name = "boot-volume"
    size = 50
    tags = {
      my-tag = "my-tag-value"
    }
    type_name = "ssd_hiiops"
  }]
  allow_app_ports = true
  configuration = {
    foo = "bar"
  }
  name = "my-instance"
  name_template = "name_template"
  password = "password"
  security_groups = [{
    id = "ae74714c-c380-48b4-87f8-758d656cdad6"
  }]
  servergroup_id = "servergroup_id"
  ssh_key_name = "my-ssh-key"
  tags = {
    my-tag = "my-tag-value"
  }
  user_data = "user_data"
  username = "username"
}
