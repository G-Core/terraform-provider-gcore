terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

variable "zone_name" {
  default = "test-rrset-simplified-1765284111.dev"
}

variable "ttl" {
  default = 300
}

variable "ip_address" {
  default = "1.2.3.4"
}

# Test: Create A record
resource "gcore_dns_zone_rrset" "test_a" {
  zone_name  = var.zone_name
  rrset_name = "www.${var.zone_name}"
  rrset_type = "A"
  ttl        = var.ttl
  
  resource_records = [
    {
      content = [jsonencode(var.ip_address)]
      enabled = true
    }
  ]
}

# Test: Create TXT record
resource "gcore_dns_zone_rrset" "test_txt" {
  zone_name  = var.zone_name
  rrset_name = "txt.${var.zone_name}"
  rrset_type = "TXT"
  ttl        = 600
  
  resource_records = [
    {
      content = [jsonencode("v=spf1 include:_spf.google.com ~all")]
      enabled = true
    }
  ]
}

output "a_record" {
  value = {
    zone_name  = gcore_dns_zone_rrset.test_a.zone_name
    rrset_name = gcore_dns_zone_rrset.test_a.rrset_name
    rrset_type = gcore_dns_zone_rrset.test_a.rrset_type
    ttl        = gcore_dns_zone_rrset.test_a.ttl
    name       = gcore_dns_zone_rrset.test_a.name
    type       = gcore_dns_zone_rrset.test_a.type
  }
}

output "txt_record" {
  value = {
    zone_name  = gcore_dns_zone_rrset.test_txt.zone_name
    rrset_name = gcore_dns_zone_rrset.test_txt.rrset_name
    rrset_type = gcore_dns_zone_rrset.test_txt.rrset_type
    ttl        = gcore_dns_zone_rrset.test_txt.ttl
  }
}
