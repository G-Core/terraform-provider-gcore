package cloud_region_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func TestAccCloudRegionDataSource_byID(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRegionDataSourceConfigByID(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_region.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_region.test",
						tfjsonpath.New("display_name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_region.test",
						tfjsonpath.New("state"), knownvalue.StringExact("ACTIVE")),
					statecheck.ExpectKnownValue("data.gcore_cloud_region.test",
						tfjsonpath.New("endpoint_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_region.test",
						tfjsonpath.New("zone"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_region.test",
						tfjsonpath.New("country"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccCloudRegionDataSource_withVolumeTypes(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRegionDataSourceConfigWithVolumeTypes(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_region.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_region.test",
						tfjsonpath.New("available_volume_types"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func testAccCloudRegionDataSourceConfigByID() string {
	return fmt.Sprintf(`
data "gcore_cloud_region" "test" {
  region_id = %s
}`, acctest.RegionID())
}

func testAccCloudRegionDataSourceConfigWithVolumeTypes() string {
	return fmt.Sprintf(`
data "gcore_cloud_region" "test" {
  region_id        = %s
  show_volume_types = true
}`, acctest.RegionID())
}
