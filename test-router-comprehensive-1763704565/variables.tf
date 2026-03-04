variable "project_id" {
  type = number
}

variable "region_id" {
  type = number
}

variable "router_name" {
  default = "test-router"
}

variable "interfaces" {
  type    = list(string)
  default = []
}

variable "routes" {
  type = list(object({
    destination = string
    nexthop     = string
  }))
  default = []
}

variable "enable_external_gateway" {
  default = false
}

variable "external_gateway_snat" {
  default = true
}
