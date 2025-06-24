variable "app_secret_slot_0" {
  description = "Slot 0 secret value"
  type        = string
  sensitive   = true
  default     = ""
}

variable "app_secret_slot_1" {
  description = "Slot 1 secret value"
  type        = string
  sensitive   = true
  default     = ""
}

resource "gcore_fastedge_secret" "test_secret" {
  name = "terraform_test_secret"
  comment = "My test secret"
  slot {
    id = 0
    value = var.app_secret_slot_0
  }
  slot {
    id = 1
    value = var.app_secret_slot_1
  }
}
