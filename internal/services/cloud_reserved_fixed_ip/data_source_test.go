package cloud_reserved_fixed_ip_test

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

func TestAccCloudReservedFixedIPDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudReservedFixedIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudReservedFixedIPDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_reserved_fixed_ip.test",
						tfjsonpath.New("port_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_reserved_fixed_ip.test",
						tfjsonpath.New("fixed_ip_address"), knownvalue.NotNull()),
					statecheck.CompareValuePairs(
						"gcore_cloud_reserved_fixed_ip.test", tfjsonpath.New("port_id"),
						"data.gcore_cloud_reserved_fixed_ip.test", tfjsonpath.New("port_id"),
						compare.ValuesSame(),
					),
					statecheck.CompareValuePairs(
						"gcore_cloud_reserved_fixed_ip.test", tfjsonpath.New("fixed_ip_address"),
						"data.gcore_cloud_reserved_fixed_ip.test", tfjsonpath.New("fixed_ip_address"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func testAccCloudReservedFixedIPDataSourceConfig() string {
	return fmt.Sprintf(`
resource "gcore_cloud_reserved_fixed_ip" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  type       = "external"
}

data "gcore_cloud_reserved_fixed_ip" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  port_id    = gcore_cloud_reserved_fixed_ip.test.port_id
}`, acctest.ProjectID(), acctest.RegionID())
}
