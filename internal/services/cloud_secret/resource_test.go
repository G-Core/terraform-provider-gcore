package cloud_secret_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cloud"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

const (
	testCertificate = `-----BEGIN CERTIFICATE-----
MIIDCTCCAfGgAwIBAgIUXZ1cl3k+XxbudA4k5NawJ6gTUOYwDQYJKoZIhvcNAQEL
BQAwFDESMBAGA1UEAwwJbG9jYWxob3N0MB4XDTI2MDIyNDIxNTkzNFoXDTI3MDIy
NDIxNTkzNFowFDESMBAGA1UEAwwJbG9jYWxob3N0MIIBIjANBgkqhkiG9w0BAQEF
AAOCAQ8AMIIBCgKCAQEAul7ibwUKpn/ORrX4AZLBNCb4juoSoJ58uplJPti6dMR7
5/jMPJd4o1SDL1MXVeaVNH7ipRjrebjyLRQQNQ2ymLooixPFcKJnYvxrSv/XQLNZ
jAy9uFS9pbDsETnr19e3s4YFxTO8M5TJ4yNI+W2vakdlLe91KL5KaQ0vXlo3JIHa
zsCCE6QwNJBTL64H1rIeM0jKJNgoSRAEeejusZkWB1j0Dzx2tdt6Z/+kSI0Y1nsA
IERFqSIbIvtv9qqhw8s94uzhVvwplaWJVfxZjQb5X3LkfDFenH3zMDcOIjcO62mx
UI3/mm+EoCPfw2g/4hplLBaUFCrvKp5ATNHzgLNjuwIDAQABo1MwUTAdBgNVHQ4E
FgQUavcynOpeZEBfnQKMm9w8Mwsu7dUwHwYDVR0jBBgwFoAUavcynOpeZEBfnQKM
m9w8Mwsu7dUwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEAKh8Y
Qp97288rZlcb9vycTuxuejujhlBss3H7IX/qvAuW9FNUfi/K8ZTIwEmNxa7poXwG
Ybcr+RjmCAen4tgzD1awDRZl8lJdYP+VZGZDe4ZhjFLDUS5O6BqXeTyqOUJUIugk
hL/VKbi6dSwv+VLQzr9RSP63CL/BUOAr/lWeixOHfgsuyim3v1GzuNm98iEfgMaE
NoTvUi41KDby9FWdIl/NKh+BQRWX+XO7LhHlPfjCAqAQiCWggp6fNkwlnicoL9Bj
P3dqnCLmFKq/asMEixZ9Be3P4g+Gxv+evrqodfcTHIzXQKx6oCmioy2Ibh7thFqd
T30UxTjWgQ5M6GFs1w==
-----END CERTIFICATE-----`

	testPrivateKey = `-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQC6XuJvBQqmf85G
tfgBksE0JviO6hKgnny6mUk+2Lp0xHvn+Mw8l3ijVIMvUxdV5pU0fuKlGOt5uPIt
FBA1DbKYuiiLE8Vwomdi/GtK/9dAs1mMDL24VL2lsOwROevX17ezhgXFM7wzlMnj
I0j5ba9qR2Ut73UovkppDS9eWjckgdrOwIITpDA0kFMvrgfWsh4zSMok2ChJEAR5
6O6xmRYHWPQPPHa123pn/6RIjRjWewAgREWpIhsi+2/2qqHDyz3i7OFW/CmVpYlV
/FmNBvlfcuR8MV6cffMwNw4iNw7rabFQjf+ab4SgI9/DaD/iGmUsFpQUKu8qnkBM
0fOAs2O7AgMBAAECggEAV9W2NAqVPWlIp2zFiMBIjDK9vGU3CDoJoMvziEoOfk/H
8ckBQKFGuvtupFQt0E9PDKGsYZEAEasHgBVPmiFthatexkU4LWBtB2rdikhPg2/D
iUzL8V8Gzls2ttuselpxeot0lr9OOKUsDP+pOdzm6ljhp6eOePhOC3qqU3aDPJGj
SyZbWL1bE1PXEg+QKidb/JxhBCvtrVNkBYYZx55NKgA+8rPIFfeXSEcSaGN4U5mh
gTIqbhx2UKHPG7omcuAfssWGsYSPNOm0GFvinD83nT42Vlq/2At55nezfCPboZQU
2GPUBW2VIZOhAT7fI8Lp6k3CXpMR4jjvOY03oK0ZAQKBgQD2q2eexvUPwzQGuNC3
FFkEK43J5c1XYC8pbLG5Qet13OBJ5ZYAgDp5rXQtdwSJo8zqnZKQ91WH/kMYXKys
MnFlSnsWHk/zQacIO5JFwF+uIrLnKgMHixVL6BVTLul9DdfwFUZ9isoE6EE7yZg/
j9WaLO1Rq6nGm4hbFxEZTC1fOwKBgQDBa5UbPH/ilO0dz3gfPTNp6bDTd+q3gQaS
nv9ExZokR02Ibgn/V1LNBwrqYutKq3+iukUQ8vRDiKR8UouBRlZaGcXdwwsbYyVX
1dh6OZAYnfKxQTsT2g6aDdw1Zxh+We3VhHFRqethihP3sD8tUCvdT+AG4yi6HX8H
R2/VPcrFgQKBgQDzosr/Ja5Jelm5xfPI0N03ZDlw0HzoL4WFmelUfQqvaJFUC+MD
7aNUKYGVonel51bv6OEqNFGTuAzXVDns/wnHrTAz4Y7ASvlLBWPtZxxaJ8Wi03kY
i0Rmq/3cInrWXMULSkhMmbf97tT305+AMHYfCP8CataO53Jf3kGyRe6OnwKBgQCn
2yCEiYWGcq4w/7sEiU1ULh1p+Bi3df7pQZjQ6xfxQfv0WWLNuM+/5MvBS2Vc4Oac
p0CHDAGVlkEBL3WoFA2eld1Urg62jt16k7gRomD+LBzRXYXSnZuscDjaE4V7Kbow
YYciUu9WL8lSXB8HyRq4LriB4aOXmT+DZqiUC9MsgQKBgQCYodb7UvuMfZ7+Fi/p
V6FyQL/MyJIyPw5HistV0Y4VZ6ry0ogfBsGCHtvHfoj4ODqmzky/m8W85exfgIWX
peJ6HQI+06E+BCe4IDgZYLuqho8IJ24A+yAIqcPEWa+DEEubVZLQBPFzaJA/YoJf
1EGWZTD5RZ8QpF0/ykUdbt5cLg==
-----END PRIVATE KEY-----`
)

