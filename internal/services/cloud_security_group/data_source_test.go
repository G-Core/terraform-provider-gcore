package cloud_security_group_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/stainless-sdks/gcore-terraform/internal/acctest"
)

func TestAccCloudSecurityGroupDataSource_byID(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSecurityGroupDataSourceConfigByID(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_security_group.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.CompareValuePairs(
						"gcore_cloud_security_group.test", tfjsonpath.New("id"),
						"data.gcore_cloud_security_group.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func TestAccCloudSecurityGroupDataSource_byName(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSecurityGroupDataSourceConfigByName(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_security_group.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.CompareValuePairs(
						"gcore_cloud_security_group.test", tfjsonpath.New("id"),
						"data.gcore_cloud_security_group.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func testAccCloudSecurityGroupDataSourceConfigByID(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_security_group" "test" {
  project_id  = %[1]s
  region_id   = %[2]s
  name        = %[3]q
  description = "test description"
}

data "gcore_cloud_security_group" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  group_id   = gcore_cloud_security_group.test.id
}`, acctest.ProjectID(), acctest.RegionID(), name)
}

func testAccCloudSecurityGroupDataSourceConfigByName(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_security_group" "test" {
  project_id  = %[1]s
  region_id   = %[2]s
  name        = %[3]q
  description = "test description"
}

data "gcore_cloud_security_group" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  find_one_by = {
    name = gcore_cloud_security_group.test.name
  }
}`, acctest.ProjectID(), acctest.RegionID(), name)
}
