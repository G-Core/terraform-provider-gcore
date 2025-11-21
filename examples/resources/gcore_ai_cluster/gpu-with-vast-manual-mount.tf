resource "gcore_file_share" "vast" {
  name       = "tf-file-share-vast"
  size       = 10
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
  type_name  = "vast"
  protocol   = "NFS"
  share_settings {
    allowed_characters = "LCD"
    path_length = "LCD"
    root_squash = true
  }
}

resource "gcore_ai_cluster" "gpu_cluster" {
  flavor          = "bm3-ai-1xlarge-h200-141-8"
  image_id        = "aab83c98-7c9c-4942-a488-6c8b63dd42bd"
  cluster_name    = "cluster-for-vast"
  keypair_name    = "my-keypair"
  instances_count = 1
  user_data = base64encode(<<-EOT
  #cloud-config
  runcmd:
    - mkdir -p /mnt/vast
    - apt-get update -y
    - apt-get install -y nfs-common
    - mount -o vers=3,nconnect=56,remoteports=dns,spread_reads,spread_writes,noextend ${data.gcore_file_share.vast.connection_point} /mnt/vast
  EOT
    )

  interface {
    type = "external"
  }

  // This interface is required to ensure that the AI cluster
  // is connected to the same network as the file share.
  // Without it, mounting the NFS share will fail.
  interface {
    type       = "subnet"
    network_id = gcore_file_share.vast.network[0].network_id
    subnet_id  = gcore_file_share.vast.network[0].subnet_id
  }

  cluster_metadata = {
    my-metadata-key = "my-metadata-value"
  }

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}