func TestAccCloudSecret_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSecretConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_secret.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_secret.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_secret.test",
						tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_secret.test",
						tfjsonpath.New("secret_type"), knownvalue.NotNull()),
				},
			},
		},
	})
}

// TestAccCloudSecret_importBlockAdoptionApplies verifies that a real import
// block applies the adoption without replacing the secret.
func TestAccCloudSecret_importBlockAdoptionApplies(t *testing.T) {
	rName := acctest.RandomName()
	var secretID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSecretConfig(rName),
				Check: func(s *terraform.State) error {
					rs, ok := s.RootModule().Resources["gcore_cloud_secret.test"]
					if !ok {
						return fmt.Errorf("gcore_cloud_secret.test not found in state")
					}
					secretID = rs.Primary.ID
					return nil
				},
			},
			{
				Config: testAccCloudSecretConfigForgotten(),
			},
			{
				ConfigFile: testAccCloudSecretImportConfigFile(t, rName, &secretID),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_secret.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_secret.test",
						tfjsonpath.New("payload_wo_version"), knownvalue.Int64Exact(1)),
				},
				Check: func(s *terraform.State) error {
					rs, ok := s.RootModule().Resources["gcore_cloud_secret.test"]
					if !ok {
						return fmt.Errorf("gcore_cloud_secret.test not found in state")
					}
					if rs.Primary.ID != secretID {
						return fmt.Errorf("secret id changed from %q to %q during import adoption", secretID, rs.Primary.ID)
					}
					return nil
				},
			},
			{
				Config: testAccCloudSecretConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: func(s *terraform.State) error {
					rs, ok := s.RootModule().Resources["gcore_cloud_secret.test"]
					if !ok {
						return fmt.Errorf("gcore_cloud_secret.test not found in state")
					}
					if rs.Primary.ID != secretID {
						return fmt.Errorf("secret id changed from %q to %q after import adoption", secretID, rs.Primary.ID)
					}
					return nil
				},
			},
		},
	})
}

