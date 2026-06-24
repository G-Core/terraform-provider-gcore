package cdn_logs_uploader_config_test

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

func TestAccCDNLogsUploaderConfigDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()
	policyID, targetID := testAccLogsUploaderConfigDeps(t, rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNLogsUploaderConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNLogsUploaderConfigDataSourceConfig(rName, policyID, targetID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cdn_logs_uploader_config.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("data.gcore_cdn_logs_uploader_config.test",
						tfjsonpath.New("policy"), knownvalue.Int64Exact(policyID)),
					statecheck.ExpectKnownValue("data.gcore_cdn_logs_uploader_config.test",
						tfjsonpath.New("target"), knownvalue.Int64Exact(targetID)),
					// data source and resource resolve to the same config id
					statecheck.CompareValuePairs(
						"gcore_cdn_logs_uploader_config.test", tfjsonpath.New("id"),
						"data.gcore_cdn_logs_uploader_config.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func testAccCDNLogsUploaderConfigDataSourceConfig(name string, policyID, targetID int64) string {
	return fmt.Sprintf(`
resource "gcore_cdn_logs_uploader_config" "test" {
  name   = %[1]q
  policy = %[2]d
  target = %[3]d
}

data "gcore_cdn_logs_uploader_config" "test" {
  id = gcore_cdn_logs_uploader_config.test.id
}`, name, policyID, targetID)
}
