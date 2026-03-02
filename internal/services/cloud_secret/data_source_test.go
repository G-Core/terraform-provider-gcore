package cloud_secret_test

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

func TestAccCloudSecretDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSecretDataSourceConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_secret.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.CompareValuePairs(
						"gcore_cloud_secret.test", tfjsonpath.New("id"),
						"data.gcore_cloud_secret.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
					statecheck.CompareValuePairs(
						"gcore_cloud_secret.test", tfjsonpath.New("status"),
						"data.gcore_cloud_secret.test", tfjsonpath.New("status"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func testAccCloudSecretDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_secret" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  payload_wo_version = 1
  payload = {
    certificate_wo       = <<EOT
%[4]s
EOT
    certificate_chain_wo = <<EOT
%[4]s
EOT
    private_key_wo       = <<EOT
%[5]s
EOT
  }
}

data "gcore_cloud_secret" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  secret_id  = gcore_cloud_secret.test.id
}`, acctest.ProjectID(), acctest.RegionID(), name, testCertificate, testPrivateKey)
}
