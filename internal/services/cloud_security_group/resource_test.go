package cloud_security_group_test

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
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/stainless-sdks/gcore-terraform/internal/acctest"
)

func TestAccCloudSecurityGroup_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSecurityGroupConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_security_group.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_security_group.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_security_group.test",
						tfjsonpath.New("description"), knownvalue.StringExact("test description")),
					statecheck.ExpectKnownValue("gcore_cloud_security_group.test",
						tfjsonpath.New("region"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_security_group.test",
						tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccCloudSecurityGroup_update(t *testing.T) {
	rName := acctest.RandomName()

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSecurityGroupConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_security_group.test",
						tfjsonpath.New("description"), knownvalue.StringExact("test description")),
					compareIDSame.AddStateValue(
						"gcore_cloud_security_group.test",
						tfjsonpath.New("id"),
					),
				},
			},
			{
				Config: testAccCloudSecurityGroupConfigUpdated(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_security_group.test",
						tfjsonpath.New("description"), knownvalue.StringExact("updated description")),
					compareIDSame.AddStateValue(
						"gcore_cloud_security_group.test",
						tfjsonpath.New("id"),
					),
				},
			},
		},
	})
}

func TestAccCloudSecurityGroup_import(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSecurityGroupConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_security_group.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
				},
			},
			{
				ResourceName:      "gcore_cloud_security_group.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_security_group.test", "project_id", "region_id", "id"),
			},
		},
	})
}

func testAccCheckCloudSecurityGroupDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cloud_security_group" {
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

		_, err = client.Cloud.SecurityGroups.Get(context.Background(), rs.Primary.ID, cloud.SecurityGroupGetParams{
			ProjectID: param.NewOpt(projectID),
			RegionID:  param.NewOpt(regionID),
		})

		if err == nil {
			return fmt.Errorf("security group %s still exists", rs.Primary.ID)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking security group deletion: %w", err)
		}
	}
	return nil
}

func testAccCloudSecurityGroupConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_security_group" "test" {
  project_id  = %[1]s
  region_id   = %[2]s
  name        = %[3]q
  description = "test description"
}`, acctest.ProjectID(), acctest.RegionID(), name)
}

func testAccCloudSecurityGroupConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_security_group" "test" {
  project_id  = %[1]s
  region_id   = %[2]s
  name        = %[3]q
  description = "updated description"
}`, acctest.ProjectID(), acctest.RegionID(), name)
}

func TestAccCloudSecurityGroup_tags(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSecurityGroupConfigWithTags(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_security_group.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_security_group.test",
						tfjsonpath.New("tags"), knownvalue.MapExact(map[string]knownvalue.Check{
							"env":  knownvalue.StringExact("test"),
							"team": knownvalue.StringExact("infra"),
						})),
				},
			},
			{
				Config: testAccCloudSecurityGroupConfigTagsUpdated(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_security_group.test",
						tfjsonpath.New("tags"), knownvalue.MapExact(map[string]knownvalue.Check{
							"env": knownvalue.StringExact("prod"),
						})),
				},
			},
		},
	})
}

func testAccCloudSecurityGroupConfigWithTags(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_security_group" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  tags = {
    env  = "test"
    team = "infra"
  }
}`, acctest.ProjectID(), acctest.RegionID(), name)
}

func testAccCloudSecurityGroupConfigTagsUpdated(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_security_group" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  tags = {
    env = "prod"
  }
}`, acctest.ProjectID(), acctest.RegionID(), name)
}
