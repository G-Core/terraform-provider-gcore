package dns_zone_rrset_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/G-Core/gcore-go/dns"
	"github.com/G-Core/gcore-go/option"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func TestAccDNSZoneRrset_basic(t *testing.T) {
	zoneName := fmt.Sprintf("tf-test-%s.com", acctest.RandomName())
	rrsetName := fmt.Sprintf("www.%s.", zoneName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDNSZoneRrsetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSZoneRrsetConfig(zoneName, rrsetName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_dns_zone_rrset.test",
						tfjsonpath.New("zone_name"), knownvalue.StringExact(zoneName)),
					statecheck.ExpectKnownValue("gcore_dns_zone_rrset.test",
						tfjsonpath.New("rrset_name"), knownvalue.StringExact(rrsetName)),
					statecheck.ExpectKnownValue("gcore_dns_zone_rrset.test",
						tfjsonpath.New("rrset_type"), knownvalue.StringExact("A")),
					statecheck.ExpectKnownValue("gcore_dns_zone_rrset.test",
						tfjsonpath.New("ttl"), knownvalue.Int64Exact(300)),
				},
			},
			// Step 2: no drift
			{
				Config: testAccDNSZoneRrsetConfig(zoneName, rrsetName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccDNSZoneRrset_update(t *testing.T) {
	zoneName := fmt.Sprintf("tf-test-%s.com", acctest.RandomName())
	rrsetName := fmt.Sprintf("www.%s.", zoneName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDNSZoneRrsetDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with single record
			{
				Config: testAccDNSZoneRrsetConfig(zoneName, rrsetName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_dns_zone_rrset.test",
						tfjsonpath.New("ttl"), knownvalue.Int64Exact(300)),
					statecheck.ExpectKnownValue("gcore_dns_zone_rrset.test",
						tfjsonpath.New("resource_records"),
						knownvalue.ListSizeExact(1)),
				},
			},
			// Step 2: Update TTL and add second record
			{
				Config: testAccDNSZoneRrsetConfigUpdated(zoneName, rrsetName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_dns_zone_rrset.test",
						tfjsonpath.New("ttl"), knownvalue.Int64Exact(600)),
					statecheck.ExpectKnownValue("gcore_dns_zone_rrset.test",
						tfjsonpath.New("resource_records"),
						knownvalue.ListSizeExact(2)),
				},
			},
			// Step 3: no drift after update
			{
				Config: testAccDNSZoneRrsetConfigUpdated(zoneName, rrsetName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccDNSZoneRrset_withMeta(t *testing.T) {
	zoneName := fmt.Sprintf("tf-test-%s.com", acctest.RandomName())
	rrsetName := fmt.Sprintf("geo.%s.", zoneName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDNSZoneRrsetDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with meta on resource records
			{
				Config: testAccDNSZoneRrsetConfigWithMeta(zoneName, rrsetName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_dns_zone_rrset.test",
						tfjsonpath.New("resource_records"),
						knownvalue.ListSizeExact(2)),
				},
			},
			// Step 2: no drift — meta values preserved through API normalization
			{
				Config: testAccDNSZoneRrsetConfigWithMeta(zoneName, rrsetName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// Step 3: Update meta values
			{
				Config: testAccDNSZoneRrsetConfigWithMetaUpdated(zoneName, rrsetName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_dns_zone_rrset.test",
						tfjsonpath.New("resource_records"),
						knownvalue.ListSizeExact(2)),
				},
			},
			// Step 4: no drift after meta update
			{
				Config: testAccDNSZoneRrsetConfigWithMetaUpdated(zoneName, rrsetName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccDNSZoneRrset_import(t *testing.T) {
	zoneName := fmt.Sprintf("tf-test-%s.com", acctest.RandomName())
	rrsetName := fmt.Sprintf("www.%s.", zoneName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDNSZoneRrsetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSZoneRrsetConfig(zoneName, rrsetName),
			},
			{
				ResourceName:                         "gcore_dns_zone_rrset.test",
				ImportState:                          true,
				ImportStateId:                        fmt.Sprintf("%s/%s/A", zoneName, rrsetName),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "zone_name",
				ImportStateVerifyIgnore:              []string{"resource_records", "ttl", "meta", "pickers"},
			},
		},
	})
}

func TestAccDNSZoneRrset_multipleRecordsOrdering(t *testing.T) {
	zoneName := fmt.Sprintf("tf-test-%s.com", acctest.RandomName())
	rrsetName := fmt.Sprintf("multi.%s.", zoneName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDNSZoneRrsetDestroy,
		Steps: []resource.TestStep{
			// Create with 3 records to test ordering stability
			{
				Config: testAccDNSZoneRrsetConfigMultipleRecords(zoneName, rrsetName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_dns_zone_rrset.test",
						tfjsonpath.New("resource_records"),
						knownvalue.ListSizeExact(3)),
				},
			},
			// Verify no drift with multiple records (API may return in different order)
			{
				Config: testAccDNSZoneRrsetConfigMultipleRecords(zoneName, rrsetName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

// --- CheckDestroy ---

func testAccCheckDNSZoneRrsetDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "gcore_dns_zone" {
			// Zone destroy is checked separately; skip here
			res := new(http.Response)
			_, err := client.DNS.Zones.Get(
				context.Background(),
				rs.Primary.Attributes["name"],
				option.WithResponseBodyInto(&res),
			)
			if res != nil && res.StatusCode == 404 {
				continue
			}
			if err == nil {
				return fmt.Errorf("dns zone %s still exists", rs.Primary.Attributes["name"])
			}
			if !acctest.IsNotFoundError(err) {
				return fmt.Errorf("error checking dns zone deletion: %w", err)
			}
			continue
		}

		if rs.Type != "gcore_dns_zone_rrset" {
			continue
		}

		rrsetType := rs.Primary.Attributes["rrset_type"]
		zoneName := rs.Primary.Attributes["zone_name"]
		rrsetName := rs.Primary.Attributes["rrset_name"]

		res := new(http.Response)
		_, err := client.DNS.Zones.Rrsets.Get(
			context.Background(),
			rrsetType,
			dns.ZoneRrsetGetParams{
				ZoneName:  zoneName,
				RrsetName: rrsetName,
			},
			option.WithResponseBodyInto(&res),
		)
		if res != nil && res.StatusCode == 404 {
			continue // successfully deleted
		}
		if err == nil {
			return fmt.Errorf("dns zone rrset %s/%s/%s still exists", zoneName, rrsetName, rrsetType)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking dns zone rrset deletion: %w", err)
		}
	}
	return nil
}

// --- Config Helpers ---

func testAccDNSZoneRrsetConfig(zoneName, rrsetName string) string {
	return fmt.Sprintf(`
resource "gcore_dns_zone" "test" {
  name = %[1]q
}

resource "gcore_dns_zone_rrset" "test" {
  zone_name  = gcore_dns_zone.test.name
  rrset_name = %[2]q
  rrset_type = "A"
  ttl        = 300

  resource_records = [
    {
      content = ["192.168.1.1"]
      enabled = true
    },
  ]
}`, zoneName, rrsetName)
}

func testAccDNSZoneRrsetConfigUpdated(zoneName, rrsetName string) string {
	return fmt.Sprintf(`
resource "gcore_dns_zone" "test" {
  name = %[1]q
}

resource "gcore_dns_zone_rrset" "test" {
  zone_name  = gcore_dns_zone.test.name
  rrset_name = %[2]q
  rrset_type = "A"
  ttl        = 600

  resource_records = [
    {
      content = ["192.168.1.1"]
      enabled = true
    },
    {
      content = ["10.0.0.1"]
      enabled = true
    },
  ]
}`, zoneName, rrsetName)
}

func testAccDNSZoneRrsetConfigWithMeta(zoneName, rrsetName string) string {
	return fmt.Sprintf(`
resource "gcore_dns_zone" "test" {
  name = %[1]q
}

resource "gcore_dns_zone_rrset" "test" {
  zone_name  = gcore_dns_zone.test.name
  rrset_name = %[2]q
  rrset_type = "A"
  ttl        = 300

  resource_records = [
    {
      content = ["192.168.1.1"]
      enabled = true
      meta = {
        latlong = "[51.5,31.5]"
        asn     = "[12345]"
        ip      = "[\"192.168.1.0/24\"]"
      }
    },
    {
      content = ["10.0.0.1"]
      enabled = true
      meta = {
        latlong = "[40.0,20.0]"
        asn     = "[67890]"
      }
    },
  ]

  pickers = [
    {
      type = "geodns"
    },
    {
      type = "default"
    },
  ]
}`, zoneName, rrsetName)
}

func testAccDNSZoneRrsetConfigWithMetaUpdated(zoneName, rrsetName string) string {
	return fmt.Sprintf(`
resource "gcore_dns_zone" "test" {
  name = %[1]q
}

resource "gcore_dns_zone_rrset" "test" {
  zone_name  = gcore_dns_zone.test.name
  rrset_name = %[2]q
  rrset_type = "A"
  ttl        = 300

  resource_records = [
    {
      content = ["192.168.1.1"]
      enabled = true
      meta = {
        latlong = "[52.0,32.0]"
        asn     = "[12345]"
        ip      = "[\"192.168.1.0/24\"]"
      }
    },
    {
      content = ["10.0.0.1"]
      enabled = true
      meta = {
        latlong = "[41.0,21.0]"
        asn     = "[67890]"
      }
    },
  ]

  pickers = [
    {
      type = "geodns"
    },
    {
      type = "default"
    },
  ]
}`, zoneName, rrsetName)
}

func testAccDNSZoneRrsetConfigMultipleRecords(zoneName, rrsetName string) string {
	return fmt.Sprintf(`
resource "gcore_dns_zone" "test" {
  name = %[1]q
}

resource "gcore_dns_zone_rrset" "test" {
  zone_name  = gcore_dns_zone.test.name
  rrset_name = %[2]q
  rrset_type = "A"
  ttl        = 300

  resource_records = [
    {
      content = ["192.168.1.1"]
      enabled = true
    },
    {
      content = ["10.0.0.1"]
      enabled = true
    },
    {
      content = ["172.16.0.1"]
      enabled = true
    },
  ]
}`, zoneName, rrsetName)
}
