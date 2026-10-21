package waap_domain_firewall_rule_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/packages/respjson"
	"github.com/G-Core/gcore-go/waap"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

const resourceName = "gcore_waap_domain_firewall_rule.test"

// TestAccWaapDomainFirewallRule_basic creates a rule with only the required
// attributes. The post-apply plan the framework runs after every step must be
// empty, which pins the optional fields the API may fill in on its own
// (description defaults to "", block.status_code to 403). The import step uses an import block so a
// forced replacement or drift on import fails the test. The first steps check
// that invalid configs fail at plan time.
func TestAccWaapDomainFirewallRule_basic(t *testing.T) {
	domainID := waapFirewallDomainID(t)
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWaapDomainFirewallRuleDestroy,
		Steps: []resource.TestStep{
			// Plan-time validation rejects a rule with both actions before any API call.
			{
				Config: strings.Replace(testAccWaapDomainFirewallRuleMinimalConfig(domainID, rName),
					"block = {}", "block = {}\n    allow = {}", 1),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`must set exactly one of allow or block`),
			},
			// Values the API would rewrite fail at plan time, so they never cause a
			// change on every plan. None of these steps calls the API.
			{
				Config: strings.Replace(testAccWaapDomainFirewallRuleMinimalConfig(domainID, rName),
					"192.0.2.10", "2001:DB8::20", 1),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Use\s+"2001:db8::20"`),
			},
			{
				Config: strings.Replace(testAccWaapDomainFirewallRuleMinimalConfig(domainID, rName),
					fmt.Sprintf("%q", rName), fmt.Sprintf("%q", " "+rName+" "), 1),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`must\s+not\s+start\s+or\s+end\s+with\s+whitespace`),
			},
			{
				Config: testAccWaapDomainFirewallRuleMinimalConfig(domainID, rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("conditions").AtSliceIndex(0).AtMapKey("ip").AtMapKey("negation"),
						knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("action").AtMapKey("block").AtMapKey("status_code"), knownvalue.Int64Exact(403)),
				},
				Check: testAccCheckWaapDomainFirewallRuleBlock(resourceName, 403, ""),
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

// TestAccWaapDomainFirewallRule_ipv6 applies IPv6 values in the form the
// validator accepts and a name with inner whitespace. The post-apply plan must
// be empty and import must show no changes, which proves the API stores these
// values byte for byte as written (it rewrites IPv6 to RFC 5952 form and trims
// the name, so any mismatch would show a change on every plan).
func TestAccWaapDomainFirewallRule_ipv6(t *testing.T) {
	domainID := waapFirewallDomainID(t)
	rName := acctest.RandomName() + "  inner"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWaapDomainFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: strings.Replace(testAccWaapDomainFirewallRuleMinimalConfig(domainID, rName),
					"192.0.2.10", "2001:db8::20", 1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("conditions").AtSliceIndex(0).AtMapKey("ip").AtMapKey("ip_address"),
						knownvalue.StringExact("2001:db8::20")),
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				ImportStateIdFunc: acctest.BuildImportID(resourceName, "domain_id", "id"),
			},
			{
				Config: strings.NewReplacer("198.51.100.1", "2001:db8::1", "198.51.100.50", "2001:db8::ff").
					Replace(testAccWaapDomainFirewallRuleAllowRangeConfig(domainID, rName, "ipv6 range", true, false)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("conditions").AtSliceIndex(0).AtMapKey("ip_range").AtMapKey("lower_bound"),
						knownvalue.StringExact("2001:db8::1")),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("conditions").AtSliceIndex(0).AtMapKey("ip_range").AtMapKey("upper_bound"),
						knownvalue.StringExact("2001:db8::ff")),
				},
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

// TestAccWaapDomainFirewallRule_update walks one rule through in-place updates.
// Every step must plan as an update, never a replacement. Step 2 imports the
// allow rule, which must read back as allow = {} with no changes.
//
// Step 4 changes only block.status_code while action_duration stays set, and
// checks the API: the PATCH replaces the whole action object, so an update body
// with only the changed field would drop action_duration. Step 5 removes
// description; the API ignores null there, so it must be cleared with "".
// Step 6 removes action_duration: the API must clear it and the next plan must
// be empty.
func TestAccWaapDomainFirewallRule_update(t *testing.T) {
	domainID := waapFirewallDomainID(t)
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWaapDomainFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWaapDomainFirewallRuleAllowRangeConfig(domainID, rName, "office range", false, true),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("description"), knownvalue.StringExact("office range")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("enabled"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("conditions").AtSliceIndex(0).AtMapKey("ip_range").AtMapKey("negation"),
						knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("action").AtMapKey("allow"), knownvalue.MapExact(map[string]knownvalue.Check{})),
				},
			},
			// An allow rule must import as allow = {}, not null.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				ImportStateIdFunc: acctest.BuildImportID(resourceName, "domain_id", "id"),
			},
			{
				Config: testAccWaapDomainFirewallRuleBlockConfig(domainID, rName+"-v2", `description = "blocked range"`, 429, "5m"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("name"), knownvalue.StringExact(rName+"-v2")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("action").AtMapKey("allow"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("action").AtMapKey("block").AtMapKey("status_code"), knownvalue.Int64Exact(429)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("action").AtMapKey("block").AtMapKey("action_duration"), knownvalue.StringExact("5m")),
				},
			},
			{
				Config: testAccWaapDomainFirewallRuleBlockConfig(domainID, rName+"-v2", `description = "blocked range"`, 405, "5m"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: testAccCheckWaapDomainFirewallRuleBlock(resourceName, 405, "5m"),
			},
			{
				Config: testAccWaapDomainFirewallRuleBlockConfig(domainID, rName+"-v2", "", 405, "5m"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("description"), knownvalue.StringExact("")),
				},
			},
			{
				Config: testAccWaapDomainFirewallRuleBlockConfig(domainID, rName+"-v2", "", 405, ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("action").AtMapKey("block").AtMapKey("action_duration"), knownvalue.Null()),
				},
				Check: testAccCheckWaapDomainFirewallRuleBlock(resourceName, 405, ""),
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

// waapFirewallDomainID returns the shared WAAP test domain and skips when its
// plan has no room for the rules this package holds at once.
func waapFirewallDomainID(t *testing.T) int64 {
	t.Helper()

	domainID := acctest.WaapDomainID(t)

	client, err := acctest.NewGcoreClient()
	if err != nil {
		t.Fatal(err)
	}
	overview, err := client.Waap.GetAccountOverview(context.Background())
	if err != nil {
		t.Fatalf("reading WAAP account overview: %s", err)
	}
	// The five acceptance tests in this package run in parallel and each holds
	// one rule. Skip early only when the plan reports a firewall rule quota
	// without room.
	const rulesNeeded = 5
	for slug, quota := range overview.Quotas {
		if strings.Contains(slug, "firewall") && quota.Allowed-quota.Current < rulesNeeded {
			t.Skipf("WAAP quota %q has no room for %d firewall rules", slug, rulesNeeded)
		}
	}

	return domainID
}

func testAccCheckWaapDomainFirewallRuleDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_waap_domain_firewall_rule" {
			continue
		}
		_, err := getFirewallRule(client, rs)
		if err == nil {
			return fmt.Errorf("WAAP firewall rule %s still exists", rs.Primary.ID)
		}
		if !acctest.IsNotFoundError(err) {
			return err
		}
	}
	return nil
}

// testAccCheckWaapDomainFirewallRuleBlock reads the rule from the API, not from
// state, so a PATCH that dropped an unchanged nested field is caught even when
// state still carries the planned value. An empty actionDuration means the key
// must be absent from the API response.
func testAccCheckWaapDomainFirewallRuleBlock(name string, statusCode int64, actionDuration string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource %s not found in state", name)
		}
		client, err := acctest.NewGcoreClient()
		if err != nil {
			return err
		}
		rule, err := getFirewallRule(client, rs)
		if err != nil {
			return err
		}
		block := rule.Action.Block
		if block.StatusCode != statusCode || block.ActionDuration != actionDuration {
			return fmt.Errorf("API has block {status_code=%d action_duration=%q}, want {%d %q}",
				block.StatusCode, block.ActionDuration, statusCode, actionDuration)
		}
		// absent, null and "" all decode to "", so check the raw JSON: an
		// unset action_duration must be left out of the response entirely.
		if raw := block.JSON.ActionDuration.Raw(); actionDuration == "" && raw != respjson.Omitted {
			return fmt.Errorf("API still returns block.action_duration as %s, want it absent", raw)
		}
		return nil
	}
}

func getFirewallRule(client *gcore.Client, rs *terraform.ResourceState) (*waap.WaapFirewallRule, error) {
	ruleID, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
	if err != nil {
		return nil, err
	}
	domainID, err := strconv.ParseInt(rs.Primary.Attributes["domain_id"], 10, 64)
	if err != nil {
		return nil, err
	}
	return client.Waap.Domains.FirewallRules.Get(context.Background(), ruleID, waap.DomainFirewallRuleGetParams{
		DomainID: domainID,
	})
}

func testAccWaapDomainFirewallRuleMinimalConfig(domainID int64, name string) string {
	return fmt.Sprintf(`
resource "gcore_waap_domain_firewall_rule" "test" {
  domain_id = %[1]d
  name      = %[2]q
  enabled   = true
  action = {
    block = {}
  }
  conditions = [{
    ip = { ip_address = "192.0.2.10" }
  }]
}
`, domainID, name)
}

func testAccWaapDomainFirewallRuleAllowRangeConfig(domainID int64, name, description string, enabled, negation bool) string {
	return fmt.Sprintf(`
resource "gcore_waap_domain_firewall_rule" "test" {
  domain_id   = %[1]d
  name        = %[2]q
  description = %[3]q
  enabled     = %[4]t
  action = {
    allow = {}
  }
  conditions = [{
    ip_range = {
      lower_bound = "198.51.100.1"
      upper_bound = "198.51.100.50"
      negation    = %[5]t
    }
  }]
}
`, domainID, name, description, enabled, negation)
}

// testAccWaapDomainFirewallRuleBlockConfig omits action_duration when it is "".
func testAccWaapDomainFirewallRuleBlockConfig(domainID int64, name, descriptionLine string, statusCode int64, actionDuration string) string {
	durationLine := ""
	if actionDuration != "" {
		durationLine = fmt.Sprintf("action_duration = %q", actionDuration)
	}
	return fmt.Sprintf(`
resource "gcore_waap_domain_firewall_rule" "test" {
  domain_id = %[1]d
  name      = %[2]q
  %[3]s
  enabled   = true
  action = {
    block = {
      status_code = %[4]d
      %[5]s
    }
  }
  conditions = [{
    ip_range = {
      lower_bound = "198.51.100.1"
      upper_bound = "198.51.100.50"
    }
  }]
}
`, domainID, name, descriptionLine, statusCode, durationLine)
}
