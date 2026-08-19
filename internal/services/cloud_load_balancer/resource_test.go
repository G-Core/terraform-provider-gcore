package cloud_load_balancer_test

import (
	"context"
	"fmt"
	"regexp"
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

// TestAccCloudLoadBalancer_importNoReplace guards the import behaviour of the
// write-only attributes the API never returns in GET responses (`no_refresh`):
// vip_network_id and vip_subnet_id. They are null in state after
// an import, and with the built-in RequiresReplace() plan modifier a configured
// value planned a destroy+recreate of a live load balancer.
//
// The expected plan after import is a one-time update-in-place that populates
// both attributes in state from the config — never a replace.
func TestAccCloudLoadBalancer_importNoReplace(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerConfigWithVIP(rName, "test"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("vip_network_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("vip_subnet_id"), knownvalue.NotNull()),
				},
			},
			{
				Config:            testAccCloudLoadBalancerConfigWithVIP(rName, "test"),
				ResourceName:      "gcore_cloud_load_balancer.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_load_balancer.test", "project_id", "region_id", "id"),
				// The API cannot return these two attributes, so the plan right
				// after import is non-empty by design: it populates them in state.
				ExpectNonEmptyPlan: true,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Fails on Replace/DestroyBeforeCreate — the regression this guards.
						plancheck.ExpectResourceAction("gcore_cloud_load_balancer.test",
							plancheck.ResourceActionUpdate),
					},
				},
			},
		},
	})
}

// TestAccCloudLoadBalancer_importAdoptionApplies covers the half of the fix a
// plannable import block cannot reach. ImportBlockWithID steps only plan, so
// they never exercise the adopting Update - the one that must strip the
// create-only fields from the PATCH body instead of sending them to an endpoint
// that rejects them, verify the adopted values against the API, and only then
// write them to state.
//
// Command-mode import with ImportStatePersist puts the imported state into the
// working state, so the following step actually applies. An ImportStatePersist
// step cannot follow a create step for the same address (`terraform import`
// refuses when the address is already in state), so the load balancer is first
// forgotten with a `removed` block (Terraform >= 1.7) - it stays alive, only
// its state entry goes away.
func TestAccCloudLoadBalancer_importAdoptionApplies(t *testing.T) {
	rName := acctest.RandomName()
	var lbID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerConfigWithVIP(rName, "test"),
				Check: func(s *terraform.State) error {
					rs, ok := s.RootModule().Resources["gcore_cloud_load_balancer.test"]
					if !ok {
						return fmt.Errorf("gcore_cloud_load_balancer.test not found in state")
					}
					lbID = rs.Primary.ID
					return nil
				},
			},
			{
				// Forget the load balancer without destroying it, so the
				// import below targets an address absent from state.
				Config: testAccCloudLoadBalancerConfigWithVIPForgotten(rName),
			},
			{
				Config:             testAccCloudLoadBalancerConfigWithVIP(rName, "test"),
				ResourceName:       "gcore_cloud_load_balancer.test",
				ImportState:        true,
				ImportStateKind:    resource.ImportCommandWithID,
				ImportStatePersist: true,
				ImportStateIdFunc: func(*terraform.State) (string, error) {
					return fmt.Sprintf("%s/%s/%s", acctest.ProjectID(), acctest.RegionID(), lbID), nil
				},
			},
			{
				// Adopts vip_network_id/vip_subnet_id into state. Must not
				// replace the load balancer and must leave a clean plan.
				Config: testAccCloudLoadBalancerConfigWithVIP(rName, "test"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_load_balancer.test",
							plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("vip_network_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("vip_subnet_id"), knownvalue.NotNull()),
				},
			},
		},
	})
}

