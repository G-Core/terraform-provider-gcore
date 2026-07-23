package cloud_gpu_baremetal_cluster_image_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

// Public GPU bare metal images exist in the region, so the list data source is
// read directly without creating a resource. (The image-upload resource is not
// exercised here; prerequisites, when needed, are provisioned out-of-band.)

func TestAccCloudGPUBaremetalClusterImagesDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudGPUBaremetalClusterImagesDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_gpu_baremetal_cluster_images.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_gpu_baremetal_cluster_images.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_gpu_baremetal_cluster_images.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("status"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccCloudGPUBaremetalClusterImagesDataSource_maxItems(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudGPUBaremetalClusterImagesDataSourceConfigMaxItems(1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_gpu_baremetal_cluster_images.test",
						tfjsonpath.New("max_items"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("data.gcore_cloud_gpu_baremetal_cluster_images.test",
						tfjsonpath.New("items"), knownvalue.ListSizeExact(1)),
				},
			},
		},
	})
}

func testAccCloudGPUBaremetalClusterImagesDataSourceConfig() string {
	return fmt.Sprintf(`
data "gcore_cloud_gpu_baremetal_cluster_images" "test" {
  project_id = %[1]s
  region_id  = %[2]s
}
`, acctest.ProjectID(), acctest.RegionID())
}

func testAccCloudGPUBaremetalClusterImagesDataSourceConfigMaxItems(n int) string {
	return fmt.Sprintf(`
data "gcore_cloud_gpu_baremetal_cluster_images" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  max_items  = %[3]d
}
`, acctest.ProjectID(), acctest.RegionID(), n)
}
