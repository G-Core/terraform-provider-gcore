package cloud_file_share_test

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

func TestAccCloudFileShare_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudFileShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFileShareConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_file_share.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_file_share.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_file_share.test",
						tfjsonpath.New("protocol"), knownvalue.StringExact("NFS")),
					statecheck.ExpectKnownValue("gcore_cloud_file_share.test",
						tfjsonpath.New("size"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("gcore_cloud_file_share.test",
						tfjsonpath.New("type_name"), knownvalue.StringExact("standard")),
				},
			},
		},
	})
}

func TestAccCloudFileShare_update(t *testing.T) {
	rName := acctest.RandomName()
	newName := acctest.RandomName()

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudFileShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFileShareConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_file_share.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					compareIDSame.AddStateValue(
						"gcore_cloud_file_share.test",
						tfjsonpath.New("id"),
					),
				},
			},
			{
				Config: testAccCloudFileShareConfigUpdated(newName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_file_share.test",
						tfjsonpath.New("name"), knownvalue.StringExact(newName)),
					statecheck.ExpectKnownValue("gcore_cloud_file_share.test",
						tfjsonpath.New("size"), knownvalue.Int64Exact(2)),
					compareIDSame.AddStateValue(
						"gcore_cloud_file_share.test",
						tfjsonpath.New("id"),
					),
				},
			},
		},
	})
}

func TestAccCloudFileShare_import(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudFileShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFileShareConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_file_share.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
				},
			},
			{
				ResourceName:      "gcore_cloud_file_share.test",
				ImportState:       true,
				ImportStateVerify: true,
				// network is a write-only input (no_refresh) not returned by the API on read,
				// so it cannot round-trip through import.
				ImportStateVerifyIgnore: []string{"network"},
				ImportStateIdFunc:       acctest.BuildImportID("gcore_cloud_file_share.test", "project_id", "region_id", "id"),
			},
		},
	})
}

func testAccCheckCloudFileShareDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cloud_file_share" {
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

		_, err = client.Cloud.FileShares.Get(context.Background(), rs.Primary.ID, cloud.FileShareGetParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})

		if err == nil {
			return fmt.Errorf("file share %s still exists", rs.Primary.ID)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking file share deletion: %w", err)
		}
	}
	return nil
}

func testAccCloudFileShareConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_network" "test" {
  project_id    = %[1]s
  region_id     = %[2]s
  name          = "%[3]s-network"
  create_router = true
  type          = "vxlan"
}

resource "gcore_cloud_network_subnet" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  network_id = gcore_cloud_network.test.id
  name       = "%[3]s-subnet"
  cidr       = "10.0.0.0/24"
}

resource "gcore_cloud_file_share" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  protocol   = "NFS"
  size       = 1
  type_name  = "standard"
  network = {
    network_id = gcore_cloud_network.test.id
    subnet_id  = gcore_cloud_network_subnet.test.id
  }
}`, acctest.ProjectID(), acctest.RegionID(), name)
}

func testAccCloudFileShareConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_network" "test" {
  project_id    = %[1]s
  region_id     = %[2]s
  name          = "%[3]s-network"
  create_router = true
  type          = "vxlan"
}

resource "gcore_cloud_network_subnet" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  network_id = gcore_cloud_network.test.id
  name       = "%[3]s-subnet"
  cidr       = "10.0.0.0/24"
}

resource "gcore_cloud_file_share" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  protocol   = "NFS"
  size       = 2
  type_name  = "standard"
  network = {
    network_id = gcore_cloud_network.test.id
    subnet_id  = gcore_cloud_network_subnet.test.id
  }
}`, acctest.ProjectID(), acctest.RegionID(), name)
}
