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

// TestAccCloudVolume_updateCreateOnlyField: adding type_name = "standard" to
// a volume created without it is a no-op for that attribute. type_name is read
// back from the API's volume_type (Computed+Optional), so the configured value
// equals state and nothing is adopted; the volume is never replaced and state
// keeps type_name == volume_type. Before the fix this planned an in-place
// update that only rewrote state - the shape that recorded a type the volume
// did not have.
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
						acctest.ExpectResourceNotReplaced("gcore_cloud_volume.test"),
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
						tfjsonpath.New("volume_type"), knownvalue.StringExact("standard")),
					statecheck.ExpectKnownValue("gcore_cloud_volume.test",
						tfjsonpath.New("status"), knownvalue.StringExact("available")),
					statecheck.ExpectKnownValue("gcore_cloud_volume.test",
						tfjsonpath.New("bootable"), knownvalue.Bool(false)),
				},
			},
		},
	})
}

// TestAccCloudVolume_updateCreateOnlyFieldWithUnknownTags is the same no-op
// with tags left unconfigured: the computed-optional tags value must not turn
// the step into a tag PATCH, and type_name still equals volume_type.
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
						acctest.ExpectResourceNotReplaced("gcore_cloud_volume.test"),
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
						tfjsonpath.New("volume_type"), knownvalue.StringExact("standard")),
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

// TestAccCloudVolume_typeNameChangeReplaces: a real type change on a known
// volume still plans a replacement. The API has an in-place retype for
// standard<->ssd_hiiops (the v0.x provider used it); honouring it in place is a
// follow-up, so today the honest action is a loud replacement, never a silent
// state-only update.
func TestAccCloudVolume_typeNameChangeReplaces(t *testing.T) {
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

// TestAccCloudVolume_importAdoptionApplies covers the apply half of the
// adoption window. Command-mode import with ImportStatePersist puts the
// imported state into the working state, so later steps really plan and apply
// against it. The standalone plan-only step in the middle is deliberate:
// standalone plans never persist private state, so any number of them must
// leave the window open for the apply that follows.
//
// source is adopted (no_refresh, never returned); type_name is NOT adopted -
// ImportState back-filled it from volume_type, so the configured "standard"
// equals state.
func TestAccCloudVolume_importAdoptionApplies(t *testing.T) {
	rName := acctest.RandomName()
	var volumeID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudVolumeConfig(rName),
				Check:  testAccCaptureVolumeID("gcore_cloud_volume.test", &volumeID),
			},
			{
				// Forget the volume without destroying it, so the import below
				// targets an address absent from state.
				Config: testAccCloudVolumeConfigForgotten(rName),
			},
			{
				Config:             testAccCloudVolumeConfig(rName),
				ResourceName:       "gcore_cloud_volume.test",
				ImportState:        true,
				ImportStateKind:    resource.ImportCommandWithID,
				ImportStatePersist: true,
				ImportStateIdFunc: func(*terraform.State) (string, error) {
					return fmt.Sprintf("%s/%s/%s", acctest.ProjectID(), acctest.RegionID(), volumeID), nil
				},
			},
			{
				// A standalone plan of the adopting config: in-place update
				// planned, window not consumed - plans persist nothing.
				Config:             testAccCloudVolumeConfigTyped(rName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acctest.ExpectResourceNotReplaced("gcore_cloud_volume.test"),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						acctest.ExpectResourceNotReplaced("gcore_cloud_volume.test"),
					},
				},
			},
			{
				// The adopting apply: source lands in state without the volume
				// being replaced, type_name comes from the API, and the marker is
				// retired.
				Config: testAccCloudVolumeConfigTyped(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_volume.test", plancheck.ResourceActionUpdate),
						acctest.ExpectResourceNotReplaced("gcore_cloud_volume.test"),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_volume.test",
						tfjsonpath.New("source"), knownvalue.StringExact("new-volume")),
					statecheck.ExpectKnownValue("gcore_cloud_volume.test",
						tfjsonpath.New("type_name"), knownvalue.StringExact("standard")),
					statecheck.ExpectKnownValue("gcore_cloud_volume.test",
						tfjsonpath.New("volume_type"), knownvalue.StringExact("standard")),
				},
			},
			{
				// Adoption is stable: the same config plans nothing at all.
				Config:   testAccCloudVolumeConfigTyped(rName),
				PlanOnly: true,
			},
		},
	})
}

