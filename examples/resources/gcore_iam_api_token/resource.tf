resource "gcore_iam_api_token" "example_iam_api_token" {
  client_id = 0
  client_user = {
    role = {
      id = 1
      name = "Administrators"
    }
  }
  exp_date = "2021-01-01 12:00:00+00:00"
  name = "My token"
  description = "It\'s my token"
}
