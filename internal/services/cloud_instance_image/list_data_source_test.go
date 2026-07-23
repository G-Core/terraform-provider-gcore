package cloud_instance_image_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

// Public instance images exist in every region, so the list data source is
// read directly (filtered to public) without creating a resource.

func TestAccCloudInstanceImagesDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstanceImagesDataSourceConfigVisibility("public"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_instance_images.test",
						tfjsonpath.New("visibility"), knownvalue.StringExact("public")),
					statecheck.ExpectKnownValue("data.gcore_cloud_instance_images.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_instance_images.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_instance_images.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_instance_images.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("visibility"), knownvalue.StringExact("public")),
				},
			},
		},
	})
}

func TestAccCloudInstanceImagesDataSource_maxItems(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstanceImagesDataSourceConfigMaxItems(1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_instance_images.test",
						tfjsonpath.New("max_items"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("data.gcore_cloud_instance_images.test",
						tfjsonpath.New("items"), knownvalue.ListSizeExact(1)),
				},
			},
		},
	})
}

func testAccCloudInstanceImagesDataSourceConfigVisibility(visibility string) string {
	return fmt.Sprintf(`
data "gcore_cloud_instance_images" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  visibility = %[3]q
}
`, acctest.ProjectID(), acctest.RegionID(), visibility)
}

func testAccCloudInstanceImagesDataSourceConfigMaxItems(n int) string {
	return fmt.Sprintf(`
data "gcore_cloud_instance_images" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  visibility = "public"
  max_items  = %[3]d
}
`, acctest.ProjectID(), acctest.RegionID(), n)
}
