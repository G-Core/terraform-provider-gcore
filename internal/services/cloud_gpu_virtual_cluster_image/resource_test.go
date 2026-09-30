package cloud_gpu_virtual_cluster_image_test

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

// The CentOS image URL used by these tests is declared once for the package as
// listTestImageURL in list_data_source_test.go.
//
// changedTestImageURL is a second, distinct publicly-accessible cloud image. It
// is only used to prove that a genuine url change still forces replacement.
const changedTestImageURL = "http://download.cirros-cloud.net/0.6.2/cirros-0.6.2-x86_64-disk.img"

func TestAccCloudGPUVirtualClusterImage_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudGPUVirtualClusterImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudGPUVirtualClusterImageConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster_image.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster_image.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster_image.test",
						tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster_image.test",
						tfjsonpath.New("url"), knownvalue.StringExact(listTestImageURL)),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster_image.test",
						tfjsonpath.New("architecture"), knownvalue.StringExact("x86_64")),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster_image.test",
						tfjsonpath.New("os_type"), knownvalue.StringExact("linux")),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster_image.test",
						tfjsonpath.New("ssh_key"), knownvalue.StringExact("allow")),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster_image.test",
						tfjsonpath.New("cow_format"), knownvalue.Bool(false)),
				},
			},
		},
	})
}

// TestAccCloudGPUVirtualClusterImage_import covers importing the resource and,
// crucially, what happens on the *next* plan/apply after the import.
//
// url is Required, create-only and never returned by the GET endpoint
// (no_refresh), so it lands in state as null after an import. While url carried
// stringplanmodifier.RequiresReplace(), that null-vs-configured difference made
// the follow-up plan destroy and recreate live infrastructure. A plain
// ImportState + ImportStateVerify step cannot see that: it only compares the
// imported state and never plans afterwards.
//
// The second import step therefore uses ImportStateKind:
// resource.ImportBlockWithID — config-driven import, which plans an `import`
// block against the configuration. url is null in the freshly imported state and
// set in config, so the plan is legitimately non-empty (hence
// ExpectNonEmptyPlan); the regression assertion is that the plan must be an
// in-place UPDATE, never a replace. Applying that adoption is served by this
// resource's real GET-backed Update — the old no-op Update stub failed such an
// apply with "Provider produced inconsistent result after apply".
func TestAccCloudGPUVirtualClusterImage_import(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudGPUVirtualClusterImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudGPUVirtualClusterImageConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster_image.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
				},
			},
			{
				ResourceName:      "gcore_cloud_gpu_virtual_cluster_image.test",
				ImportState:       true,
				ImportStateVerify: true,
				// url is the only create-only input the API never returns;
				// cow_format is recovered from disk_format on import, and
				// hw_firmware_type and tags are echoed by GET.
				ImportStateVerifyIgnore: []string{"url"},
				ImportStateIdFunc: acctest.BuildImportID(
					"gcore_cloud_gpu_virtual_cluster_image.test",
					"project_id", "region_id", "id",
				),
			},
			{
				// Post-import convergence. ExpectNonEmptyPlan is required
				// because url is null after import and set in config;
				// ExpectResourceAction(...Update) is the actual regression
				// assertion — with RequiresReplace() on url this plan was a
				// replace and the step fails.
				ResourceName:    "gcore_cloud_gpu_virtual_cluster_image.test",
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithID,
				ImportStateIdFunc: acctest.BuildImportID(
					"gcore_cloud_gpu_virtual_cluster_image.test",
					"project_id", "region_id", "id",
				),
				ExpectNonEmptyPlan: true,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_gpu_virtual_cluster_image.test",
							plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster_image.test",
							tfjsonpath.New("name"), knownvalue.StringExact(rName)),
						plancheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster_image.test",
							tfjsonpath.New("url"), knownvalue.StringExact(listTestImageURL)),
					},
				},
			},
			{
				// The managed resource must still be drift-free afterwards.
				Config: testAccCloudGPUVirtualClusterImageConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

// TestAccCloudGPUVirtualClusterImage_importAdoptionApplies runs the leg the
// import-block test cannot: the adopting apply itself. The import-block step
// above is plan-only, so the GET-backed Update that persists url after an
// import is not executed there. This test imports with `terraform import`
// semantics (ImportCommandWithID + ImportStatePersist), applies the config and
// checks that url lands in state, the image is not replaced, and the next plan
// is empty.
//
// The image is uploaded with cow_format = true. The API never returns
// cow_format, so before it was derived from disk_format on import the same
// config planned a replacement over cow_format after every import.
func TestAccCloudGPUVirtualClusterImage_importAdoptionApplies(t *testing.T) {
	rName := acctest.RandomName()
	var imageID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudGPUVirtualClusterImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudGPUVirtualClusterImageConfigCow(rName),
				Check: func(s *terraform.State) error {
					rs, ok := s.RootModule().Resources["gcore_cloud_gpu_virtual_cluster_image.test"]
					if !ok {
						return fmt.Errorf("gcore_cloud_gpu_virtual_cluster_image.test not found in state")
					}
					imageID = rs.Primary.ID
					return nil
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster_image.test",
						tfjsonpath.New("disk_format"), knownvalue.StringExact("raw")),
				},
			},
			{
				// Forget the image without destroying it, so the import below
				// targets an address absent from state.
				Config: testAccCloudGPUVirtualClusterImageConfigForgotten(),
			},
			{
				Config:             testAccCloudGPUVirtualClusterImageConfigCow(rName),
				ResourceName:       "gcore_cloud_gpu_virtual_cluster_image.test",
				ImportState:        true,
				ImportStateKind:    resource.ImportCommandWithID,
				ImportStatePersist: true,
				ImportStateIdFunc: func(*terraform.State) (string, error) {
					return fmt.Sprintf("%s/%s/%s", acctest.ProjectID(), acctest.RegionID(), imageID), nil
				},
			},
			{
				// The adopting apply: url lands in state in place, cow_format
				// was already recovered by the import, nothing is replaced.
				Config: testAccCloudGPUVirtualClusterImageConfigCow(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_gpu_virtual_cluster_image.test",
							plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster_image.test",
						tfjsonpath.New("id"), knownvalue.StringFunc(func(v string) error {
							if v != imageID {
								return fmt.Errorf("image was replaced: %s != %s", v, imageID)
							}
							return nil
						})),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster_image.test",
						tfjsonpath.New("url"), knownvalue.StringExact(listTestImageURL)),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster_image.test",
						tfjsonpath.New("cow_format"), knownvalue.Bool(true)),
				},
			},
			{
				// Adoption is stable: the same config plans nothing.
				Config:   testAccCloudGPUVirtualClusterImageConfigCow(rName),
				PlanOnly: true,
			},
			{
				// Dropping cow_format from config keeps the recovered value:
				// "not configured" must not turn into a silent in-place update
				// to false on a raw image.
				Config: testAccCloudGPUVirtualClusterImageConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster_image.test",
						tfjsonpath.New("cow_format"), knownvalue.Bool(true)),
				},
			},
			{
				// A value that really differs from the image (false on a raw
				// image) is a genuine change and still forces replacement;
				// the replacement is stored as qcow2.
				Config: testAccCloudGPUVirtualClusterImageConfigCowFormat(rName, false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_gpu_virtual_cluster_image.test",
							plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster_image.test",
						tfjsonpath.New("cow_format"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster_image.test",
						tfjsonpath.New("disk_format"), knownvalue.StringExact("qcow2")),
				},
			},
		},
	})
}

