package cloud_load_balancer_test

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

func TestAccCloudLoadBalancer_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("vip_address"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("admin_state_up"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("provisioning_status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("operating_status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("region"), knownvalue.NotNull()),
				},
			},
		},
	})
}

// TestAccCloudLoadBalancer_update renames a load balancer in place and changes
// nothing else.
//
// The rename step also carries the regression assertion for the computed-field
// drift: a rename is a genuine in-place update, so the framework re-marks any
// computed attribute that has no plan modifier as "(known after apply)" — the
// drift users still saw after the inline listeners attribute was removed.
//
// provisioning_status is deliberately excluded: it genuinely changes during an
// update, and pinning it previously produced "Provider produced inconsistent
// result after apply".
func TestAccCloudLoadBalancer_update(t *testing.T) {
	rName := acctest.RandomName()
	rNameUpdated := acctest.RandomName()

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	renamePlanChecks := []plancheck.PlanCheck{
		plancheck.ExpectResourceAction("gcore_cloud_load_balancer.test", plancheck.ResourceActionUpdate),
		plancheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
			tfjsonpath.New("name"), knownvalue.StringExact(rNameUpdated)),
	}
	for _, attribute := range []string{"vip_fqdn", "vip_ip_family"} {
		renamePlanChecks = append(renamePlanChecks,
			acctest.ExpectNotUnknownValue("gcore_cloud_load_balancer.test", tfjsonpath.New(attribute)))
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("vip_ip_family"), knownvalue.NotNull()),
					compareIDSame.AddStateValue(
						"gcore_cloud_load_balancer.test",
						tfjsonpath.New("id"),
					),
				},
			},
			{
				Config: testAccCloudLoadBalancerConfig(rNameUpdated),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: renamePlanChecks,
					// The framework already fails the step on a non-empty plan
					// after apply; asserting it here names the drift explicitly.
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rNameUpdated)),
					compareIDSame.AddStateValue(
						"gcore_cloud_load_balancer.test",
						tfjsonpath.New("id"),
					),
				},
			},
		},
	})
}

func TestAccCloudLoadBalancer_import(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
				},
			},
			{
				ResourceName:      "gcore_cloud_load_balancer.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_load_balancer.test", "project_id", "region_id", "id"),
			},
		},
	})
}

// TestAccCloudLoadBalancer_privateVipPort covers the only way to give a load
// balancer a private VIP: a gcore_cloud_reserved_fixed_ip in the private subnet,
// passed as vip_port_id. Unlike the create-only network/subnet inputs it
// replaced, vip_port_id is returned by the API, so the plan after create and
// after import is empty, and the reserved fixed IP outlives the load balancer.
func TestAccCloudLoadBalancer_privateVipPort(t *testing.T) {
	rName := acctest.RandomName()

	comparePortID := statecheck.CompareValuePairs(
		"gcore_cloud_load_balancer.test", tfjsonpath.New("vip_port_id"),
		"gcore_cloud_reserved_fixed_ip.test", tfjsonpath.New("port_id"),
		compare.ValuesSame(),
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerConfigWithVipPort(rName, "test"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("vip_address"), knownvalue.NotNull()),
					comparePortID,
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				// An import block plans against the config, so this fails if
				// importing a private load balancer would replace it.
				Config:            testAccCloudLoadBalancerConfigWithVipPort(rName, "test"),
				ResourceName:      "gcore_cloud_load_balancer.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_load_balancer.test", "project_id", "region_id", "id"),
			},
		},
	})
}

// TestAccCloudLoadBalancer_privateVipPortDualStack places the VIP on a
// dual-stack reserved fixed IP (one IPv4 and one IPv6 address from the
// network's subnets).
func TestAccCloudLoadBalancer_privateVipPortDualStack(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerConfigDualStackVipPort(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("vip_address"), knownvalue.NotNull()),
					statecheck.CompareValuePairs(
						"gcore_cloud_load_balancer.test", tfjsonpath.New("vip_port_id"),
						"gcore_cloud_reserved_fixed_ip.test", tfjsonpath.New("port_id"),
						compare.ValuesSame(),
					),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

// TestAccCloudLoadBalancer_vipPortChangeReplaces: the update API cannot move a
// load balancer to another port, so pointing vip_port_id at a different
// reserved fixed IP has to destroy and recreate it.
func TestAccCloudLoadBalancer_vipPortChangeReplaces(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerConfigWithTwoVipPorts(rName, "test"),
			},
			{
				Config: testAccCloudLoadBalancerConfigWithTwoVipPorts(rName, "other"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_load_balancer.test",
							plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
			},
		},
	})
}

func testAccCheckCloudLoadBalancerDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cloud_load_balancer" {
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

		_, err = client.Cloud.LoadBalancers.Get(context.Background(), rs.Primary.ID, cloud.LoadBalancerGetParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})

		if err == nil {
			return fmt.Errorf("load balancer %s still exists", rs.Primary.ID)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking load balancer deletion: %w", err)
		}
	}
	return nil
}

func TestAccCloudLoadBalancer_resize(t *testing.T) {
	rName := acctest.RandomName()
	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerConfigWithFlavor(rName, "lb1-1-2"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("flavor"), knownvalue.StringExact("lb1-1-2")),
					compareIDSame.AddStateValue("gcore_cloud_load_balancer.test", tfjsonpath.New("id")),
				},
			},
			{
				Config: testAccCloudLoadBalancerConfigWithFlavor(rName, "lb1-2-4"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("flavor"), knownvalue.StringExact("lb1-2-4")),
					compareIDSame.AddStateValue("gcore_cloud_load_balancer.test", tfjsonpath.New("id")),
				},
			},
		},
	})
}

// TestAccCloudLoadBalancer_tags covers in-place tag updates.
//
// The resource no longer exposes the read-only server view of tags at all - it
// is a data source attribute now, named `tags` there - so `tags` here is the
// write-shape map and the only thing to assert.
//
// Each step's post-apply plan must come back empty (the test framework enforces
// this unless ExpectNonEmptyPlan is set), which is what rules out a perpetual
// diff on the read-only tags the API adds on its own.
func TestAccCloudLoadBalancer_tags(t *testing.T) {
	rName := acctest.RandomName()
	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerConfigWithTags(rName, `
    env = "test"
`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("tags"), knownvalue.MapExact(map[string]knownvalue.Check{
							"env": knownvalue.StringExact("test"),
						})),
					compareIDSame.AddStateValue("gcore_cloud_load_balancer.test", tfjsonpath.New("id")),
				},
			},
			// Adding a key grows the server-side tag list. This is the step
			// that used to reproduce "new element N has appeared" back when the
			// resource mirrored that list and pinned it to prior state.
			{
				Config: testAccCloudLoadBalancerConfigWithTags(rName, `
    env   = "test"
    owner = "qa"
`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("tags"), knownvalue.MapExact(map[string]knownvalue.Check{
							"env":   knownvalue.StringExact("test"),
							"owner": knownvalue.StringExact("qa"),
						})),
					compareIDSame.AddStateValue("gcore_cloud_load_balancer.test", tfjsonpath.New("id")),
				},
			},
			// Change a value and drop a key in one step: exercises merge-patch
			// deletion as well as mutation.
			{
				Config: testAccCloudLoadBalancerConfigWithTags(rName, `
    env = "prod"
`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("tags"), knownvalue.MapExact(map[string]knownvalue.Check{
							"env": knownvalue.StringExact("prod"),
						})),
					compareIDSame.AddStateValue("gcore_cloud_load_balancer.test", tfjsonpath.New("id")),
				},
			},
		},
	})
}

func testAccCloudLoadBalancerConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_load_balancer" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  flavor     = "lb1-1-2"
}`, acctest.ProjectID(), acctest.RegionID(), name)
}

// testAccCloudLoadBalancerNetworkBlocks emits one network+subnet pair under
// the given resource label, so the network resources stay identical across
// steps and are never replaced mid-test.
func testAccCloudLoadBalancerNetworkBlocks(name, network, cidr string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_network" %[4]q {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-%[4]s-net"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" %[4]q {
  project_id  = %[1]s
  region_id   = %[2]s
  name        = "%[3]s-%[4]s-subnet"
  network_id  = gcore_cloud_network.%[4]s.id
  cidr        = %[5]q
  enable_dhcp = true
}`, acctest.ProjectID(), acctest.RegionID(), name, network, cidr)
}

