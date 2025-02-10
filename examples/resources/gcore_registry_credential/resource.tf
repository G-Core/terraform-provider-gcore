provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

data "gcore_project" "project" {
  name = "Default"
}

resource "gcore_registry_credential" "creds" {
  name = "docker-io"
  username = "username"
  password = "passwd"
  registry_url = "docker.io"
}