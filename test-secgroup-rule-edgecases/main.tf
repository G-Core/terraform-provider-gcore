terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {
  # Uses environment variables
}

# Test security group for comprehensive edge case testing
resource "gcore_cloud_security_group" "test" {
  project_id = 379987
  region_id  = 76

  # ============================================================================
  # LIMITATION: Nested block required (not ideal UX)
  # ============================================================================
  # Current structure (REQUIRED):
  security_group = {
    name        = "test-secgroup-edgecases"
    description = "Comprehensive edge case testing of security group rules"
  }

  # Desired structure (DOESN'T WORK YET):
  # name        = "test-secgroup-edgecases"
  # description = "Comprehensive edge case testing of security group rules"

  # ----------------------------------------------------------------------------
  # WHY IS IT NESTED?
  # ----------------------------------------------------------------------------
  # The API request schema requires nesting:
  # POST /cloud/v1/securitygroups/{project_id}/{region_id}
  # Body: {
  #   "security_group": {           ← Nested object
  #     "name": "...",
  #     "description": "..."
  #   }
  # }
  #
  # But the response is flat:
  # Response: {
  #   "id": "...",
  #   "name": "...",               ← Flat
  #   "description": "...",         ← Flat
  #   "project_id": 379987
  # }
  #
  # ----------------------------------------------------------------------------
  # HOW TO FIX IT
  # ----------------------------------------------------------------------------
  # File: /Users/user/repos/gcore-config/openapi.yml
  #
  # Current (causes nesting):
  #   POST /cloud/v1/securitygroups/{project_id}/{region_id}:
  #     requestBody:
  #       content:
  #         application/json:
  #           schema:
  #             $ref: '#/components/schemas/CreateSecurityGroupSerializer'
  #
  # This schema has structure:
  #   CreateSecurityGroupSerializer:
  #     properties:
  #       security_group:                    ← Nesting here!
  #         $ref: '#/components/schemas/SingleCreateSecurityGroupSerializer'
  #
  # Change to (will flatten):
  #   POST /cloud/v1/securitygroups/{project_id}/{region_id}:
  #     requestBody:
  #       content:
  #         application/json:
  #           schema:
  #             $ref: '#/components/schemas/SingleCreateSecurityGroupSerializer'
  #                   ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
  #                   Use this directly (already flat)
  #
  # SingleCreateSecurityGroupSerializer already has flat structure:
  #   properties:
  #     name:
  #       type: string
  #     description:
  #       type: string
  #
  # ----------------------------------------------------------------------------
  # AFTER THE FIX
  # ----------------------------------------------------------------------------
  # 1. Update openapi.yml as shown above
  # 2. Regenerate provider: stainless generate
  # 3. New Terraform syntax will be:
  #
  #    resource "gcore_cloud_security_group" "test" {
  #      project_id  = 379987
  #      region_id   = 76
  #      name        = "test-secgroup-edgecases"  ← Flat!
  #      description = "Comprehensive testing"    ← Flat!
  #    }
  #
  # 4. Much cleaner and consistent with other cloud providers (AWS, Azure, etc.)
  # ============================================================================
}

# Edge Case 1: Basic TCP rule with port range (UPDATED for testing)
resource "gcore_cloud_security_group_rule" "tcp_range" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = gcore_cloud_security_group.test.project_id  # ✅ Reference parent
  region_id  = gcore_cloud_security_group.test.region_id   # ✅ Reference parent

  direction      = "ingress"
  ethertype      = "IPv4"
  protocol       = "tcp"
  port_range_min = 9000  # Updated from 8000
  port_range_max = 9100  # Updated from 8100
  remote_ip_prefix = "192.168.1.0/24"
  description    = "TCP port range 9000-9100 (updated)"
}

# Edge Case 2: UDP single port (DNS)
resource "gcore_cloud_security_group_rule" "udp_single" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = gcore_cloud_security_group.test.project_id
  region_id  = gcore_cloud_security_group.test.region_id

  direction      = "ingress"
  ethertype      = "IPv4"
  protocol       = "udp"
  port_range_min = 53
  port_range_max = 53
  remote_ip_prefix = "0.0.0.0/0"
  description    = "DNS UDP"
}

# Edge Case 3: ICMP (no port specification)
resource "gcore_cloud_security_group_rule" "icmp" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = gcore_cloud_security_group.test.project_id
  region_id  = gcore_cloud_security_group.test.region_id

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "icmp"
  remote_ip_prefix = "0.0.0.0/0"
  description      = "ICMP ping"
}

# Edge Case 4: IPv6 rule
resource "gcore_cloud_security_group_rule" "ipv6" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = gcore_cloud_security_group.test.project_id
  region_id  = gcore_cloud_security_group.test.region_id

  direction        = "ingress"
  ethertype        = "IPv6"
  protocol         = "tcp"
  port_range_min   = 443
  port_range_max   = 443
  remote_ip_prefix = "::/0"
  description      = "HTTPS IPv6"
}

# Edge Case 5: Egress rule
resource "gcore_cloud_security_group_rule" "egress" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = gcore_cloud_security_group.test.project_id
  region_id  = gcore_cloud_security_group.test.region_id

  direction        = "egress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = 443
  port_range_max   = 443
  remote_ip_prefix = "0.0.0.0/0"
  description      = "Outbound HTTPS"
}

# Edge Case 6: Rule without description (optional field)
resource "gcore_cloud_security_group_rule" "no_desc" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = gcore_cloud_security_group.test.project_id
  region_id  = gcore_cloud_security_group.test.region_id

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = 22
  port_range_max   = 22
  remote_ip_prefix = "10.0.0.0/8"
}

# Edge Case 7: Protocol variations (SCTP)
resource "gcore_cloud_security_group_rule" "sctp" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = gcore_cloud_security_group.test.project_id
  region_id  = gcore_cloud_security_group.test.region_id

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "sctp"
  port_range_min   = 3868
  port_range_max   = 3868
  remote_ip_prefix = "0.0.0.0/0"
  description      = "SCTP protocol test"
}

# Edge Case 8: Wide port range
resource "gcore_cloud_security_group_rule" "wide_range" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = gcore_cloud_security_group.test.project_id
  region_id  = gcore_cloud_security_group.test.region_id

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = 10000
  port_range_max   = 20000
  remote_ip_prefix = "192.168.0.0/16"
  description      = "Wide port range"
}

# Edge Case 9: Remote group ID (instead of IP prefix)
# Note: This requires an existing security group, so we'll use the test group itself
resource "gcore_cloud_security_group_rule" "remote_group" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = gcore_cloud_security_group.test.project_id
  region_id  = gcore_cloud_security_group.test.region_id

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = 3306
  port_range_max   = 3306
  remote_group_id  = gcore_cloud_security_group.test.id
  description      = "MySQL from same security group"
}

# Outputs for verification
output "security_group_id" {
  value = gcore_cloud_security_group.test.id
}

output "rule_ids" {
  value = {
    tcp_range    = gcore_cloud_security_group_rule.tcp_range.id
    udp_single   = gcore_cloud_security_group_rule.udp_single.id
    icmp         = gcore_cloud_security_group_rule.icmp.id
    ipv6         = gcore_cloud_security_group_rule.ipv6.id
    egress       = gcore_cloud_security_group_rule.egress.id
    no_desc      = gcore_cloud_security_group_rule.no_desc.id
    sctp         = gcore_cloud_security_group_rule.sctp.id
    wide_range   = gcore_cloud_security_group_rule.wide_range.id
    remote_group = gcore_cloud_security_group_rule.remote_group.id
  }
}

output "rule_count" {
  value = 9
}