// TestAccCloudSecret_payloadChangeForcesReplacement: a genuine
// payload_wo_version bump must still force replacement (the adoption
// window must not weaken this for a managed secret), and the recreated
// secret must get a new ID.
func TestAccCloudSecret_payloadChangeForcesReplacement(t *testing.T) {
	rName := acctest.RandomName()
	var firstID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSecretConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_secret.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
				},
				Check: func(s *terraform.State) error {
					rs, ok := s.RootModule().Resources["gcore_cloud_secret.test"]
					if !ok {
						return fmt.Errorf("gcore_cloud_secret.test not found in state")
					}
					firstID = rs.Primary.ID
					return nil
				},
			},
			{
				Config: testAccCloudSecretConfigVersion(rName, 2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_secret.test", plancheck.ResourceActionReplace),
					},
				},
				Check: func(s *terraform.State) error {
					rs, ok := s.RootModule().Resources["gcore_cloud_secret.test"]
					if !ok {
						return fmt.Errorf("gcore_cloud_secret.test not found in state")
					}
					if rs.Primary.ID == firstID {
						return fmt.Errorf("expected a new secret id after payload_wo_version bump, got the same id %q", firstID)
					}
					return nil
				},
			},
		},
	})
}

// TestAccCloudSecret_importCommandAdoptionApplies covers the half of the
// window the import-block test cannot reach: the apply that actually
// adopts. Command-mode import with ImportStatePersist puts the imported
// state into the working state, so the following steps really plan and
// apply against it, and the same config plans nothing afterward.
func TestAccCloudSecret_importCommandAdoptionApplies(t *testing.T) {
	rName := acctest.RandomName()
	var secretID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSecretConfig(rName),
				Check: func(s *terraform.State) error {
					rs, ok := s.RootModule().Resources["gcore_cloud_secret.test"]
					if !ok {
						return fmt.Errorf("gcore_cloud_secret.test not found in state")
					}
					secretID = rs.Primary.ID
					return nil
				},
			},
			{
				// Forget the secret without destroying it, so the import
				// below targets an address absent from state.
				Config: testAccCloudSecretConfigForgotten(),
			},
			{
				Config:             testAccCloudSecretConfig(rName),
				ResourceName:       "gcore_cloud_secret.test",
				ImportState:        true,
				ImportStateKind:    resource.ImportCommandWithID,
				ImportStatePersist: true,
				ImportStateIdFunc: func(*terraform.State) (string, error) {
					return fmt.Sprintf("%s/%s/%s", acctest.ProjectID(), acctest.RegionID(), secretID), nil
				},
			},
			{
				// The adopting apply: payload/payload_wo_version land in
				// state without the secret being replaced, and the marker
				// is retired.
				Config: testAccCloudSecretConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_secret.test", plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_secret.test",
						tfjsonpath.New("payload_wo_version"), knownvalue.Int64Exact(1)),
				},
			},
			{
				// Adoption is stable: with the values in state, the same
				// config plans nothing at all.
				Config:   testAccCloudSecretConfig(rName),
				PlanOnly: true,
			},
		},
	})
}

func testAccCheckCloudSecretDestroy(s *terraform.State) error {
	return acctest.CheckResourceDestroyed(s, "gcore_cloud_secret", func(client *gcore.Client, id string) error {
		_, err := client.Cloud.Secrets.Get(context.Background(), id, cloud.SecretGetParams{})
		return err
	})
}

func testAccCloudSecretConfig(name string) string {
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
}`, acctest.ProjectID(), acctest.RegionID(), name, testCertificate, testPrivateKey)
}

func testAccCloudSecretImportConfigFile(t *testing.T, name string, secretID *string) config.TestStepConfigFunc {
	t.Helper()
	dir := t.TempDir()

	return func(req config.TestStepConfigRequest) string {
		file := filepath.Join(dir, fmt.Sprintf("step_%d.tf", req.StepNumber))
		importID := fmt.Sprintf("%s/%s/%s", acctest.ProjectID(), acctest.RegionID(), *secretID)
		content := fmt.Sprintf(`
import {
  to = gcore_cloud_secret.test
  id = %q
}

%s`, importID, testAccCloudSecretConfig(name))
		if err := os.WriteFile(file, []byte(content), 0o600); err != nil {
			t.Fatalf("writing %s: %s", file, err)
		}
		return file
	}
}

func testAccCloudSecretConfigVersion(name string, version int) string {
	return fmt.Sprintf(`
