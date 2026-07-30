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
// drift: a rename is a genuine in-place update, so the framework re-marks every
// computed attribute that has no plan modifier as "(known after apply)" — the
// drift users still saw after the inline listeners attribute was removed. Every
// computed attribute except provisioning_status must therefore plan as known.
// The expected-known set is derived from the resource schema itself, see
// stableComputedAttributeNames in schema_guard_test.go.
func TestAccCloudLoadBalancer_update(t *testing.T) {
	rName := acctest.RandomName()
	rNameUpdated := acctest.RandomName()

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	renamePlanChecks := []plancheck.PlanCheck{
		plancheck.ExpectResourceAction("gcore_cloud_load_balancer.test", plancheck.ResourceActionUpdate),
		plancheck.ExpectKnownValue("gcore_cloud_load_balancer.test",
			tfjsonpath.New("name"), knownvalue.StringExact(rNameUpdated)),
	}
	for _, attribute := range stableComputedAttributeNames() {
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

func testAccCloudLoadBalancerConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_load_balancer" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  flavor     = "lb1-1-2"
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
