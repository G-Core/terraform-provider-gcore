package cloud_k8s_flavor_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

// K8s-eligible flavors are region-global and always present, so the list data
// source is read directly without creating any resource.

func TestAccCloudK8SFlavorsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudK8SFlavorsDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_k8s_flavors.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("flavor_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_k8s_flavors.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("flavor_name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_k8s_flavors.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("vcpus"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_k8s_flavors.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("ram"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccCloudK8SFlavorsDataSource_maxItems(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudK8SFlavorsDataSourceConfigMaxItems(1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_k8s_flavors.test",
						tfjsonpath.New("max_items"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("data.gcore_cloud_k8s_flavors.test",
						tfjsonpath.New("items"), knownvalue.ListSizeExact(1)),
				},
			},
		},
	})
}

func testAccCloudK8SFlavorsDataSourceConfig() string {
	return fmt.Sprintf(`
data "gcore_cloud_k8s_flavors" "test" {
  project_id = %[1]s
  region_id  = %[2]s
}
`, acctest.ProjectID(), acctest.RegionID())
}

func testAccCloudK8SFlavorsDataSourceConfigMaxItems(n int) string {
	return fmt.Sprintf(`
data "gcore_cloud_k8s_flavors" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  max_items  = %[3]d
}
`, acctest.ProjectID(), acctest.RegionID(), n)
}
