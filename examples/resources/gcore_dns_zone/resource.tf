provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

# Basic DNS Zone example - minimal configuration
resource "gcore_dns_zone" "example_zone" {
  name = "example_zone.com"
}

# Advanced DNS Zone example - showcasing all available options
resource "gcore_dns_zone" "advanced_zone" {
  name    = "advanced-example.com"
  dnssec  = true
  enabled = true
  
  # SOA record fields
  contact        = "admin@advanced-example.com"
  expiry         = 604800    # 1 week
  nx_ttl         = 3600      # 1 hour
  primary_server = "ns1.advanced-example.com."
  refresh        = 3600      # 1 hour
  retry          = 1800      # 30 minutes
  serial         = 2024010100 # YYYYMMDDNN format
  
  # Meta configuration with webhook
  meta = {
    webhook_url    = "https://hooks.example.com/dns-changes"
    webhook_method = "POST"
    environment    = "production"
    managed_by     = "terraform"
  }
}

# DNS Zone with DNSSEC disabled (explicit)
resource "gcore_dns_zone" "simple_zone" {
  name   = "simple-example.org"
  dnssec = false
  
  # Basic SOA configuration
  contact = "hostmaster@simple-example.org"
  expiry  = 1209600  # 2 weeks
  refresh = 7200     # 2 hours
  retry   = 3600     # 1 hour
}

# DNS Zone for development/testing (disabled)
resource "gcore_dns_zone" "test_zone" {
  name    = "test-example.dev"
  enabled = false  # Zone disabled, records won't resolve
  
  contact = "devops@test-example.dev"
  
  # Custom meta for testing
  meta = {
    environment = "testing"
    auto_delete = "true"
    owner       = "development-team"
  }
}
