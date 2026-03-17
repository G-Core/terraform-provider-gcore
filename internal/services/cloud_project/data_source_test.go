package cloud_project_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func TestAccCloudProjectDataSource_byProjectID(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudProjectDataSourceConfigByProjectID(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_project.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_project.test",
						tfjsonpath.New("name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_project.test",
						tfjsonpath.New("state"), knownvalue.StringExact("ACTIVE")),
					statecheck.ExpectKnownValue("data.gcore_cloud_project.test",
						tfjsonpath.New("is_default"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_project.test",
						tfjsonpath.New("created_at"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccCloudProjectDataSource_defaultProject(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudProjectDataSourceConfigDefault(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_project.default",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_project.default",
						tfjsonpath.New("name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cloud_project.default",
						tfjsonpath.New("state"), knownvalue.StringExact("ACTIVE")),
				},
			},
		},
	})
}

func testAccCloudProjectDataSourceConfigByProjectID() string {
	return fmt.Sprintf(`
data "gcore_cloud_project" "test" {
  project_id = %s
}`, acctest.ProjectID())
}

func testAccCloudProjectDataSourceConfigDefault() string {
	return `
data "gcore_cloud_project" "default" {}
`
}