// TestAccCloudGPUVirtualClusterImage_urlChangeForcesReplacement is the negative
// counterpart of the import test: when the prior url in state is known and the
// configured url actually changes, the resource must still be replaced.
//
// This is the failure mode of the earlier url_wo (write-only) attempt, which
// silently did nothing on a url change.
func TestAccCloudGPUVirtualClusterImage_urlChangeForcesReplacement(t *testing.T) {
	rName := acctest.RandomName()

	compareIDDiffer := statecheck.CompareValue(compare.ValuesDiffer())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudGPUVirtualClusterImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudGPUVirtualClusterImageConfigWithURL(rName, listTestImageURL),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster_image.test",
						tfjsonpath.New("url"), knownvalue.StringExact(listTestImageURL)),
					compareIDDiffer.AddStateValue(
						"gcore_cloud_gpu_virtual_cluster_image.test",
						tfjsonpath.New("id"),
					),
				},
			},
			{
				Config: testAccCloudGPUVirtualClusterImageConfigWithURL(rName, changedTestImageURL),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_gpu_virtual_cluster_image.test",
							plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_virtual_cluster_image.test",
						tfjsonpath.New("url"), knownvalue.StringExact(changedTestImageURL)),
					// A different ID proves the image was really recreated.
					compareIDDiffer.AddStateValue(
						"gcore_cloud_gpu_virtual_cluster_image.test",
						tfjsonpath.New("id"),
					),
				},
			},
		},
	})
}

func testAccCheckCloudGPUVirtualClusterImageDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cloud_gpu_virtual_cluster_image" {
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

		_, err = client.Cloud.GPUVirtual.Clusters.Images.Get(
			context.Background(),
			rs.Primary.ID,
			cloud.GPUVirtualClusterImageGetParams{
				ProjectID: param.NewOpt(projectID),
				RegionID:  param.NewOpt(regionID),
			},
		)

		if err == nil {
			return fmt.Errorf("gpu virtual cluster image %s still exists", rs.Primary.ID)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking gpu virtual cluster image deletion: %w", err)
		}
	}
	return nil
}

func testAccCloudGPUVirtualClusterImageConfig(name string) string {
	return testAccCloudGPUVirtualClusterImageConfigWithURL(name, listTestImageURL)
}

// testAccCloudGPUVirtualClusterImageConfigCow uploads the image with
// cow_format = true, which the API stores as disk_format raw.
func testAccCloudGPUVirtualClusterImageConfigCow(name string) string {
	return testAccCloudGPUVirtualClusterImageConfigCowFormat(name, true)
}

func testAccCloudGPUVirtualClusterImageConfigCowFormat(name string, cowFormat bool) string {
	return fmt.Sprintf(`
resource "gcore_cloud_gpu_virtual_cluster_image" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  url        = %[4]q
  cow_format = %[5]t
}`, acctest.ProjectID(), acctest.RegionID(), name, listTestImageURL, cowFormat)
}

// testAccCloudGPUVirtualClusterImageConfigForgotten replaces the resource block
// with a removed block (Terraform >= 1.7): the state entry goes away, the live
// image stays, and a later step can import it again.
func testAccCloudGPUVirtualClusterImageConfigForgotten() string {
	return `
removed {
  from = gcore_cloud_gpu_virtual_cluster_image.test

  lifecycle {
    destroy = false
  }
}`
}

func testAccCloudGPUVirtualClusterImageConfigWithURL(name, url string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_gpu_virtual_cluster_image" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  url        = %[4]q
}`, acctest.ProjectID(), acctest.RegionID(), name, url)
}
