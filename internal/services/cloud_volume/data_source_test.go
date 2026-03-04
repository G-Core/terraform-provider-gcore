package cloud_volume_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func TestAccCloudVolumeDataSource_byID(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudVolumeDataSourceConfigByID(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_volume.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.CompareValuePairs(
						"gcore_cloud_volume.test", tfjsonpath.New("id"),
						"data.gcore_cloud_volume.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
					statecheck.CompareValuePairs(
						"gcore_cloud_volume.test", tfjsonpath.New("size"),
						"data.gcore_cloud_volume.test", tfjsonpath.New("size"),
						compare.ValuesSame(),
					),
					statecheck.CompareValuePairs(
						"gcore_cloud_volume.test", tfjsonpath.New("status"),
						"data.gcore_cloud_volume.test", tfjsonpath.New("status"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func testAccCloudVolumeDataSourceConfigByID(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_volume" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  source     = "new-volume"
  name       = %[3]q
  size       = 1
}

data "gcore_cloud_volume" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  volume_id  = gcore_cloud_volume.test.id
}`, acctest.ProjectID(), acctest.RegionID(), name)
}
