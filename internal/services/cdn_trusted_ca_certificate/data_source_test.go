package cdn_trusted_ca_certificate_test

import (
	"fmt"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCDNTrustedCaCertificateDataSource_basic(t *testing.T) {
	t.Parallel()
	rName := acctest.RandomName()
	resourceName := "gcore_cdn_trusted_ca_certificate.test"
	dataSourceName := "data.gcore_cdn_trusted_ca_certificate.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNTrustedCaCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNTrustedCaCertificateDataSourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "cert_issuer", resourceName, "cert_issuer"),
					resource.TestCheckResourceAttrPair(dataSourceName, "cert_subject_cn", resourceName, "cert_subject_cn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "validity_not_after", resourceName, "validity_not_after"),
					resource.TestCheckResourceAttrPair(dataSourceName, "validity_not_before", resourceName, "validity_not_before"),
				),
			},
		},
	})
}

func testAccCDNTrustedCaCertificateDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_trusted_ca_certificate" "test" {
  name            = %[1]q
  ssl_certificate = <<-EOT
%[2]s
EOT
}

data "gcore_cdn_trusted_ca_certificate" "test" {
  id = gcore_cdn_trusted_ca_certificate.test.id
}
`, name, testCACertPEM)
}
