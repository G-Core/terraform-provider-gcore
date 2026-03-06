package cdn_certificate_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func TestAccCDNCertificate_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNCertificateAutomatedConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_certificate.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cdn_certificate.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cdn_certificate.test",
						tfjsonpath.New("ssl_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cdn_certificate.test",
						tfjsonpath.New("automated"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("gcore_cdn_certificate.test",
						tfjsonpath.New("validate_root_ca"), knownvalue.Bool(false)),
				},
			},
		},
	})
}

func TestAccCDNCertificate_update(t *testing.T) {
	rName := acctest.RandomName()
	rNameUpdated := acctest.RandomName()

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNCertificateAutomatedConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_certificate.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					compareIDSame.AddStateValue(
						"gcore_cdn_certificate.test",
						tfjsonpath.New("id"),
					),
				},
			},
			{
				Config: testAccCDNCertificateAutomatedConfig(rNameUpdated),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_certificate.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rNameUpdated)),
					// ID should not change — in-place update, not recreate
					compareIDSame.AddStateValue(
						"gcore_cdn_certificate.test",
						tfjsonpath.New("id"),
					),
				},
			},
		},
	})
}

func TestAccCDNCertificate_import(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNCertificateAutomatedConfig(rName),
			},
			{
				ResourceName:            "gcore_cdn_certificate.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccCDNCertificateImportStateIDFunc("gcore_cdn_certificate.test"),
				ImportStateVerifyIgnore: []string{"validate_root_ca", "deleted"},
			},
		},
	})
}

// testAccCDNCertificateImportStateIDFunc returns the ssl_id for import (which equals the id).
func testAccCDNCertificateImportStateIDFunc(resourceName string) func(*terraform.State) (string, error) {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		sslID := rs.Primary.Attributes["ssl_id"]
		if sslID == "" {
			return "", fmt.Errorf("ssl_id not set in state for %s", resourceName)
		}
		return sslID, nil
	}
}

func testAccCheckCDNCertificateDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_cdn_certificate" {
			continue
		}

		sslID, err := strconv.ParseInt(rs.Primary.Attributes["ssl_id"], 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing ssl_id: %w", err)
		}

		_, err = client.CDN.Certificates.Get(context.Background(), sslID)
		if err == nil {
			return fmt.Errorf("CDN certificate %d still exists", sslID)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking CDN certificate deletion: %w", err)
		}
	}

	return nil
}

func testAccCDNCertificateAutomatedConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_certificate" "test" {
  name      = %[1]q
  automated = true
}`, name)
}
