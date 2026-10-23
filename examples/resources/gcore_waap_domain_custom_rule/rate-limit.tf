# Show a captcha to clients that send more than 100 API requests in 10 seconds.
# A request_rate condition must be the only condition of its rule.
resource "gcore_waap_domain_custom_rule" "rate_limit" {
  domain_id = 12345
  name      = "API rate limit"
  enabled   = true

  action = {
    captcha = {}
  }

  conditions = [
    {
      request_rate = {
        path_pattern = "^/api/"
        requests     = 100
        time         = 10
        http_methods = ["GET", "POST"]
      }
    },
  ]
}
