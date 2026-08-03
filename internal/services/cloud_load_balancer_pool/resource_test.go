package cloud_load_balancer_pool_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

// TestAccCloudLoadBalancerPool_noHealthmonitor is a regression test: a pool
// created WITHOUT a healthmonitor block previously failed at apply with
// "Value Conversion Error ... Path: healthmonitor", because the omitted
// Computed+Optional block plans as unknown and the model stores it as a plain
// Go pointer that cannot hold an unknown value.
func TestAccCloudLoadBalancerPool_noHealthmonitor(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerPoolDestroy,
		Steps: []resource.TestStep{
			{
				// Step 1: create with NO healthmonitor block — the bug repro.
				Config: testAccCloudLoadBalancerPoolConfigNoHealthmonitor(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gcore_cloud_load_balancer_pool.test", "name", fmt.Sprintf("%s-pool", rName)),
					resource.TestCheckNoResourceAttr("gcore_cloud_load_balancer_pool.test", "healthmonitor.type"),
				),
			},
			{
				// Step 2: same config — must be a no-drift plan.
				Config:   testAccCloudLoadBalancerPoolConfigNoHealthmonitor(rName),
				PlanOnly: true,
			},
			{
				// Step 3: add a healthmonitor block — in-place update, no replace.
				Config: testAccCloudLoadBalancerPoolConfigWithHealthmonitor(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gcore_cloud_load_balancer_pool.test", "healthmonitor.type", "TCP"),
					resource.TestCheckResourceAttr("gcore_cloud_load_balancer_pool.test", "healthmonitor.delay", "10"),
				),
			},
			{
				// Step 4: remove the block again. Documents current semantics:
				// Computed retention keeps the monitor; the plan must be empty
				// (removal-by-omission is a pre-existing no-op, tracked separately).
				Config:   testAccCloudLoadBalancerPoolConfigNoHealthmonitor(rName),
				PlanOnly: true,
			},
		},
	})
}

func testAccCheckCloudLoadBalancerPoolDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cloud_load_balancer_pool" {
			continue
		}

		projectID, err := strconv.ParseInt(rs.Primary.Attributes["project_id"], 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing project_id: %w", err)
		}
		regionID, err := strconv.ParseInt(rs.Primary.Attributes["region_id"], 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing region_id: %w", err)
		}

		_, err = client.Cloud.LoadBalancers.Pools.Get(context.Background(), rs.Primary.ID, cloud.LoadBalancerPoolGetParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})
		if err != nil {
			if acctest.IsNotFoundError(err) {
				continue
			}
			return fmt.Errorf("error getting pool: %w", err)
		}
		return fmt.Errorf("pool %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccCloudLoadBalancerPoolConfigNoHealthmonitor(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_load_balancer" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-lb"
  flavor     = "lb1-1-2"
}

resource "gcore_cloud_load_balancer_pool" "test" {
  project_id       = %[1]s
  region_id        = %[2]s
  name             = "%[3]s-pool"
  lb_algorithm     = "ROUND_ROBIN"
  protocol         = "HTTP"
  load_balancer_id = gcore_cloud_load_balancer.test.id
}`, acctest.ProjectID(), acctest.RegionID(), name)
}

func testAccCloudLoadBalancerPoolConfigWithHealthmonitor(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_load_balancer" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-lb"
  flavor     = "lb1-1-2"
}

resource "gcore_cloud_load_balancer_pool" "test" {
  project_id       = %[1]s
  region_id        = %[2]s
  name             = "%[3]s-pool"
  lb_algorithm     = "ROUND_ROBIN"
  protocol         = "HTTP"
  load_balancer_id = gcore_cloud_load_balancer.test.id

  healthmonitor = {
    type        = "TCP"
    delay       = 10
    max_retries = 3
    timeout     = 5
  }
}`, acctest.ProjectID(), acctest.RegionID(), name)
}
