provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_cdn_rule_template" "cdn_example_com_rule_template_1" {
  name        = "All PNG images template"
  rule        = "/folder/images/*.png"
  rule_type   = 0
  weight      = 1
  override_origin_protocol = "HTTPS"

  options {
    edge_cache_settings {
      default = "14d"
    }
    gzip_on {
      value = true
    }
    ignore_query_string {
      value = true
    }
  }
}
