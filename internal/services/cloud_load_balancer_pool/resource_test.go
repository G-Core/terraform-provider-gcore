package cloud_load_balancer_pool_test

import (
	"context"
	"fmt"
	"strconv"
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

const hmAddr = "gcore_cloud_load_balancer_pool_health_monitor.test"

// testAccCloudLoadBalancerPoolLegacyVersion is the last release in which the
// pool managed its health monitor through a nested healthmonitor attribute.
const testAccCloudLoadBalancerPoolLegacyVersion = "2.0.0-rc.2"

// testAccRequireGcoreProvider pins configs that follow a legacy-release step to
// G-Core/gcore, the provider the legacy state records. Without it the config
// implies hashicorp/gcore, and init still has to resolve G-Core/gcore for the
// state, which the locked pre-release version blocks.
const testAccRequireGcoreProvider = `
terraform {
  required_providers {
    gcore = {
      source = "G-Core/gcore"
    }
  }
}
`

// TestAccCloudLoadBalancerPool_reimportKeepsHealthMonitor: re-importing a pool
// whose monitor belongs to gcore_cloud_load_balancer_pool_health_monitor must
// be a no-op. While the pool had a nested healthmonitor, the import read the
// monitor into it, the configuration had none, and the plan proposed
// healthmonitor = null, which deleted the monitor the other resource owned.
func TestAccCloudLoadBalancerPool_reimportKeepsHealthMonitor(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerPoolConfigWithHealthMonitorResource(rName),
			},
			{
				ResourceName:      poolAddr,
				ImportState:       true,
				ImportStateIdFunc: acctest.BuildImportID(poolAddr, "project_id", "region_id", "id"),
				ImportStateKind:   resource.ImportBlockWithID,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(poolAddr, plancheck.ResourceActionNoop),
						plancheck.ExpectResourceAction(hmAddr, plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

// TestAccCloudLoadBalancerPool_migrateHealthmonitorToResource follows the
// documented migration from the nested healthmonitor attribute: in one apply,
// drop the block from the pool and import the existing monitor into
// gcore_cloud_load_balancer_pool_health_monitor. The monitor must survive with
// the same ID; the pool is never replaced and no longer touches the monitor.
func TestAccCloudLoadBalancerPool_migrateHealthmonitorToResource(t *testing.T) {
	rName := acctest.RandomName()
	var monitorID string

	// Serve the provider under test as G-Core/gcore so it replaces the legacy
	// release in the same state. t.Setenv rules out a parallel test.
	t.Setenv("TF_ACC_PROVIDER_NAMESPACE", "G-Core")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		CheckDestroy: testAccCheckCloudLoadBalancerPoolDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"gcore": {Source: "G-Core/gcore", VersionConstraint: testAccCloudLoadBalancerPoolLegacyVersion},
				},
				Config: testAccCloudLoadBalancerPoolConfigLegacyHealthmonitor(rName),
				Check:  testAccCaptureHealthMonitorID(poolAddr, &monitorID),
			},
			{
				ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
				Config:                   testAccRequireGcoreProvider + testAccCloudLoadBalancerPoolConfigImportHealthMonitor(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acctest.ExpectResourceNotReplaced(poolAddr),
						plancheck.ExpectResourceAction(hmAddr, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHealthMonitorIDUnchanged(poolAddr, &monitorID),
					resource.TestCheckResourceAttrPtr(hmAddr, "id", &monitorID),
				),
			},
			{
				ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
				Config:                   testAccRequireGcoreProvider + testAccCloudLoadBalancerPoolConfigWithHealthMonitorResource(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
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
				Config: testAccCloudLoadBalancerPoolConfigBasic(rName),
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
				Config: testAccCloudLoadBalancerPoolConfigBasic(rName),
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

func testAccCloudLoadBalancerPoolConfigBasic(name string) string {
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

func testAccCloudLoadBalancerPoolConfigWithHealthMonitorResource(name string) string {
	return testAccCloudLoadBalancerPoolConfigBasic(name) + fmt.Sprintf(`

resource "gcore_cloud_load_balancer_pool_health_monitor" "test" {
  project_id  = %[1]s
  region_id   = %[2]s
  pool_id     = gcore_cloud_load_balancer_pool.test.id
  type        = "TCP"
  delay       = 10
  timeout     = 5
  max_retries = 3
}`, acctest.ProjectID(), acctest.RegionID())
}

// testAccCloudLoadBalancerPoolConfigImportHealthMonitor is the migration
// apply: the pool without its healthmonitor block, and an import block that
// adopts the monitor the legacy release created.
func testAccCloudLoadBalancerPoolConfigImportHealthMonitor(name string) string {
	return testAccCloudLoadBalancerPoolConfigWithHealthMonitorResource(name) + fmt.Sprintf(`

import {
  to = gcore_cloud_load_balancer_pool_health_monitor.test
  id = "%[1]s/%[2]s/${gcore_cloud_load_balancer_pool.test.id}"
}`, acctest.ProjectID(), acctest.RegionID())
}

// testAccCloudLoadBalancerPoolConfigLegacyHealthmonitor uses the nested
// healthmonitor attribute, so it only applies with the legacy release. The
// monitor matches testAccCloudLoadBalancerPoolConfigWithHealthMonitorResource
// so the import that adopts it plans no change.
func testAccCloudLoadBalancerPoolConfigLegacyHealthmonitor(name string) string {
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

// getPoolHealthMonitor reads the health monitor of the pool at addr from the
// API.
func getPoolHealthMonitor(s *terraform.State, addr string) (*cloud.HealthMonitor, error) {
	rs, ok := s.RootModule().Resources[addr]
	if !ok {
		return nil, fmt.Errorf("%s not found in state", addr)
	}
	projectID, err := strconv.ParseInt(rs.Primary.Attributes["project_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing project_id: %w", err)
	}
	regionID, err := strconv.ParseInt(rs.Primary.Attributes["region_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing region_id: %w", err)
	}
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return nil, err
	}
	return client.Cloud.LoadBalancers.Pools.HealthMonitors.Get(context.Background(), rs.Primary.ID,
		cloud.LoadBalancerPoolHealthMonitorGetParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})
}

// testAccCaptureHealthMonitorID records the ID of the pool's health monitor.
func testAccCaptureHealthMonitorID(addr string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		hm, err := getPoolHealthMonitor(s, addr)
		if err != nil {
			return fmt.Errorf("error getting health monitor of %s: %w", addr, err)
		}
		*id = hm.ID
		return nil
	}
}

// testAccCheckHealthMonitorIDUnchanged asserts the pool still has the health
// monitor recorded by testAccCaptureHealthMonitorID, not a recreated one.
func testAccCheckHealthMonitorIDUnchanged(addr string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		hm, err := getPoolHealthMonitor(s, addr)
		if err != nil {
			return fmt.Errorf("error getting health monitor of %s: %w", addr, err)
		}
		if hm.ID != *id {
			return fmt.Errorf("expected health monitor %s on %s, found %s", *id, addr, hm.ID)
		}
		return nil
	}
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
