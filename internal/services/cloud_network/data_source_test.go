package cloud_network_test

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

func TestAccCloudNetworkDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudNetworkDataSourceConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_network.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.CompareValuePairs(
						"gcore_cloud_network.test", tfjsonpath.New("id"),
						"data.gcore_cloud_network.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func testAccCloudNetworkDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_network" "test" {
  project_id    = %[1]s
  region_id     = %[2]s
  name          = %[3]q
  create_router = false
  type          = "vxlan"
}

data "gcore_cloud_network" "test" {
  project_id  = %[1]s
  region_id   = %[2]s
  network_id  = gcore_cloud_network.test.id
}`, acctest.ProjectID(), acctest.RegionID(), name)
}
