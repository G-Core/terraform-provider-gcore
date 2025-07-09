# G-Core IAM User Resource Example
#
# IMPORTANT: Due to API limitations, the user creation process works in two steps:
#
# 1. First "terraform apply":
#    - Creates user with basic fields: email, name, lang, user_role
#    - Ignores: phone, company, auth_types (not supported by invite API)
#    - User will be created but phone/company/auth_types will be empty/default
#
# 2. Second "terraform apply" (same config):
#    - Detects differences between desired config and actual state
#    - Updates phone, company, auth_types using PATCH API
#    - User will have all specified fields
#
# This behavior is intentional per G-Core API design.

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
