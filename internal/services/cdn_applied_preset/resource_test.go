package cdn_applied_preset_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cdn"
	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func TestAccCDNAppliedPreset_basic(t *testing.T) {
	presetID := cdnResourcePresetID(t)
	rName := acctest.RandomName()
	cname := rName + ".example.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNAppliedPresetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNAppliedPresetConfig(rName, cname, presetID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_applied_preset.test",
						tfjsonpath.New("preset_id"), knownvalue.Int64Exact(presetID)),
					// object_id is the CDN resource the preset is applied to, and
					// id mirrors it - the apply response carries no id of its own.
					statecheck.CompareValuePairs(
						"gcore_cdn_resource.test", tfjsonpath.New("id"),
						"gcore_cdn_applied_preset.test", tfjsonpath.New("object_id"),
						compare.ValuesSame(),
					),
					statecheck.CompareValuePairs(
						"gcore_cdn_applied_preset.test", tfjsonpath.New("object_id"),
						"gcore_cdn_applied_preset.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
				},
			},
			{
				ResourceName:      "gcore_cdn_applied_preset.test",
				ImportState:       true,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cdn_applied_preset.test", "preset_id", "id"),
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccCDNAppliedPreset_replace covers the resource having no update endpoint:
// every attribute is RequiresReplace, so pointing object_id at another CDN
// resource must unapply the preset and re-apply it to the new object. The last
// step then removes the resource from the config, so both the replacement's
// destroy phase and a plain destroy are checked against a CDN resource that is
// still there to hold a stale association.
func TestAccCDNAppliedPreset_replace(t *testing.T) {
	presetID := cdnResourcePresetID(t)
	rName := acctest.RandomName()
	cname := rName + ".example.com"

	compareIDDiffer := statecheck.CompareValue(compare.ValuesDiffer())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNAppliedPresetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNAppliedPresetConfigTwoResources(rName, cname, presetID, "test"),
				ConfigStateChecks: []statecheck.StateCheck{
					compareIDDiffer.AddStateValue("gcore_cdn_applied_preset.test", tfjsonpath.New("id")),
					statecheck.CompareValuePairs(
						"gcore_cdn_resource.test", tfjsonpath.New("id"),
						"gcore_cdn_applied_preset.test", tfjsonpath.New("object_id"),
						compare.ValuesSame(),
					),
				},
			},
			{
				Config: testAccCDNAppliedPresetConfigTwoResources(rName, cname, presetID, "second"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cdn_applied_preset.test", plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					compareIDDiffer.AddStateValue("gcore_cdn_applied_preset.test", tfjsonpath.New("id")),
					statecheck.CompareValuePairs(
						"gcore_cdn_resource.second", tfjsonpath.New("id"),
						"gcore_cdn_applied_preset.test", tfjsonpath.New("object_id"),
						compare.ValuesSame(),
					),
				},
				// The replacement must also leave the first CDN resource without
				// the preset. Asserting it here, rather than leaving it to
				// CheckDestroy, is what makes the unapply observable at all.
				Check: checkPresetNotAppliedTo("gcore_cdn_resource.test", presetID),
			},
			// Dropping the resource from the config exercises the plain delete
			// path while the CDN resource it was applied to is still alive.
			{
				Config: testAccCDNAppliedPresetConfigCDNResourcesOnly(rName, cname),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cdn_applied_preset.test", plancheck.ResourceActionDestroy),
					},
				},
				Check: checkPresetNotAppliedTo("gcore_cdn_resource.second", presetID),
			},
		},
	})
}

// TestAccCDNAppliedPreset_externalDrift covers Read mapping a 404 onto
// RemoveResource: unapplying the preset outside terraform must leave the
// resource planned for creation rather than erroring the plan.
func TestAccCDNAppliedPreset_externalDrift(t *testing.T) {
	presetID := cdnResourcePresetID(t)
	rName := acctest.RandomName()
	cname := rName + ".example.com"
	config := testAccCDNAppliedPresetConfig(rName, cname, presetID)

	var objectID int64

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNAppliedPresetDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  captureAppliedPresetObjectID("gcore_cdn_applied_preset.test", &objectID),
			},
			{
				PreConfig: func() { unapplyPreset(t, presetID, objectID) },
				Config:    config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cdn_applied_preset.test", plancheck.ResourceActionCreate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_applied_preset.test",
						tfjsonpath.New("preset_id"), knownvalue.Int64Exact(presetID)),
				},
			},
		},
	})
}

func testAccCDNAppliedPresetConfig(name, cname string, presetID int64) string {
	return fmt.Sprintf(`
resource "gcore_cdn_origin_group" "test" {
  name = %[1]q
  sources = [
    {
      source = "example.com"
    }
  ]
}

resource "gcore_cdn_resource" "test" {
  cname        = %[2]q
  origin_group = gcore_cdn_origin_group.test.id
}

resource "gcore_cdn_applied_preset" "test" {
  preset_id = %[3]d
  object_id = gcore_cdn_resource.test.id
}`, name, cname, presetID)
}

