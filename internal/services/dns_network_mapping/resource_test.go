package dns_network_mapping_test

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

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func TestAccDNSNetworkMapping_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDNSNetworkMappingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSNetworkMappingConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_dns_network_mapping.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_dns_network_mapping.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_dns_network_mapping.test",
						tfjsonpath.New("mapping"), knownvalue.ListSizeExact(2)),
				},
			},
		},
	})
}

func TestAccDNSNetworkMapping_full(t *testing.T) {
	rName := acctest.RandomName()

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDNSNetworkMappingDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with 2 mapping entries
			{
				Config: testAccDNSNetworkMappingConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_dns_network_mapping.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_dns_network_mapping.test",
						tfjsonpath.New("mapping"), knownvalue.ListSizeExact(2)),
					statecheck.ExpectKnownValue("gcore_dns_network_mapping.test",
						tfjsonpath.New("mapping").AtSliceIndex(0).AtMapKey("tags"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("dev"),
							knownvalue.StringExact("dc1"),
						})),
					statecheck.ExpectKnownValue("gcore_dns_network_mapping.test",
						tfjsonpath.New("mapping").AtSliceIndex(0).AtMapKey("cidr4"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("10.0.0.0/24"),
						})),
					statecheck.ExpectKnownValue("gcore_dns_network_mapping.test",
						tfjsonpath.New("mapping").AtSliceIndex(0).AtMapKey("cidr6"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("2001:db8:1::/48"),
						})),
					compareIDSame.AddStateValue(
						"gcore_dns_network_mapping.test",
						tfjsonpath.New("id"),
					),
				},
			},
			// Step 2: Update to 1 mapping entry with modified values
			{
				Config: testAccDNSNetworkMappingConfigUpdated(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_dns_network_mapping.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_dns_network_mapping.test",
						tfjsonpath.New("mapping"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue("gcore_dns_network_mapping.test",
						tfjsonpath.New("mapping").AtSliceIndex(0).AtMapKey("tags"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("dev"),
							knownvalue.StringExact("dc1"),
							knownvalue.StringExact("updated"),
						})),
					statecheck.ExpectKnownValue("gcore_dns_network_mapping.test",
						tfjsonpath.New("mapping").AtSliceIndex(0).AtMapKey("cidr4"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("10.0.0.0/24"),
							knownvalue.StringExact("10.0.1.0/24"),
						})),
					statecheck.ExpectKnownValue("gcore_dns_network_mapping.test",
						tfjsonpath.New("mapping").AtSliceIndex(0).AtMapKey("cidr6"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("2001:db8:2::/48"),
						})),
					// ID should remain the same after update (in-place)
					compareIDSame.AddStateValue(
						"gcore_dns_network_mapping.test",
						tfjsonpath.New("id"),
					),
				},
			},
			// Step 3: Import
			{
				ResourceName:    "gcore_dns_network_mapping.test",
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithID,
			},
		},
	})
}

func testAccCheckDNSNetworkMappingDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_dns_network_mapping" {
			continue
		}

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing ID %q: %w", rs.Primary.ID, err)
		}

		_, err = client.DNS.NetworkMappings.Get(context.Background(), id)
		if err == nil {
			return fmt.Errorf("dns_network_mapping %s still exists", rs.Primary.ID)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking dns_network_mapping deletion: %w", err)
		}
	}
	return nil
}

func testAccDNSNetworkMappingConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_dns_network_mapping" "test" {
  name = %[1]q
  mapping = [
    {
      tags  = ["dev", "dc1"]
      cidr4 = ["10.0.0.0/24"]
      cidr6 = ["2001:db8:1::/48"]
    },
    {
      tags  = ["prod", "dc2"]
      cidr4 = ["10.1.0.0/24"]
    }
  ]
}`, name)
}

func testAccDNSNetworkMappingConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "gcore_dns_network_mapping" "test" {
  name = %[1]q
  mapping = [
    {
      tags  = ["dev", "dc1", "updated"]
      cidr4 = ["10.0.0.0/24", "10.0.1.0/24"]
      cidr6 = ["2001:db8:2::/48"]
    }
  ]
}`, name)
}

func testAccDNSNetworkMappingDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_dns_network_mapping" "test" {
  name = %[1]q
  mapping = [
    {
      tags  = ["ds-test"]
      cidr4 = ["10.2.0.0/24"]
    }
  ]
}

data "gcore_dns_network_mapping" "test" {
  id = gcore_dns_network_mapping.test.id
}`, name)
}
