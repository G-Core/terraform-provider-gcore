package waap_domain_custom_rule_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func condition(i int, kind string) tfjsonpath.Path {
	return tfjsonpath.New("conditions").AtSliceIndex(i).AtMapKey(kind)
}

func TestAccWaapDomainCustomRule_basic(t *testing.T) {
	rName := acctest.RandomName()
	domainID := acctest.WaapDomainID(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWaapDomainCustomRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWaapDomainCustomRuleConfig(domainID, rName, false, `
  description = "block rule"
  action = {
    block = {
      status_code     = 429
      action_duration = "5m"
    }
  }
  conditions = [
    { ip = { ip_address = "192.0.2.10", negation = true } },
    { url = { url = "/admin", match_type = "Exact" } },
    { country = { country_code = ["US", "CN"] } },
  ]`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("enabled"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("description"), knownvalue.StringExact("block rule")),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("action").AtMapKey("block").AtMapKey("status_code"), knownvalue.Int64Exact(429)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("action").AtMapKey("block").AtMapKey("action_duration"), knownvalue.StringExact("5m")),
					statecheck.ExpectKnownValue(resourceName, condition(0, "ip").AtMapKey("negation"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(resourceName, condition(1, "url").AtMapKey("match_type"), knownvalue.StringExact("Exact")),
					statecheck.ExpectKnownValue(resourceName, condition(2, "country").AtMapKey("country_code"),
						knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("US"), knownvalue.StringExact("CN")})),
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
// with 403, the owner types with COMMERCIAL, every negation with false and
// every match type with Contains. Removing a set owner_types list plans the
// COMMERCIAL default, which is what the API then stores.
func TestAccWaapDomainCustomRule_defaults(t *testing.T) {
	rName := acctest.RandomName()
	domainID := acctest.WaapDomainID(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWaapDomainCustomRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWaapDomainCustomRuleConfig(domainID, rName, true, `
  action = {
    block = {}
  }
  conditions = [
    { owner_types = {} },
    { user_agent = { user_agent = "curl/" } },
  ]`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("action").AtMapKey("block").AtMapKey("status_code"), knownvalue.Int64Exact(403)),
					statecheck.ExpectKnownValue(resourceName, condition(0, "owner_types").AtMapKey("owner_types"),
						knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("COMMERCIAL")})),
					statecheck.ExpectKnownValue(resourceName, condition(0, "owner_types").AtMapKey("negation"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(resourceName, condition(1, "user_agent").AtMapKey("match_type"), knownvalue.StringExact("Contains")),
				},
			},
			// setting owner_types and then removing it returns to the API default
			{
				Config: testAccWaapDomainCustomRuleConfig(domainID, rName, true, `
  action = {
    block = {}
  }
  conditions = [
    { owner_types = { owner_types = ["ISP", "HOSTING_SERVICES"] } },
    { user_agent = { user_agent = "curl/" } },
  ]`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, condition(0, "owner_types").AtMapKey("owner_types"),
						knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("ISP"), knownvalue.StringExact("HOSTING_SERVICES")})),
				},
			},
			{
				Config: testAccWaapDomainCustomRuleConfig(domainID, rName, true, `
  action = {
    block = {}
  }
  conditions = [
    { owner_types = {} },
    { user_agent = { user_agent = "curl/" } },
  ]`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(resourceName, condition(0, "owner_types").AtMapKey("owner_types"),
							knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("COMMERCIAL")})),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, condition(0, "owner_types").AtMapKey("owner_types"),
						knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("COMMERCIAL")})),
				},
			},
		},
	})
}

