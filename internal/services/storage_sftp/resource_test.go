package storage_sftp_test

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

const testAccStorageSftpLocation = "lux"

func TestAccStorageSftp_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageSftpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageSftpConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("address"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("full_name"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("provisioning_status"), knownvalue.StringExact("active")),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("location_name"), knownvalue.StringExact(testAccStorageSftpLocation)),
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("password"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccStorageSftpConfig(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("password"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccStorageSftp_update(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageSftpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageSftpConfigWithHTTPDisabled(rName, false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("is_http_disabled"), knownvalue.Bool(false)),
				},
			},
			{
				Config: testAccStorageSftpConfigWithHTTPDisabled(rName, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_storage_sftp.test",
							plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("is_http_disabled"), knownvalue.Bool(true)),
				},
			},
			{
				Config: testAccStorageSftpConfigWithHTTPDisabled(rName, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_storage_sftp.test",
						tfjsonpath.New("password"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func testAccCheckStorageSftpDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_storage_sftp" {
			continue
		}

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse id %q: %w", rs.Primary.ID, err)
		}

		_, err = client.Storage.SftpStorages.Get(context.Background(), id)
		if err == nil {
			return fmt.Errorf("storage sftp %d still exists", id)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking sftp storage deletion: %w", err)
		}
	}

	return nil
}

func testAccStorageSftpConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_storage_sftp" "test" {
  name          = %[1]q
  location_name = %[2]q
  password_mode = "auto"
}`, name, testAccStorageSftpLocation)
}

func testAccStorageSftpConfigWithHTTPDisabled(name string, httpDisabled bool) string {
	return fmt.Sprintf(`
resource "gcore_storage_sftp" "test" {
  name             = %[1]q
  location_name    = %[2]q
  password_mode    = "auto"
  is_http_disabled = %[3]t
}`, name, testAccStorageSftpLocation, httpDisabled)
}
