package cdn_certificate_test

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

func TestAccCDNCertificateDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNCertificateDataSourceConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					// Verify data source reads the same certificate
					statecheck.ExpectKnownValue("data.gcore_cdn_certificate.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("data.gcore_cdn_certificate.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.gcore_cdn_certificate.test",
						tfjsonpath.New("automated"), knownvalue.Bool(true)),
					// Verify data source and resource have same ssl_id
					statecheck.CompareValuePairs(
						"gcore_cdn_certificate.test", tfjsonpath.New("ssl_id"),
						"data.gcore_cdn_certificate.test", tfjsonpath.New("ssl_id"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func testAccCDNCertificateDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_certificate" "test" {
  name      = %[1]q
  automated = true
}

data "gcore_cdn_certificate" "test" {
  ssl_id = gcore_cdn_certificate.test.ssl_id
}`, name)
}
