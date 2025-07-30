resource "gcore_cloud_user_role_assignment" "example_cloud_user_role_assignment" {
  role = "ClientAdministrator"
  user_id = 777
  client_id = 8
  project_id = 0
}
