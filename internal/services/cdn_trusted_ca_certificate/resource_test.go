package cdn_trusted_ca_certificate_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func TestAccCDNTrustedCaCertificate_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNTrustedCaCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNTrustedCaCertificateConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_trusted_ca_certificate.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cdn_trusted_ca_certificate.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cdn_trusted_ca_certificate.test",
						tfjsonpath.New("cert_subject_cn"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cdn_trusted_ca_certificate.test",
						tfjsonpath.New("validity_not_after"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cdn_trusted_ca_certificate.test",
						tfjsonpath.New("validity_not_before"), knownvalue.NotNull()),
				},
			},
			// Import test
			{
				ResourceName:            "gcore_cdn_trusted_ca_certificate.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccCDNTrustedCaCertificateImportStateIDFunc("gcore_cdn_trusted_ca_certificate.test"),
				ImportStateVerifyIgnore: []string{"ssl_certificate", "deleted"},
			},
		},
	})
}

// testAccCDNTrustedCaCertificateImportStateIDFunc returns the ID as a string for import.
func testAccCDNTrustedCaCertificateImportStateIDFunc(resourceName string) func(*terraform.State) (string, error) {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		id := rs.Primary.Attributes["id"]
		if id == "" {
			return "", fmt.Errorf("id not set in state for %s", resourceName)
		}
		return id, nil
	}
}

func testAccCheckCDNTrustedCaCertificateDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cdn_trusted_ca_certificate" {
			continue
		}

		id, err := strconv.ParseInt(rs.Primary.Attributes["id"], 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing id: %w", err)
		}

		_, err = client.CDN.TrustedCaCertificates.Get(context.Background(), id)
		if err == nil {
			return fmt.Errorf("CDN trusted CA certificate %d still exists", id)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking CDN trusted CA certificate deletion: %w", err)
		}
	}

	return nil
}

// testCACertPEM is a self-signed CA certificate for testing.
// Generated with: openssl req -x509 -newkey rsa:2048 -keyout /dev/null -out /dev/stdout -days 3650 -nodes -subj "/CN=tf-test-ca"
const testCACertPEM = "-----BEGIN CERTIFICATE-----\n" +
	"MIIDCjCCAfKgAwIBAgITL/pYrT1kgOYtbe40xPcRZbmqNTANBgkqhkiG9w0BAQsF\n" +
	"ADAVMRMwEQYDVQQDDAp0Zi10ZXN0LWNhMB4XDTI2MDMxNzE5MTYxOFoXDTM2MDMx\n" +
	"NDE5MTYxOFowFTETMBEGA1UEAwwKdGYtdGVzdC1jYTCCASIwDQYJKoZIhvcNAQEB\n" +
	"BQADggEPADCCAQoCggEBALJWJYInR6FCkr1W3fZ06P6GNBq1HSojChWa5Rpp4Fh8\n" +
	"z+eZKsKJkTFoYSwvfUUN0m8HxywUYeBSYG2fLaAko0ebPe4cOHt7VStvpOnJ+fYS\n" +
	"KO01C7P0/WzN0ROehOvDW3CRVTByCXA8deqWzV9otTL07LLdU9FBYl0Hyz/eNl40\n" +
	"2QtSHeBsQ/vOXYr097v5ReFvyVttCXYez7znZOqQuJycy/s8rCuGyzGyGOI8iF0m\n" +
	"yVM3lFFcEEGItGkKPPzF19csTUAs9UvRrMVJBM0Lhv9VcOfeVPhqwCIkiY9CEz44\n" +
	"FjE7yFpErqPLobBW3ZOumtD/st5P+tOEaUpTYRKzQU0CAwEAAaNTMFEwHQYDVR0O\n" +
	"BBYEFE6AOvgVff0u3N3FmY3u71wVphNmMB8GA1UdIwQYMBaAFE6AOvgVff0u3N3F\n" +
	"mY3u71wVphNmMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEBAEx5\n" +
	"LyMm1UJDtIpTNjRdpxftOAXXeucX6BZzZxx79ohMjuqWhnk6hKhjfxUmz4XnZKTN\n" +
	"rxh/m6SaC4ABOM8j1aNhySD4nWwRkg9JYq7Fu6Nh6ePs3dVUmJWNfOtCcCoz1JGU\n" +
	"xOs5QTjkreoat6iNvVjJ8okuExhaOHElbZtI0Z+d0LbOdziStC/dMSf31I6thuzC\n" +
	"sGa5c3sMTa28Kix0fJOKykNwZHlZ98It8Rp2QuhsdGSUkF3O6P/KS6dvkk5aTiqX\n" +
	"k9w55B/5O8Q6b2Xo+T0tgN+paHiq1F0Q3YiE+Ptlb6O9Auh6e868q/+uBYhEbKkf\n" +
	"sXBOQryEQ4bdqc8CRGw=\n" +
	"-----END CERTIFICATE-----\n"

func testAccCDNTrustedCaCertificateConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_trusted_ca_certificate" "test" {
  name            = %[1]q
  ssl_certificate = %[2]q
}`, name, testCACertPEM)
}
