package dns_zone_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/G-Core/gcore-go"
	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
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