// testAccCloudLoadBalancerVipPortBlocks emits a network+subnet pair and a
// reserved fixed IP in that subnet, all under the given resource label.
func testAccCloudLoadBalancerVipPortBlocks(name, label, cidr string) string {
	return testAccCloudLoadBalancerNetworkBlocks(name, label, cidr) + fmt.Sprintf(`

resource "gcore_cloud_reserved_fixed_ip" %[3]q {
  project_id = %[1]s
  region_id  = %[2]s
  type       = "subnet"
  network_id = gcore_cloud_network.%[3]s.id
  subnet_id  = gcore_cloud_network_subnet.%[3]s.id
  is_vip     = false
}`, acctest.ProjectID(), acctest.RegionID(), label)
}

// testAccCloudLoadBalancerConfigWithVipPort places the load balancer VIP on a
// reserved fixed IP in a private subnet.
func testAccCloudLoadBalancerConfigWithVipPort(name, label string) string {
	return testAccCloudLoadBalancerVipPortBlocks(name, label, "10.0.0.0/24") + fmt.Sprintf(`

resource "gcore_cloud_load_balancer" "test" {
  project_id  = %[1]s
  region_id   = %[2]s
  name        = %[3]q
  flavor      = "lb1-1-2"
  vip_port_id = gcore_cloud_reserved_fixed_ip.%[4]s.port_id
}`, acctest.ProjectID(), acctest.RegionID(), name, label)
}

// testAccCloudLoadBalancerConfigWithTwoVipPorts declares two reserved fixed IPs
// and points the load balancer at whichever one `label` names, so a step can
// swap the VIP port without the ports themselves being replaced.
func testAccCloudLoadBalancerConfigWithTwoVipPorts(name, label string) string {
	return testAccCloudLoadBalancerVipPortBlocks(name, "test", "10.0.0.0/24") +
		testAccCloudLoadBalancerVipPortBlocks(name, "other", "10.1.0.0/24") + fmt.Sprintf(`

resource "gcore_cloud_load_balancer" "test" {
  project_id  = %[1]s
  region_id   = %[2]s
  name        = %[3]q
  flavor      = "lb1-1-2"
  vip_port_id = gcore_cloud_reserved_fixed_ip.%[4]s.port_id
}`, acctest.ProjectID(), acctest.RegionID(), name, label)
}

// testAccCloudLoadBalancerConfigDualStackVipPort builds a network with an IPv4
// and an IPv6 subnet and a dual-stack reserved fixed IP as the VIP port.
func testAccCloudLoadBalancerConfigDualStackVipPort(name string) string {
	return testAccCloudLoadBalancerNetworkBlocks(name, "test", "10.0.0.0/24") + fmt.Sprintf(`

resource "gcore_cloud_network_subnet" "test_v6" {
  project_id  = %[1]s
  region_id   = %[2]s
  name        = "%[3]s-test-subnet-v6"
  network_id  = gcore_cloud_network.test.id
  cidr        = "fd00::/120"
  enable_dhcp = true
}

resource "gcore_cloud_reserved_fixed_ip" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  type       = "any_subnet"
  network_id = gcore_cloud_network.test.id
  ip_family  = "dual"
  is_vip     = false

  depends_on = [
    gcore_cloud_network_subnet.test,
    gcore_cloud_network_subnet.test_v6,
  ]
}

resource "gcore_cloud_load_balancer" "test" {
  project_id  = %[1]s
  region_id   = %[2]s
  name        = %[3]q
  flavor      = "lb1-1-2"
  vip_port_id = gcore_cloud_reserved_fixed_ip.test.port_id
}`, acctest.ProjectID(), acctest.RegionID(), name)
}

func testAccCloudLoadBalancerConfigWithFlavor(name, flavor string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_load_balancer" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  flavor     = %[4]q
}`, acctest.ProjectID(), acctest.RegionID(), name, flavor)
}

func testAccCloudLoadBalancerConfigWithTags(name, tags string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_load_balancer" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  flavor     = "lb1-1-2"

  tags = {
%[4]s  }
}`, acctest.ProjectID(), acctest.RegionID(), name, tags)
}

