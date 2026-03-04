variable "project_name" {
  description = "Project name"
  type        = string
  default     = "default"
}

variable "region_id" {
  description = "Region to deploy in"
  type        = number
  default     = 76
}

variable "lb_name" {
  description = "Load balancer name"
  type        = string
  default     = "test-lb-udp-source"
}

variable "lb_flavor" {
  description = "Flavor for the UDP load balancer"
  type        = string
  default     = "lb1-1-2"
}

variable "vip_network_cidr" {
  description = "CIDR for VIP network"
  type        = string
  default     = "10.0.52.0/24"
}

variable "backend_subnet_cidr" {
  description = "CIDR for backend subnet"
  type        = string
  default     = "10.0.53.0/24"
}

variable "udp_listener_port" {
  description = "Frontend UDP port"
  type        = number
  default     = 5353
}

variable "backend_port" {
  description = "Backend UDP port"
  type        = number
  default     = 5353
}

variable "persistence_timeout" {
  description = "SOURCE_IP session persistence timeout"
  type        = number
  default     = 60
}

variable "timeout_client_data" {
  description = "Client idle timeout"
  type        = number
  default     = 20000
}

variable "timeout_member_connect" {
  description = "Backend connect timeout"
  type        = number
  default     = 4000
}

variable "timeout_member_data" {
  description = "Backend idle timeout"
  type        = number
  default     = 30000
}

variable "lb_tags" {
  description = "Tags applied to the UDP LB"
  type        = map(string)
  default     = {
    "environment" = "qa"
    "scenario"    = "udp-source"
  }
}
