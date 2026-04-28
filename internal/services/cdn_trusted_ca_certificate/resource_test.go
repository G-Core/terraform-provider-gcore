package cdn_trusted_ca_certificate_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
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
				ImportStateIdFunc:       acctest.BuildImportID("gcore_cdn_trusted_ca_certificate.test", "id"),
				ImportStateVerifyIgnore: []string{"ssl_certificate", "deleted"},
			},
		},
	})
}

func TestAccCDNTrustedCaCertificate_update(t *testing.T) {
	t.Parallel()
	rName := acctest.RandomName()
	rNameUpdated := acctest.RandomName()
	resourceName := "gcore_cdn_trusted_ca_certificate.test"

	var certID string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNTrustedCaCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNTrustedCaCertificateConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
				},
				Check: func(s *terraform.State) error {
					rs := s.RootModule().Resources[resourceName]
					certID = rs.Primary.ID
					return nil
				},
			},
			{
				Config: testAccCDNTrustedCaCertificateConfig(rNameUpdated),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("name"), knownvalue.StringExact(rNameUpdated)),
				},
				Check: func(s *terraform.State) error {
					rs := s.RootModule().Resources[resourceName]
					if rs.Primary.ID != certID {
						return fmt.Errorf("resource was recreated: ID changed from %s to %s", certID, rs.Primary.ID)
					}
					return nil
				},
			},
		},
	})
}

func TestAccCDNTrustedCaCertificate_importDrift(t *testing.T) {
	t.Parallel()
	rName := acctest.RandomName()
	resourceName := "gcore_cdn_trusted_ca_certificate.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNTrustedCaCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNTrustedCaCertificateConfig(rName),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ssl_certificate", "deleted"},
				ImportStateIdFunc:       acctest.BuildImportID(resourceName, "id"),
			},
			{
				// After import, ssl_certificate is null in state because the API
				// doesn't return the PEM. Re-applying the same config should NOT
				// trigger replacement.
				Config: testAccCDNTrustedCaCertificateConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
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
