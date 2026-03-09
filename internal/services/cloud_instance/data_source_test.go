package cloud_instance_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func TestAccCloudInstanceDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()
	imageID := latestUbuntuImageID(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInstanceDataSourceConfig(rName, imageID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_instance.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.CompareValuePairs(
						"gcore_cloud_instance.test", tfjsonpath.New("id"),
						"data.gcore_cloud_instance.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func testAccCloudInstanceDataSourceConfig(name, imageID string) string {
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
}

data "gcore_cloud_instance" "test" {
  project_id  = %[1]s
  region_id   = %[2]s
  instance_id = gcore_cloud_instance.test.id
}`, acctest.ProjectID(), acctest.RegionID(), name, imageID)
}
