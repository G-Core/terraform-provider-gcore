resource "gcore_waap_domain_custom_rule" "example_waap_domain_custom_rule" {
  domain_id = 1
  action = {
    allow = {

    }
    block = {
      action_duration = "12h"
      status_code = 403
    }
    captcha = {

    }
    handshake = {

    }
    monitor = {

    }
    tag = {
      tags = ["string"]
    }
  }
  conditions = [{
    content_type = {
      content_type = ["application/xml"]
      negation = true
    }
    country = {
      country_code = ["Mv"]
      negation = true
    }
    file_extension = {
      file_extension = ["pdf"]
      negation = true
    }
    header = {
      header = "Origin"
      value = "value"
      match_type = "Exact"
      negation = true
    }
    header_exists = {
      header = "Origin"
      negation = true
    }
    http_method = {
      http_method = "CONNECT"
      negation = true
    }
    ip = {
      ip_address = "ip_address"
      negation = true
    }
    ip_range = {
      lower_bound = "lower_bound"
      upper_bound = "upper_bound"
      negation = true
    }
    organization = {
      organization = "UptimeRobot s.r.o"
      negation = true
    }
    owner_types = {
      negation = true
      owner_types = ["COMMERCIAL"]
    }
    request_rate = {
      path_pattern = "/"
      requests = 20
      time = 1
      http_methods = ["CONNECT"]
      ips = ["string"]
      user_defined_tag = "SQfNklznVLBBpr"
    }
    response_header = {
      header = "header"
      value = "value"
      match_type = "Exact"
      negation = true
    }
    response_header_exists = {
      header = "header"
      negation = true
    }
    session_request_count = {
      request_count = 1
      negation = true
    }
    tags = {
      tags = ["string"]
      negation = true
    }
    url = {
      url = "/wp-admin/"
      match_type = "Exact"
      negation = true
    }
    user_agent = {
      user_agent = "curl/"
      match_type = "Exact"
      negation = true
    }
    user_defined_tags = {
      tags = ["SQfNklznVLBBpr"]
      negation = true
    }
  }]
  enabled = true
  name = "Block foobar bot"
  description = "description"
}
