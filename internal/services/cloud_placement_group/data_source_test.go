package cloud_placement_group_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/stainless-sdks/gcore-terraform/internal/acctest"
)

func TestAccCloudPlacementGroupDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudPlacementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudPlacementGroupDataSourceConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_placement_group.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("data.gcore_cloud_placement_group.test",
						tfjsonpath.New("policy"), knownvalue.StringExact("affinity")),
					statecheck.CompareValuePairs(
						"gcore_cloud_placement_group.test", tfjsonpath.New("id"),
						"data.gcore_cloud_placement_group.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func testAccCloudPlacementGroupDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_placement_group" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  policy     = "affinity"
}

data "gcore_cloud_placement_group" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  group_id   = gcore_cloud_placement_group.test.servergroup_id
}`, acctest.ProjectID(), acctest.RegionID(), name)
}
