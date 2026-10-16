package waap_domain_advanced_rule_test

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

// Reads one rule by rule_id and by find_one_by name, and lists rules filtered
// by name.
func TestAccWaapDomainAdvancedRuleDataSource_basic(t *testing.T) {
	domainID := acctest.WaapDomainID(t)
	rName := acctest.RandomName()

	byID := "data.gcore_waap_domain_advanced_rule.by_id"
	byName := "data.gcore_waap_domain_advanced_rule.by_name"
	list := "data.gcore_waap_domain_advanced_rules.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWaapDomainAdvancedRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWaapDomainAdvancedRuleConfigEnabled(domainID, rName, false, `
  description = "data source test"
  action = {
    block = { status_code = 418 }
  }`) + fmt.Sprintf(`
data "gcore_waap_domain_advanced_rule" "by_id" {
  domain_id = %[1]d
  rule_id   = gcore_waap_domain_advanced_rule.test.id
}

data "gcore_waap_domain_advanced_rule" "by_name" {
  domain_id = %[1]d
  find_one_by = {
    name = gcore_waap_domain_advanced_rule.test.name
  }
}

data "gcore_waap_domain_advanced_rules" "test" {
  domain_id = %[1]d
  name      = gcore_waap_domain_advanced_rule.test.name
}
`, domainID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs(resourceName, tfjsonpath.New("id"), byID, tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue(byID, tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue(byID, tfjsonpath.New("description"), knownvalue.StringExact("data source test")),
					statecheck.ExpectKnownValue(byID,
						tfjsonpath.New("action").AtMapKey("block").AtMapKey("status_code"), knownvalue.Int64Exact(418)),
					statecheck.CompareValuePairs(resourceName, tfjsonpath.New("id"), byName, tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue(list, tfjsonpath.New("items"), knownvalue.ListSizeExact(1)),
					statecheck.CompareValuePairs(resourceName, tfjsonpath.New("id"),
						list, tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("id"), compare.ValuesSame()),
				},
			},
		},
	})
}
