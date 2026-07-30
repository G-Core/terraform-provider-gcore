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
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
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

// TestAccCloudVolume_importNoReplacement verifies that importing a volume whose
// configuration declares the write-only create-only fields (source, type_name)
// never proposes a destroy/recreate.
//
// The API does not return source or type_name on GET, so they land in state as
// null after import. Before the fix these were plain RequiresReplace attributes,
// so the first plan after import proposed replacing the volume — silently
// destroying data. They now use the import-safe plan modifier, which adopts the
// config value instead. The plan is a non-destructive update-in-place that
// reconciles the write-only fields into state; a second plan is fully empty.
func TestAccCloudVolume_importNoReplacement(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudVolumeConfigTyped(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_volume.test",
						tfjsonpath.New("source"), knownvalue.StringExact("new-volume")),
					statecheck.ExpectKnownValue("gcore_cloud_volume.test",
						tfjsonpath.New("type_name"), knownvalue.StringExact("standard")),
				},
			},
			// Import. source and type_name come back null from the API, so the
			// plan reconciles them in place. The critical assertion is that the
			// volume is updated, never replaced.
			{
				ResourceName:       "gcore_cloud_volume.test",
				ImportState:        true,
				ImportStateIdFunc:  acctest.BuildImportID("gcore_cloud_volume.test", "project_id", "region_id", "id"),
				ImportStateKind:    resource.ImportBlockWithID,
				ExpectNonEmptyPlan: true,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_volume.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue("gcore_cloud_volume.test",
							tfjsonpath.New("source"), knownvalue.StringExact("new-volume")),
						plancheck.ExpectKnownValue("gcore_cloud_volume.test",
							tfjsonpath.New("type_name"), knownvalue.StringExact("standard")),
						// The volume itself must survive: its id is preserved,
						// not recomputed as "known after apply".
						plancheck.ExpectKnownValue("gcore_cloud_volume.test",
							tfjsonpath.New("id"), knownvalue.NotNull()),
					},
				},
			},
			// After the one reconciliation apply the plan must be fully empty.
			{
				Config: testAccCloudVolumeConfigTyped(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

// TestAccCloudVolume_importFromSnapshotNoReplacement covers the source value the
// API cannot be distinguished from a blank volume: a volume created from a
// snapshot returns volume_image_metadata=null and empty snapshot_ids, exactly
// like source="new-volume". Import must still avoid replacement.
func TestAccCloudVolume_importFromSnapshotNoReplacement(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudVolumeConfigFromSnapshot(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_volume.from_snapshot",
						tfjsonpath.New("source"), knownvalue.StringExact("snapshot")),
				},
			},
			{
				ResourceName:       "gcore_cloud_volume.from_snapshot",
				ImportState:        true,
				ImportStateIdFunc:  acctest.BuildImportID("gcore_cloud_volume.from_snapshot", "project_id", "region_id", "id"),
				ImportStateKind:    resource.ImportBlockWithID,
				ExpectNonEmptyPlan: true,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_volume.from_snapshot", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue("gcore_cloud_volume.from_snapshot",
							tfjsonpath.New("source"), knownvalue.StringExact("snapshot")),
					},
				},
			},
			{
				Config: testAccCloudVolumeConfigFromSnapshot(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

// TestAccCloudVolume_sourceChangeStillReplaces guards the other side of the fix:
// suppressing replacement after import must not suppress it for a genuine
// config change, where the prior state value is known.
func TestAccCloudVolume_sourceChangeStillReplaces(t *testing.T) {
	rName := acctest.RandomName()

	compareIDDifferent := statecheck.CompareValue(compare.ValuesDiffer())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudVolumeConfigTyped(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					compareIDDifferent.AddStateValue(
						"gcore_cloud_volume.test",
						tfjsonpath.New("id"),
					),
				},
			},
			// Changing type_name on a resource whose state holds a known value
			// must still force replacement.
			{
				Config: testAccCloudVolumeConfigTypeName(rName, "ssd_hiiops"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_volume.test", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					compareIDDifferent.AddStateValue(
						"gcore_cloud_volume.test",
						tfjsonpath.New("id"),
					),
				},
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

func testAccCloudVolumeConfigTyped(name string) string {
	return testAccCloudVolumeConfigTypeName(name, "standard")
}

func testAccCloudVolumeConfigTypeName(name, typeName string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_volume" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  source     = "new-volume"
  name       = %[3]q
  size       = 1
  type_name  = %[4]q
}`, acctest.ProjectID(), acctest.RegionID(), name, typeName)
}

func testAccCloudVolumeConfigFromSnapshot(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_volume" "source" {
  project_id = %[1]s
  region_id  = %[2]s
  source     = "new-volume"
  name       = "%[3]s-src"
  size       = 1
}

resource "gcore_cloud_volume_snapshot" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  volume_id  = gcore_cloud_volume.source.id
  name       = "%[3]s-snap"
}

resource "gcore_cloud_volume" "from_snapshot" {
  project_id  = %[1]s
  region_id   = %[2]s
  source      = "snapshot"
  snapshot_id = gcore_cloud_volume_snapshot.test.id
  name        = "%[3]s-fromsnap"
}`, acctest.ProjectID(), acctest.RegionID(), name)
}