// TestAccCloudLoadBalancer_staleMarkerMismatchErrors is the reviewer's
// stale-marker scenario end to end: import with a config that sets none of the
// create-only attributes (nothing to adopt, no update ever runs, the marker
// stays in private state), then point vip_network_id/vip_subnet_id at a
// different network on day two.
//
// The plan legitimately shows an in-place update - no plan-time hook can
// retire the marker, because Terraform discards private-state writes made
// during no-op plans - but the apply must fail loudly against the API
// evidence instead of silently writing a network the load balancer is not on
// into state. The final plan-only step proves nothing was adopted.
func TestAccCloudLoadBalancer_staleMarkerMismatchErrors(t *testing.T) {
	rName := acctest.RandomName()
	var lbID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				// The load balancer genuinely lives on network "test": created
				// with the VIP attributes, which only exist in config for this
				// step.
				Config: testAccCloudLoadBalancerConfigWithTwoNetworks(rName, "test"),
				Check: func(s *terraform.State) error {
					rs, ok := s.RootModule().Resources["gcore_cloud_load_balancer.test"]
					if !ok {
						return fmt.Errorf("gcore_cloud_load_balancer.test not found in state")
					}
					lbID = rs.Primary.ID
					return nil
				},
			},
			{
				// Forget the load balancer without destroying it (Terraform >= 1.7).
				Config: testAccCloudLoadBalancerConfigTwoNetworksForgotten(rName),
			},
			{
				// Re-import with a config that omits the create-only
				// attributes: the adoption marker is written and nothing ever
				// clears it, because every following plan is a no-op.
				Config:             testAccCloudLoadBalancerConfigTwoNetworksNoVIP(rName),
				ResourceName:       "gcore_cloud_load_balancer.test",
				ImportState:        true,
				ImportStateKind:    resource.ImportCommandWithID,
				ImportStatePersist: true,
				ImportStateIdFunc: func(*terraform.State) (string, error) {
					return fmt.Sprintf("%s/%s/%s", acctest.ProjectID(), acctest.RegionID(), lbID), nil
				},
			},
			{
				// A no-op apply between import and the day-two edit: the
				// marker survives it, which is exactly why apply-time
				// verification has to exist.
				Config: testAccCloudLoadBalancerConfigTwoNetworksNoVIP(rName),
			},
			{
				// Day two: the stale marker suppresses replacement, so the
				// plan is an in-place update - and the apply must then refuse
				// the adoption because the API says the load balancer lives on
				// network "test", not "other".
				Config: testAccCloudLoadBalancerConfigWithTwoNetworks(rName, "other"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_load_balancer.test",
							plancheck.ResourceActionUpdate),
					},
				},
				ExpectError: regexp.MustCompile(`does not match the imported load balancer`),
			},
			{
				// Nothing may have been adopted: with the VIP attributes gone
				// from config again, state must still hold nulls and the plan
				// must be empty.
				Config:   testAccCloudLoadBalancerConfigTwoNetworksNoVIP(rName),
				PlanOnly: true,
			},
		},
	})
}

