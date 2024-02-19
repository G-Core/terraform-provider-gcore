variable "dns_nameservers" {
  type = list(string)
  default = ["8.8.4.4", "1.1.1.1"]
}

variable "host_routes" {
  type = list(object({
    destination = string
    nexthop = string
   }))
  default = [
    {
      destination = "10.0.3.0/24"
      nexthop = "10.0.0.13"
    },
    {
      destination = "10.0.4.0/24"
      nexthop = "10.0.0.14"
    },
  ]
}