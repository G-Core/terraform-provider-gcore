package cloud_gpu_virtual_cluster_image_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

// listTestImageURL is a small, publicly-accessible Linux cloud image used to
// upload a GPU virtual image for the list data source test.
const listTestImageURL = "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud-2003.qcow2.xz"

// GPU virtual images are custom uploads, so create one first and then read the
// list data source (deferred via depends_on) to assert it is returned.

func TestAccCloudGPUVirtualClusterImagesDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudGPUVirtualClusterImagesDataSourceConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_gpu_virtual_cluster_images.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_gpu_virtual_cluster_images.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_gpu_virtual_cluster_images.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("status"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func testAccCloudGPUVirtualClusterImagesDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_gpu_virtual_cluster_image" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  url_wo     = %[4]q
}

data "gcore_cloud_gpu_virtual_cluster_images" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  depends_on = [gcore_cloud_gpu_virtual_cluster_image.test]
}
`, acctest.ProjectID(), acctest.RegionID(), name, listTestImageURL)
}
