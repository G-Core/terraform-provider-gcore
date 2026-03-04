variable "project_id" {
  description = "Gcore project ID"
  type        = number
  default     = 379987
}

variable "region_id" {
  description = "Gcore region ID"
  type        = number
  default     = 76
}

variable "instance_name" {
  description = "Instance name"
  type        = string
  default     = "test-resize-fix"
}

variable "flavor_id" {
  description = "Instance flavor ID"
  type        = string
  default     = "g1-standard-1-2"  # Start with small flavor
}

variable "volume_size" {
  description = "Boot volume size in GiB"
  type        = number
  default     = 10  # Start with 10GB
}
