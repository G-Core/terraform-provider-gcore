package cloud_secret_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cloud"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
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

// Note: Import test is not implemented because the payload field is required
// and forces replacement, but it's not returned by the API during read/import.
// This is a known limitation of the secrets API.

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