// Update in place across action kinds and conditions, and clear the
// description by removing it from the config.
func TestAccWaapDomainCustomRule_update(t *testing.T) {
	rName := acctest.RandomName()
	rNameUpdated := acctest.RandomName()
	domainID := acctest.WaapDomainID(t)

	idCheck := statecheck.CompareValue(compare.ValuesSame())
	updateInPlace := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWaapDomainCustomRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWaapDomainCustomRuleConfig(domainID, rName, false, `
  description = "before"
  action = {
    block = {
      status_code     = 429
      action_duration = "5m"
    }
  }
  conditions = [
    { ip = { ip_address = "192.0.2.10" } },
  ]`),
				ConfigStateChecks: []statecheck.StateCheck{
					idCheck.AddStateValue(resourceName, tfjsonpath.New("id")),
				},
			},
			// block -> monitor: the API rejects "block": null, so only the new kind
			// may be sent
			{
				Config: testAccWaapDomainCustomRuleConfig(domainID, rNameUpdated, true, `
  action = {
    monitor = {}
  }
  conditions = [
    { http_method = { http_method = "POST", negation = true } },
    { ip_range = { lower_bound = "2001:db8::1", upper_bound = "2001:db8::ff" } },
  ]`),
				ConfigPlanChecks: updateInPlace,
				ConfigStateChecks: []statecheck.StateCheck{
					idCheck.AddStateValue(resourceName, tfjsonpath.New("id")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("name"), knownvalue.StringExact(rNameUpdated)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("action").AtMapKey("block"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("action").AtMapKey("monitor"), knownvalue.MapExact(map[string]knownvalue.Check{})),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("conditions"), knownvalue.ListSizeExact(2)),
				},
			},
			// monitor -> tag; the API sorts the tag action, user-defined tags and ASNs, so they are sets
			{
				Config: testAccWaapDomainCustomRuleConfig(domainID, rNameUpdated, true, `
  description = "tagged"
  action = {
    tag = {
      tags = ["zeta", "alpha"]
    }
  }
  conditions = [
    { asn = { asns = [15169, 13335] } },
    { ja3 = { ja3_fingerprints = ["e7d705a3286e19ea42f587b344ee6865"] } },
    { user_defined_tags = { tags = ["Zeta", "Alpha"] } },
  ]`),
				ConfigPlanChecks: updateInPlace,
				ConfigStateChecks: []statecheck.StateCheck{
					idCheck.AddStateValue(resourceName, tfjsonpath.New("id")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("description"), knownvalue.StringExact("tagged")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("action").AtMapKey("monitor"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, condition(0, "asn").AtMapKey("asns"),
						knownvalue.SetExact([]knownvalue.Check{knownvalue.Int64Exact(13335), knownvalue.Int64Exact(15169)})),
					statecheck.ExpectKnownValue(resourceName, condition(2, "user_defined_tags").AtMapKey("tags"),
						knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact("Alpha"), knownvalue.StringExact("Zeta")})),
				},
			},
			// a reordered set is the same value
			{
				Config: testAccWaapDomainCustomRuleConfig(domainID, rNameUpdated, true, `
  description = "tagged"
  action = {
    tag = {
      tags = ["alpha", "zeta"]
    }
  }
  conditions = [
    { asn = { asns = [13335, 15169] } },
    { ja3 = { ja3_fingerprints = ["e7d705a3286e19ea42f587b344ee6865"] } },
    { user_defined_tags = { tags = ["Alpha", "Zeta"] } },
  ]`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			// request_rate must be the only condition of a rule
			{
				Config: testAccWaapDomainCustomRuleConfig(domainID, rNameUpdated, true, `
  action = {
    captcha = {}
  }
  conditions = [
    { request_rate = {
      path_pattern = "/api"
      requests     = 20
      time         = 10
      http_methods = ["POST", "GET"]
      ips          = ["192.0.2.1", "2001:db8::1"]
    } },
  ]`),
				ConfigPlanChecks: updateInPlace,
				ConfigStateChecks: []statecheck.StateCheck{
					idCheck.AddStateValue(resourceName, tfjsonpath.New("id")),
					statecheck.ExpectKnownValue(resourceName, condition(0, "request_rate").AtMapKey("http_methods"),
						knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("POST"), knownvalue.StringExact("GET")})),
				},
			},
		},
	})
}

// An empty-object action kind comes back as {}, which the decoder would read as
// null, so import must keep it and plan a no-op. A rule without a description
// cannot be created through Terraform (the "" default always sends it); that
// read path is covered by TestUnmarshalReadImportedRule.
func TestAccWaapDomainCustomRule_importEmptyAction(t *testing.T) {
	rName := acctest.RandomName()
	domainID := acctest.WaapDomainID(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWaapDomainCustomRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWaapDomainCustomRuleConfig(domainID, rName, false, `
  action = {
    handshake = {}
  }
  conditions = [
    { header_exists = { header = "X-Test" } },
  ]`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("action").AtMapKey("handshake"), knownvalue.MapExact(map[string]knownvalue.Check{})),
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

// Values the API would store in another form, and shapes it would reject, fail
// at plan time.
func TestAccWaapDomainCustomRule_invalidConfig(t *testing.T) {
	rName := acctest.RandomName()
	domainID := acctest.WaapDomainID(t)

	// terraform wraps diagnostics, so a space in want matches any whitespace
	step := func(body string, want string) resource.TestStep {
		return resource.TestStep{
			Config:      testAccWaapDomainCustomRuleConfig(domainID, rName, false, body),
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(strings.ReplaceAll(want, " ", `\s+`)),
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			step(`
  action = { allow = {} }
  conditions = [
    { ip = { ip_address = "192.0.2.10" }, url = { url = "/admin" } },
  ]`, `Each condition must set exactly one condition type, got 2`),
			step(`
  action = { allow = {} }
  conditions = [
    { ip = { ip_address = "2001:DB8::1" } },
  ]`, `Non-canonical IP address`),
			step(`
  action = { allow = {} }
  conditions = [
    { ja3 = { ja3_fingerprints = ["E7D705A3286E19EA42F587B344EE6865"] } },
  ]`, `lowercase 32-character hexadecimal JA3 fingerprint`),
			step(`
  action = { allow = {} }
  conditions = [
    { country = { country_code = ["us"] } },
  ]`, `uppercase ISO 3166-1 alpha-2 country code`),
			step(`
  action = { allow = {} }
  conditions = [
    { http_method = { http_method = "post" } },
  ]`, `value must be one of`),
			step(`
  action = { allow = {} }
  conditions = [
    { header = { header = "X-Foo", value = " bar " } },
  ]`, `Leading or trailing whitespace`),
			step(`
  action = { tag = { tags = ["lead "] } }
  conditions = [
    { content_type = { content_type = ["text/html "] } },
    { tags = { tags = ["xss "] } },
    { user_defined_tags = { tags = ["qa "] } },
  ]`, `Leading or trailing whitespace`),
			step(`
  action = {
    allow   = {}
    monitor = {}
  }
  conditions = [
    { ip = { ip_address = "192.0.2.10" } },
  ]`, `Invalid Attribute Combination`),
			{
				Config: testAccWaapDomainCustomRuleConfig(domainID, " "+rName, false, `
  action = { allow = {} }
  conditions = [
    { ip = { ip_address = "192.0.2.10" } },
  ]`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Leading or trailing whitespace`),
			},
		},
	})
}
