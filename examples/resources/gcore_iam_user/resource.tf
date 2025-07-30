provider "gcore" {
  permanent_api_token = "your-api-token-here"
  api_endpoint        = "https://api.gcore.com"
}

# Complete user example with all available fields
resource "gcore_iam_user" "example_user" {
  email     = "user@example.com"
  name      = "Example User"
  lang      = "en"
  phone     = "+1234567890"  # Set on second apply
  company   = "Example Corp" # Set on second apply
  client_id = 12345          # Replace with your client ID

  user_role {
    id   = 2
    name = "Users"
  }

  auth_types = ["password", "sso"] # Set on second apply

}
