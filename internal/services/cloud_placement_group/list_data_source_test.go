package cloud_placement_group_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

// Server groups do not exist by default, so create one first and then read the
// list data source (deferred via depends_on) to assert it is returned.

func TestAccCloudPlacementGroupsDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudPlacementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudPlacementGroupsDataSourceConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_placement_groups.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("servergroup_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_placement_groups.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_placement_groups.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("policy"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func testAccCloudPlacementGroupsDataSourceConfig(name string) string {
	return fmt.Sprintf(`
%[1]s

data "gcore_cloud_placement_groups" "test" {
  project_id = %[2]s
  region_id  = %[3]s
  depends_on = [gcore_cloud_placement_group.test]
}
`, testAccCloudPlacementGroupConfig(name, "anti-affinity"), acctest.ProjectID(), acctest.RegionID())
}
