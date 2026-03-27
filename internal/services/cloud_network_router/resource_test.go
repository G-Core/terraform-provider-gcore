package cloud_network_router_test

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
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func TestAccCloudNetworkRouter_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudNetworkRouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudNetworkRouterConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_network_router.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_network_router.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_network_router.test",
						tfjsonpath.New("status"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccCloudNetworkRouter_update(t *testing.T) {
	rName := acctest.RandomName()
	rNameUpdated := acctest.RandomName()

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudNetworkRouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudNetworkRouterConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_network_router.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					compareIDSame.AddStateValue("gcore_cloud_network_router.test", tfjsonpath.New("id")),
				},
			},
			{
				Config: testAccCloudNetworkRouterConfig(rNameUpdated),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_network_router.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rNameUpdated)),
					compareIDSame.AddStateValue("gcore_cloud_network_router.test", tfjsonpath.New("id")),
				},
			},
		},
	})
}

func TestAccCloudNetworkRouter_import(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudNetworkRouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudNetworkRouterConfig(rName),
			},
			{
				ResourceName:      "gcore_cloud_network_router.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_network_router.test", "project_id", "region_id", "id"),
			},
		},
	})
}

func testAccCheckCloudNetworkRouterDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cloud_network_router" {
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

		_, err = client.Cloud.Networks.Routers.Get(context.Background(), rs.Primary.ID, cloud.NetworkRouterGetParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})

		if err == nil {
			return fmt.Errorf("router %s still exists", rs.Primary.ID)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking router deletion: %w", err)
		}
	}
	return nil
}

func TestAccCloudNetworkRouter_updateRoutes(t *testing.T) {
	rName := acctest.RandomName()
	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudNetworkRouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudNetworkRouterWithSubnetConfig(rName, ""),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_network_router.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					compareIDSame.AddStateValue("gcore_cloud_network_router.test", tfjsonpath.New("id")),
				},
			},
			{
				Config: testAccCloudNetworkRouterWithSubnetConfig(rName, `
  routes = [{
    destination = "192.168.100.0/24"
    nexthop     = "10.0.0.1"
  }]`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_network_router.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_network_router.test",
						tfjsonpath.New("routes"), knownvalue.ListSizeExact(1)),
					compareIDSame.AddStateValue("gcore_cloud_network_router.test", tfjsonpath.New("id")),
				},
			},
		},
	})
}

func testAccCloudNetworkRouterConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_network_router" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
}`, acctest.ProjectID(), acctest.RegionID(), name)
}

func testAccCloudNetworkRouterWithSubnetConfig(name string, routes string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_network" "test" {
  project_id    = %[1]s
  region_id     = %[2]s
  name          = "tf-test-net-%[3]s"
  create_router = false
  type          = "vxlan"
}

resource "gcore_cloud_network_subnet" "test" {
  project_id                = %[1]s
  region_id                 = %[2]s
  name                      = "tf-test-subnet-%[3]s"
  network_id                = gcore_cloud_network.test.id
  cidr                      = "10.0.0.0/24"
  enable_dhcp               = true
  connect_to_network_router = false
}

resource "gcore_cloud_network_router" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  interfaces = [{
    type      = "subnet"
    subnet_id = gcore_cloud_network_subnet.test.id
  }]
  %[4]s
}`, acctest.ProjectID(), acctest.RegionID(), name, routes)
}
