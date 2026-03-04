terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

resource "gcore_fastedge_secret" "test" {
  name    = var.secret_name
  comment = var.secret_comment

  secret_slots = var.secret_slots

  secret_slots_wo_version = var.secret_slots_wo_version
}

output "secret_id" {
  value = gcore_fastedge_secret.test.id
}

output "secret_name" {
  value = gcore_fastedge_secret.test.name
}

output "secret_app_count" {
  value = gcore_fastedge_secret.test.app_count
}

output "secret_slots_info" {
  value     = gcore_fastedge_secret.test.secret_slots
  sensitive = true
}
