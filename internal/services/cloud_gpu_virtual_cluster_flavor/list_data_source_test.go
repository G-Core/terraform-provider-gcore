package cloud_gpu_virtual_cluster_flavor_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

// GPU virtual flavors are region-global and always present, so the list data
// source is read directly without creating any resource.

func TestAccCloudGPUVirtualClusterFlavorsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudGPUVirtualClusterFlavorsDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_gpu_virtual_cluster_flavors.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_gpu_virtual_cluster_flavors.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("architecture"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_gpu_virtual_cluster_flavors.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("capacity"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccCloudGPUVirtualClusterFlavorsDataSource_maxItems(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudGPUVirtualClusterFlavorsDataSourceConfigMaxItems(1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_gpu_virtual_cluster_flavors.test",
						tfjsonpath.New("max_items"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("data.gcore_cloud_gpu_virtual_cluster_flavors.test",
						tfjsonpath.New("items"), knownvalue.ListSizeExact(1)),
				},
			},
		},
	})
}

func testAccCloudGPUVirtualClusterFlavorsDataSourceConfig() string {
	return fmt.Sprintf(`
data "gcore_cloud_gpu_virtual_cluster_flavors" "test" {
  project_id = %[1]s
  region_id  = %[2]s
}
`, acctest.ProjectID(), acctest.RegionID())
}

func testAccCloudGPUVirtualClusterFlavorsDataSourceConfigMaxItems(n int) string {
	return fmt.Sprintf(`
data "gcore_cloud_gpu_virtual_cluster_flavors" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  max_items  = %[3]d
}
`, acctest.ProjectID(), acctest.RegionID(), n)
}
