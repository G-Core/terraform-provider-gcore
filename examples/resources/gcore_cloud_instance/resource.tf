# Create a boot volume from an image
resource "gcore_cloud_volume" "boot" {
  project_id = 1
  region_id  = 1
  name       = "boot-volume"
  source     = "image"
  image_id   = "your-image-id"
  size       = 20
  type_name  = "ssd_hiiops"
  tags = {
    my-tag = "my-tag-value"
  }
}

# Create an instance with the existing volume
resource "gcore_cloud_instance" "example_cloud_instance" {
  project_id = 1
  region_id  = 1
  name       = "my-instance"
  flavor     = "g2-standard-32-64"

  # Attach existing volumes by ID (only volume_id is required)
  volumes = [{ volume_id = gcore_cloud_volume.boot.id }]

  interfaces = [{
    type           = "external"
    interface_name = "eth0"
    ip_family      = "ipv4"
    security_groups = [{
      id = "ae74714c-c380-48b4-87f8-758d656cdad6"
    }]
  }]

  security_groups = [{
    id = "ae74714c-c380-48b4-87f8-758d656cdad6"
  }]

  ssh_key_name = "my-ssh-key"
  user_data    = "user_data"
  username     = "username"
  password     = "password"

  tags = {
    my-tag = "my-tag-value"
  }
}
