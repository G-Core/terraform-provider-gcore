package cloud_load_balancer_pool_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
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

// TestAccCloudLoadBalancerPool_importLBOnly: a pool created through
// load_balancer_id imports with the same config as a no-op plan.
// load_balancer_id is read back from loadbalancers[] and listener_id stays
// null because the pool has no listener.
func TestAccCloudLoadBalancerPool_importLBOnly(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerPoolConfigNoHealthmonitor(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(poolAddr, tfjsonpath.New("listener_id"), knownvalue.Null()),
					statecheck.CompareValuePairs(
						"gcore_cloud_load_balancer.test", tfjsonpath.New("id"),
						poolAddr, tfjsonpath.New("load_balancer_id"),
						compare.ValuesSame(),
					),
				},
			},
			{
				ResourceName:      poolAddr,
				ImportState:       true,
				ImportStateIdFunc: acctest.BuildImportID(poolAddr, "project_id", "region_id", "id"),
				ImportStateKind:   resource.ImportBlockWithID,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(poolAddr, plancheck.ResourceActionNoop),
						plancheck.ExpectKnownValue(poolAddr, tfjsonpath.New("listener_id"), knownvalue.Null()),
						plancheck.ExpectKnownValue(poolAddr, tfjsonpath.New("load_balancer_id"), knownvalue.NotNull()),
					},
				},
			},
		},
	})
}

// TestAccCloudLoadBalancerPool_importAttachedViaListener is the case a literal
// replace-unless-adopting on listener_id would have destroyed: a pool attached
// through a listener's default_pool_id, imported with a config that carries
// only load_balancer_id. listener_id is back-filled from listeners[] and, being
// Computed+Optional, is never compared against the missing config - the plan
// is a no-op and nothing cascades to the pool members.
func TestAccCloudLoadBalancerPool_importAttachedViaListener(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerPoolConfigDefaultPoolOfListener(rName),
			},
			{
				// The listener attachment happened after the pool's Create, so
				// state has to catch up through Read before the import compares.
				RefreshState: true,
			},
			{
				Config: testAccCloudLoadBalancerPoolConfigDefaultPoolOfListener(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs(
						"gcore_cloud_load_balancer_listener.test", tfjsonpath.New("id"),
						poolAddr, tfjsonpath.New("listener_id"),
						compare.ValuesSame(),
					),
				},
			},
			{
				ResourceName:      poolAddr,
				ImportState:       true,
				ImportStateIdFunc: acctest.BuildImportID(poolAddr, "project_id", "region_id", "id"),
				ImportStateKind:   resource.ImportBlockWithID,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(poolAddr, plancheck.ResourceActionNoop),
						plancheck.ExpectKnownValue(poolAddr, tfjsonpath.New("listener_id"), knownvalue.NotNull()),
						plancheck.ExpectKnownValue(poolAddr, tfjsonpath.New("load_balancer_id"), knownvalue.NotNull()),
					},
				},
			},
		},
	})
}

// TestAccCloudLoadBalancerPool_detachedListenerClearsListenerID: refresh
// records the listener that uses the pool as its default pool. Once that
// listener is gone, refresh clears listener_id and the plan stays empty.
func TestAccCloudLoadBalancerPool_detachedListenerClearsListenerID(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerPoolConfigDefaultPoolOfListener(rName),
			},
			{
				RefreshState: true,
				Check: resource.TestCheckResourceAttrPair(
					poolAddr, "listener_id", "gcore_cloud_load_balancer_listener.test", "id"),
			},
			{
				Config: testAccCloudLoadBalancerPoolConfigNoHealthmonitor(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				RefreshState: true,
				Check:        resource.TestCheckNoResourceAttr(poolAddr, "listener_id"),
			},
		},
	})
}

// TestAccCloudLoadBalancerPool_dayTwoListenerIDForcesReplace is the ticket's
// assertion on the pool: setting listener_id on a pool created without one
// must plan the replacement that actually attaches it - no PATCH re-parents a
// pool - with and without a refresh. There is no marker and no window here.
func TestAccCloudLoadBalancerPool_dayTwoListenerIDForcesReplace(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerPoolConfigWithDetachedListener(rName, false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(poolAddr, tfjsonpath.New("listener_id"), knownvalue.Null()),
				},
			},
			{
				Config:             testAccCloudLoadBalancerPoolConfigWithDetachedListener(rName, true),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(poolAddr, plancheck.ResourceActionReplace),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(poolAddr, plancheck.ResourceActionReplace),
					},
				},
			},
		},
	})
}

// TestAccCloudLoadBalancerPool_addTrueLoadBalancerIDIsNoop: a pool created
// through listener_id has load_balancer_id back-filled at Create, so adding
// the true value to the config later plans nothing.
func TestAccCloudLoadBalancerPool_addTrueLoadBalancerIDIsNoop(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerPoolConfigViaListener(rName, false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs(
						"gcore_cloud_load_balancer.test", tfjsonpath.New("id"),
						poolAddr, tfjsonpath.New("load_balancer_id"),
						compare.ValuesSame(),
					),
				},
			},
			{
				Config:   testAccCloudLoadBalancerPoolConfigViaListener(rName, true),
				PlanOnly: true,
			},
		},
	})
}

// TestAccCloudLoadBalancerPool_removeListenerIDIsNoop: dropping listener_id
// from the config of a pool created through it is a no-op - Computed keeps the
// state value.
func TestAccCloudLoadBalancerPool_removeListenerIDIsNoop(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerPoolConfigViaListener(rName, true),
			},
			{
				Config:   testAccCloudLoadBalancerPoolConfigWithDetachedListener(rName, false),
				PlanOnly: true,
			},
		},
	})
}

