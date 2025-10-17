resource "gcore_file_share" "file_share_vast" {
  name       = "tf-file-share-vast"
  size       = 1
  project_id = data.gcore_project.pr.id
  region_id  = data.gcore_region.rg.id
  type_name  = "vast"
  protocol   = "NFS"
}

variable "user_userdata" {
 description = "This is a variable of type string"
 type        = string
 default     = <<EOF
#cloud-config
runcmd:
  - mkdir -p /mount/path
  - apt-get update -y
  - apt-get install -y nfs-common
  - mount -o vers=3,nconnect=56,remoteports=dns,spread_reads,spread_writes,noextend ${resource.gcore_file_share.file_share_vast.connection_point} /mount/path
EOF
}

resource "gcore_ai_cluster" "gpu_cluster" {
  flavor = "bm3-ai-1xlarge-h200-141-8"
  image_id = "18126da2-261a-4e56-a059-82e71477bada"
  cluster_name = "my-gpu-cluster-for-vast"
  keypair_name = "qa-prod-tk-def"
  instances_count = 1

  interface {
    type = "external"
  }

  interface {
    type = "subnet"
    network_id = gcore_file_share.file_share_vast.network[0].network_id
    subnet_id = gcore_file_share.file_share_vast.network[0].subnet_id
  }

  cluster_metadata = {
    my-metadata-key = "my-metadata-value"
  }

  user_data = base64encode(var.user_userdata)

  project_id = data.gcore_project.pr.id
  region_id  = data.gcore_region.rg.id
}
