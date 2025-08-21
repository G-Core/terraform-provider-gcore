provider gcore {
  permanent_api_token = "768660$.............a43f91f"
}

resource "gcore_iam_api_token" "iam_api_token" {
  client_id   = 9262944
  name        = "Test Api Token"
  description = "API token for Terraform test"
  exp_date    =  "2026-01-01T12:00:00.000000Z"
  client_user {
    role {
      id    = 2
      name  = "Users"
    }
  }
}