resource "gcore_cloud_secret" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
  payload_wo_version = %[6]d
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
}`, acctest.ProjectID(), acctest.RegionID(), name, testCertificate, testPrivateKey, version)
}

func testAccCloudSecretConfigForgotten() string {
	return `
removed {
  from = gcore_cloud_secret.test

  lifecycle {
    destroy = false
  }
}`
}

func testAccCloudSecretConfigNoPayload(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_secret" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
}`, acctest.ProjectID(), acctest.RegionID(), name)
}

func testAccCloudSecretConfigPayloadNoVersion(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_secret" "test" {
  project_id = %[1]s
  region_id  = %[2]s
  name       = %[3]q
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
}`, acctest.ProjectID(), acctest.RegionID(), name, testCertificate, testPrivateKey)
}

// TestAccCloudSecret_importWithoutPayload covers importing a secret without
// supplying its material at all: the import apply must adopt nothing
// (payload/payload_wo_version stay null), a later plan with the same
// no-payload config must be empty, replacing the secret (e.g. renaming it)
// without payload must fail up front, and the secret must otherwise survive
// untouched.
func TestAccCloudSecret_importWithoutPayload(t *testing.T) {
	rName := acctest.RandomName()
	var secretID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSecretConfig(rName),
				Check: func(s *terraform.State) error {
					rs, ok := s.RootModule().Resources["gcore_cloud_secret.test"]
					if !ok {
						return fmt.Errorf("gcore_cloud_secret.test not found in state")
					}
					secretID = rs.Primary.ID
					return nil
				},
			},
			{
				// Forget the secret without destroying it, so the import
				// below targets an address absent from state.
				Config: testAccCloudSecretConfigForgotten(),
			},
			{
				Config:             testAccCloudSecretConfigNoPayload(rName),
				ResourceName:       "gcore_cloud_secret.test",
				ImportState:        true,
				ImportStateKind:    resource.ImportCommandWithID,
				ImportStatePersist: true,
				ImportStateIdFunc: func(*terraform.State) (string, error) {
					return fmt.Sprintf("%s/%s/%s", acctest.ProjectID(), acctest.RegionID(), secretID), nil
				},
			},
			{
				// Nothing was adopted, so the same no-payload config plans
				// nothing at all.
				Config:   testAccCloudSecretConfigNoPayload(rName),
				PlanOnly: true,
			},
			{
				// Renaming without payload forces replacement, which needs
				// the material this config still does not have.
				Config:      testAccCloudSecretConfigNoPayload(rName + "-renamed"),
				ExpectError: regexp.MustCompile(`Missing secret payload`),
			},
			{
				Config: testAccCloudSecretConfigNoPayload(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: func(s *terraform.State) error {
					rs, ok := s.RootModule().Resources["gcore_cloud_secret.test"]
					if !ok {
						return fmt.Errorf("gcore_cloud_secret.test not found in state")
					}
					if rs.Primary.ID != secretID {
						return fmt.Errorf("secret id changed from %q to %q; it should have survived untouched", secretID, rs.Primary.ID)
					}
					return nil
				},
			},
		},
	})
}

// TestAccCloudSecret_removePayloadInPlace covers removing payload and
// payload_wo_version from the config of a managed secret: the secret must
// survive with both attributes recorded as null in state, a repeat plan must
// be empty, and re-adding the material must replace the secret with a new
// id (no marker survives a removal).
func TestAccCloudSecret_removePayloadInPlace(t *testing.T) {
	rName := acctest.RandomName()
	var secretID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSecretConfig(rName),
				Check: func(s *terraform.State) error {
					rs, ok := s.RootModule().Resources["gcore_cloud_secret.test"]
					if !ok {
						return fmt.Errorf("gcore_cloud_secret.test not found in state")
					}
					secretID = rs.Primary.ID
					return nil
				},
			},
			{
				Config: testAccCloudSecretConfigNoPayload(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_secret.test", plancheck.ResourceActionUpdate),
					},
				},
				// carryComputedAttributesForward (plan_modifiers.go) keeps the
				// 7 other computed attributes at their prior known values on
				// this plan; without it they would show "(known after
				// apply)" here even though nothing about them changed, which
				// would break a downstream reference (a `count`, another
				// resource's RequiresReplace) evaluated at plan time. Only
				// the removed attributes are asserted below because the acc
				// test framework does not expose a "known, unchanged" plan
				// check for the others.
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_secret.test",
						tfjsonpath.New("payload"), knownvalue.Null()),
					statecheck.ExpectKnownValue("gcore_cloud_secret.test",
						tfjsonpath.New("payload_wo_version"), knownvalue.Null()),
				},
				Check: func(s *terraform.State) error {
					rs, ok := s.RootModule().Resources["gcore_cloud_secret.test"]
					if !ok {
						return fmt.Errorf("gcore_cloud_secret.test not found in state")
					}
					if rs.Primary.ID != secretID {
						return fmt.Errorf("secret id changed from %q to %q; removal must not replace the secret", secretID, rs.Primary.ID)
					}
					return nil
				},
			},
			{
				Config:   testAccCloudSecretConfigNoPayload(rName),
				PlanOnly: true,
			},
			{
				// No marker exists after a removal, so re-adding the
				// material is a real change and must replace the secret.
				Config: testAccCloudSecretConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_secret.test", plancheck.ResourceActionReplace),
					},
				},
				Check: func(s *terraform.State) error {
					rs, ok := s.RootModule().Resources["gcore_cloud_secret.test"]
					if !ok {
						return fmt.Errorf("gcore_cloud_secret.test not found in state")
					}
					if rs.Primary.ID == secretID {
						return fmt.Errorf("expected a new secret id after re-adding payload, got the same id %q", secretID)
					}
					return nil
				},
			},
		},
	})
}

// TestAccCloudSecret_importWithoutPayloadThenSetReplaces is the do-nothing-
// then-act path: import without material, apply the no-payload config
// exactly once with an empty plan (the refresh in that apply spends the
// marker), then set payload. If the marker had not been retired by the
// no-op apply, this would plan an Update instead of a Replace, which would
// mean the docs claim (one apply after import, then the window closes) is
// wrong - escalate that, do not add steps to make this test pass.
func TestAccCloudSecret_importWithoutPayloadThenSetReplaces(t *testing.T) {
	rName := acctest.RandomName()
	var secretID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSecretConfig(rName),
				Check: func(s *terraform.State) error {
					rs, ok := s.RootModule().Resources["gcore_cloud_secret.test"]
					if !ok {
						return fmt.Errorf("gcore_cloud_secret.test not found in state")
					}
					secretID = rs.Primary.ID
					return nil
				},
			},
			{
				Config: testAccCloudSecretConfigForgotten(),
			},
			{
				Config:             testAccCloudSecretConfigNoPayload(rName),
				ResourceName:       "gcore_cloud_secret.test",
				ImportState:        true,
				ImportStateKind:    resource.ImportCommandWithID,
				ImportStatePersist: true,
				ImportStateIdFunc: func(*terraform.State) (string, error) {
					return fmt.Sprintf("%s/%s/%s", acctest.ProjectID(), acctest.RegionID(), secretID), nil
				},
			},
			{
				// Exactly one no-op apply: its refresh is the one that
				// spends the marker.
				Config: testAccCloudSecretConfigNoPayload(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				Config: testAccCloudSecretConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cloud_secret.test", plancheck.ResourceActionReplace),
					},
				},
				Check: func(s *terraform.State) error {
					rs, ok := s.RootModule().Resources["gcore_cloud_secret.test"]
					if !ok {
						return fmt.Errorf("gcore_cloud_secret.test not found in state")
					}
					if rs.Primary.ID == secretID {
						return fmt.Errorf("expected a new secret id after setting payload, got the same id %q", secretID)
					}
					return nil
				},
			},
		},
	})
}

// TestAccCloudSecret_createWithoutPayloadFails: a brand-new secret cannot be
// created without its material.
func TestAccCloudSecret_createWithoutPayloadFails(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCloudSecretConfigNoPayload(rName),
				ExpectError: regexp.MustCompile(`Missing secret payload`),
			},
		},
	})
}

// TestAccCloudSecret_payloadWithoutVersionFails: payload and
// payload_wo_version must be configured together.
func TestAccCloudSecret_payloadWithoutVersionFails(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCloudSecretConfigPayloadNoVersion(rName),
				ExpectError: regexp.MustCompile(`must be set together`),
			},
		},
	})
}