// TestAccCloudLoadBalancer_lateAdoptionCorrectValue pins the deliberate
// semantics change that comes with apply-time verification: filling in the
// TRUE vip_network_id/vip_subnet_id long after an adopt-nothing import is a
// verified in-place adoption, not a destroy+recreate to the same value. It
// also pins that the successful adoption retires the marker: a subsequent
// change to another network must plan a replacement again.
func TestAccCloudLoadBalancer_lateAdoptionCorrectValue(t *testing.T) {
	rName := acctest.RandomName()
	var lbID string

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerConfigWithTwoNetworks(rName, "test"),
				Check: func(s *terraform.State) error {
					rs, ok := s.RootModule().Resources["gcore_cloud_load_balancer.test"]
					if !ok {
						return fmt.Errorf("gcore_cloud_load_balancer.test not found in state")
					}
					lbID = rs.Primary.ID
					return nil
				},
			},
			{
				Config: testAccCloudLoadBalancerConfigTwoNetworksForgotten(rName),
			},
			{
				Config:             testAccCloudLoadBalancerConfigTwoNetworksNoVIP(rName),
				ResourceName:       "gcore_cloud_load_balancer.test",
				ImportState:        true,
				ImportStateKind:    resource.ImportCommandWithID,
				ImportStatePersist: true,
				ImportStateIdFunc: func(*terraform.State) (string, error) {
					return fmt.Sprintf("%s/%s/%s", acctest.ProjectID(), acctest.RegionID(), lbID), nil
				},
			},
			{
				// Day two with the values the load balancer really has: the
				// apply verifies them against the API and adopts in place.
				Config: testAccCloudLoadBalancerConfigWithTwoNetworks(rName, "test"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_load_balancer.test",
							plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("vip_network_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
						tfjsonpath.New("vip_subnet_id"), knownvalue.NotNull()),
					compareIDSame.AddStateValue("gcore_cloud_load_balancer.test", tfjsonpath.New("id")),
				},
			},
			{
				// The adoption cleared the marker, so a genuine network change
				// is back on the replacement path. The plan check fails the
				// step before the apply if the marker were still suppressing
				// replacement.
				Config: testAccCloudLoadBalancerConfigWithTwoNetworks(rName, "other"),
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

// TestAccCloudLoadBalancer_vipNetworkChangeStillReplaces is the negative case:
// the import-safe modifiers must not weaken replacement for a genuine change.
// Pointing vip_network_id/vip_subnet_id at a different network has to destroy
// and recreate the load balancer, because the update API cannot move it.
func TestAccCloudLoadBalancer_vipNetworkChangeStillReplaces(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudLoadBalancerConfigWithTwoNetworks(rName, "test"),
			},
			{
				Config: testAccCloudLoadBalancerConfigWithTwoNetworks(rName, "other"),
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
// The state checks assert `tags`, not `tags_v2`: tags_v2 also carries read-only
// system tags in an order the API does not guarantee, so there is no stable
// value to assert. It is still what the test exercises — pinning tags_v2 to
// prior state fails step 2 during apply, before any state check runs, with
// "Provider produced inconsistent result after apply".
//
// Each step's post-apply plan must come back empty (the test framework enforces
// this unless ExpectNonEmptyPlan is set), which is what rules out the opposite
// failure of leaving tags_v2 unpinned: a perpetual diff.
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
			// Adding a key grows tags_v2. This is the step that reproduces
			// "new element N has appeared" when tags_v2 is pinned to prior state.
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
// the given resource label. Shared between the with-LB, without-LB and
// forgotten-LB configs so the network resources stay identical across steps
// and are never replaced mid-test.
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

// testAccCloudLoadBalancerRemovedBlock forgets the load balancer without
// destroying it (Terraform >= 1.7): its state entry goes away, the live
// resource stays, and a later step can legally import it again.
const testAccCloudLoadBalancerRemovedBlock = `
removed {
  from = gcore_cloud_load_balancer.test

  lifecycle {
    destroy = false
  }
}`

// testAccCloudLoadBalancerConfigWithVIP points the load balancer at a
// dedicated network+subnet pair, so vip_network_id/vip_subnet_id are set and
// can be adopted after an import.
func testAccCloudLoadBalancerConfigWithVIP(name, network string) string {
	return testAccCloudLoadBalancerNetworkBlocks(name, network, "10.0.0.0/24") + fmt.Sprintf(`

resource "gcore_cloud_load_balancer" "test" {
  project_id     = %[1]s
  region_id      = %[2]s
  name           = %[3]q
  flavor         = "lb1-1-2"
  vip_network_id = gcore_cloud_network.%[4]s.id
  vip_subnet_id  = gcore_cloud_network_subnet.%[4]s.id
}`, acctest.ProjectID(), acctest.RegionID(), name, network)
}

// testAccCloudLoadBalancerConfigWithVIPForgotten keeps the network+subnet of
// testAccCloudLoadBalancerConfigWithVIP(name, "test") and replaces the load
// balancer block with a removed block.
func testAccCloudLoadBalancerConfigWithVIPForgotten(name string) string {
	return testAccCloudLoadBalancerNetworkBlocks(name, "test", "10.0.0.0/24") +
		testAccCloudLoadBalancerRemovedBlock
}

// testAccCloudLoadBalancerTwoNetworksBlocks declares the two network+subnet
// pairs the stale-marker tests move the load balancer between.
func testAccCloudLoadBalancerTwoNetworksBlocks(name string) string {
	return testAccCloudLoadBalancerNetworkBlocks(name, "test", "10.0.0.0/24") +
		testAccCloudLoadBalancerNetworkBlocks(name, "other", "10.1.0.0/24")
}

// testAccCloudLoadBalancerConfigWithTwoNetworks declares two networks and points
// the load balancer at whichever one `network` names, so a step can swap the VIP
// network without the target network itself being replaced.
func testAccCloudLoadBalancerConfigWithTwoNetworks(name, network string) string {
	return testAccCloudLoadBalancerTwoNetworksBlocks(name) + fmt.Sprintf(`

resource "gcore_cloud_load_balancer" "test" {
  project_id     = %[1]s
  region_id      = %[2]s
  name           = %[3]q
  flavor         = "lb1-1-2"
  vip_network_id = gcore_cloud_network.%[4]s.id
  vip_subnet_id  = gcore_cloud_network_subnet.%[4]s.id
}`, acctest.ProjectID(), acctest.RegionID(), name, network)
}

// testAccCloudLoadBalancerConfigTwoNetworksNoVIP keeps both networks but omits
// the create-only VIP attributes from the load balancer: the adopt-nothing
// import config.
func testAccCloudLoadBalancerConfigTwoNetworksNoVIP(name string) string {
	return testAccCloudLoadBalancerTwoNetworksBlocks(name) + fmt.Sprintf(`

resource "gcore_cloud_load_balancer" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  flavor     = "lb1-1-2"
}`, acctest.ProjectID(), acctest.RegionID(), name)
}

// testAccCloudLoadBalancerConfigTwoNetworksForgotten keeps both networks and
// replaces the load balancer block with a removed block.
func testAccCloudLoadBalancerConfigTwoNetworksForgotten(name string) string {
	return testAccCloudLoadBalancerTwoNetworksBlocks(name) +
		testAccCloudLoadBalancerRemovedBlock
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
