# Create an MX record
resource "gcore_dns_zone_rrset" "mx_record" {
  zone_name  = "example.com"
  rrset_name = "example.com"
  rrset_type = "MX"
  ttl        = 300

  resource_records = [{
    content = [jsonencode("10 mail.example.com.")]
    enabled = true
  }]
}