// testAccCDNAppliedPresetConfigCDNResourcesOnly declares the two CDN resources
// the preset is moved between, without applying it to either. Kept byte-identical
// to the block testAccCDNAppliedPresetConfigTwoResources builds on so switching
// between them never plans a change to the CDN resources themselves.
func testAccCDNAppliedPresetConfigCDNResourcesOnly(name, cname string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_origin_group" "test" {
  name = %[1]q
  sources = [
    {
      source = "example.com"
    }
  ]
}

resource "gcore_cdn_resource" "test" {
  cname        = %[2]q
  origin_group = gcore_cdn_origin_group.test.id
}

resource "gcore_cdn_resource" "second" {
  cname        = "second-%[2]s"
  origin_group = gcore_cdn_origin_group.test.id
}`, name, cname)
}

// testAccCDNAppliedPresetConfigTwoResources declares two CDN resources and
// applies the preset to the one named by target ("test" or "second").
func testAccCDNAppliedPresetConfigTwoResources(name, cname string, presetID int64, target string) string {
	return testAccCDNAppliedPresetConfigCDNResourcesOnly(name, cname) + fmt.Sprintf(`

resource "gcore_cdn_applied_preset" "test" {
  preset_id = %[1]d
  object_id = gcore_cdn_resource.%[2]s.id
}`, presetID, target)
}

// cdnResourcePresetID returns the ID of a preset applicable to CDN resources,
// skipping the test when the account has none. Presets are read-only and
// account-scoped, so there is nothing to create and nothing to hardcode.
func cdnResourcePresetID(t *testing.T) int64 {
	t.Helper()

	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	client, err := acctest.NewGcoreClient()
	if err != nil {
		t.Fatal(err)
	}

	presets, err := client.CDN.Presets.List(context.Background(), cdn.PresetListParams{})
	if err != nil {
		t.Fatalf("failed to list CDN presets: %s", err)
	}

	for _, preset := range presets.Results {
		if preset.ObjectType == cdn.PresetDetailObjectTypeCDNResource {
			return preset.ID
		}
	}

	t.Skip("no CDN preset applicable to CDN resources available on this account")
	return 0
}

// captureAppliedPresetObjectID records the object the preset was applied to so
// a later step can unapply it out of band.
func captureAppliedPresetObjectID(resourceName string, objectID *int64) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		parsed, err := strconv.ParseInt(rs.Primary.Attributes["object_id"], 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing object_id: %w", err)
		}
		*objectID = parsed
		return nil
	}
}

// checkPresetNotAppliedTo asserts the preset is not applied to the CDN resource
// at cdnResourceAddress, which must still exist. CheckDestroy cannot make this
// assertion: it runs once the CDN resources are themselves destroyed, and the
// lookup then reports not found whether or not the preset was ever unapplied.
//
// A pass does not prove the provider's Delete did the unapplying - the API could
// drop the old association when the preset is applied elsewhere - but it does
// hold the invariant that no stale association is left behind.
func checkPresetNotAppliedTo(cdnResourceAddress string, presetID int64) resource.TestCheckFunc {
	return acctest.CheckResourceExists(cdnResourceAddress, func(client *gcore.Client, id string) error {
		objectID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing CDN resource id: %w", err)
		}

		_, err = client.CDN.Presets.Applied.Get(context.Background(), objectID, cdn.PresetAppliedGetParams{
			PresetID: presetID,
		})
		if err == nil {
			return fmt.Errorf("CDN preset %d is still applied to object %d", presetID, objectID)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("unexpected error checking CDN applied preset %d/%d: %w", presetID, objectID, err)
		}

		return nil
	})
}

func unapplyPreset(t *testing.T, presetID, objectID int64) {
	t.Helper()
	client, err := acctest.NewGcoreClient()
	if err != nil {
		t.Fatal(err)
	}
	err = client.CDN.Presets.Applied.Unapply(context.Background(), objectID, cdn.PresetAppliedUnapplyParams{
		PresetID: presetID,
	})
	if err != nil {
		t.Fatalf("failed to unapply CDN preset %d from object %d: %s", presetID, objectID, err)
	}
}

// testAccCheckCDNAppliedPresetDestroy cannot use acctest.CheckResourceDestroyed:
// reading an applied preset needs preset_id as well as the object ID.
func testAccCheckCDNAppliedPresetDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cdn_applied_preset" {
			continue
		}

		objectID, err := strconv.ParseInt(rs.Primary.Attributes["object_id"], 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing object_id: %w", err)
		}
		presetID, err := strconv.ParseInt(rs.Primary.Attributes["preset_id"], 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing preset_id: %w", err)
		}

		_, err = client.CDN.Presets.Applied.Get(context.Background(), objectID, cdn.PresetAppliedGetParams{
			PresetID: presetID,
		})
		if err == nil {
			return fmt.Errorf("CDN preset %d is still applied to object %d", presetID, objectID)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("unexpected error checking CDN applied preset %d/%d: %w", presetID, objectID, err)
		}
	}

	return nil
}