// TestAccCloudLoadBalancerPool_listenerRecreateReplacesPool verifies that
// replacing the parent listener also replaces the pool with the new listener.
func TestAccCloudLoadBalancerPool_listenerRecreateReplacesPool(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerPoolConfigViaListenerPort(rName, 80),
			},
			{
				Config: testAccCloudLoadBalancerPoolConfigViaListenerPort(rName, 81),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_load_balancer_listener.test", plancheck.ResourceActionReplace),
						plancheck.ExpectResourceAction(poolAddr, plancheck.ResourceActionReplace),
						plancheck.ExpectUnknownValue(poolAddr, tfjsonpath.New("listener_id")),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs(
						"gcore_cloud_load_balancer_listener.test", tfjsonpath.New("id"),
						poolAddr, tfjsonpath.New("listener_id"),
						compare.ValuesSame(),
					),
				},
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

// testAccCloudLoadBalancerPoolConfigDefaultPoolOfListener: the pool is created
// through the load balancer and then attached as the listener's default pool.
func testAccCloudLoadBalancerPoolConfigDefaultPoolOfListener(name string) string {
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
}

resource "gcore_cloud_load_balancer_listener" "test" {
  project_id       = %[1]s
  region_id        = %[2]s
  load_balancer_id = gcore_cloud_load_balancer.test.id
  name             = "%[3]s-listener"
  protocol         = "HTTP"
  protocol_port    = 80
  default_pool_id  = gcore_cloud_load_balancer_pool.test.id
}`, acctest.ProjectID(), acctest.RegionID(), name)
}

// testAccCloudLoadBalancerPoolConfigWithDetachedListener: LB + a listener with
// no default pool + a pool created through the LB; withListenerID adds
// listener_id = listener.id to the pool on day two.
func testAccCloudLoadBalancerPoolConfigWithDetachedListener(name string, withListenerID bool) string {
	// Without listener_id nothing orders the pool after the listener, and the
	// two creates race for the load balancer's task lock (409 Conflict).
	listenerID := "  depends_on       = [gcore_cloud_load_balancer_listener.test]\n"
	if withListenerID {
		listenerID = "  listener_id      = gcore_cloud_load_balancer_listener.test.id\n"
	}
	return fmt.Sprintf(`
resource "gcore_cloud_load_balancer" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-lb"
  flavor     = "lb1-1-2"
}

resource "gcore_cloud_load_balancer_listener" "test" {
  project_id       = %[1]s
  region_id        = %[2]s
  load_balancer_id = gcore_cloud_load_balancer.test.id
  name             = "%[3]s-listener"
  protocol         = "HTTP"
  protocol_port    = 80
}

resource "gcore_cloud_load_balancer_pool" "test" {
  project_id       = %[1]s
  region_id        = %[2]s
  name             = "%[3]s-pool"
  lb_algorithm     = "ROUND_ROBIN"
  protocol         = "HTTP"
  load_balancer_id = gcore_cloud_load_balancer.test.id
%[4]s}`, acctest.ProjectID(), acctest.RegionID(), name, listenerID)
}

// testAccCloudLoadBalancerPoolConfigViaListener creates the pool through
// listener_id only; withLoadBalancerID adds the true load_balancer_id later.
func testAccCloudLoadBalancerPoolConfigViaListener(name string, withLoadBalancerID bool) string {
	lbID := ""
	if withLoadBalancerID {
		lbID = "  load_balancer_id = gcore_cloud_load_balancer.test.id\n"
	}
	return fmt.Sprintf(`
resource "gcore_cloud_load_balancer" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-lb"
  flavor     = "lb1-1-2"
}

resource "gcore_cloud_load_balancer_listener" "test" {
  project_id       = %[1]s
  region_id        = %[2]s
  load_balancer_id = gcore_cloud_load_balancer.test.id
  name             = "%[3]s-listener"
  protocol         = "HTTP"
  protocol_port    = 80
}

resource "gcore_cloud_load_balancer_pool" "test" {
  project_id   = %[1]s
  region_id    = %[2]s
  name         = "%[3]s-pool"
  lb_algorithm = "ROUND_ROBIN"
  protocol     = "HTTP"
  listener_id  = gcore_cloud_load_balancer_listener.test.id
%[4]s}`, acctest.ProjectID(), acctest.RegionID(), name, lbID)
}

func testAccCloudLoadBalancerPoolConfigViaListenerPort(name string, port int) string {
	return fmt.Sprintf(`
resource "gcore_cloud_load_balancer" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-lb"
  flavor     = "lb1-1-2"
}

resource "gcore_cloud_load_balancer_listener" "test" {
  project_id       = %[1]s
  region_id        = %[2]s
  load_balancer_id = gcore_cloud_load_balancer.test.id
  name             = "%[3]s-listener"
  protocol         = "HTTP"
  protocol_port    = %[4]d
}

resource "gcore_cloud_load_balancer_pool" "test" {
  project_id   = %[1]s
  region_id    = %[2]s
  name         = "%[3]s-pool"
  lb_algorithm = "ROUND_ROBIN"
  protocol     = "HTTP"
  listener_id  = gcore_cloud_load_balancer_listener.test.id
}`, acctest.ProjectID(), acctest.RegionID(), name, port)
}
