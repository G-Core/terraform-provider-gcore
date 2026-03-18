package dns_network_mapping_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func TestAccDNSNetworkMappingDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDNSNetworkMappingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSNetworkMappingDataSourceConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_dns_network_mapping.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("data.gcore_dns_network_mapping.test",
						tfjsonpath.New("mapping"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue("data.gcore_dns_network_mapping.test",
						tfjsonpath.New("mapping").AtSliceIndex(0).AtMapKey("tags"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("ds-test"),
						})),
					statecheck.ExpectKnownValue("data.gcore_dns_network_mapping.test",
						tfjsonpath.New("mapping").AtSliceIndex(0).AtMapKey("cidr4"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("10.2.0.0/24"),
						})),
					// Verify data source reads the same ID as the resource
					statecheck.CompareValuePairs(
						"gcore_dns_network_mapping.test", tfjsonpath.New("id"),
						"data.gcore_dns_network_mapping.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}