// TestAccCloudLoadBalancer_separateFloatingIP covers the only supported way to
// give a load balancer a public IP: a standalone gcore_cloud_floating_ip
// attached to the load balancer's VIP port.
//
// The inline `floating_ip` block used to do this in one resource, but the API
// has no delete-time cleanup parameter on DELETE /cloud/v1/loadbalancers/...
// (unlike instances and bare metal, which take `delete_floatings` /
// `all_floating_ips`), so every destroy left the inline-created floating IP
// allocated and billable with nothing in the plan to signal it. Modelling the
// floating IP as its own resource puts it under Terraform's lifecycle instead,
// which is what CheckDestroy asserts here.
func TestAccCloudLoadBalancer_separateFloatingIP(t *testing.T) {
	rName := acctest.RandomName()

	comparePortID := statecheck.CompareValuePairs(
		"gcore_cloud_load_balancer.test", tfjsonpath.New("vip_port_id"),
		"gcore_cloud_floating_ip.test", tfjsonpath.New("port_id"),
		compare.ValuesSame(),
	)
	compareFixedIP := statecheck.CompareValuePairs(
		"gcore_cloud_load_balancer.test", tfjsonpath.New("vip_address"),
		"gcore_cloud_floating_ip.test", tfjsonpath.New("fixed_ip_address"),
		compare.ValuesSame(),
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerWithFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerConfigSeparateFloatingIP(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("vip_port_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_floating_ip.test",
						tfjsonpath.New("floating_ip_address"), knownvalue.NotNull()),
					comparePortID,
					compareFixedIP,
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				// Re-reading the load balancer after the floating IP is attached
				// must surface it in the computed floating_ips list, which is how
				// a user discovers the public address from the load balancer side.
				//
				// The floating IP's own status is asserted here rather than in the
				// step above: the create task returns before Neutron flips the
				// association to ACTIVE, so the value only settles on the next read.
				Config: testAccCloudLoadBalancerConfigSeparateFloatingIP(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_floating_ip.test",
						tfjsonpath.New("status"), knownvalue.StringExact("ACTIVE")),
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("floating_ips"), knownvalue.ListSizeExact(1)),
					statecheck.CompareValuePairs(
						"gcore_cloud_load_balancer.test",
						tfjsonpath.New("floating_ips").AtSliceIndex(0).AtMapKey("floating_ip_address"),
						"gcore_cloud_floating_ip.test", tfjsonpath.New("floating_ip_address"),
						compare.ValuesSame(),
					),
				},
			},
			{
				// Importing the attached floating IP must produce a no-op plan —
				// an import that forced replacement would tear the public address
				// away from a live load balancer.
				Config:            testAccCloudLoadBalancerConfigSeparateFloatingIP(rName),
				ResourceName:      "gcore_cloud_floating_ip.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_floating_ip.test", "project_id", "region_id", "id"),
			},
		},
	})
}

// testAccCheckCloudLoadBalancerWithFloatingIPDestroy asserts that neither the
// load balancer nor the floating IP survives `terraform destroy`. The floating
// IP half is the regression guard: an inline-created floating IP passed the
// load-balancer-only check while staying allocated.
func testAccCheckCloudLoadBalancerWithFloatingIPDestroy(s *terraform.State) error {
	if err := testAccCheckCloudLoadBalancerDestroy(s); err != nil {
		return err
	}

	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cloud_floating_ip" {
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

		_, err = client.Cloud.FloatingIPs.Get(context.Background(), rs.Primary.ID, cloud.FloatingIPGetParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})

		if err == nil {
			return fmt.Errorf("floating IP %s still exists", rs.Primary.ID)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking floating IP deletion: %w", err)
		}
	}
	return nil
}

// testAccCloudLoadBalancerConfigSeparateFloatingIP builds a load balancer whose
// VIP is a reserved fixed IP on a private subnet. A load balancer created
// without vip_port_id gets its VIP straight on the external network and the API
// then refuses to associate a floating IP with that port, so the private-VIP
// shape is the one that needs a floating IP at all.
func testAccCloudLoadBalancerConfigSeparateFloatingIP(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_network" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-net"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "test" {
  project_id  = %[1]s
  region_id   = %[2]s
  name        = "%[3]s-subnet"
  network_id  = gcore_cloud_network.test.id
  cidr        = "192.168.40.0/24"
  enable_dhcp = true
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

resource "gcore_cloud_floating_ip" "test" {
  project_id       = %[1]s
  region_id        = %[2]s
  port_id          = gcore_cloud_load_balancer.test.vip_port_id
  fixed_ip_address = gcore_cloud_load_balancer.test.vip_address
  tags = {
    %[4]s
  }
}`, acctest.ProjectID(), acctest.RegionID(), name, acctest.OwnershipTagHCL())
}
