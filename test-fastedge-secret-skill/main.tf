terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

resource "gcore_fastedge_secret" "test" {
  name    = var.name
  comment = var.comment

  secret_slots = var.secret_slots
}

variable "name" {
  default = "tf-skill-test-secret"
}

variable "comment" {
  default = "Comprehensive skill test"
}

variable "secret_slots" {
  type = list(object({
    slot  = number
    value = optional(string)
  }))
  default = [
    {
      slot  = 0
      value = "secret-value-slot0"
    }
  ]
  sensitive = true
}

output "secret_id" {
  value = gcore_fastedge_secret.test.id
}

output "secret_name" {
  value = gcore_fastedge_secret.test.name
}

output "app_count" {
  value = gcore_fastedge_secret.test.app_count
}
