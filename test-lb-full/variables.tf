variable "lb_name" {
  default = "qa-lb-test"
}

variable "lb_flavor" {
  default = "lb1-2-4"
}

variable "listener_name" {
  default = "qa-listener"
}

variable "listener_protocol" {
  default = "HTTP"
}

variable "listener_port" {
  default = 80
}

variable "timeout_client_data" {
  default = null
}

variable "timeout_member_connect" {
  default = null
}

variable "timeout_member_data" {
  default = null
}

variable "create_pool" {
  default = true
}

variable "pool_name" {
  default = "qa-pool"
}

variable "pool_algorithm" {
  default = "ROUND_ROBIN"
}

variable "pool_protocol" {
  default = "HTTP"
}

variable "pool_healthmonitor" {
  type = object({
    type             = string
    delay            = number
    max_retries      = number
    timeout          = number
    http_method      = optional(string)
    max_retries_down = optional(number)
    url_path         = optional(string)
    expected_codes   = optional(string)
  })
  default = null
}

