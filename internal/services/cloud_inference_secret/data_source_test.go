package cloud_inference_secret_test

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

func TestAccCloudInferenceSecretDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudInferenceSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudInferenceSecretDataSourceConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_inference_secret.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("data.gcore_cloud_inference_secret.test",
						tfjsonpath.New("type"), knownvalue.StringExact("aws-iam")),
					statecheck.CompareValuePairs(
						"gcore_cloud_inference_secret.test", tfjsonpath.New("id"),
						"data.gcore_cloud_inference_secret.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func testAccCloudInferenceSecretDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_inference_secret" "test" {
  project_id      = %[1]s
  name            = %[2]q
  type            = "aws-iam"
  data_wo_version = 1
  data = {
    aws_access_key_id_wo     = "AKIAIOSFODNN7EXAMPLE"
    aws_secret_access_key_wo = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
  }
}

data "gcore_cloud_inference_secret" "test" {
  project_id  = %[1]s
  secret_name = gcore_cloud_inference_secret.test.id
}`, acctest.ProjectID(), name)
}
