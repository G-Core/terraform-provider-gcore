terraform {
  required_providers {
    gcore = {
      source  = "gcore/gcore"
      version = "99.0.0"
    }
  }
}

provider "gcore" {}

resource "gcore_fastedge_secret" "test_secret" {
  name    = "terraform_jira_test_secret_v2"
  comment = "Test secret for GCLOUD2-23104 - simplified"

  secret_slots = [{
    slot  = 0
    value = "my-secret-value"
  }]
}
