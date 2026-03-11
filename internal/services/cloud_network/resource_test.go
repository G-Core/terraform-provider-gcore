package cloud_network_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func TestAccCloudNetwork_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudNetworkConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_network.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_network.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_network.test",
						tfjsonpath.New("type"), knownvalue.StringExact("vxlan")),
				},
			},
		},
	})
}

func TestAccCloudNetwork_update(t *testing.T) {
	rName := acctest.RandomName()
	rNameUpdated := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudNetworkConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_network.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
				},
			},
			{
				Config: testAccCloudNetworkConfig(rNameUpdated),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_network.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rNameUpdated)),
				},
			},
		},
	})
}

func TestAccCloudNetwork_import(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudNetworkConfig(rName),
			},
			{
				ResourceName:      "gcore_cloud_network.test",
				ImportState:       true,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_network.test", "project_id", "region_id", "id"),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCloudNetworkDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cloud_network" {
			continue
		}

		projectID, _ := strconv.ParseInt(rs.Primary.Attributes["project_id"], 10, 64)
		regionID, _ := strconv.ParseInt(rs.Primary.Attributes["region_id"], 10, 64)

		_, err := client.Cloud.Networks.Get(context.Background(), rs.Primary.ID, cloud.NetworkGetParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})

		if err == nil {
			return fmt.Errorf("network %s still exists", rs.Primary.ID)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking network deletion: %w", err)
		}
	}
	return nil
}

func testAccCloudNetworkConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_network" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  type       = "vxlan"
}`, acctest.ProjectID(), acctest.RegionID(), name)
}
