package cdn_cdn_resource_rule_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/stainless-sdks/gcore-terraform/internal/acctest"
)

func TestAccCDNCDNResourceRuleDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()
	fixture := acctest.SetupCDNResource(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNCDNResourceRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNCDNResourceRuleDataSourceConfig(fixture.ResourceID, rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cdn_cdn_resource_rule.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.CompareValuePairs(
						"gcore_cdn_cdn_resource_rule.test", tfjsonpath.New("id"),
						"data.gcore_cdn_cdn_resource_rule.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func testAccCDNCDNResourceRuleDataSourceConfig(resourceID int64, name string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_cdn_resource_rule" "test" {
  resource_id = %d
  name        = %q
  rule        = "/images/*"
  rule_type   = 0
  active      = false
  weight      = 1
}

data "gcore_cdn_cdn_resource_rule" "test" {
  rule_id     = gcore_cdn_cdn_resource_rule.test.id
  resource_id = gcore_cdn_cdn_resource_rule.test.resource_id
}`, resourceID, name)
}
