resource "gcore_cloud_k8s_cluster" "example_cloud_k8s_cluster" {
  project_id = 0
  region_id = 0
  keypair = "some_keypair"
  name = "string"
  pools = [{
    flavor_id = "g1-standard-1-2"
    min_node_count = 3
    name = "my-pool"
    auto_healing_enabled = true
    boot_volume_size = 50
    boot_volume_type = "ssd_hiiops"
    crio_config = {
      default-ulimits = "nofile=1024:2048"
    }
    is_public_ipv4 = true
    kubelet_config = {
      podMaxPids = "4096"
    }
    labels = {
      my-label = "foo"
    }
    max_node_count = 5
    servergroup_policy = "affinity"
    taints = {
      my-taint = "bar:NoSchedule"
    }
  }]
  version = "1.28.1"
  add_ons = {
    slurm = {
      enabled = true
      file_share_id = "cbc94d0e-06c6-4d12-9e86-9782ba14fc8c"
      ssh_key_ids = ["25735292-bd97-44b0-a1af-d7eab876261d", "efc01f3a-35b9-4385-89f9-e38439093ee7"]
      worker_count = 2
    }
  }
  authentication = {
    oidc = {
      client_id = "kubernetes"
      groups_claim = "groups"
      groups_prefix = "oidc:"
      issuer_url = "https://accounts.provider.example"
      required_claims = {
        claim = "value"
      }
      signing_algs = ["RS256", "RS512"]
      username_claim = "sub"
      username_prefix = "oidc:"
    }
  }
  autoscaler_config = {
    scale-down-unneeded-time = "5m"
  }
  cni = {
    cilium = {
      encryption = true
      hubble_relay = true
      hubble_ui = true
      lb_acceleration = true
      lb_mode = "snat"
      mask_size = 24
      mask_size_v6 = 120
      routing_mode = "tunnel"
      tunnel = "geneve"
    }
    cloud_k8s_cluster_provider = "cilium"
  }
  csi = {
    nfs = {
      vast_enabled = true
    }
  }
  ddos_profile = {
    enabled = true
    fields = [{
      base_field = 10
      field_value = [45046, 45047]
      value = null
    }]
    profile_template = 29
    profile_template_name = "profile_template_name"
  }
  fixed_network = "3fa85f64-5717-4562-b3fc-2c963f66afa6"
  fixed_subnet = "3fa85f64-5717-4562-b3fc-2c963f66afa6"
  is_ipv6 = true
  logging = {
    destination_region_id = 1
    enabled = true
    retention_policy = {
      period = 45
    }
    topic_name = "my-log-name"
  }
  pods_ip_pool = "172.16.0.0/18"
  pods_ipv6_pool = "2a03:90c0:88:393::/64"
  services_ip_pool = "172.24.0.0/18"
  services_ipv6_pool = "2a03:90c0:88:381::/108"
}
