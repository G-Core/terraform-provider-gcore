# Block POST requests to /admin from outside the office range for one hour
resource "gcore_waap_domain_custom_rule" "block_admin" {
  domain_id   = 12345
  name        = "Block admin from outside office"
  description = "Only the office may post to /admin"
  enabled     = true

  action = {
    block = {
      status_code     = 403
      action_duration = "1h"
    }
  }

  conditions = [
    { url = { url = "/admin", match_type = "Contains" } },
    { http_method = { http_method = "POST" } },
    { ip_range = { lower_bound = "192.0.2.0", upper_bound = "192.0.2.255", negation = true } },
  ]
}
