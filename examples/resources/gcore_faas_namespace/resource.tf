provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_faas_namespace" "ns" {
        project_id = 1
        region_id = 1
        name = "testns"
        description = "test description"
        envs = {
            BIG_ENV = "EXAMPLE"
        }
}
