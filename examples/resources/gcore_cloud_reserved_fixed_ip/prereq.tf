# Configure project and region
variable "project_id" {
  type        = number
  description = "Gcore project ID"
  default     = 1
}

variable "region_id" {
  type        = number
  description = "Gcore region ID"
  default     = 1
}
