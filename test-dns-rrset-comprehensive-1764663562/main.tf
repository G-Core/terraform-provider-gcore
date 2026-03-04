terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# Test A Record - with 2 IPs
resource "gcore_dns_zone_rrset" "test_a" {
  zone_name  = "maxima.lt"
  rrset_name = "tf-test-comprehensive-a.maxima.lt"
  rrset_type = "A"
  ttl        = 600

  resource_records = [
    {
      content = ["\"192.168.1.100\""]
      enabled = true
    },
    {
      content = ["\"192.168.1.101\""]
      enabled = true
    }
  ]
}

# Test AAAA Record
resource "gcore_dns_zone_rrset" "test_aaaa" {
  zone_name  = "maxima.lt"
  rrset_name = "tf-test-comprehensive-aaaa.maxima.lt"
  rrset_type = "AAAA"
  ttl        = 600

  resource_records = [
    {
      content = ["\"2001:db8::1\""]
      enabled = true
    }
  ]
}

# Test CNAME Record
resource "gcore_dns_zone_rrset" "test_cname" {
  zone_name  = "maxima.lt"
  rrset_name = "tf-test-comprehensive-cname.maxima.lt"
  rrset_type = "CNAME"
  ttl        = 600

  resource_records = [
    {
      content = ["\"maxima.lt.\""]
      enabled = true
    }
  ]
}

# Test MX Record
resource "gcore_dns_zone_rrset" "test_mx" {
  zone_name  = "maxima.lt"
  rrset_name = "tf-test-comprehensive-mx.maxima.lt"
  rrset_type = "MX"
  ttl        = 600

  resource_records = [
    {
      content = ["10", "\"mail.maxima.lt.\""]
      enabled = true
    }
  ]
}

# Test TXT Record
resource "gcore_dns_zone_rrset" "test_txt" {
  zone_name  = "maxima.lt"
  rrset_name = "tf-test-comprehensive-txt.maxima.lt"
  rrset_type = "TXT"
  ttl        = 600

  resource_records = [
    {
      content = ["\"v=spf1 include:_spf.google.com ~all\""]
      enabled = true
    }
  ]
}

output "a_record_ids" {
  value = [for r in gcore_dns_zone_rrset.test_a.resource_records : r.id]
}
