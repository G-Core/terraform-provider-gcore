package cloud_floating_ip_test

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

func TestAccCloudFloatingIP_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFloatingIPConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_floating_ip.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_floating_ip.test",
						tfjsonpath.New("floating_ip_address"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_floating_ip.test",
						tfjsonpath.New("status"), knownvalue.StringExact("DOWN")),
				},
			},
		},
	})
}

func TestAccCloudFloatingIP_withTags(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFloatingIPConfigWithTags("value1"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_floating_ip.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_floating_ip.test",
						tfjsonpath.New("tags").AtMapKey("test_key"), knownvalue.StringExact("value1")),
				},
			},
			{
				Config: testAccCloudFloatingIPConfigWithTags("value2"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_floating_ip.test",
						tfjsonpath.New("tags").AtMapKey("test_key"), knownvalue.StringExact("value2")),
				},
			},
		},
	})
}

func TestAccCloudFloatingIP_import(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFloatingIPConfig(),
			},
			{
				ResourceName:      "gcore_cloud_floating_ip.test",
				ImportState:       true,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_floating_ip.test", "project_id", "region_id", "id"),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCloudFloatingIPDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cloud_floating_ip" {
			continue
		}

		projectID, _ := strconv.ParseInt(rs.Primary.Attributes["project_id"], 10, 64)
		regionID, _ := strconv.ParseInt(rs.Primary.Attributes["region_id"], 10, 64)

		_, err := client.Cloud.FloatingIPs.Get(context.Background(), rs.Primary.ID, cloud.FloatingIPGetParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})

		if err == nil {
			return fmt.Errorf("floating IP %s still exists", rs.Primary.ID)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking floating IP deletion: %w", err)
		}
	}
	return nil
}

func testAccCloudFloatingIPConfig() string {
	return fmt.Sprintf(`
resource "gcore_cloud_floating_ip" "test" {
  project_id = %[1]s
  region_id  = %[2]s
}`, acctest.ProjectID(), acctest.RegionID())
}

func testAccCloudFloatingIPConfigWithTags(tagValue string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_floating_ip" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  tags = {
    test_key = %[3]q
  }
}`, acctest.ProjectID(), acctest.RegionID(), tagValue)
}
