terraform {
  required_providers {
    gcore = {
      source  = "G-Core/gcore"
      version = "~> 0.10"
    }
  }
}

provider "gcore" {}

resource "gcore_fastedge_secret" "test_secret" {
  name    = "terraform_jira_test_secret"
  comment = "Test secret for GCLOUD2-23104 - old provider"

  slot {
    id    = 0
    value = "old-provider-secret-value-slot0"
  }
}
