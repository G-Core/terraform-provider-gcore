package waap_domain_custom_rule_test

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

// The single-rule data source, by rule_id and by find_one_by, and the list data
// source must all report the rule the resource created.
func TestAccWaapDomainCustomRuleDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()
	domainID := acctest.WaapDomainID(t)

	const (
		byID   = "data.gcore_waap_domain_custom_rule.by_id"
		byName = "data.gcore_waap_domain_custom_rule.by_name"
		list   = "data.gcore_waap_domain_custom_rules.all"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWaapDomainCustomRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWaapDomainCustomRuleConfig(domainID, rName, false, `
  description = "data source"
  action = {
    block = {
      status_code = 418
    }
  }
  conditions = [
    { url = { url = "/admin", match_type = "Exact" } },
  ]`) + fmt.Sprintf(`
data "gcore_waap_domain_custom_rule" "by_id" {
  domain_id = gcore_waap_domain_custom_rule.test.domain_id
  rule_id   = gcore_waap_domain_custom_rule.test.id
}

data "gcore_waap_domain_custom_rule" "by_name" {
  domain_id   = gcore_waap_domain_custom_rule.test.domain_id
  find_one_by = { name = gcore_waap_domain_custom_rule.test.name }
}

data "gcore_waap_domain_custom_rules" "all" {
  domain_id = gcore_waap_domain_custom_rule.test.domain_id
  name      = %q
}
`, rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs(byID, tfjsonpath.New("id"), resourceName, tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.CompareValuePairs(byName, tfjsonpath.New("id"), resourceName, tfjsonpath.New("id"), compare.ValuesSame()),
					statecheck.ExpectKnownValue(byID, tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue(byID, tfjsonpath.New("description"), knownvalue.StringExact("data source")),
					statecheck.ExpectKnownValue(byID,
						tfjsonpath.New("action").AtMapKey("block").AtMapKey("status_code"), knownvalue.Int64Exact(418)),
					statecheck.ExpectKnownValue(byID,
						tfjsonpath.New("conditions").AtSliceIndex(0).AtMapKey("url").AtMapKey("url"), knownvalue.StringExact("/admin")),
					statecheck.ExpectKnownValue(list, tfjsonpath.New("items"), knownvalue.ListSizeExact(1)),
					statecheck.CompareValuePairs(list, tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("id"),
						resourceName, tfjsonpath.New("id"), compare.ValuesSame()),
				},
			},
		},
	})
}
