package cloud_load_balancer_listener_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

const listenerAddr = "gcore_cloud_load_balancer_listener.test"

// TestAccCloudLoadBalancerListener_importNoReplace guards the original bug:
// default_pool_id and insert_x_forwarded are never returned by the API, so
// they are null after an import, and a configured value used to plan a
// destroy+recreate of a live listener. The plan right after import must be
// an in-place update that adopts them.
func TestAccCloudLoadBalancerListener_importNoReplace(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerConfig(rName, listenerWithPoolAndXFF(true)),
			},
			{
				Config:             testAccListenerConfig(rName, listenerWithPoolAndXFF(true)),
				ResourceName:       listenerAddr,
				ImportState:        true,
				ImportStateKind:    resource.ImportBlockWithID,
				ImportStateIdFunc:  acctest.BuildImportID(listenerAddr, "project_id", "region_id", "id"),
				ExpectNonEmptyPlan: true,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(listenerAddr, plancheck.ResourceActionUpdate),
					},
				},
			},
		},
	})
}

// TestAccCloudLoadBalancerListener_importAdoptionApplies covers what a
// plannable import cannot reach: the adopting apply itself. It must strip the
// two attributes from the PATCH body, record them in state, leave a clean
// plan, and a real change afterwards must still replace the listener.
func TestAccCloudLoadBalancerListener_importAdoptionApplies(t *testing.T) {
	rName := acctest.RandomName()
	var listenerID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerConfig(rName, listenerWithPoolAndXFF(true)),
				Check:  captureID(listenerAddr, &listenerID),
			},
			{
				// Forget the listener without destroying it, so the import below
				// targets an address absent from state (Terraform >= 1.7).
				Config: testAccListenerConfig(rName, "") + testAccListenerRemovedBlock,
			},
			{
				Config:             testAccListenerConfig(rName, listenerWithPoolAndXFF(true)),
				ResourceName:       listenerAddr,
				ImportState:        true,
				ImportStateKind:    resource.ImportCommandWithID,
				ImportStatePersist: true,
				ImportStateIdFunc:  importIDFor(&listenerID),
			},
			{
				Config: testAccListenerConfig(rName, listenerWithPoolAndXFF(true)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(listenerAddr, plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(listenerAddr, tfjsonpath.New("default_pool_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(listenerAddr, tfjsonpath.New("insert_x_forwarded"), knownvalue.Bool(true)),
				},
			},
			{
				Config: testAccListenerConfig(rName, listenerWithPoolAndXFF(false)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(listenerAddr, plancheck.ResourceActionReplace),
					},
				},
			},
		},
	})
}

// TestAccCloudLoadBalancerListener_markerRetiresWithoutUpdate is the "do
// nothing, then act later" path. The import adopts nothing, so no Update ever
// runs to close the window; Read must retire the marker on its own. Setting
// insert_x_forwarded on day two was never applied to the listener, so it must
// plan a replacement, not be silently recorded in state.
func TestAccCloudLoadBalancerListener_markerRetiresWithoutUpdate(t *testing.T) {
	rName := acctest.RandomName()
	var listenerID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerConfig(rName, listenerPlain),
				Check:  captureID(listenerAddr, &listenerID),
			},
			{
				Config: testAccListenerConfig(rName, "") + testAccListenerRemovedBlock,
			},
			{
				Config:             testAccListenerConfig(rName, listenerPlain),
				ResourceName:       listenerAddr,
				ImportState:        true,
				ImportStateKind:    resource.ImportCommandWithID,
				ImportStatePersist: true,
				ImportStateIdFunc:  importIDFor(&listenerID),
			},
			{
				Config: testAccListenerConfig(rName, listenerPlain),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				Config: testAccListenerConfig(rName, listenerWithXFFOnly),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(listenerAddr, plancheck.ResourceActionReplace),
					},
				},
			},
		},
	})
}

