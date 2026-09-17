package cloud_load_balancer_pool_member_test

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

func TestAccCloudLoadBalancerPoolMember_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerPoolMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerPoolMemberConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("address"), knownvalue.StringExact("192.168.1.10")),
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("protocol_port"), knownvalue.Int64Exact(80)),
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("admin_state_up"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("backup"), knownvalue.Bool(false)),
				},
			},
			{
				Config: testAccCloudLoadBalancerPoolMemberConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccCloudLoadBalancerPoolMember_optionalFields(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerPoolMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerPoolMemberConfigWithOptionalFields(rName, 5, true),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("address"), knownvalue.StringExact("192.168.1.20")),
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("protocol_port"), knownvalue.Int64Exact(8080)),
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("weight"), knownvalue.Int64Exact(5)),
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("backup"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("admin_state_up"), knownvalue.Bool(true)),
				},
			},
			{
				Config: testAccCloudLoadBalancerPoolMemberConfigWithOptionalFields(rName, 5, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccCloudLoadBalancerPoolMember_updateWeight(t *testing.T) {
	rName := acctest.RandomName()
	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerPoolMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerPoolMemberConfigWithWeight(rName, 1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("weight"), knownvalue.Int64Exact(1)),
					compareIDSame.AddStateValue(
						"gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("id"),
					),
				},
			},
			{
				Config: testAccCloudLoadBalancerPoolMemberConfigWithWeight(rName, 10),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("weight"), knownvalue.Int64Exact(10)),
					compareIDSame.AddStateValue(
						"gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("id"),
					),
				},
			},
		},
	})
}

func testAccCheckCloudLoadBalancerPoolMemberDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cloud_load_balancer_pool_member" {
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
		poolID := rs.Primary.Attributes["pool_id"]

		pool, err := client.Cloud.LoadBalancers.Pools.Get(context.Background(), poolID, cloud.LoadBalancerPoolGetParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})
		if err != nil {
			if acctest.IsNotFoundError(err) {
				continue
			}
			return fmt.Errorf("error getting pool: %w", err)
		}

		for _, member := range pool.Members {
			if member.ID == rs.Primary.ID {
				return fmt.Errorf("pool member %s still exists in pool %s", rs.Primary.ID, poolID)
			}
		}
	}
	return nil
}

func testAccCloudLoadBalancerPoolMemberConfig(name string) string {
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

resource "gcore_cloud_load_balancer_pool_member" "test" {
  project_id    = %[1]s
  region_id     = %[2]s
  pool_id       = gcore_cloud_load_balancer_pool.test.id
  address       = "192.168.1.10"
  protocol_port = 80
}`, acctest.ProjectID(), acctest.RegionID(), name)
}

func testAccCloudLoadBalancerPoolMemberConfigWithOptionalFields(name string, weight int, backup bool) string {
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

resource "gcore_cloud_load_balancer_pool_member" "test" {
  project_id    = %[1]s
  region_id     = %[2]s
  pool_id       = gcore_cloud_load_balancer_pool.test.id
  address       = "192.168.1.20"
  protocol_port = 8080
  weight        = %[4]d
  backup        = %[5]t
}`, acctest.ProjectID(), acctest.RegionID(), name, weight, backup)
}

