# Import test resource
resource "gcore_cdn_origin_group" "imported" {
  name      = "import-test-origin"
  auth_type = "none"
  use_next  = true  # API default

  sources = [
    {
      source  = "example.org"
      enabled = true
      backup  = false
    }
  ]
}

output "imported_id" {
  value = gcore_cdn_origin_group.imported.id
}

output "imported_sources" {
  value = gcore_cdn_origin_group.imported.sources
}
