package cloud_instance_test

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

	"github.com/stainless-sdks/gcore-terraform/internal/acctest"
)

func testAccCheckCloudInstanceDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cloud_instance" {
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

		_, err = client.Cloud.Instances.Get(
			context.Background(),
			rs.Primary.ID,
			cloud.InstanceGetParams{
				ProjectID: param.NewOpt(projectID),
				RegionID:  param.NewOpt(regionID),
			},
		)
		if err == nil {
			return fmt.Errorf("cloud instance %s still exists", rs.Primary.ID)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking instance deletion: %w", err)
		}
	}
	return nil
}

func TestAccCloudInstance_basic(t *testing.T) {
	rName := acctest.RandomName()
	imageID := acctest.LatestUbuntuImageID(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstanceConfig(rName, imageID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_instance.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_instance.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_instance.test",
						tfjsonpath.New("status"), knownvalue.NotNull()),
				},
			},
			{
				ResourceName:      "gcore_cloud_instance.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password_wo",
					"password_wo_version",
					"allow_app_ports",
					"user_data",
					"volumes",
				},
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_instance.test", "project_id", "region_id", "id"),
			},
		},
	})
}

func TestAccCloudInstance_update(t *testing.T) {
	rName := acctest.RandomName()
	rNameUpdated := acctest.RandomName()
	imageID := acctest.LatestUbuntuImageID(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstanceConfig(rName, imageID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_instance.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
				},
			},
			{
				Config: testAccCloudInstanceConfig(rNameUpdated, imageID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_instance.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rNameUpdated)),
				},
			},
		},
	})
}

func testAccCloudInstanceConfig(name, imageID string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_volume" "boot" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = "%[3]s-vol"
  size       = 10
  type_name  = "ssd_hiiops"
  source     = "image"
  image_id   = %[4]q
}

resource "gcore_cloud_instance" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  flavor     = "g1-standard-1-2"

  volumes = [
    {
      volume_id  = gcore_cloud_volume.boot.id
      boot_index = 0
    }
  ]

  interfaces = [
    {
      type = "external"
    }
  ]
}`, acctest.ProjectID(), acctest.RegionID(), name, imageID)
}
