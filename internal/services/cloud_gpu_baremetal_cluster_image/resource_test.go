package cloud_gpu_baremetal_cluster_image_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
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

// testImageURLAlternate is a second, different image URL used to prove that a
// genuine change of the create-only `url` attribute still forces replacement.
const testImageURLAlternate = "http://download.cirros-cloud.net/0.6.2/cirros-0.6.2-x86_64-disk.img"

func TestAccCloudGPUBaremetalClusterImage_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudGPUBaremetalClusterImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudGPUBaremetalClusterImageConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_baremetal_cluster_image.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_baremetal_cluster_image.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_baremetal_cluster_image.test",
						tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_baremetal_cluster_image.test",
						tfjsonpath.New("architecture"), knownvalue.StringExact("x86_64")),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_baremetal_cluster_image.test",
						tfjsonpath.New("os_type"), knownvalue.StringExact("linux")),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_baremetal_cluster_image.test",
						tfjsonpath.New("ssh_key"), knownvalue.StringExact("allow")),
					statecheck.ExpectKnownValue("gcore_cloud_gpu_baremetal_cluster_image.test",
						tfjsonpath.New("cow_format"), knownvalue.Bool(false)),
				},
			},
		},
	})
}

// TestAccCloudGPUBaremetalClusterImage_import covers importing the resource and,
// crucially, what the *next* plan after the import looks like.
//
// url is Required, create-only and never returned by the GET endpoint
// (no_refresh), so it lands in state as null after an import. While url carried
// stringplanmodifier.RequiresReplace(), that null-in-state vs set-in-config
// difference made the follow-up plan destroy and recreate live infrastructure.
// A plain ImportState + ImportStateVerify step cannot see that: it only compares
// the imported state against the pre-import state and never plans afterwards.
//
// The third step closes that gap with ImportStateKind: resource.ImportBlockWithID
// — config-driven import, which plans an `import` block against the real config.
// url is null in the freshly imported state and set in config, so the plan is
// legitimately non-empty (hence ExpectNonEmptyPlan); the regression assertion is
// that it must be an in-place UPDATE, never a replace. When a user applies that
// update, the adoption is served by this resource's real GET-backed Update — a
// no-op Update stub fails such an apply with "Provider produced inconsistent
// result after apply", because the framework pre-populates UpdateResponse.State
// with the prior (url-null) state.
func TestAccCloudGPUBaremetalClusterImage_import(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudGPUBaremetalClusterImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudGPUBaremetalClusterImageConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_baremetal_cluster_image.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
				},
			},
			{
				ResourceName:      "gcore_cloud_gpu_baremetal_cluster_image.test",
				ImportState:       true,
				ImportStateVerify: true,
				// url, hw_firmware_type, tags, and cow_format are create-only
				// inputs not returned by the API on read.
				ImportStateVerifyIgnore: []string{"url", "hw_firmware_type", "tags", "cow_format"},
				ImportStateIdFunc: acctest.BuildImportID(
					"gcore_cloud_gpu_baremetal_cluster_image.test",
					"project_id", "region_id", "id",
				),
			},
			{
				// Post-import convergence. ImportStateVerify is not supported
				// with plannable import blocks, so this is a separate step from
				// the one above.
				//
				// ExpectNonEmptyPlan is required because url is null after import
				// and set in config; ExpectResourceAction(...Update) is the
				// actual regression assertion — with RequiresReplace() on url
				// this plan was a replace and the step fails.
				ResourceName:    "gcore_cloud_gpu_baremetal_cluster_image.test",
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithID,
				ImportStateIdFunc: acctest.BuildImportID(
					"gcore_cloud_gpu_baremetal_cluster_image.test",
					"project_id", "region_id", "id",
				),
				ExpectNonEmptyPlan: true,
				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_gpu_baremetal_cluster_image.test",
							plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue("gcore_cloud_gpu_baremetal_cluster_image.test",
							tfjsonpath.New("name"), knownvalue.StringExact(rName)),
						plancheck.ExpectKnownValue("gcore_cloud_gpu_baremetal_cluster_image.test",
							tfjsonpath.New("url"), knownvalue.StringExact(testImageURL)),
					},
				},
			},
			{
				// The managed resource must still be drift-free afterwards.
				Config: testAccCloudGPUBaremetalClusterImageConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

// TestAccCloudGPUBaremetalClusterImage_urlChangeForcesReplacement is the negative
// counterpart to the import test: RequiresReplaceIfPriorValueKnown() must keep
// replacing the image when the prior `url` in state is known and the config
// changes it. The write-only `url_wo` hack this replaced silently did nothing
// here, so assert the Replace action explicitly.
func TestAccCloudGPUBaremetalClusterImage_urlChangeForcesReplacement(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudGPUBaremetalClusterImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudGPUBaremetalClusterImageConfigWithURL(rName, testImageURL),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_baremetal_cluster_image.test",
						tfjsonpath.New("url"), knownvalue.StringExact(testImageURL)),
				},
			},
			{
				Config: testAccCloudGPUBaremetalClusterImageConfigWithURL(rName, testImageURLAlternate),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_gpu_baremetal_cluster_image.test",
							plancheck.ResourceActionReplace),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_gpu_baremetal_cluster_image.test",
						tfjsonpath.New("url"), knownvalue.StringExact(testImageURLAlternate)),
				},
			},
		},
	})
}

func testAccCheckCloudGPUBaremetalClusterImageDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cloud_gpu_baremetal_cluster_image" {
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

		_, err = client.Cloud.GPUBaremetal.Clusters.Images.Get(
			context.Background(),
			rs.Primary.ID,
			cloud.GPUBaremetalClusterImageGetParams{
				ProjectID: param.NewOpt(projectID),
				RegionID:  param.NewOpt(regionID),
			},
		)

		if err == nil {
			return fmt.Errorf("gpu baremetal cluster image %s still exists", rs.Primary.ID)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking gpu baremetal cluster image deletion: %w", err)
		}
	}
	return nil
}

func testAccCloudGPUBaremetalClusterImageConfig(name string) string {
	return testAccCloudGPUBaremetalClusterImageConfigWithURL(name, testImageURL)
}

func testAccCloudGPUBaremetalClusterImageConfigWithURL(name, url string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_gpu_baremetal_cluster_image" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  url        = %[4]q
}`, acctest.ProjectID(), acctest.RegionID(), name, url)
}
