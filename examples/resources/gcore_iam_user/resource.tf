resource "gcore_iam_user" "example_iam_user" {
  user_id = 0
  auth_types = ["password"]
  company = "company"
  email = "dev@stainless.com"
  groups = [{
    id = 1
    name = "Administrators"
  }]
  lang = "de"
  name = "name"
  phone = "phone"
}
