package storage_ssh_key_test

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

func TestAccStorageSSHKeyDataSource_byID(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageSSHKeyDataSourceConfigByID(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_storage_ssh_key.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("data.gcore_storage_ssh_key.test",
						tfjsonpath.New("public_key"), knownvalue.StringExact(testAccStorageSSHKeyPublicKey)),
					statecheck.CompareValuePairs(
						"gcore_storage_ssh_key.test", tfjsonpath.New("id"),
						"data.gcore_storage_ssh_key.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func TestAccStorageSSHKeyDataSource_byName(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageSSHKeyDataSourceConfigByName(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_storage_ssh_key.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.CompareValuePairs(
						"gcore_storage_ssh_key.test", tfjsonpath.New("id"),
						"data.gcore_storage_ssh_key.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func TestAccStorageSSHKeysDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageSSHKeysDataSourceConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_storage_ssh_keys.test",
						tfjsonpath.New("items"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue("data.gcore_storage_ssh_keys.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(rName)),
					statecheck.CompareValuePairs(
						"gcore_storage_ssh_key.test", tfjsonpath.New("id"),
						"data.gcore_storage_ssh_keys.test", tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("id"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func testAccStorageSSHKeyDataSourceConfigByID(name string) string {
	return fmt.Sprintf(`
resource "gcore_storage_ssh_key" "test" {
  name       = %[1]q
  public_key = %[2]q
}

data "gcore_storage_ssh_key" "test" {
  ssh_key_id = gcore_storage_ssh_key.test.id
}`, name, testAccStorageSSHKeyPublicKey)
}

func testAccStorageSSHKeyDataSourceConfigByName(name string) string {
	return fmt.Sprintf(`
resource "gcore_storage_ssh_key" "test" {
  name       = %[1]q
  public_key = %[2]q
}

data "gcore_storage_ssh_key" "test" {
  find_one_by = {
    name = gcore_storage_ssh_key.test.name
  }
}`, name, testAccStorageSSHKeyPublicKey)
}

func testAccStorageSSHKeysDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_storage_ssh_key" "test" {
  name       = %[1]q
  public_key = %[2]q
}

data "gcore_storage_ssh_keys" "test" {
  name = gcore_storage_ssh_key.test.name
}`, name, testAccStorageSSHKeyPublicKey)
}
