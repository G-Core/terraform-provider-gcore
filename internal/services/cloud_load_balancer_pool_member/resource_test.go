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