// TestAccCloudVolume_importThenSetForcesReplace is the ticket's headline
// assertion on the reproduced resource: after the adopting apply consumed the
// window, setting image_id on day two must plan the replacement that actually
// delivers the value - both with and without a refresh, because Update cleared
// the marker. Plan-only: actually replacing proves nothing more.
func TestAccCloudVolume_importThenSetForcesReplace(t *testing.T) {
	rName := acctest.RandomName()
	imageID := acctest.LatestUbuntuImageID(t)
	var volumeID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudVolumeConfig(rName),
				Check:  testAccCaptureVolumeID("gcore_cloud_volume.test", &volumeID),
			},
			{
				Config: testAccCloudVolumeConfigForgotten(rName),
			},
			{
				Config:             testAccCloudVolumeConfig(rName),
				ResourceName:       "gcore_cloud_volume.test",
				ImportState:        true,
				ImportStateKind:    resource.ImportCommandWithID,
				ImportStatePersist: true,
				ImportStateIdFunc: func(*terraform.State) (string, error) {
					return fmt.Sprintf("%s/%s/%s", acctest.ProjectID(), acctest.RegionID(), volumeID), nil
				},
			},
			{
				// The adopting apply (source only).
				Config: testAccCloudVolumeConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_volume.test", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				// Day two: the window is closed, so image_id must plan a replacement.
				Config:             testAccCloudVolumeConfigNewVolumeWithImageID(rName, imageID),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_volume.test", plancheck.ResourceActionReplace),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_volume.test", plancheck.ResourceActionReplace),
					},
				},
			},
		},
	})
}

// TestAccCloudVolume_importMismatchFails: an adopted instance_id_to_attach_to
// is checked against the volume's attachments at the top of the adopting
// Update; a volume that is not attached to that instance fails the apply
// loudly. The follow-up apply with the true config succeeds, which proves the
// failed apply re-armed the window instead of consuming it.
func TestAccCloudVolume_importMismatchFails(t *testing.T) {
	rName := acctest.RandomName()
	var volumeID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudVolumeConfig(rName),
				Check:  testAccCaptureVolumeID("gcore_cloud_volume.test", &volumeID),
			},
			{
				Config: testAccCloudVolumeConfigForgotten(rName),
			},
			{
				Config:             testAccCloudVolumeConfig(rName),
				ResourceName:       "gcore_cloud_volume.test",
				ImportState:        true,
				ImportStateKind:    resource.ImportCommandWithID,
				ImportStatePersist: true,
				ImportStateIdFunc: func(*terraform.State) (string, error) {
					return fmt.Sprintf("%s/%s/%s", acctest.ProjectID(), acctest.RegionID(), volumeID), nil
				},
			},
			{
				Config:      testAccCloudVolumeConfigAttachedTo(rName, "00000000-0000-4000-8000-000000000000"),
				ExpectError: regexp.MustCompile(`contradicts the imported volume`),
			},
			{
				Config: testAccCloudVolumeConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_volume.test", plancheck.ResourceActionUpdate),
						acctest.ExpectResourceNotReplaced("gcore_cloud_volume.test"),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_volume.test",
						tfjsonpath.New("source"), knownvalue.StringExact("new-volume")),
				},
			},
		},
	})
}

// TestAccCloudVolume_importAdoptsWithoutRefresh_documentsRefreshFalseResidual
// pins a documented residual, not a requirement: with -refresh=false Read
// never runs, so the marker written at import is never advanced and the
// adopting plan stays an in-place update for as long as the user keeps
// skipping refresh. Every plan in this test case runs with -refresh=false.
func TestAccCloudVolume_importAdoptsWithoutRefresh_documentsRefreshFalseResidual(t *testing.T) {
	rName := acctest.RandomName()
	var volumeID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudVolumeDestroy,
		AdditionalCLIOptions: &resource.AdditionalCLIOptions{
			Plan: resource.PlanOptions{NoRefresh: true},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCloudVolumeConfig(rName),
				Check:  testAccCaptureVolumeID("gcore_cloud_volume.test", &volumeID),
			},
			{
				Config: testAccCloudVolumeConfigForgotten(rName),
			},
			{
				Config:             testAccCloudVolumeConfig(rName),
				ResourceName:       "gcore_cloud_volume.test",
				ImportState:        true,
				ImportStateKind:    resource.ImportCommandWithID,
				ImportStatePersist: true,
				ImportStateIdFunc: func(*terraform.State) (string, error) {
					return fmt.Sprintf("%s/%s/%s", acctest.ProjectID(), acctest.RegionID(), volumeID), nil
				},
			},
			{
				Config:             testAccCloudVolumeConfigTyped(rName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_volume.test", plancheck.ResourceActionUpdate),
						acctest.ExpectResourceNotReplaced("gcore_cloud_volume.test"),
					},
				},
			},
		},
	})
}

