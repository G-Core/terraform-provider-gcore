package cdn_certificate_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
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

// TestAccCDNCertificate_custom uploads a certificate without setting automated.
// The API reports automated as false, and the plan must accept that instead of
// replacing the certificate on every apply. A replacement creates the new
// certificate under the same name, which the API rejects, or deletes one that
// a CDN resource still uses.
func TestAccCDNCertificate_custom(t *testing.T) {
	rName := acctest.RandomName()
	cname := rName + ".example.com"
	cert, key := testAccCDNCertificateSelfSigned(t, cname)

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNCertificateCustomConfig(rName, cert, key),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_certificate.test",
						tfjsonpath.New("automated"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("gcore_cdn_certificate.test",
						tfjsonpath.New("cert_subject_cn"), knownvalue.StringExact(cname)),
					compareIDSame.AddStateValue("gcore_cdn_certificate.test", tfjsonpath.New("id")),
				},
			},
			{
				Config: testAccCDNCertificateCustomConfig(rName, cert, key) +
					testAccCDNCertificateAttachedResourceConfig(rName, cname),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cdn_certificate.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					compareIDSame.AddStateValue("gcore_cdn_certificate.test", tfjsonpath.New("id")),
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

func testAccCDNCertificateCustomConfig(name, cert, key string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_certificate" "test" {
  name                       = %[1]q
  ssl_certificate_wo         = %[2]q
  ssl_private_key_wo         = %[3]q
  ssl_certificate_wo_version = 1
}`, name, cert, key)
}

func testAccCDNCertificateAttachedResourceConfig(name, cname string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_origin_group" "test" {
  name    = %[1]q
  sources = [{ source = "example.com" }]
}

resource "gcore_cdn_resource" "test" {
  cname        = %[2]q
  origin_group = gcore_cdn_origin_group.test.id
  ssl_enabled  = true
  ssl_data     = gcore_cdn_certificate.test.ssl_id
}`, name, cname)
}

// testAccCDNCertificateSelfSigned returns a PEM certificate and key for cn.
func testAccCDNCertificateSelfSigned(t *testing.T, cn string) (string, string) {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generating key: %v", err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject:      pkix.Name{CommonName: cn},
		DNSNames:     []string{cn},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(30 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("creating certificate: %v", err)
	}
	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		t.Fatalf("encoding key: %v", err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})),
		string(pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER}))
}
