package cloud_network_subnet_test

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

func TestAccCloudNetworkSubnet_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudNetworkSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudNetworkSubnetConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_network_subnet.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_network_subnet.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_network_subnet.test",
						tfjsonpath.New("cidr"), knownvalue.StringExact("10.0.0.0/24")),
					statecheck.ExpectKnownValue("gcore_cloud_network_subnet.test",
						tfjsonpath.New("enable_dhcp"), knownvalue.Bool(true)),
				},
			},
		},
	})
}

func TestAccCloudNetworkSubnet_update(t *testing.T) {
	rName := acctest.RandomName()
	rNameUpdated := acctest.RandomName()
	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudNetworkSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudNetworkSubnetConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_network_subnet.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					compareIDSame.AddStateValue("gcore_cloud_network_subnet.test", tfjsonpath.New("id")),
				},
			},
			{
				Config: testAccCloudNetworkSubnetConfigWithDNS(rNameUpdated, `dns_nameservers = ["8.8.8.8"]`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_network_subnet.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rNameUpdated)),
					statecheck.ExpectKnownValue("gcore_cloud_network_subnet.test",
						tfjsonpath.New("dns_nameservers"), knownvalue.ListSizeExact(1)),
					compareIDSame.AddStateValue("gcore_cloud_network_subnet.test", tfjsonpath.New("id")),
				},
			},
		},
	})
}

func TestAccCloudNetworkSubnet_import(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudNetworkSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudNetworkSubnetConfig(rName),
			},
			{
				ResourceName:      "gcore_cloud_network_subnet.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_network_subnet.test", "project_id", "region_id", "id"),
			},
		},
	})
}

func testAccCheckCloudNetworkSubnetDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cloud_network_subnet" {
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

		_, err = client.Cloud.Networks.Subnets.Get(context.Background(), rs.Primary.ID, cloud.NetworkSubnetGetParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})

		if err == nil {
			return fmt.Errorf("subnet %s still exists", rs.Primary.ID)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking subnet deletion: %w", err)
		}
	}
	return nil
}

func testAccCloudNetworkSubnetConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_network" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "tf-test-net-%[3]s"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "test" {
  project_id  = %[1]s
  region_id   = %[2]s
  name        = %[3]q
  network_id  = gcore_cloud_network.test.id
  cidr        = "10.0.0.0/24"
  enable_dhcp = true
}`, acctest.ProjectID(), acctest.RegionID(), name)
}

func testAccCloudNetworkSubnetConfigWithDNS(name string, dnsServers string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_network" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "tf-test-net-%[3]s"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "test" {
  project_id  = %[1]s
  region_id   = %[2]s
  name        = %[3]q
  network_id  = gcore_cloud_network.test.id
  cidr        = "10.0.0.0/24"
  enable_dhcp = true
  %[4]s
}`, acctest.ProjectID(), acctest.RegionID(), name, dnsServers)
}
