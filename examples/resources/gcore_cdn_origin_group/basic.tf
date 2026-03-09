# Create a CDN origin group with primary and backup origins
resource "gcore_cdn_origin_group" "example" {
  name     = "my-origin-group"
  use_next = true

  sources = [
    {
      source  = "example.com"
      enabled = true
      backup  = false
    },
    {
      source  = "mirror.example.com"
      enabled = true
      backup  = true
    },
  ]
}
