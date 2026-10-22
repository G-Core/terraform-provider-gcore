resource "gcore_waap_domain_custom_rule" "example_waap_domain_custom_rule" {
  domain_id = 1
  action = {
    allow = {
      foo = "bar"
    }
    block = {
      action_duration = "12h"
      status_code = 403
    }
    captcha = {
      foo = "bar"
    }
    handshake = {
      foo = "bar"
    }
    monitor = {
      foo = "bar"
    }
    tag = {
      tags = ["string"]
    }
  }
  conditions = [{
    asn = {
      asns = [1]
      negation = true
    }
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
    ja3 = {
      ja3_fingerprints = ["e7d705a3286e19ea42f587b344ee6865"]
      negation = true
    }
    ja4 = {
      ja4_fingerprints = ["t13d3113h2_e8f1e7e78f70_ce5650b735ce"]
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
