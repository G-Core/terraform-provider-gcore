package cloud_volume_snapshot_test

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

func TestAccCloudVolumeSnapshot_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudVolumeSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudVolumeSnapshotConfig(rName, "test"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_volume_snapshot.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_volume_snapshot.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_volume_snapshot.test",
						tfjsonpath.New("status"), knownvalue.StringExact("available")),
					statecheck.ExpectKnownValue("gcore_cloud_volume_snapshot.test",
						tfjsonpath.New("size"), knownvalue.Int64Exact(1)),
					statecheck.CompareValuePairs(
						"gcore_cloud_volume.test", tfjsonpath.New("id"),
						"gcore_cloud_volume_snapshot.test", tfjsonpath.New("volume_id"),
						compare.ValuesSame(),
					),
					// The API adds bootable / volume_name system tags to every
					// snapshot; the resource filters them out, so the map holds
					// exactly the user-supplied tags.
					statecheck.ExpectKnownValue("gcore_cloud_volume_snapshot.test",
						tfjsonpath.New("tags"), knownvalue.MapExact(map[string]knownvalue.Check{
							"environment": knownvalue.StringExact("test"),
						})),
				},
			},
		},
	})
}

func TestAccCloudVolumeSnapshot_update(t *testing.T) {
	rName := acctest.RandomName()
	newName := acctest.RandomName()

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudVolumeSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudVolumeSnapshotConfig(rName, "test"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_volume_snapshot.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					compareIDSame.AddStateValue(
						"gcore_cloud_volume_snapshot.test",
						tfjsonpath.New("id"),
					),
				},
			},
			{
				// Name and tags update in place: same snapshot ID, no replacement.
				Config: testAccCloudVolumeSnapshotConfig(newName, "staging"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_volume_snapshot.test",
						tfjsonpath.New("name"), knownvalue.StringExact(newName)),
					statecheck.ExpectKnownValue("gcore_cloud_volume_snapshot.test",
						tfjsonpath.New("tags"), knownvalue.MapExact(map[string]knownvalue.Check{
							"environment": knownvalue.StringExact("staging"),
						})),
					compareIDSame.AddStateValue(
						"gcore_cloud_volume_snapshot.test",
						tfjsonpath.New("id"),
					),
				},
			},
		},
	})
}

func TestAccCloudVolumeSnapshot_import(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudVolumeSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudVolumeSnapshotConfig(rName, "test"),
			},
			{
				ResourceName:      "gcore_cloud_volume_snapshot.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_volume_snapshot.test", "project_id", "region_id", "id"),
			},
		},
	})
}

func testAccCheckCloudVolumeSnapshotDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cloud_volume_snapshot" {
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

		_, err = client.Cloud.VolumeSnapshots.Get(context.Background(), rs.Primary.ID, cloud.VolumeSnapshotGetParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})

		if err == nil {
			return fmt.Errorf("volume snapshot %s still exists", rs.Primary.ID)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking volume snapshot deletion: %w", err)
		}
	}
	return nil
}

func testAccCloudVolumeSnapshotConfig(name, environmentTag string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_volume" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  source     = "new-volume"
  name       = %[3]q
  size       = 1
}

resource "gcore_cloud_volume_snapshot" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  volume_id  = gcore_cloud_volume.test.id
  name       = %[3]q

  tags = {
    environment = %[4]q
  }
}`, acctest.ProjectID(), acctest.RegionID(), name, environmentTag)
}
