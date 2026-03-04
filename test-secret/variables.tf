variable "secret_name" {
  type    = string
  default = "tf-test-secret"
}

variable "secret_comment" {
  type    = string
  default = null
}

variable "secret_slots" {
  type = list(object({
    slot  = number
    value = optional(string)
  }))
  default = null
}

variable "secret_slots_wo_version" {
  type    = number
  default = null
}
