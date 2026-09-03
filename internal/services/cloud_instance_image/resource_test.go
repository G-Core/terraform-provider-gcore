package cloud_instance_image_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

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

// testImageURL is a small, publicly-accessible Linux cloud image used for acceptance tests.
// Using a known-good minimal image to keep upload times short.
const testImageURL = "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud-2003.qcow2.xz"

// changedTestImageURL is a second, equally small image used to prove that a
// genuine url change still forces a replacement.
const changedTestImageURL = "http://download.cirros-cloud.net/0.6.2/cirros-0.6.2-x86_64-disk.img"

func TestAccCloudInstanceImage_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudInstanceImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstanceImageConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_instance_image.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_instance_image.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_instance_image.test",
						tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_instance_image.test",
						tfjsonpath.New("visibility"), knownvalue.NotNull()),
				},
			},
		},
	})
}

// TestAccCloudInstanceImage_import covers importing the resource and, crucially,
// what the *next* plan after the import looks like.
//
// url is Required, create-only and never returned by the GET endpoint
// (no_refresh), so it lands in state as null after an import. While url carried
// stringplanmodifier.RequiresReplace(), that null-in-state vs set-in-config
// difference made the follow-up plan destroy and recreate live infrastructure.
// A plain ImportState + ImportStateVerify step cannot see that: it only compares
// the imported state against the pre-import state and never plans afterwards.
//
// The third step closes that gap with ImportStateKind: resource.ImportBlockWithID
// — config-driven import, which drives the import from an `import` block against
// the configuration and plans it. url (plus the other create-only attributes) is
// null in the freshly imported state, so the plan is legitimately non-empty; the
// point of the test is that it must be an in-place UPDATE and never a
// replace/destroy. When a user applies that update, the adoption is served by
// this resource's real Update (PATCH) implementation, which
// TestAccCloudInstanceImage_update exercises end to end.
func TestAccCloudInstanceImage_import(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudInstanceImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstanceImageConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_instance_image.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
				},
			},
			{
				ResourceName:      "gcore_cloud_instance_image.test",
				ImportState:       true,
				ImportStateVerify: true,
				// url, hw_firmware_type, tags, and cow_format are create-only inputs
				// not returned by the API on read.
				ImportStateVerifyIgnore: []string{"url", "hw_firmware_type", "tags", "cow_format"},
				ImportStateIdFunc: acctest.BuildImportID(
					"gcore_cloud_instance_image.test",
					"project_id", "region_id", "id",
				),
			},
			{
				// Post-import convergence. ImportStateVerify is not supported with
				// plannable import blocks, so this is a separate step from the one
				// above.
				//
				// ExpectNonEmptyPlan is required because url is null after import
				// and set in config; ExpectResourceAction(...Update) is the actual
				// regression assertion — with RequiresReplace() on url this plan was
				// a replace and the step fails.
				ResourceName:    "gcore_cloud_instance_image.test",
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithID,
				ImportStateIdFunc: acctest.BuildImportID(
					"gcore_cloud_instance_image.test",
					"project_id", "region_id", "id",
				),
				ExpectNonEmptyPlan: true,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_instance_image.test",
							plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue("gcore_cloud_instance_image.test",
							tfjsonpath.New("name"), knownvalue.StringExact(rName)),
						plancheck.ExpectKnownValue("gcore_cloud_instance_image.test",
							tfjsonpath.New("url"), knownvalue.StringExact(testImageURL)),
					},
				},
			},
			{
				// The managed resource must still be drift-free afterwards.
				Config: testAccCloudInstanceImageConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

// TestAccCloudInstanceImage_urlChangeForcesReplacement is the negative
// counterpart of the import test: when the prior url in state is known and the
// configured url actually changes, the resource must still be replaced.
// RequiresReplaceIfPriorValueKnown() only relaxes the replacement when the prior
// value is null (i.e. post-import).
//
// This is the failure mode of the earlier url_wo (write-only) attempt, which
// silently did nothing on a url change.
//
// The change step is a real apply rather than PlanOnly: the framework skips
// ConfigPlanChecks.PreApply entirely for PlanOnly steps, and asserting the
// replacement is the whole point of the test. The replacement image
// (cirros, ~20MB) is deliberately the smallest one available.
func TestAccCloudInstanceImage_urlChangeForcesReplacement(t *testing.T) {
	rName := acctest.RandomName()

	compareIDDiffer := statecheck.CompareValue(compare.ValuesDiffer())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudInstanceImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstanceImageConfigWithURL(rName, testImageURL),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_instance_image.test",
						tfjsonpath.New("url"), knownvalue.StringExact(testImageURL)),
					compareIDDiffer.AddStateValue(
						"gcore_cloud_instance_image.test",
						tfjsonpath.New("id"),
					),
				},
			},
			{
				Config: testAccCloudInstanceImageConfigWithURL(rName, changedTestImageURL),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_instance_image.test",
							plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_instance_image.test",
						tfjsonpath.New("url"), knownvalue.StringExact(changedTestImageURL)),
					// A different ID proves the image was really recreated.
					compareIDDiffer.AddStateValue(
						"gcore_cloud_instance_image.test",
						tfjsonpath.New("id"),
					),
				},
			},
		},
	})
}

func testAccCheckCloudInstanceImageDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cloud_instance_image" {
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

		// Deletion is asynchronous: the image record can still be returned by
		// GET for a short window after the delete task completes. Poll until it
		// 404s (deleted) or the timeout elapses, instead of asserting once.
		deadline := time.Now().Add(30 * time.Second)
		for {
			_, err = client.Cloud.Instances.Images.Get(
				context.Background(),
				rs.Primary.ID,
				cloud.InstanceImageGetParams{
					ProjectID: param.NewOpt(projectID),
					RegionID:  param.NewOpt(regionID),
				},
			)
			if err != nil {
				if acctest.IsNotFoundError(err) {
					break // deleted as expected
				}
				return fmt.Errorf("error checking instance image deletion: %w", err)
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("instance image %s still exists", rs.Primary.ID)
			}
			time.Sleep(3 * time.Second)
		}
	}
	return nil
}

func testAccCloudInstanceImageConfig(name string) string {
	return testAccCloudInstanceImageConfigWithURL(name, testImageURL)
}

func testAccCloudInstanceImageConfigWithURL(name, url string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_instance_image" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  url        = %[4]q
}`, acctest.ProjectID(), acctest.RegionID(), name, url)
}

func TestAccCloudInstanceImage_update(t *testing.T) {
	rName := acctest.RandomName()
	rNameUpdated := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudInstanceImageDestroy,
		Steps: []resource.TestStep{
			// Create with optional attributes (os_type, ssh_key) and tags set.
			{
				Config: testAccCloudInstanceImageConfigOptional(rName, `{
    env  = "acctest"
    team = "dx"
  }`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_instance_image.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_instance_image.test",
						tfjsonpath.New("os_type"), knownvalue.StringExact("linux")),
					statecheck.ExpectKnownValue("gcore_cloud_instance_image.test",
						tfjsonpath.New("ssh_key"), knownvalue.StringExact("allow")),
					statecheck.ExpectKnownValue("gcore_cloud_instance_image.test",
						tfjsonpath.New("tags"), knownvalue.MapExact(map[string]knownvalue.Check{
							"env":  knownvalue.StringExact("acctest"),
							"team": knownvalue.StringExact("dx"),
						})),
				},
			},
			// In-place update: change name and modify tags (drop "team", change "env").
			{
				Config: testAccCloudInstanceImageConfigOptional(rNameUpdated, `{
    env = "acctest-updated"
  }`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_instance_image.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rNameUpdated)),
					statecheck.ExpectKnownValue("gcore_cloud_instance_image.test",
						tfjsonpath.New("tags"), knownvalue.MapExact(map[string]knownvalue.Check{
							"env": knownvalue.StringExact("acctest-updated"),
						})),
				},
			},
			// Remove all tags - exercises the Update tag-removal path (resource.go
			// rewrites "tags":null to "tags":{} so the API accepts removal).
			{
				Config: testAccCloudInstanceImageConfigOptional(rNameUpdated, `{}`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_instance_image.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rNameUpdated)),
					statecheck.ExpectKnownValue("gcore_cloud_instance_image.test",
						tfjsonpath.New("tags"), knownvalue.MapSizeExact(0)),
				},
			},
		},
	})
}

func testAccCloudInstanceImageConfigOptional(name, tags string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_instance_image" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  url        = %[4]q
  os_type    = "linux"
  ssh_key    = "allow"
  tags       = %[5]s
}`, acctest.ProjectID(), acctest.RegionID(), name, testImageURL, tags)
}