func testAccCloudLoadBalancerPoolMemberConfigWithWeight(name string, weight int) string {
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

resource "gcore_cloud_load_balancer_pool_member" "test" {
  project_id    = %[1]s
  region_id     = %[2]s
  pool_id       = gcore_cloud_load_balancer_pool.test.id
  address       = "192.168.1.30"
  protocol_port = 80
  weight        = %[4]d
}`, acctest.ProjectID(), acctest.RegionID(), name, weight)
}

// Import became possible once `get` was mapped. The trailing plan step catches
// attributes Read writes into state that the configuration does not carry.
func TestAccCloudLoadBalancerPoolMember_import(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerPoolMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerPoolMemberConfig(rName),
			},
			{
				ResourceName: "gcore_cloud_load_balancer_pool_member.test",
				ImportState:  true,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_load_balancer_pool_member.test",
					"project_id", "region_id", "pool_id", "id"),
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudLoadBalancerPoolMemberConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

// admin_state_up and backup are in the PATCH body, so they update in place
// instead of forcing replacement.
func TestAccCloudLoadBalancerPoolMember_updateInPlace(t *testing.T) {
	rName := acctest.RandomName()
	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerPoolMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerPoolMemberConfigPatchable(rName, false, true, 1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("backup"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("admin_state_up"), knownvalue.Bool(true)),
					compareIDSame.AddStateValue(
						"gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("id"),
					),
				},
			},
			{
				Config: testAccCloudLoadBalancerPoolMemberConfigPatchable(rName, true, false, 7),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acctest.ExpectResourceNotReplaced("gcore_cloud_load_balancer_pool_member.test"),
						plancheck.ExpectResourceAction("gcore_cloud_load_balancer_pool_member.test",
							plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("backup"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("admin_state_up"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("weight"), knownvalue.Int64Exact(7)),
					compareIDSame.AddStateValue(
						"gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("id"),
					),
				},
			},
		},
	})
}

// address is absent from the PATCH body, so changing it must replace the member
// rather than promise a change the API cannot make.
func TestAccCloudLoadBalancerPoolMember_replaceOnAddress(t *testing.T) {
	rName := acctest.RandomName()
	compareIDDiffer := statecheck.CompareValue(compare.ValuesDiffer())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerPoolMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerPoolMemberConfigAddress(rName, "192.168.1.40"),
				ConfigStateChecks: []statecheck.StateCheck{
					compareIDDiffer.AddStateValue(
						"gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("id"),
					),
				},
			},
			{
				Config: testAccCloudLoadBalancerPoolMemberConfigAddress(rName, "192.168.1.41"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_load_balancer_pool_member.test",
							plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("address"), knownvalue.StringExact("192.168.1.41")),
					compareIDDiffer.AddStateValue(
						"gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("id"),
					),
				},
			},
		},
	})
}

// The two read-only fields the per-member GET returns must not leave the plan
// dirty.
func TestAccCloudLoadBalancerPoolMember_computedStatuses(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerPoolMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerPoolMemberConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("operating_status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("provisioning_status"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccCloudLoadBalancerPoolMemberConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

// A populated subnet_id in the response while the configuration leaves it null
// used to plan a destroy and recreate on every plan: Read writes every field not
// tagged no_refresh, and RequiresReplace turned the resulting null proposal into
// a replacement. Naming the subnet and then dropping it reaches that state, as
// does importing a member created elsewhere from an instance, where the API
// resolved the subnet itself.
func TestAccCloudLoadBalancerPoolMember_subnetIDDroppedFromConfig(t *testing.T) {
	rName := acctest.RandomName()
	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerPoolMemberDestroy,
		Steps: []resource.TestStep{
			{
				// subnet_id named, so the API certainly returns it.
				Config: testAccCloudLoadBalancerPoolMemberConfigWithSubnet(rName, true),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("subnet_id"), knownvalue.NotNull()),
					compareIDSame.AddStateValue(
						"gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("id"),
					),
				},
			},
			{
				// Same member, subnet_id gone from the configuration. Optional+Computed
				// falls back to prior state; a plain optional would plan null here and
				// destroy the member.
				Config: testAccCloudLoadBalancerPoolMemberConfigWithSubnet(rName, false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
						acctest.ExpectResourceNotReplaced("gcore_cloud_load_balancer_pool_member.test"),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					compareIDSame.AddStateValue(
						"gcore_cloud_load_balancer_pool_member.test",
						tfjsonpath.New("id"),
					),
				},
			},
			{
				// And it stays settled across a further refresh.
				Config: testAccCloudLoadBalancerPoolMemberConfigWithSubnet(rName, false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

// A configured subnet must win over the API value and survive a refresh.
func TestAccCloudLoadBalancerPoolMember_subnetIDSet(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerPoolMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerPoolMemberConfigWithSubnet(rName, true),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs(
						"gcore_cloud_network_subnet.test", tfjsonpath.New("id"),
						"gcore_cloud_load_balancer_pool_member.test", tfjsonpath.New("subnet_id"),
						compare.ValuesSame(),
					),
				},
			},
			{
				Config: testAccCloudLoadBalancerPoolMemberConfigWithSubnet(rName, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func testAccCloudLoadBalancerPoolMemberConfigPatchable(name string, backup, adminStateUp bool, weight int) string {
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

resource "gcore_cloud_load_balancer_pool_member" "test" {
  project_id     = %[1]s
  region_id      = %[2]s
  pool_id        = gcore_cloud_load_balancer_pool.test.id
  address        = "192.168.1.50"
  protocol_port  = 80
  backup         = %[4]t
  admin_state_up = %[5]t
  weight         = %[6]d
}`, acctest.ProjectID(), acctest.RegionID(), name, backup, adminStateUp, weight)
}

func testAccCloudLoadBalancerPoolMemberConfigAddress(name, address string) string {
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

resource "gcore_cloud_load_balancer_pool_member" "test" {
  project_id    = %[1]s
  region_id     = %[2]s
  pool_id       = gcore_cloud_load_balancer_pool.test.id
  address       = %[4]q
  protocol_port = 80
}`, acctest.ProjectID(), acctest.RegionID(), name, address)
}

func testAccCloudLoadBalancerPoolMemberConfigWithSubnet(name string, withSubnetID bool) string {
	subnetAttr := ""
	if withSubnetID {
		subnetAttr = "subnet_id     = gcore_cloud_network_subnet.test.id"
	}
	return fmt.Sprintf(`
resource "gcore_cloud_network" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-net"
}

resource "gcore_cloud_network_subnet" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-subnet"
  network_id = gcore_cloud_network.test.id
  cidr       = "192.168.77.0/24"
}

resource "gcore_cloud_load_balancer" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-lb"
  flavor     = "lb1-1-2"

  # Octavia allocates member ports in the subnet without the balancer's
  # configuration referencing it, so the dependency has to be explicit or the
  # destroy removes the subnet while those ports still hold it.
  depends_on = [gcore_cloud_network_subnet.test]
}

resource "gcore_cloud_load_balancer_pool" "test" {
  project_id       = %[1]s
  region_id        = %[2]s
  name             = "%[3]s-pool"
  lb_algorithm     = "ROUND_ROBIN"
  protocol         = "HTTP"
  load_balancer_id = gcore_cloud_load_balancer.test.id
}

resource "gcore_cloud_load_balancer_pool_member" "test" {
  project_id    = %[1]s
  region_id     = %[2]s
  pool_id       = gcore_cloud_load_balancer_pool.test.id
  address       = "192.168.77.10"
  protocol_port = 80
  %[4]s
}`, acctest.ProjectID(), acctest.RegionID(), name, subnetAttr)
}
