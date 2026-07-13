package cdn_logs_uploader_target_test

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

func TestAccCDNLogsUploaderTargetDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNLogsUploaderTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNLogsUploaderTargetDataSourceConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("data.gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("description"), knownvalue.StringExact("data source acctest")),
					statecheck.ExpectKnownValue("data.gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("storage_type"), knownvalue.StringExact("s3_other")),
					statecheck.ExpectKnownValue("data.gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("client_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("created"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("access_key_id"), knownvalue.StringExact("ds-access-key")),
					statecheck.ExpectKnownValue("data.gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("bucket_name"), knownvalue.StringExact("ds-bucket")),
					statecheck.ExpectKnownValue("data.gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("region"), knownvalue.StringExact("us-east-1")),
					statecheck.ExpectKnownValue("data.gcore_cdn_logs_uploader_target.test",
						tfjsonpath.New("config").AtMapKey("endpoint"), knownvalue.StringExact("https://s3-ds.example.com")),
					statecheck.CompareValuePairs(
						"gcore_cdn_logs_uploader_target.test", tfjsonpath.New("id"),
						"data.gcore_cdn_logs_uploader_target.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func testAccCDNLogsUploaderTargetDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_logs_uploader_target" "test" {
  name         = %[1]q
  description  = "data source acctest"
  storage_type = "s3_other"

  config = {
    access_key_id     = "ds-access-key"
    secret_access_key = "ds-secret-key-12345"
    bucket_name       = "ds-bucket"
    region            = "us-east-1"
    endpoint          = "https://s3-ds.example.com"
    use_path_style    = true
  }
}

data "gcore_cdn_logs_uploader_target" "test" {
  id = gcore_cdn_logs_uploader_target.test.id
}`, name)
}
