resource "gcore_k8sv2" "cluster" {
  project_id    = data.gcore_project.project.id
  region_id     = data.gcore_region.region.id
  name          = "my-k8s-cluster"
  keypair       = gcore_keypair.my_keypair.sshkey_name
  version       = "v1.31.9"
  pool {
    name             = "my-k8s-pool"
    flavor_id        = "g1-standard-2-4"
    servergroup_policy = "soft-anti-affinity"
    min_node_count   = 1
    max_node_count   = 1
    boot_volume_size = 10
    boot_volume_type = "standard"
    is_public_ipv4 = true
  }
  ddos_profile {
    enabled = true
    fields {
      base_field = 1353
      field_value = jsonencode(["AF"])
    }
    fields {
      base_field = 1354
      field_value = jsonencode(50)
    }
    fields {
      base_field = 1355
      field_value = jsonencode(150)
    }
    fields {
      base_field = 1356
      field_value = jsonencode(300)
    }
    fields {
      base_field = 1357
      field_value = jsonencode(300)
    }

    fields {
      base_field = 1352
      field_value = jsonencode([
        {
          "sip_list":["192.168.0.1","10.10.0.1"],
          "dport_list": ["27015","27025"],
          "proto_list": ["udp"],
          "sport_list": ["27025"],
          "policy": "DROP"
        }
      ])
    }
    profile_template = 1128
  }
}

data "gcore_k8sv2_kubeconfig" "config" {
  cluster_name       = gcore_k8sv2.cluster.name
  region_id          = data.gcore_region.region.id
  project_id         = data.gcore_project.project.id
}

// to store kubeconfig in a file pls use
// terraform output -raw kubeconfig > config.yaml
output "kubeconfig" {
  value = data.gcore_k8sv2_kubeconfig.config.kubeconfig
}
