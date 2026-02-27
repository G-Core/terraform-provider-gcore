package cloud_placement_group_test

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

func TestAccCloudPlacementGroup_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudPlacementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudPlacementGroupConfig(rName, "affinity"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_placement_group.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_placement_group.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_placement_group.test",
						tfjsonpath.New("servergroup_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_placement_group.test",
						tfjsonpath.New("policy"), knownvalue.StringExact("affinity")),
					statecheck.ExpectKnownValue("gcore_cloud_placement_group.test",
						tfjsonpath.New("region"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccCloudPlacementGroup_antiAffinity(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudPlacementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudPlacementGroupConfig(rName, "anti-affinity"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_placement_group.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_placement_group.test",
						tfjsonpath.New("policy"), knownvalue.StringExact("anti-affinity")),
				},
			},
		},
	})
}

func TestAccCloudPlacementGroup_import(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudPlacementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudPlacementGroupConfig(rName, "affinity"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_placement_group.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
				},
			},
			{
				ResourceName:      "gcore_cloud_placement_group.test",
				ImportState:       true,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_placement_group.test", "project_id", "region_id", "servergroup_id"),
				ImportStateVerify: true,
			},
		},
	})
}

// testAccCheckCloudPlacementGroupDestroy verifies the placement group is deleted.
// This resource requires project_id and region_id in the API call.
func testAccCheckCloudPlacementGroupDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cloud_placement_group" {
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

		_, err = client.Cloud.PlacementGroups.Get(context.Background(), rs.Primary.ID, cloud.PlacementGroupGetParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})

		if err == nil {
			return fmt.Errorf("placement group %s still exists", rs.Primary.ID)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking placement group deletion: %w", err)
		}
	}
	return nil
}

func testAccCloudPlacementGroupConfig(name, policy string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_placement_group" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  policy     = %[4]q
}`, acctest.ProjectID(), acctest.RegionID(), name, policy)
}
