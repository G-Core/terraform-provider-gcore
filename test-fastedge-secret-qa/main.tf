terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

resource "gcore_fastedge_secret" "test" {
  name    = "terraform-qa-test-secret"
  comment = "QA test for GCLOUD2-23104"
  secret_slots = [
    {
      slot  = 0
      value = "my-secret-value-0"
    },
    {
      slot  = 1
      value = "my-secret-value-1"
    }
  ]
}

output "secret_id" {
  value = gcore_fastedge_secret.test.id
}
