package cdn_logs_uploader_policy_test

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

func TestAccCDNLogsUploaderPolicyDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNLogsUploaderPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNLogsUploaderPolicyDataSourceConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("data.gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("description"), knownvalue.StringExact("acctest policy")),
					statecheck.ExpectKnownValue("data.gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("format_type"), knownvalue.StringExact("json")),
					statecheck.ExpectKnownValue("data.gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("field_delimiter"), knownvalue.StringExact(",")),
					statecheck.ExpectKnownValue("data.gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("fields"), knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("remote_addr"),
							knownvalue.StringExact("status"),
						})),
					statecheck.ExpectKnownValue("data.gcore_cdn_logs_uploader_policy.test",
						tfjsonpath.New("tags"), knownvalue.MapExact(map[string]knownvalue.Check{
							"env": knownvalue.StringExact("test"),
						})),
					statecheck.CompareValuePairs(
						"gcore_cdn_logs_uploader_policy.test", tfjsonpath.New("id"),
						"data.gcore_cdn_logs_uploader_policy.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func testAccCDNLogsUploaderPolicyDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_logs_uploader_policy" "test" {
  name                     = %[1]q
  description              = "acctest policy"
  format_type              = "json"
  fields                   = ["remote_addr", "status"]
  field_delimiter          = ","
  field_separator          = "|"
  include_shield_logs      = true
  log_sample_rate          = 0.5
  tags = {
    env = "test"
  }
}

data "gcore_cdn_logs_uploader_policy" "test" {
  id = gcore_cdn_logs_uploader_policy.test.id
}`, name)
}
