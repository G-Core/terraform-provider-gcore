package cdn_resource_rule_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/gcore-go/cdn"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func testAccCheckCDNCDNResourceRuleDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cdn_cdn_resource_rule" {
			continue
		}
		ruleIDStr := rs.Primary.ID
		resourceIDStr := rs.Primary.Attributes["resource_id"]

		ruleID, err := strconv.ParseInt(ruleIDStr, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse rule ID: %w", err)
		}
		resourceID, err := strconv.ParseInt(resourceIDStr, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse resource_id: %w", err)
		}

		_, err = client.CDN.CDNResources.Rules.Get(ctx, ruleID, cdn.CDNResourceRuleGetParams{
			ResourceID: resourceID,
		})
		if err == nil {
			return fmt.Errorf("CDN rule %d still exists", ruleID)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking CDN rule deletion: %w", err)
		}
	}
	return nil
}

func TestAccCDNCDNResourceRule_basic(t *testing.T) {
	rName := acctest.RandomName()
	fixture := setupCDNResource(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNCDNResourceRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNCDNResourceRuleConfig(fixture.ResourceID, rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_cdn_resource_rule.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cdn_cdn_resource_rule.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cdn_cdn_resource_rule.test",
						tfjsonpath.New("rule"), knownvalue.StringExact("/images/*")),
				},
			},
			{
				ResourceName:            "gcore_cdn_cdn_resource_rule.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       acctest.BuildImportID("gcore_cdn_cdn_resource_rule.test", "resource_id", "id"),
				ImportStateVerifyIgnore: []string{"deleted"},
			},
		},
	})
}

func TestAccCDNCDNResourceRule_update(t *testing.T) {
	rName := acctest.RandomName()
	rNameUpdated := acctest.RandomName()
	fixture := setupCDNResource(t)

	idCheck := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNCDNResourceRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNCDNResourceRuleConfig(fixture.ResourceID, rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_cdn_resource_rule.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					idCheck.AddStateValue("gcore_cdn_cdn_resource_rule.test", tfjsonpath.New("id")),
				},
			},
			{
				Config: testAccCDNCDNResourceRuleConfig(fixture.ResourceID, rNameUpdated),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_cdn_resource_rule.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rNameUpdated)),
					idCheck.AddStateValue("gcore_cdn_cdn_resource_rule.test", tfjsonpath.New("id")),
				},
			},
		},
	})
}

func testAccCDNCDNResourceRuleConfig(resourceID int64, name string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_cdn_resource_rule" "test" {
  resource_id = %d
  name        = %q
  rule        = "/images/*"
  rule_type   = 0
  active      = false
  weight      = 1
}`, resourceID, name)
}