// Listener variants, as the extra attributes of the listener block.
// admin_state_up is set explicitly because an import that leaves it unset
// hits a separate API bug on apply (a 400 on admin_state_up).
func listenerWithPoolAndXFF(xff bool) string {
	return fmt.Sprintf(`
  default_pool_id    = gcore_cloud_load_balancer_pool.test.id
  insert_x_forwarded = %t`, xff)
}

// user_list is set explicitly: the API returns [] and an import otherwise plans
// a one-off `[] -> null` update, which would run Update and close the adoption
// window - hiding the no-Update path the marker test exists for.
const (
	listenerPlain       = "\n  user_list = []"
	listenerWithXFFOnly = "\n  user_list = []\n  insert_x_forwarded = true"
)

// testAccListenerConfig declares network, subnet, reserved fixed IP, load
// balancer and pool, plus the listener when listenerAttrs is non-empty.
func testAccListenerConfig(name, listenerAttrs string) string {
	config := fmt.Sprintf(`
resource "gcore_cloud_network" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
}

resource "gcore_cloud_network_subnet" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  network_id = gcore_cloud_network.test.id
  cidr       = "10.30.0.0/24"
}

resource "gcore_cloud_reserved_fixed_ip" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  type       = "subnet"
  network_id = gcore_cloud_network.test.id
  subnet_id  = gcore_cloud_network_subnet.test.id
  is_vip     = false
}

resource "gcore_cloud_load_balancer" "test" {
  project_id  = %[1]s
  region_id   = %[2]s
  name        = %[3]q
  flavor      = "lb1-1-2"
  vip_port_id = gcore_cloud_reserved_fixed_ip.test.port_id
}

resource "gcore_cloud_load_balancer_pool" "test" {
  project_id       = %[1]s
  region_id        = %[2]s
  name             = %[3]q
  load_balancer_id = gcore_cloud_load_balancer.test.id
  lb_algorithm     = "ROUND_ROBIN"
  protocol         = "HTTP"
}
`, acctest.ProjectID(), acctest.RegionID(), name)

	if listenerAttrs == "" {
		return config
	}

	return config + fmt.Sprintf(`
resource "gcore_cloud_load_balancer_listener" "test" {
  project_id       = %[1]s
  region_id        = %[2]s
  load_balancer_id = gcore_cloud_load_balancer.test.id
  name             = %[3]q
  protocol         = "HTTP"
  protocol_port    = 80
  admin_state_up   = true%[4]s

  # The API locks the load balancer while a child is created, so the pool and
  # listener must not be created concurrently even when no pool is attached.
  depends_on = [gcore_cloud_load_balancer_pool.test]
}
`, acctest.ProjectID(), acctest.RegionID(), name, listenerAttrs)
}

const testAccListenerRemovedBlock = `
removed {
  from = gcore_cloud_load_balancer_listener.test

  lifecycle {
    destroy = false
  }
}`

func captureID(addr string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[addr]
		if !ok {
			return fmt.Errorf("%s not found in state", addr)
		}
		*id = rs.Primary.ID
		return nil
	}
}

func importIDFor(id *string) resource.ImportStateIdFunc {
	return func(*terraform.State) (string, error) {
		return fmt.Sprintf("%s/%s/%s", acctest.ProjectID(), acctest.RegionID(), *id), nil
	}
}

func testAccCheckCloudLoadBalancerListenerDestroy(s *terraform.State) error {
	regionID, err := strconv.ParseInt(acctest.RegionID(), 10, 64)
	if err != nil {
		return fmt.Errorf("parsing region id: %w", err)
	}

	return acctest.CheckResourceDestroyed(s, "gcore_cloud_load_balancer_listener", func(client *gcore.Client, id string) error {
		_, err := client.Cloud.LoadBalancers.Listeners.Get(context.Background(), id, cloud.LoadBalancerListenerGetParams{
			RegionID: param.NewOpt(regionID),
		})
		return err
	})
}