// TestAccCloudVolume_removeImageIDRejectedAtPlan pins the guard in front of
// the removal semantics: dropping image_id from a source="image" volume would
// plan a replacement, and Terraform destroys before it creates, so the
// prerequisite check in ModifyPlan rejects the plan instead (the API would
// refuse the create). This test exercises that check, not the replace
// modifier - TestAccCloudVolume_removeImageIDFromNewVolumeReplaces does.
func TestAccCloudVolume_removeImageIDRejectedAtPlan(t *testing.T) {
	rName := acctest.RandomName()
	imageID := acctest.LatestUbuntuImageID(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudVolumeConfigImageWithSize(rName, imageID, 5),
			},
			{
				Config:      testAccCloudVolumeConfigImageNoImageID(rName, 5),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}

// TestAccCloudVolume_removeImageIDFromNewVolumeReplaces isolates the removal
// semantics of an Optional-only create-only attribute at the replace modifier:
// a new-volume created with a stray image_id (the API ignores it) has the
// attribute dropped with `source` unchanged, so nothing but the modifier's
// "removal replaces" row can produce the replacement - loud, the same as the
// generated base and the v0.x provider, never the silent `X -> null` state
// rewrite of 2.0.0-rc.1. This is the only removal that reaches the modifier
// without a live instance: dropping image_id under source="image" or
// snapshot_id under source="snapshot" is stopped by the prerequisite check
// first, by design.
func TestAccCloudVolume_removeImageIDFromNewVolumeReplaces(t *testing.T) {
	rName := acctest.RandomName()
	imageID := acctest.LatestUbuntuImageID(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudVolumeConfigNewVolumeWithImageID(rName, imageID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_volume.test",
						tfjsonpath.New("image_id"), knownvalue.StringExact(imageID)),
				},
			},
			{
				Config:             testAccCloudVolumeConfig(rName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_volume.test", plancheck.ResourceActionReplace),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_volume.test", plancheck.ResourceActionReplace),
					},
				},
			},
		},
	})
}

// TestAccCloudVolume_removeSnapshotIDReplaces: the snapshot-sourced volume is
// switched to a plain new-volume config (source changes, snapshot_id is
// dropped). Both transitions plan a replacement, with and without refresh.
// `source` alone drives the Replace here; snapshot_id removal under
// source="snapshot" is guarded by the prerequisite check and cannot reach the
// modifier on its own, so this test does not isolate it.
func TestAccCloudVolume_removeSnapshotIDReplaces(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudVolumeConfigSizeFromSnapshot(rName),
			},
			{
				Config:             testAccCloudVolumeConfigFromSnapshotDetached(rName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_volume.from_snapshot", plancheck.ResourceActionReplace),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_volume.from_snapshot", plancheck.ResourceActionReplace),
					},
				},
			},
		},
	})
}

func testAccCaptureVolumeID(address string, into *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[address]
		if !ok {
			return fmt.Errorf("%s not found in state", address)
		}
		*into = rs.Primary.ID
		return nil
	}
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

// testAccCloudVolumeConfigForgotten removes the volume from state without
// destroying it (Terraform >= 1.7), so a following import targets an address
// absent from state.
func testAccCloudVolumeConfigForgotten(string) string {
	return `
removed {
  from = gcore_cloud_volume.test

  lifecycle {
    destroy = false
  }
}`
}

func testAccCloudVolumeConfigNewVolumeWithImageID(name, imageID string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_volume" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  source     = "new-volume"
  image_id   = %[4]q
  name       = %[3]q
  size       = 1
}`, acctest.ProjectID(), acctest.RegionID(), name, imageID)
}

func testAccCloudVolumeConfigAttachedTo(name, instanceID string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_volume" "test" {
  project_id               = %[1]s
  region_id                = %[2]s
  source                   = "new-volume"
  name                     = %[3]q
  size                     = 1
  instance_id_to_attach_to = %[4]q
}`, acctest.ProjectID(), acctest.RegionID(), name, instanceID)
}

// testAccCloudVolumeConfigImageNoImageID is an image-sourced volume with the
// image_id dropped - rejected at plan time by the config validator.
func testAccCloudVolumeConfigImageNoImageID(name string, size int64) string {
	return fmt.Sprintf(`
resource "gcore_cloud_volume" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  source     = "image"
  name       = %[3]q
  size       = %[4]d
}`, acctest.ProjectID(), acctest.RegionID(), name, size)
}

// testAccCloudVolumeConfigFromSnapshotDetached keeps the origin volume and its
// snapshot but turns from_snapshot into a plain new volume.
func testAccCloudVolumeConfigFromSnapshotDetached(name string) string {
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
  project_id = %[1]s
  region_id  = %[2]s
  source     = "new-volume"
  name       = "%[3]s-fromsnap"
  size       = 2
}`, acctest.ProjectID(), acctest.RegionID(), name)
}
