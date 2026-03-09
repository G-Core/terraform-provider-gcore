# DNS Zone with SOA configuration
resource "gcore_dns_zone" "advanced" {
  name    = "advanced-example.com"
  enabled = true

  # SOA record fields
  contact        = "admin@advanced-example.com"
  expiry         = 604800
  nx_ttl         = 3600
  primary_server = "ns1.advanced-example.com."
  refresh        = 3600
  retry          = 1800
}
