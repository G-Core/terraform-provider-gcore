package waap_domain_advanced_rule_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/gcore-go/waap"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

const resourceName = "gcore_waap_domain_advanced_rule.test"

func testAccCheckWaapDomainAdvancedRuleDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_waap_domain_advanced_rule" {
			continue
		}
		ruleID, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse rule ID: %w", err)
		}
		domainID, err := strconv.ParseInt(rs.Primary.Attributes["domain_id"], 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse domain_id: %w", err)
		}
		_, err = client.Waap.Domains.AdvancedRules.Get(context.Background(), ruleID, waap.DomainAdvancedRuleGetParams{
			DomainID: domainID,
		})
		if err == nil {
			return fmt.Errorf("WAAP advanced rule %d still exists", ruleID)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking WAAP advanced rule deletion: %w", err)
		}
	}
	return nil
}

func TestAccWaapDomainAdvancedRule_basic(t *testing.T) {
	rName := acctest.RandomName()
	domainID := acctest.WaapDomainID(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWaapDomainAdvancedRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWaapDomainAdvancedRuleConfig(domainID, rName, `
  description = "block rule"
  phase       = "access"
  action = {
    block = {
      status_code     = 429
      action_duration = "5m"
    }
  }`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("description"), knownvalue.StringExact("block rule")),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("action").AtMapKey("block").AtMapKey("status_code"), knownvalue.Int64Exact(429)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("action").AtMapKey("block").AtMapKey("action_duration"), knownvalue.StringExact("5m")),
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: acctest.BuildImportID(resourceName, "domain_id", "id"),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				ImportStateIdFunc: acctest.BuildImportID(resourceName, "domain_id", "id"),
			},
		},
	})
}

// Omitted optional fields must not drift: the API fills the block status code
// with 403 and phase with "access".
func TestAccWaapDomainAdvancedRule_defaults(t *testing.T) {
	rName := acctest.RandomName()
	domainID := acctest.WaapDomainID(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWaapDomainAdvancedRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWaapDomainAdvancedRuleConfig(domainID, rName, `
  action = {
    block = {}
  }`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("action").AtMapKey("block").AtMapKey("status_code"), knownvalue.Int64Exact(403)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("phase"), knownvalue.StringExact("access")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("description"), knownvalue.StringExact("")),
				},
			},
		},
	})
}

// Update in place across action kinds, and clear the description by removing
// it from the config.
func TestAccWaapDomainAdvancedRule_update(t *testing.T) {
	rName := acctest.RandomName()
	rNameUpdated := acctest.RandomName()
	domainID := acctest.WaapDomainID(t)

	idCheck := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWaapDomainAdvancedRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWaapDomainAdvancedRuleConfig(domainID, rName, `
  description = "before"
  action = {
    block = {
      action_duration = "5m"
    }
  }`),
				ConfigStateChecks: []statecheck.StateCheck{
					idCheck.AddStateValue(resourceName, tfjsonpath.New("id")),
				},
			},
			{
				Config: testAccWaapDomainAdvancedRuleConfigEnabled(domainID, rNameUpdated, false, `
  phase = "header_filter"
  action = {
    allow = {}
  }`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					idCheck.AddStateValue(resourceName, tfjsonpath.New("id")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("name"), knownvalue.StringExact(rNameUpdated)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("enabled"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("phase"), knownvalue.StringExact("header_filter")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("action").AtMapKey("block"), knownvalue.Null()),
				},
			},
			{
				Config: testAccWaapDomainAdvancedRuleConfig(domainID, rNameUpdated, `
  description = "tagged"
  action = {
    tag = {
      tags = ["zeta", "Alpha", "beta", "beta"]
    }
  }`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					idCheck.AddStateValue(resourceName, tfjsonpath.New("id")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("description"), knownvalue.StringExact("tagged")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("action").AtMapKey("allow"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("action").AtMapKey("tag").AtMapKey("tags"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("Alpha"),
							knownvalue.StringExact("beta"),
							knownvalue.StringExact("zeta"),
						})),
				},
			},
			// The API stores tags as a set and returns them in any order, so a
			// reordered config is the same value.
			{
				Config: testAccWaapDomainAdvancedRuleConfig(domainID, rNameUpdated, `
  description = "tagged"
  action = {
    tag = {
      tags = ["beta", "zeta", "Alpha"]
    }
  }`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

// A rule created outside Terraform without a description and with an
// empty-object action kind must import with a no-op plan.
func TestAccWaapDomainAdvancedRule_importExternal(t *testing.T) {
	rName := acctest.RandomName()
	domainID := acctest.WaapDomainID(t)

	client, err := acctest.NewGcoreClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	rule, err := client.Waap.Domains.AdvancedRules.New(context.Background(), domainID, waap.DomainAdvancedRuleNewParams{
		Name:    rName,
		Enabled: false,
		Source:  "request.ip == '192.0.2.10'",
		Action: waap.DomainAdvancedRuleNewParamsAction{
			Monitor: map[string]any{},
		},
	})
	if err != nil {
		t.Fatalf("failed to create rule: %v", err)
	}
	t.Cleanup(func() {
		_ = client.Waap.Domains.AdvancedRules.Delete(context.Background(), rule.ID, waap.DomainAdvancedRuleDeleteParams{
			DomainID: domainID,
		})
	})

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWaapDomainAdvancedRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
import {
  to = gcore_waap_domain_advanced_rule.test
  id = "%[1]d/%[2]d"
}
`, domainID, rule.ID) + testAccWaapDomainAdvancedRuleConfigEnabled(domainID, rName, false, `
  action = {
    monitor = {}
  }`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("action").AtMapKey("monitor"), knownvalue.MapExact(map[string]knownvalue.Check{})),
				},
			},
		},
	})
}

func testAccWaapDomainAdvancedRuleConfig(domainID int64, name, extra string) string {
	return testAccWaapDomainAdvancedRuleConfigEnabled(domainID, name, true, extra)
}

func testAccWaapDomainAdvancedRuleConfigEnabled(domainID int64, name string, enabled bool, extra string) string {
	return fmt.Sprintf(`
resource "gcore_waap_domain_advanced_rule" "test" {
  domain_id = %d
  name      = %q
  enabled   = %t
  source    = "request.ip == '192.0.2.10'"
%s
}
`, domainID, name, enabled, extra)
}
