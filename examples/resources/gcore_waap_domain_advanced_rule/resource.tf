# Block an IP address for 5 minutes. Exactly one action can be set per rule:
# allow = {}, block = {...}, captcha = {}, handshake = {}, monitor = {} or tag = {...}.
resource "gcore_waap_domain_advanced_rule" "example_waap_domain_advanced_rule" {
  domain_id   = 1
  name        = "Block foobar bot"
  description = "Block requests from a single IP address"
  enabled     = true
  source      = "request.ip == '117.20.32.55'"
  phase       = "access"
  action = {
    block = {
      status_code     = 403
      action_duration = "5m"
    }
  }
}
