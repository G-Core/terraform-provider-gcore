package cloud_instance_image_test

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

		_, err = client.Cloud.Instances.Images.Get(
			context.Background(),
			rs.Primary.ID,
			cloud.InstanceImageGetParams{
				ProjectID: param.NewOpt(projectID),
				RegionID:  param.NewOpt(regionID),
			},
		)

		if err == nil {
			return fmt.Errorf("instance image %s still exists", rs.Primary.ID)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking instance image deletion: %w", err)
		}
	}
	return nil
}

func testAccCloudInstanceImageConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_instance_image" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  url        = %[4]q
}`, acctest.ProjectID(), acctest.RegionID(), name, testImageURL)
}
