package cloud_inference_registry_credential_test

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

func TestAccCloudInferenceRegistryCredentialDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudInferenceRegistryCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInferenceRegistryCredentialDataSourceConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_inference_registry_credential.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.CompareValuePairs(
						"gcore_cloud_inference_registry_credential.test", tfjsonpath.New("id"),
						"data.gcore_cloud_inference_registry_credential.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
					statecheck.CompareValuePairs(
						"gcore_cloud_inference_registry_credential.test", tfjsonpath.New("username"),
						"data.gcore_cloud_inference_registry_credential.test", tfjsonpath.New("username"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func testAccCloudInferenceRegistryCredentialDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_inference_registry_credential" "test" {
  project_id          = %[1]s
  name                = %[2]q
  registry_url        = "https://registry.example.com"
  username            = "testuser"
  password_wo         = "testpassword"
  password_wo_version = 1
}

data "gcore_cloud_inference_registry_credential" "test" {
  project_id      = %[1]s
  credential_name = gcore_cloud_inference_registry_credential.test.id
}`, acctest.ProjectID(), name)
}
