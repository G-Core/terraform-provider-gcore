package dns_zone_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/G-Core/gcore-go"
	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func TestAccDNSZone_basic(t *testing.T) {
	rName := acctest.RandomName() + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDNSZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSZoneConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_dns_zone.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_dns_zone.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_dns_zone.test",
						tfjsonpath.New("enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("gcore_dns_zone.test",
						tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_dns_zone.test",
						tfjsonpath.New("serial"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccDNSZone_update(t *testing.T) {
	rName := acctest.RandomName() + ".com"

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDNSZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSZoneConfigContact(rName, ""),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_dns_zone.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					compareIDSame.AddStateValue(
						"gcore_dns_zone.test",
						tfjsonpath.New("id"),
					),
				},
			},
			{
				Config: testAccDNSZoneConfigContact(rName, "admin.example.com"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_dns_zone.test",
						tfjsonpath.New("contact"), knownvalue.StringExact("admin.example.com")),
					compareIDSame.AddStateValue(
						"gcore_dns_zone.test",
						tfjsonpath.New("id"),
					),
				},
			},
		},
	})
}

func TestAccDNSZone_withMeta(t *testing.T) {
	rName := acctest.RandomName() + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDNSZoneDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with meta
			{
				Config: testAccDNSZoneConfigWithMeta(rName, "https://example.com/webhook", "POST"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_dns_zone.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_dns_zone.test",
						tfjsonpath.New("meta").AtMapKey("webhook"),
						knownvalue.StringExact("https://example.com/webhook")),
					statecheck.ExpectKnownValue("gcore_dns_zone.test",
						tfjsonpath.New("meta").AtMapKey("webhook_method"),
						knownvalue.StringExact("POST")),
				},
			},
			// Step 2: no drift
			{
				Config: testAccDNSZoneConfigWithMeta(rName, "https://example.com/webhook", "POST"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Step 3: Update meta
			{
				Config: testAccDNSZoneConfigWithMeta(rName, "https://example.com/webhook-v2", "PATCH"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_dns_zone.test",
						tfjsonpath.New("meta").AtMapKey("webhook"),
						knownvalue.StringExact("https://example.com/webhook-v2")),
					statecheck.ExpectKnownValue("gcore_dns_zone.test",
						tfjsonpath.New("meta").AtMapKey("webhook_method"),
						knownvalue.StringExact("PATCH")),
				},
			},
			// Step 4: no drift after update
			{
				Config: testAccDNSZoneConfigWithMeta(rName, "https://example.com/webhook-v2", "PATCH"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Step 5: Import
			{
				ResourceName:      "gcore_dns_zone.test",
				ImportState:       true,
				ImportStateIdFunc: acctest.BuildImportID("gcore_dns_zone.test", "name"),
				ImportStateKind:   resource.ImportBlockWithID,
				// meta is cleared on import to avoid drift from server-injected keys
				ImportStateVerifyIgnore: []string{"meta"},
			},
		},
	})
}

func TestAccDNSZone_import(t *testing.T) {
	rName := acctest.RandomName() + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDNSZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSZoneConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_dns_zone.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
				},
			},
			{
				ResourceName:      "gcore_dns_zone.test",
				ImportState:       true,
				ImportStateIdFunc: acctest.BuildImportID("gcore_dns_zone.test", "name"),
				ImportStateKind:   resource.ImportBlockWithID,
			},
		},
	})
}

func TestAccDNSZone_dnssec(t *testing.T) {
	rName := acctest.RandomName() + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDNSZoneDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with DNSSEC enabled
			{
				Config: testAccDNSZoneConfigDnssec(rName, true),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_dns_zone.test",
						tfjsonpath.New("dnssec_enabled"), knownvalue.Bool(true)),
				},
			},
			// Step 2: no drift
			{
				Config: testAccDNSZoneConfigDnssec(rName, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Step 3: Disable DNSSEC
			{
				Config: testAccDNSZoneConfigDnssec(rName, false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_dns_zone.test",
						tfjsonpath.New("dnssec_enabled"), knownvalue.Bool(false)),
				},
			},
		},
	})
}

func testAccCheckDNSZoneDestroy(s *terraform.State) error {
	return acctest.CheckResourceDestroyed(s, "gcore_dns_zone", func(client *gcore.Client, id string) error {
		_, err := client.DNS.Zones.Get(context.Background(), id)
		return err
	})
}

func testAccDNSZoneConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_dns_zone" "test" {
  name = %[1]q
}`, name)
}

func testAccDNSZoneConfigContact(name string, contact string) string {
	if contact == "" {
		return fmt.Sprintf(`
resource "gcore_dns_zone" "test" {
  name = %[1]q
}`, name)
	}
	return fmt.Sprintf(`
resource "gcore_dns_zone" "test" {
  name    = %[1]q
  contact = %[2]q
}`, name, contact)
}

func testAccDNSZoneConfigWithMeta(name, webhook, method string) string {
	return fmt.Sprintf(`
resource "gcore_dns_zone" "test" {
  name = %[1]q
  meta = {
    webhook        = %[2]q
    webhook_method = %[3]q
  }
}`, name, webhook, method)
}

func testAccDNSZoneConfigDnssec(name string, enabled bool) string {
	return fmt.Sprintf(`
resource "gcore_dns_zone" "test" {
  name           = %[1]q
  dnssec_enabled = %[2]t
}`, name, enabled)
}
