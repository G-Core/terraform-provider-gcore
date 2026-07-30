package cloud_volume_test

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

// TestAccCloudVolume_updateCreateOnlyField verifies that reconciling a
// create-only field missing from state completes without returning unknown
// computed values after apply.
func TestAccCloudVolume_updateCreateOnlyField(t *testing.T) {
	rName := acctest.RandomName()

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudVolumeConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					compareIDSame.AddStateValue(
						"gcore_cloud_volume.test",
						tfjsonpath.New("id"),
					),
				},
			},
			{
				Config: testAccCloudVolumeConfigTypedWithTags(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_volume.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					compareIDSame.AddStateValue(
						"gcore_cloud_volume.test",
						tfjsonpath.New("id"),
					),
					statecheck.ExpectKnownValue("gcore_cloud_volume.test",
						tfjsonpath.New("type_name"), knownvalue.StringExact("standard")),
					statecheck.ExpectKnownValue("gcore_cloud_volume.test",
						tfjsonpath.New("status"), knownvalue.StringExact("available")),
					statecheck.ExpectKnownValue("gcore_cloud_volume.test",
						tfjsonpath.New("bootable"), knownvalue.Bool(false)),
				},
			},
		},
	})
}

// TestAccCloudVolume_updateCreateOnlyFieldWithUnknownTags verifies that an
// unconfigured computed-optional tags value is not treated as a tag change
// during a state-only reconciliation update.
func TestAccCloudVolume_updateCreateOnlyFieldWithUnknownTags(t *testing.T) {
	rName := acctest.RandomName()

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudVolumeConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					compareIDSame.AddStateValue(
						"gcore_cloud_volume.test",
						tfjsonpath.New("id"),
					),
				},
			},
			{
				Config: testAccCloudVolumeConfigTyped(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_volume.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectUnknownValue("gcore_cloud_volume.test", tfjsonpath.New("tags")),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					compareIDSame.AddStateValue(
						"gcore_cloud_volume.test",
						tfjsonpath.New("id"),
					),
					statecheck.ExpectKnownValue("gcore_cloud_volume.test",
						tfjsonpath.New("type_name"), knownvalue.StringExact("standard")),
					statecheck.ExpectKnownValue("gcore_cloud_volume.test",
						tfjsonpath.New("tags"), knownvalue.MapExact(map[string]knownvalue.Check{})),
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

// TestAccCloudVolume_sizeDerivedFromSnapshot verifies that a volume restored
// from a snapshot settles after the first apply when the configuration leaves
// size out. The size comes from the snapshot, so there is nothing for the user
// to declare.
//
// size is computed-optional: with no value in the configuration the size the
// platform derived stays in state. Before the fix size was optional-only, so
// every subsequent plan proposed unsetting it (`size = 1 -> null`) and the
// configuration never converged.
func TestAccCloudVolume_sizeDerivedFromSnapshot(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudVolumeConfigSizeFromSnapshot(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					// The restored volume inherits the 2 GiB of the source
					// volume without declaring a size of its own.
					statecheck.ExpectKnownValue("gcore_cloud_volume.from_snapshot",
						tfjsonpath.New("size"), knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue("gcore_cloud_volume.from_snapshot",
						tfjsonpath.New("status"), knownvalue.StringExact("available")),
				},
			},
			// Plan again on the unchanged configuration: no diff.
			{
				Config: testAccCloudVolumeConfigSizeFromSnapshot(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

// TestAccCloudVolume_sizeFromImage pins down the image source, where size is
// computed-optional in the schema but still mandatory at the API: only the
// snapshot source derives a size server-side.
//
// Terraform cannot express per-source requiredness on a single attribute, so the
// RequireSizeUnlessDerived plan modifier reports the missing size at plan time
// instead of letting the apply fail on a 400. The test asserts that, then that
// an image volume with an explicit size applies and converges.
func TestAccCloudVolume_sizeFromImage(t *testing.T) {
	rName := acctest.RandomName()
	image := acctest.LatestUbuntuImage(t)
	size := image.MinDisk
	if size < 1 {
		size = 5
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudVolumeDestroy,
		Steps: []resource.TestStep{
			// An image volume must declare a size, and the plan says so.
			{
				Config:      testAccCloudVolumeConfigImageNoSize(rName, image.ID),
				ExpectError: regexp.MustCompile(`(?s)size is required when source is`),
			},
			{
				Config: testAccCloudVolumeConfigImageWithSize(rName, image.ID, size),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_volume.test",
						tfjsonpath.New("size"), knownvalue.Int64Exact(size)),
					statecheck.ExpectKnownValue("gcore_cloud_volume.test",
						tfjsonpath.New("bootable"), knownvalue.Bool(true)),
				},
			},
			// Declaring size explicitly must stay diff-free too.
			{
				Config: testAccCloudVolumeConfigImageWithSize(rName, image.ID, size),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

// TestAccCloudVolume_sizeRequiredForNewVolume covers the same guard for a blank
// volume, the other source the API cannot derive a size for. Nothing is created:
// the plan fails before any request is sent.
func TestAccCloudVolume_sizeRequiredForNewVolume(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCloudVolumeConfigNoSize(rName),
				ExpectError: regexp.MustCompile(`(?s)size is required when source is`),
			},
		},
	})
}

// TestAccCloudVolume_resize checks the other half of making size
// computed-optional: an explicitly declared size still drives a resize, and
// dropping size from the configuration afterwards must not resize anything
// back.
func TestAccCloudVolume_resize(t *testing.T) {
	rName := acctest.RandomName()

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudVolumeConfigWithSize(rName, 1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_volume.test",
						tfjsonpath.New("size"), knownvalue.Int64Exact(1)),
					compareIDSame.AddStateValue(
						"gcore_cloud_volume.test",
						tfjsonpath.New("id"),
					),
				},
			},
			// Growing the declared size resizes in place, no replacement.
			{
				Config: testAccCloudVolumeConfigWithSize(rName, 2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_volume.test",
							plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_volume.test",
						tfjsonpath.New("size"), knownvalue.Int64Exact(2)),
					compareIDSame.AddStateValue(
						"gcore_cloud_volume.test",
						tfjsonpath.New("id"),
					),
				},
			},
			// Removing size from the configuration is a no-op: the value in
			// state stands, so there is no plan and no spurious resize.
			{
				Config: testAccCloudVolumeConfigNoSize(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_volume.test",
						tfjsonpath.New("size"), knownvalue.Int64Exact(2)),
					compareIDSame.AddStateValue(
						"gcore_cloud_volume.test",
						tfjsonpath.New("id"),
					),
				},
			},
		},
	})
}

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
					statecheck.ExpectKnownValue("gcore_cloud_volume.test", tfjsonpath.New("source"), knownvalue.StringExact("new-volume")),
					statecheck.ExpectKnownValue("gcore_cloud_volume.test", tfjsonpath.New("type_name"), knownvalue.StringExact("standard")),
				},
			},
			{
				ResourceName:       "gcore_cloud_volume.test",
				ImportState:        true,
				ImportStateIdFunc:  acctest.BuildImportID("gcore_cloud_volume.test", "project_id", "region_id", "id"),
				ImportStateKind:    resource.ImportBlockWithID,
				ExpectNonEmptyPlan: true,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_volume.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue("gcore_cloud_volume.test", tfjsonpath.New("source"), knownvalue.StringExact("new-volume")),
						plancheck.ExpectKnownValue("gcore_cloud_volume.test", tfjsonpath.New("type_name"), knownvalue.StringExact("standard")),
						plancheck.ExpectKnownValue("gcore_cloud_volume.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					},
				},
			},
			{
				Config: testAccCloudVolumeConfigTyped(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

func TestAccCloudVolume_importFromSnapshotNoReplacement(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudVolumeConfigSizeFromSnapshot(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_volume.from_snapshot", tfjsonpath.New("source"), knownvalue.StringExact("snapshot")),
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
						plancheck.ExpectKnownValue("gcore_cloud_volume.from_snapshot", tfjsonpath.New("source"), knownvalue.StringExact("snapshot")),
					},
				},
			},
			{
				Config: testAccCloudVolumeConfigSizeFromSnapshot(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

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
					compareIDDifferent.AddStateValue("gcore_cloud_volume.test", tfjsonpath.New("id")),
				},
			},
			{
				Config: testAccCloudVolumeConfigTypeName(rName, "ssd_hiiops"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_volume.test", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					compareIDDifferent.AddStateValue("gcore_cloud_volume.test", tfjsonpath.New("id")),
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
	return testAccCloudVolumeConfigWithSize(name, 1)
}

func testAccCloudVolumeConfigWithSize(name string, size int64) string {
	return fmt.Sprintf(`
resource "gcore_cloud_volume" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  source     = "new-volume"
  name       = %[3]q
  size       = %[4]d
}`, acctest.ProjectID(), acctest.RegionID(), name, size)
}

// testAccCloudVolumeConfigNoSize is the same volume with size left out of the
// configuration entirely.
func testAccCloudVolumeConfigNoSize(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_volume" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  source     = "new-volume"
  name       = %[3]q
}`, acctest.ProjectID(), acctest.RegionID(), name)
}

func testAccCloudVolumeConfigTyped(name string) string {
	return testAccCloudVolumeConfigTypeName(name, "standard")
}

func testAccCloudVolumeConfigTypedWithTags(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_volume" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  source     = "new-volume"
  name       = %[3]q
  size       = 1
  tags       = {}
  type_name  = "standard"
}`, acctest.ProjectID(), acctest.RegionID(), name)
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

// testAccCloudVolumeConfigSizeFromSnapshot restores a volume from a snapshot of
// a 2 GiB volume, without declaring a size on the restored volume.
func testAccCloudVolumeConfigSizeFromSnapshot(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_volume" "origin" {
  project_id = %[1]s
  region_id  = %[2]s
  source     = "new-volume"
  name       = "%[3]s-origin"
  size       = 2
}

resource "gcore_cloud_volume_snapshot" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  volume_id  = gcore_cloud_volume.origin.id
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

// testAccCloudVolumeConfigImageNoSize creates a bootable volume from an image
// without declaring a size.
func testAccCloudVolumeConfigImageNoSize(name, imageID string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_volume" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  source     = "image"
  image_id   = %[4]q
  name       = %[3]q
}`, acctest.ProjectID(), acctest.RegionID(), name, imageID)
}

func testAccCloudVolumeConfigImageWithSize(name, imageID string, size int64) string {
	return fmt.Sprintf(`
resource "gcore_cloud_volume" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  source     = "image"
  image_id   = %[4]q
  name       = %[3]q
  size       = %[5]d
}`, acctest.ProjectID(), acctest.RegionID(), name, imageID, size)
}
