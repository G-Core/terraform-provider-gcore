provider gcore {
  permanent_api_token = "768660$.............a43f91f"
}

resource "gcore_cdn_resource" "cdn_resource" {
  cname  = "api.example.com"
  origin = "origin.example.com"
  options {
    waap { value = true }
  }
}

resource "gcore_waap_domain" "domain" {
  name   = gcore_cdn_resource.cdn_resource.cname
  status = "monitor"
}

resource "gcore_waap_api_path" "api_path" {
  domain_id = gcore_waap_domain.domain.id
  path = "/v1/items/{id}"
  method = "GET"
  http_scheme = "HTTPS"
  api_version = "v1"
  tags = ["tag1", "tag2"]
  api_groups = ["group1", "group2"]
}
