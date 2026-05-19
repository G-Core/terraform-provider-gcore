package storage_sftp_test

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

func TestAccStorageSftpDataSource_byID(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageSftpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageSftpDataSourceConfigByID(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_storage_sftp.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("data.gcore_storage_sftp.test",
						tfjsonpath.New("location_name"), knownvalue.StringExact(testAccStorageSftpLocation)),
					statecheck.ExpectKnownValue("data.gcore_storage_sftp.test",
						tfjsonpath.New("full_name"), knownvalue.NotNull()),
					statecheck.CompareValuePairs(
						"gcore_storage_sftp.test", tfjsonpath.New("id"),
						"data.gcore_storage_sftp.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func TestAccStorageSftpDataSource_byName(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageSftpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageSftpDataSourceConfigByName(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_storage_sftp.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.CompareValuePairs(
						"gcore_storage_sftp.test", tfjsonpath.New("id"),
						"data.gcore_storage_sftp.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func TestAccStorageSftpsDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageSftpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageSftpsDataSourceConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_storage_sftps.test",
						tfjsonpath.New("items"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue("data.gcore_storage_sftps.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("data.gcore_storage_sftps.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("full_name"), knownvalue.NotNull()),
					statecheck.CompareValuePairs(
						"gcore_storage_sftp.test", tfjsonpath.New("id"),
						"data.gcore_storage_sftps.test", tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("id"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func testAccStorageSftpDataSourceConfigByID(name string) string {
	return fmt.Sprintf(`
resource "gcore_storage_sftp" "test" {
  name          = %[1]q
  location_name = %[2]q
  password_mode = "auto"
}

data "gcore_storage_sftp" "test" {
  storage_id = gcore_storage_sftp.test.id
}`, name, testAccStorageSftpLocation)
}

func testAccStorageSftpDataSourceConfigByName(name string) string {
	return fmt.Sprintf(`
resource "gcore_storage_sftp" "test" {
  name          = %[1]q
  location_name = %[2]q
  password_mode = "auto"
}

data "gcore_storage_sftp" "test" {
  find_one_by = {
    name = gcore_storage_sftp.test.name
  }
}`, name, testAccStorageSftpLocation)
}

func testAccStorageSftpsDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_storage_sftp" "test" {
  name          = %[1]q
  location_name = %[2]q
  password_mode = "auto"
}

data "gcore_storage_sftps" "test" {
  name = gcore_storage_sftp.test.name
}`, name, testAccStorageSftpLocation)
}
