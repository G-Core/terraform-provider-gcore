resource "gcore_cloud_gpu_baremetal_cluster" "example_cloud_gpu_baremetal_cluster" {
  project_id = 0
  region_id = 0
  flavor = "bm3-ai-1xlarge-h100-80-8"
  image_id = "f01fd9a0-9548-48ba-82dc-a8c8b2d6f2f1"
  interfaces = [{
    network_id = "024a29e9-b4b7-4c91-9a46-505be123d9f8"
    subnet_id = "91200a6c-07e0-42aa-98da-32d1f6545ae7"
    type = "subnet"
    floating_ip = {
      source = "new"
    }
    interface_name = "interface_name"
  }]
  name = "my-gpu-cluster"
  instances_count = 1
  password = "password"
  security_groups = [{
    id = "ae74714c-c380-48b4-87f8-758d656cdad6"
  }]
  ssh_key_name = "my-ssh-key"
  tags = {
    my-tag = "my-tag-value"
  }
  user_data = "user_data"
  username = "username"
}
