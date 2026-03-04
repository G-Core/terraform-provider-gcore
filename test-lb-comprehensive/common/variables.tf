variable "project_id" {
  description = "GCore project ID"
  type        = number
  default     = 379987
}

variable "region_id" {
  description = "GCore region ID"
  type        = number
  default     = 76
}

variable "subnet_id" {
  description = "Subnet ID for pool members"
  type        = string
  default     = "replace-with-actual-subnet-id"
}

variable "network_id" {
  description = "Network ID for load balancer VIP"
  type        = string
  default     = "replace-with-actual-network-id"
}
