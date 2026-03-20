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
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

// testImageURL is a small, publicly-accessible Linux cloud image used for acceptance tests.
// Using a known-good minimal image to keep upload times short.
const testImageURL = "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud-2003.qcow2.xz"

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
				// url, hw_firmware_type, tags, and cow_format are not returned by the API (no_refresh)
				ImportStateVerifyIgnore: []string{"url", "hw_firmware_type", "tags", "cow_format"},
				ImportStateIdFunc: acctest.BuildImportID(
					"gcore_cloud_gpu_baremetal_cluster_image.test",
					"project_id", "region_id", "id",
				),
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
	return fmt.Sprintf(`
resource "gcore_cloud_gpu_baremetal_cluster_image" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  url        = %[4]q
}`, acctest.ProjectID(), acctest.RegionID(), name, testImageURL)
}
