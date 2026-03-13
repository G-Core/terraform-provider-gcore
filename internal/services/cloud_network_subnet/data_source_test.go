package cloud_network_subnet_test

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

func TestAccCloudNetworkSubnetDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudNetworkSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudNetworkSubnetDataSourceConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_network_subnet.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("data.gcore_cloud_network_subnet.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_network_subnet.test",
						tfjsonpath.New("cidr"), knownvalue.StringExact("10.0.0.0/24")),
					statecheck.CompareValuePairs(
						"gcore_cloud_network_subnet.test", tfjsonpath.New("id"),
						"data.gcore_cloud_network_subnet.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func testAccCloudNetworkSubnetDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_network" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "tf-test-net-%[3]s"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "test" {
  project_id  = %[1]s
  region_id   = %[2]s
  name        = %[3]q
  network_id  = gcore_cloud_network.test.id
  cidr        = "10.0.0.0/24"
  enable_dhcp = true
}

data "gcore_cloud_network_subnet" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  subnet_id  = gcore_cloud_network_subnet.test.id
}`, acctest.ProjectID(), acctest.RegionID(), name)
}
