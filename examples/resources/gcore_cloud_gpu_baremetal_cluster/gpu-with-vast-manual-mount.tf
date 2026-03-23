# GPU bare metal cluster with VAST file share manual mount
resource "gcore_cloud_file_share" "vast" {
  project_id = 1
  region_id  = 1
  name       = "tf-file-share-vast"
  size       = 10
  type_name  = "vast"
  protocol   = "NFS"
  share_settings = {
    allowed_characters = "LCD"
    path_length        = "LCD"
    root_squash        = true
  }
}

resource "gcore_cloud_gpu_baremetal_cluster" "gpu_cluster" {
  project_id    = 1
  region_id     = 1
  flavor        = "bm3-ai-ndp2-1xlarge-h100-80-8"
  image_id      = "234c133c-b37e-4744-8a26-dc32fe407066"
  name          = "cluster-for-vast"
  servers_count = 1

  servers_settings = {
    interfaces = [
      {
        type = "external"
      },
      {
        # This interface is required to ensure that the AI cluster
        # is connected to the same network as the file share.
        # Without it, mounting the NFS share will fail.
        type       = "subnet"
        network_id = gcore_cloud_file_share.vast.network_id
        subnet_id  = gcore_cloud_file_share.vast.subnet_id
      },
    ]
    credentials = {
      ssh_key_name = "my-keypair"
    }
    file_shares = [{
      id         = gcore_cloud_file_share.vast.id
      mount_path = "/mnt/vast"
    }]
    user_data = base64encode(<<-EOT
      #cloud-config
      runcmd:
        - mkdir -p /mnt/vast
        - apt-get update -y
        - apt-get install -y nfs-common
        - mount -o vers=3,nconnect=56,remoteports=dns,spread_reads,spread_writes,noextend ${data.gcore_cloud_file_share.vast.connection_point} /mnt/vast
    EOT
    )
  }

  tags = {
    my-tag-key = "my-tag-value"
  }
}
