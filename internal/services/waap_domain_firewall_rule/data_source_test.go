package waap_domain_firewall_rule_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

// TestAccWaapDomainFirewallRuleDataSource_basic reads one rule by rule_id and
// by find_one_by name. The find_one_by path lists rules under domain_id, which
// the list call takes as a positional argument rather than a params field.
func TestAccWaapDomainFirewallRuleDataSource_basic(t *testing.T) {
	domainID := waapFirewallDomainID(t)
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWaapDomainFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWaapDomainFirewallRuleDataSourceConfig(domainID, rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs(
						resourceName, tfjsonpath.New("id"),
						"data.gcore_waap_domain_firewall_rule.by_id", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
					statecheck.ExpectKnownValue("data.gcore_waap_domain_firewall_rule.by_id",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("data.gcore_waap_domain_firewall_rule.by_id",
						tfjsonpath.New("description"), knownvalue.StringExact("data source test")),
					statecheck.ExpectKnownValue("data.gcore_waap_domain_firewall_rule.by_id",
						tfjsonpath.New("action").AtMapKey("block").AtMapKey("status_code"), knownvalue.Int64Exact(418)),
					statecheck.ExpectKnownValue("data.gcore_waap_domain_firewall_rule.by_id",
						tfjsonpath.New("conditions").AtSliceIndex(0).AtMapKey("ip").AtMapKey("ip_address"),
						knownvalue.StringExact("192.0.2.20")),
					statecheck.CompareValuePairs(
						resourceName, tfjsonpath.New("id"),
						"data.gcore_waap_domain_firewall_rule.by_name", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

// TestAccWaapDomainFirewallRulesDataSource_basic filters the list by name so
// the result holds exactly the rule this test created.
func TestAccWaapDomainFirewallRulesDataSource_basic(t *testing.T) {
	domainID := waapFirewallDomainID(t)
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWaapDomainFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWaapDomainFirewallRulesDataSourceConfig(domainID, rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_waap_domain_firewall_rules.test",
						tfjsonpath.New("items"), knownvalue.ListSizeExact(1)),
					statecheck.CompareValuePairs(
						resourceName, tfjsonpath.New("id"),
						"data.gcore_waap_domain_firewall_rules.test", tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("id"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func testAccWaapDomainFirewallRuleDataSourceConfig(domainID int64, name string) string {
	return fmt.Sprintf(`
resource "gcore_waap_domain_firewall_rule" "test" {
  domain_id   = %[1]d
  name        = %[2]q
  description = "data source test"
  enabled     = false
  action = {
    block = { status_code = 418 }
  }
  conditions = [{
    ip = { ip_address = "192.0.2.20" }
  }]
}

data "gcore_waap_domain_firewall_rule" "by_id" {
  domain_id = %[1]d
  rule_id   = gcore_waap_domain_firewall_rule.test.id
}

data "gcore_waap_domain_firewall_rule" "by_name" {
  domain_id = %[1]d
  find_one_by = {
    name = gcore_waap_domain_firewall_rule.test.name
  }
}
`, domainID, name)
}

func testAccWaapDomainFirewallRulesDataSourceConfig(domainID int64, name string) string {
	return fmt.Sprintf(`
resource "gcore_waap_domain_firewall_rule" "test" {
  domain_id = %[1]d
  name      = %[2]q
  enabled   = false
  action = {
    allow = {}
  }
  conditions = [{
    ip = { ip_address = "192.0.2.30" }
  }]
}

data "gcore_waap_domain_firewall_rules" "test" {
  domain_id = %[1]d
  name      = gcore_waap_domain_firewall_rule.test.name
}
`, domainID, name)
}
