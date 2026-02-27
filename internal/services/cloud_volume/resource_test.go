package cloud_volume_test

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

	"github.com/stainless-sdks/gcore-terraform/internal/acctest"
)

func TestAccCloudVolume_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudVolumeConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_volume.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_volume.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_volume.test",
						tfjsonpath.New("status"), knownvalue.StringExact("available")),
					statecheck.ExpectKnownValue("gcore_cloud_volume.test",
						tfjsonpath.New("size"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("gcore_cloud_volume.test",
						tfjsonpath.New("bootable"), knownvalue.Bool(false)),
				},
			},
		},
	})
}

func TestAccCloudVolume_update(t *testing.T) {
	rName := acctest.RandomName()
	newName := acctest.RandomName()

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudVolumeConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_volume.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					compareIDSame.AddStateValue(
						"gcore_cloud_volume.test",
						tfjsonpath.New("id"),
					),
				},
			},
			{
				Config: testAccCloudVolumeConfig(newName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_volume.test",
						tfjsonpath.New("name"), knownvalue.StringExact(newName)),
					compareIDSame.AddStateValue(
						"gcore_cloud_volume.test",
						tfjsonpath.New("id"),
					),
				},
			},
		},
	})
}

func TestAccCloudVolume_import(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudVolumeConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_volume.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
				},
			},
			{
				ResourceName:            "gcore_cloud_volume.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source"},
				ImportStateIdFunc:       acctest.BuildImportID("gcore_cloud_volume.test", "project_id", "region_id", "id"),
			},
		},
	})
}

func testAccCheckCloudVolumeDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cloud_volume" {
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

		_, err = client.Cloud.Volumes.Get(context.Background(), rs.Primary.ID, cloud.VolumeGetParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})

		if err == nil {
			return fmt.Errorf("volume %s still exists", rs.Primary.ID)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking volume deletion: %w", err)
		}
	}
	return nil
}

func testAccCloudVolumeConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_volume" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  source     = "new-volume"
  name       = %[3]q
  size       = 1
}`, acctest.ProjectID(), acctest.RegionID(), name)
}
