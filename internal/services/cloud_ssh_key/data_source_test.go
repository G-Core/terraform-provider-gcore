package cloud_ssh_key_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/stainless-sdks/gcore-terraform/internal/acctest"
)

func TestAccCloudSSHKeyDataSource_byID(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSSHKeyDataSourceConfigByID(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_ssh_key.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.CompareValuePairs(
						"gcore_cloud_ssh_key.test", tfjsonpath.New("id"),
						"data.gcore_cloud_ssh_key.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
					statecheck.CompareValuePairs(
						"gcore_cloud_ssh_key.test", tfjsonpath.New("fingerprint"),
						"data.gcore_cloud_ssh_key.test", tfjsonpath.New("fingerprint"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func TestAccCloudSSHKeyDataSource_byName(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSSHKeyDataSourceConfigByName(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cloud_ssh_key.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.CompareValuePairs(
						"gcore_cloud_ssh_key.test", tfjsonpath.New("id"),
						"data.gcore_cloud_ssh_key.test", tfjsonpath.New("id"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func testAccCloudSSHKeyDataSourceConfigByID(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_ssh_key" "test" {
  project_id = %[1]s
  name       = %[2]q
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFo9Yv8MkEsmt8mUm4gMRpOuIiLPkJfSdmR3lKA3GsOG tf-test@example.com"
}

data "gcore_cloud_ssh_key" "test" {
  project_id = %[1]s
  ssh_key_id = gcore_cloud_ssh_key.test.id
}`, acctest.ProjectID(), name)
}

func testAccCloudSSHKeyDataSourceConfigByName(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_ssh_key" "test" {
  project_id = %[1]s
  name       = %[2]q
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFo9Yv8MkEsmt8mUm4gMRpOuIiLPkJfSdmR3lKA3GsOG tf-test@example.com"
}

data "gcore_cloud_ssh_key" "test" {
  project_id = %[1]s
  find_one_by = {
    name = gcore_cloud_ssh_key.test.name
  }
}`, acctest.ProjectID(), name)
}
