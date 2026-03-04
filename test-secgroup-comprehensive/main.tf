terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

# -------------------------------------------------------
# Security group — rules managed via separate resources only
# Extended test plan covering JIRA bugs and edge cases
# -------------------------------------------------------
resource "gcore_cloud_security_group" "test" {
  name        = var.sg_name
  description = var.sg_description
}

# --- Rule 1: TCP port range (standard SSH) ---
resource "gcore_cloud_security_group_rule" "ingress_ssh" {
  count = var.create_rules ? 1 : 0

  group_id   = gcore_cloud_security_group.test.id
  project_id = gcore_cloud_security_group.test.project_id
  region_id  = gcore_cloud_security_group.test.region_id

  direction        = "ingress"
  protocol         = "tcp"
  port_range_min   = 22
  port_range_max   = 22
  ethertype        = "IPv4"
  remote_ip_prefix = "0.0.0.0/0"
  description      = "Allow SSH"
}

# --- Rule 2: UDP single port (DNS - GCLOUD2-22026 edge case) ---
resource "gcore_cloud_security_group_rule" "ingress_dns" {
  count = var.create_rules ? 1 : 0

  group_id   = gcore_cloud_security_group.test.id
  project_id = gcore_cloud_security_group.test.project_id
  region_id  = gcore_cloud_security_group.test.region_id

  direction        = "ingress"
  protocol         = "udp"
  port_range_min   = 53
  port_range_max   = 53
  ethertype        = "IPv4"
  remote_ip_prefix = "0.0.0.0/0"
  description      = "Allow DNS UDP"
}

# --- Rule 3: ICMP (no port range - protocol without ports) ---
resource "gcore_cloud_security_group_rule" "ingress_icmp" {
  count = var.create_rules ? 1 : 0

  group_id   = gcore_cloud_security_group.test.id
  project_id = gcore_cloud_security_group.test.project_id
  region_id  = gcore_cloud_security_group.test.region_id

  direction        = "ingress"
  protocol         = "icmp"
  ethertype        = "IPv4"
  remote_ip_prefix = "0.0.0.0/0"
  description      = "Allow ICMP ping"
}

# --- Rule 4: IPv6 rule (ipv6-icmp - GCLOUD2-4818 edge case) ---
resource "gcore_cloud_security_group_rule" "ingress_icmpv6" {
  count = var.create_rules ? 1 : 0

  group_id   = gcore_cloud_security_group.test.id
  project_id = gcore_cloud_security_group.test.project_id
  region_id  = gcore_cloud_security_group.test.region_id

  direction        = "ingress"
  protocol         = "ipv6-icmp"
  ethertype        = "IPv6"
  remote_ip_prefix = "::/0"
  description      = "Allow ICMPv6"
}

# --- Rule 5: Wide TCP port range (GCLOUD2-22026 edge case) ---
resource "gcore_cloud_security_group_rule" "ingress_high_ports" {
  count = var.create_rules ? 1 : 0

  group_id   = gcore_cloud_security_group.test.id
  project_id = gcore_cloud_security_group.test.project_id
  region_id  = gcore_cloud_security_group.test.region_id

  direction        = "ingress"
  protocol         = "tcp"
  port_range_min   = 10000
  port_range_max   = 20000
  ethertype        = "IPv4"
  remote_ip_prefix = "10.0.0.0/8"
  description      = "Allow high TCP ports from private range"
}

# --- Rule 6: Egress rule ---
resource "gcore_cloud_security_group_rule" "egress_https" {
  count = var.create_rules ? 1 : 0

  group_id   = gcore_cloud_security_group.test.id
  project_id = gcore_cloud_security_group.test.project_id
  region_id  = gcore_cloud_security_group.test.region_id

  direction        = "egress"
  protocol         = "tcp"
  port_range_min   = 443
  port_range_max   = 443
  ethertype        = "IPv4"
  remote_ip_prefix = "0.0.0.0/0"
  description      = "Allow HTTPS egress"
}

# --- Rule 7: SCTP protocol (uncommon protocol test) ---
resource "gcore_cloud_security_group_rule" "ingress_sctp" {
  count = var.create_rules ? 1 : 0

  group_id   = gcore_cloud_security_group.test.id
  project_id = gcore_cloud_security_group.test.project_id
  region_id  = gcore_cloud_security_group.test.region_id

  direction        = "ingress"
  protocol         = "sctp"
  port_range_min   = 5060
  port_range_max   = 5060
  ethertype        = "IPv4"
  remote_ip_prefix = "0.0.0.0/0"
  description      = "Allow SCTP signaling"
}

# --- Rule 8: Rule without description (optional field omitted) ---
resource "gcore_cloud_security_group_rule" "ingress_http_no_desc" {
  count = var.create_rules ? 1 : 0

  group_id   = gcore_cloud_security_group.test.id
  project_id = gcore_cloud_security_group.test.project_id
  region_id  = gcore_cloud_security_group.test.region_id

  direction        = "ingress"
  protocol         = "tcp"
  port_range_min   = 80
  port_range_max   = 80
  ethertype        = "IPv4"
  remote_ip_prefix = "0.0.0.0/0"
}

# --- Rule 9: Remote group self-reference (GCLOUD2-5064 edge case) ---
resource "gcore_cloud_security_group_rule" "ingress_self_ref" {
  count = var.create_rules ? 1 : 0

  group_id   = gcore_cloud_security_group.test.id
  project_id = gcore_cloud_security_group.test.project_id
  region_id  = gcore_cloud_security_group.test.region_id

  direction       = "ingress"
  protocol        = "tcp"
  port_range_min  = 1
  port_range_max  = 65535
  ethertype       = "IPv4"
  remote_group_id = gcore_cloud_security_group.test.id
  description     = "Allow all TCP from self"
}

# ---- Second SG for import test ----
resource "gcore_cloud_security_group" "import_test" {
  count       = var.create_import_sg ? 1 : 0
  name        = "tf-test-sg-import"
  description = "SG for import testing"
}

# -------------------------------------------------------
# Outputs for verification
# -------------------------------------------------------
output "sg_id" {
  value = gcore_cloud_security_group.test.id
}

output "sg_security_group_rules" {
  value = gcore_cloud_security_group.test.security_group_rules
}

output "sg_revision" {
  value = gcore_cloud_security_group.test.revision_number
}

output "rule_count" {
  value = var.create_rules ? 9 : 0
}

output "rule_ids" {
  value = var.create_rules ? {
    ssh         = gcore_cloud_security_group_rule.ingress_ssh[0].id
    dns         = gcore_cloud_security_group_rule.ingress_dns[0].id
    icmp        = gcore_cloud_security_group_rule.ingress_icmp[0].id
    icmpv6      = gcore_cloud_security_group_rule.ingress_icmpv6[0].id
    high_ports  = gcore_cloud_security_group_rule.ingress_high_ports[0].id
    egress_https = gcore_cloud_security_group_rule.egress_https[0].id
    sctp        = gcore_cloud_security_group_rule.ingress_sctp[0].id
    http_no_desc = gcore_cloud_security_group_rule.ingress_http_no_desc[0].id
    self_ref    = gcore_cloud_security_group_rule.ingress_self_ref[0].id
  } : {}
}

output "import_sg_id" {
  value = var.create_import_sg ? gcore_cloud_security_group.import_test[0].id : "not created"
}
