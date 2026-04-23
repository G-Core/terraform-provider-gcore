package storage_location_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func TestAccStorageLocationsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageLocationsDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_storage_locations.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("address"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_storage_locations.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_storage_locations.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("technical_name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_storage_locations.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("type"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccStorageLocationsDataSource_maxItems(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageLocationsDataSourceConfigMaxItems(1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_storage_locations.test",
						tfjsonpath.New("max_items"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("data.gcore_storage_locations.test",
						tfjsonpath.New("items"), knownvalue.ListSizeExact(1)),
				},
			},
		},
	})
}

func testAccStorageLocationsDataSourceConfig() string {
	return `
data "gcore_storage_locations" "test" {}
`
}

func testAccStorageLocationsDataSourceConfigMaxItems(n int) string {
	return fmt.Sprintf(`
data "gcore_storage_locations" "test" {
  max_items = %d
}
`, n)
}
