data "gcore_project" "multi_sg_project" {
  name = "Default"
}

data "gcore_region" "multi_sg_region" {
  name = "Luxembourg-2"
}

# Get the default security group
data "gcore_securitygroup" "multi_sg_default" {
  name       = "default"
  project_id = data.gcore_project.multi_sg_project.id
  region_id  = data.gcore_region.multi_sg_region.id
}

# Create a custom security group for GPU workloads
resource "gcore_securitygroup" "gpu_sg" {
  name        = "gpu-workload-sg"
  description = "Security group for GPU cluster workloads"

  # Allow SSH access
  security_group_rules {
    direction        = "ingress"
    ethertype        = "IPv4"
    protocol         = "tcp"
    port_range_min   = 22
    port_range_max   = 22
    remote_ip_prefix = "0.0.0.0/0"
  }

  # Allow Jupyter/JupyterLab access
  security_group_rules {
    direction        = "ingress"
    ethertype        = "IPv4"
    protocol         = "tcp"
    port_range_min   = 8888
    port_range_max   = 8888
    remote_ip_prefix = "0.0.0.0/0"
  }

  # Allow HTTPS outbound
  security_group_rules {
    direction        = "egress"
    ethertype        = "IPv4"
    protocol         = "tcp"
    port_range_min   = 443
    port_range_max   = 443
    remote_ip_prefix = "0.0.0.0/0"
  }

  project_id = data.gcore_project.multi_sg_project.id
  region_id  = data.gcore_region.multi_sg_region.id
}

# GPU cluster with multiple security groups
resource "gcore_ai_cluster" "multi_sg_gpu_cluster" {
  flavor          = "bm3-ai-ndp-1xlarge-h100-80-8"         # Updated GPU flavor for region 164
  image_id        = "a0e6e6f2-1d23-4841-81cc-fc0038f1ccb9" # Latest Ubuntu 24.04 GPU image
  cluster_name    = "multi-sg-gpu-cluster"
  keypair_name    = "my-keypair"
  instances_count = 1

  interface {
    type = "external"
  }

  # Multiple security groups: input uses ID field
  security_group {
    id = data.gcore_securitygroup.multi_sg_default.id
  }

  security_group {
    id = gcore_securitygroup.gpu_sg.id
  }

  cluster_metadata = {
    environment = "production"
    workload    = "gpu-ml"
    security    = "multi-sg-demo"
  }

  project_id = data.gcore_project.multi_sg_project.id
  region_id  = data.gcore_region.multi_sg_region.id
}

# Data source to read the cluster back (output will have security group names)
data "gcore_ai_cluster" "multi_sg_cluster_read" {
  cluster_id = gcore_ai_cluster.multi_sg_gpu_cluster.id
  project_id = data.gcore_project.multi_sg_project.id
  region_id  = data.gcore_region.multi_sg_region.id
}

# Outputs showing the security group structure
output "multi_sg_cluster_info" {
  description = "Cluster information"
  value = {
    cluster_id = gcore_ai_cluster.multi_sg_gpu_cluster.id
    status     = gcore_ai_cluster.multi_sg_gpu_cluster.cluster_status
  }
}

output "multi_sg_security_groups_input" {
  description = "Security groups as configured (with IDs)"
  value       = gcore_ai_cluster.multi_sg_gpu_cluster.security_group
}

output "multi_sg_security_groups_output" {
  description = "Security groups from data source (with names)"
  value       = data.gcore_ai_cluster.multi_sg_cluster_read.security_group
}

output "multi_sg_poplar_servers_security_groups" {
  description = "Security groups per server (with names)"
  value       = data.gcore_ai_cluster.multi_sg_cluster_read.poplar_servers[*].security_groups
}
