package cloud_floating_ip_test

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

func TestAccCloudFloatingIPDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFloatingIPDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_floating_ip.test",
						tfjsonpath.New("floating_ip_address"), knownvalue.NotNull()),
					statecheck.CompareValuePairs(
						"gcore_cloud_floating_ip.test", tfjsonpath.New("id"),
						"data.gcore_cloud_floating_ip.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
					statecheck.CompareValuePairs(
						"gcore_cloud_floating_ip.test", tfjsonpath.New("floating_ip_address"),
						"data.gcore_cloud_floating_ip.test", tfjsonpath.New("floating_ip_address"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func testAccCloudFloatingIPDataSourceConfig() string {
	return fmt.Sprintf(`
resource "gcore_cloud_floating_ip" "test" {
  project_id = %[1]s
  region_id  = %[2]s
}

data "gcore_cloud_floating_ip" "test" {
  project_id     = %[1]s
  region_id      = %[2]s
  floating_ip_id = gcore_cloud_floating_ip.test.id
}`, acctest.ProjectID(), acctest.RegionID())
}
