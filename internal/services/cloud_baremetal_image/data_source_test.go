package cloud_baremetal_image_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

// Bare metal images are pre-existing (public/shared) images in a bare-metal
// capable region, so these tests read them directly without creating a resource.

func TestAccCloudBaremetalImagesDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudBaremetalImagesDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_baremetal_images.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_baremetal_images.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_baremetal_images.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("architecture"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_baremetal_images.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("is_baremetal"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccCloudBaremetalImagesDataSource_visibilityPublic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudBaremetalImagesDataSourceConfigVisibility("public"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_baremetal_images.test",
						tfjsonpath.New("visibility"), knownvalue.StringExact("public")),
					statecheck.ExpectKnownValue("data.gcore_cloud_baremetal_images.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("visibility"), knownvalue.StringExact("public")),
				},
			},
		},
	})
}

// The ticket's primary use case: look up an image by name instead of hardcoding
// its ID. The `name` filter is a case-insensitive substring match, so assert the
// echoed filter and that the first returned item's name actually contains it.
func TestAccCloudBaremetalImagesDataSource_filterByName(t *testing.T) {
	const name = "ubuntu"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudBaremetalImagesDataSourceConfigName(name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_baremetal_images.test",
						tfjsonpath.New("name"), knownvalue.StringExact(name)),
					// items[0].name must match the requested substring (case-insensitive)
					statecheck.ExpectKnownValue("data.gcore_cloud_baremetal_images.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringRegexp(regexp.MustCompile("(?i)"+name))),
				},
			},
		},
	})
}

func TestAccCloudBaremetalImagesDataSource_maxItems(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudBaremetalImagesDataSourceConfigMaxItems(1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_baremetal_images.test",
						tfjsonpath.New("max_items"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("data.gcore_cloud_baremetal_images.test",
						tfjsonpath.New("items"), knownvalue.ListSizeExact(1)),
				},
			},
		},
	})
}

func testAccCloudBaremetalImagesDataSourceConfig() string {
	return fmt.Sprintf(`
data "gcore_cloud_baremetal_images" "test" {
  project_id = %[1]s
  region_id  = %[2]s
}
`, acctest.ProjectID(), acctest.RegionID())
}

func testAccCloudBaremetalImagesDataSourceConfigVisibility(visibility string) string {
	return fmt.Sprintf(`
data "gcore_cloud_baremetal_images" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  visibility = %[3]q
}
`, acctest.ProjectID(), acctest.RegionID(), visibility)
}

func testAccCloudBaremetalImagesDataSourceConfigName(name string) string {
	return fmt.Sprintf(`
data "gcore_cloud_baremetal_images" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
}
`, acctest.ProjectID(), acctest.RegionID(), name)
}

func testAccCloudBaremetalImagesDataSourceConfigMaxItems(n int) string {
	return fmt.Sprintf(`
data "gcore_cloud_baremetal_images" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  max_items  = %[3]d
}
`, acctest.ProjectID(), acctest.RegionID(), n)
}
