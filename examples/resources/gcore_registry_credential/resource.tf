resource "gcore_registry_credential" "creds" {
        project_id = 184550
        name = "docker-io"
        username = "username"
        password = "passwd"
        registry_url = "docker.io"
}