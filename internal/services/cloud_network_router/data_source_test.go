package cloud_network_router_test

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

func TestAccCloudNetworkRouterDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudNetworkRouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudNetworkRouterDataSourceConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_network_router.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.CompareValuePairs(
						"gcore_cloud_network_router.test", tfjsonpath.New("id"),
						"data.gcore_cloud_network_router.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func testAccCloudNetworkRouterDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_network_router" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
}

data "gcore_cloud_network_router" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  router_id  = gcore_cloud_network_router.test.id
}`, acctest.ProjectID(), acctest.RegionID(), name)
}
