package cloud_load_balancer_pool_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

const poolAddr = "gcore_cloud_load_balancer_pool.test"

// healthmonitorIsNull asserts the whole attribute is null, not merely that a
// leaf is missing. TestCheckNoResourceAttr("healthmonitor.type") cannot tell
// the two apart: the state shim omits null-valued leaves from the flatmap, so
// it passes just as happily for an object whose every attribute is null — which
// is exactly the value this test has to reject.
var healthmonitorIsNull = statecheck.ExpectKnownValue(poolAddr, tfjsonpath.New("healthmonitor"), knownvalue.Null())

// checkHealthmonitorIsNull is the same assertion for steps that cannot carry a
// state check — RefreshState steps take no Config, and ConfigStateChecks is
// only valid alongside one. A null object contributes no keys to the flatmap at
// all, while an object of nulls still contributes its "healthmonitor.%" count,
// so the presence of any healthmonitor key is the phantom.
func checkHealthmonitorIsNull(s *terraform.State) error {
	rs, ok := s.RootModule().Resources[poolAddr]
	if !ok {
		return fmt.Errorf("%s not found in state", poolAddr)
	}
	for key, value := range rs.Primary.Attributes {
		if strings.HasPrefix(key, "healthmonitor.") {
			return fmt.Errorf("expected a null healthmonitor, found %s = %q in state", key, value)
		}
	}
	return nil
}

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
				Check:  resource.TestCheckResourceAttr(poolAddr, "name", fmt.Sprintf("%s-pool", rName)),
				ConfigStateChecks: []statecheck.StateCheck{
					healthmonitorIsNull,
				},
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
					resource.TestCheckResourceAttr(poolAddr, "healthmonitor.type", "TCP"),
					resource.TestCheckResourceAttr(poolAddr, "healthmonitor.delay", "10"),
				),
			},
			{
				// Step 4: remove the block again. Omitting it plans null, so Update
				// dispatches to the health monitor DELETE endpoint and the monitor
				// is actually gone rather than silently retained. The apply used to
				// fail here — the decoded response left an object of nulls in state
				// where the plan promised null.
				Config: testAccCloudLoadBalancerPoolConfigNoHealthmonitor(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					healthmonitorIsNull,
				},
			},
			{
				// Step 5: the removal must settle — no perpetual diff afterwards.
				// Re-applying used to re-issue the DELETE for the monitor that was
				// already gone, which the API answers with 400 "not found".
				Config:   testAccCloudLoadBalancerPoolConfigNoHealthmonitor(rName),
				PlanOnly: true,
			},
			{
				// Step 6: refresh has to agree. Read decodes the response into the
				// prior state, a path Update never exercises, and it is where the
				// perpetual healthmonitor = {} -> null diff lived.
				RefreshState: true,
				Check:        checkHealthmonitorIsNull,
			},
			{
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
