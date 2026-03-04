terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

variable "zone_name" {
  default = "test-comprehensive-1765285206.dev"
}

# ============================================================
# Test Group 1: Multiple Record Types
# ============================================================

# Test 1.1: A Record - basic IP address
resource "gcore_dns_zone_rrset" "a_record" {
  zone_name  = var.zone_name
  rrset_name = "www.${var.zone_name}"
  rrset_type = "A"
  ttl        = var.a_ttl
  
  resource_records = [
    {
      content = [jsonencode(var.a_ip)]
      enabled = true
    }
  ]
}

# Test 1.2: TXT Record - SPF record
resource "gcore_dns_zone_rrset" "txt_record" {
  zone_name  = var.zone_name
  rrset_name = "spf.${var.zone_name}"
  rrset_type = "TXT"
  ttl        = 600
  
  resource_records = [
    {
      content = [jsonencode(var.txt_value)]
      enabled = true
    }
  ]
}

# Test 1.3: CNAME Record - alias
resource "gcore_dns_zone_rrset" "cname_record" {
  zone_name  = var.zone_name
  rrset_name = "alias.${var.zone_name}"
  rrset_type = "CNAME"
  ttl        = 300
  
  resource_records = [
    {
      content = [jsonencode("www.${var.zone_name}.")]
      enabled = true
    }
  ]
}

# Test 1.4: MX Record - mail server with priority
resource "gcore_dns_zone_rrset" "mx_record" {
  zone_name  = var.zone_name
  rrset_name = var.zone_name
  rrset_type = "MX"
  ttl        = 3600
  
  resource_records = [
    {
      content = [10, jsonencode("mail.${var.zone_name}.")]
      enabled = true
    }
  ]
}

# ============================================================
# Test Group 2: Edge Cases - Special Names
# ============================================================

# Test 2.1: Deep subdomain
resource "gcore_dns_zone_rrset" "deep_subdomain" {
  zone_name  = var.zone_name
  rrset_name = "level1.level2.level3.${var.zone_name}"
  rrset_type = "A"
  ttl        = 300
  
  resource_records = [
    {
      content = [jsonencode("10.0.0.1")]
      enabled = true
    }
  ]
}

# Test 2.2: Wildcard record
resource "gcore_dns_zone_rrset" "wildcard" {
  zone_name  = var.zone_name
  rrset_name = "*.${var.zone_name}"
  rrset_type = "A"
  ttl        = 300
  
  resource_records = [
    {
      content = [jsonencode("10.0.0.2")]
      enabled = true
    }
  ]
}

# Test 2.3: Multiple resource records in one rrset
resource "gcore_dns_zone_rrset" "multi_record" {
  zone_name  = var.zone_name
  rrset_name = "multi.${var.zone_name}"
  rrset_type = "A"
  ttl        = 300
  
  resource_records = [
    {
      content = [jsonencode("192.168.1.1")]
      enabled = true
    },
    {
      content = [jsonencode("192.168.1.2")]
      enabled = true
    },
    {
      content = [jsonencode("192.168.1.3")]
      enabled = var.multi_third_enabled
    }
  ]
}

# ============================================================
# Variables for Update Testing
# ============================================================

variable "a_ttl" {
  description = "TTL for A record - change to test updates"
  default     = 300
}

variable "a_ip" {
  description = "IP for A record - change to test updates"
  default     = "1.2.3.4"
}

variable "txt_value" {
  description = "TXT record value - change to test updates"
  default     = "v=spf1 include:_spf.google.com ~all"
}

variable "multi_third_enabled" {
  description = "Enable/disable third record in multi_record"
  default     = true
}

# ============================================================
# Outputs for Verification
# ============================================================

output "a_record_info" {
  value = {
    zone_name  = gcore_dns_zone_rrset.a_record.zone_name
    rrset_name = gcore_dns_zone_rrset.a_record.rrset_name
    rrset_type = gcore_dns_zone_rrset.a_record.rrset_type
    ttl        = gcore_dns_zone_rrset.a_record.ttl
    name       = gcore_dns_zone_rrset.a_record.name
    type       = gcore_dns_zone_rrset.a_record.type
  }
}

output "all_records_summary" {
  value = {
    a_name      = gcore_dns_zone_rrset.a_record.name
    txt_name    = gcore_dns_zone_rrset.txt_record.name
    cname_name  = gcore_dns_zone_rrset.cname_record.name
    mx_name     = gcore_dns_zone_rrset.mx_record.name
    deep_name   = gcore_dns_zone_rrset.deep_subdomain.name
    wild_name   = gcore_dns_zone_rrset.wildcard.name
    multi_name  = gcore_dns_zone_rrset.multi_record.name
  }
}
