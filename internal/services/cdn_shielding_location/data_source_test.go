package cdn_shielding_location_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func TestAccCDNShieldingLocationsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNShieldingLocationsDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cdn_shielding_locations.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cdn_shielding_locations.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("city"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cdn_shielding_locations.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("country"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cdn_shielding_locations.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("datacenter"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccCDNShieldingLocationsDataSource_maxItems(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNShieldingLocationsDataSourceConfigMaxItems(1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cdn_shielding_locations.test",
						tfjsonpath.New("max_items"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("data.gcore_cdn_shielding_locations.test",
						tfjsonpath.New("items"), knownvalue.ListSizeExact(1)),
				},
			},
		},
	})
}

func testAccCDNShieldingLocationsDataSourceConfig() string {
	return `
data "gcore_cdn_shielding_locations" "test" {}
`
}

func testAccCDNShieldingLocationsDataSourceConfigMaxItems(n int) string {
	return fmt.Sprintf(`
data "gcore_cdn_shielding_locations" "test" {
  max_items = %d
}
`, n)
}
