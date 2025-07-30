resource "gcore_waap_domain_advanced_rule" "example_waap_domain_advanced_rule" {
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
  enabled = true
  name = "name"
  source = "request.rate_limit([], \'.*events\', 5, 200, [], [], \'\', \'ip\') and not (\'mb-web-ui\' in request.headers[\'Cookie\'] or \'mb-mobile-ios\' in request.headers[\'Cookie\'] or \'session-token\' in request.headers[\'Cookie\']) and not request.headers[\'session\']"
  description = "description"
  phase = "access"
}
