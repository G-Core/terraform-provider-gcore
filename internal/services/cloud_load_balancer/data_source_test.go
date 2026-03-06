package cloud_load_balancer_test

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

func TestAccCloudLoadBalancerDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerDataSourceConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_load_balancer.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					// verify data source reads the same ID as the resource
					statecheck.CompareValuePairs(
						"gcore_cloud_load_balancer.test", tfjsonpath.New("id"),
						"data.gcore_cloud_load_balancer.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
					// verify data source reads the same vip_address as the resource
					statecheck.CompareValuePairs(
						"gcore_cloud_load_balancer.test", tfjsonpath.New("vip_address"),
						"data.gcore_cloud_load_balancer.test", tfjsonpath.New("vip_address"),
						compare.ValuesSame(),
					),
					statecheck.ExpectKnownValue("data.gcore_cloud_load_balancer.test",
						tfjsonpath.New("admin_state_up"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_load_balancer.test",
						tfjsonpath.New("provisioning_status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_load_balancer.test",
						tfjsonpath.New("operating_status"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func testAccCloudLoadBalancerDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_load_balancer" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
}

data "gcore_cloud_load_balancer" "test" {
  project_id       = %[1]s
  region_id        = %[2]s
  load_balancer_id = gcore_cloud_load_balancer.test.id
}`, acctest.ProjectID(), acctest.RegionID(), name)
}
