provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

data "gcore_project" "pr" {
  name = "test"
}

data "gcore_region" "rg" {
  name = "ED-10 Preprod"
}

data "gcore_lblistener" "l" {
  name            = "test-listener"
  loadbalancer_id = "59b2eabc-c0a8-4545-8081-979bd963c6ab" //optional
  region_id       = data.gcore_region.rg.id
  project_id      = data.gcore_project.pr.id
}

output "view" {
  value = data.gcore_lblistener.l
}
