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

resource "gcore_fastedge_binary" "test_binary" {
  filename = "test.wasm"
}

resource "gcore_fastedge_app" "app_from_binary" {
  status = "enabled"
  name = "terraform-test1"
  comment = "Terraform test app 1"
  binary = gcore_fastedge_binary.test_binary.id
  env = {
    "foo" = "bar"
  }
}

resource "gcore_fastedge_template" "test_template" {
  name = "terraform_test_template"
  binary = gcore_fastedge_binary.test_binary.id
  short_descr = "short description"
  long_descr = "long description"
  param {
      name = "foo"
      type = "string"
      mandatory = true
      descr = "Parameter foo"
  }
  param {
      name = "bar"
      type = "number"
      descr = "Parameter bar"
  }
}

resource "gcore_fastedge_app" "app_from_template" {
  status = "enabled"
  name = "terraform-test2"
  comment = "Terraform test app 2"
  template = gcore_fastedge_template.test_template.id
  env = {
    "foo" = "foo_value"
    "bar" = 123
  }
  secrets = {
    "baz" = gcore_fastedge_secret.test_secret.id
  }
}
