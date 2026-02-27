package cloud_ssh_key_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cloud"
	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/stainless-sdks/gcore-terraform/internal/acctest"
)

func TestAccCloudSSHKey_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSSHKeyConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_ssh_key.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cloud_ssh_key.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_ssh_key.test",
						tfjsonpath.New("fingerprint"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cloud_ssh_key.test",
						tfjsonpath.New("state"), knownvalue.StringExact("ACTIVE")),
					statecheck.ExpectKnownValue("gcore_cloud_ssh_key.test",
						tfjsonpath.New("shared_in_project"), knownvalue.Bool(true)),
				},
			},
		},
	})
}

func TestAccCloudSSHKey_update(t *testing.T) {
	rName := acctest.RandomName()

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSSHKeyConfigShared(rName, false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_ssh_key.test",
						tfjsonpath.New("shared_in_project"), knownvalue.Bool(false)),
					compareIDSame.AddStateValue(
						"gcore_cloud_ssh_key.test",
						tfjsonpath.New("id"),
					),
				},
			},
			{
				Config: testAccCloudSSHKeyConfigShared(rName, true),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_ssh_key.test",
						tfjsonpath.New("shared_in_project"), knownvalue.Bool(true)),
					compareIDSame.AddStateValue(
						"gcore_cloud_ssh_key.test",
						tfjsonpath.New("id"),
					),
				},
			},
		},
	})
}

func TestAccCloudSSHKey_import(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudSSHKeyConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cloud_ssh_key.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
				},
			},
			{
				ResourceName:      "gcore_cloud_ssh_key.test",
				ImportState:       true,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cloud_ssh_key.test", "project_id", "id"),
				ImportStateKind:   resource.ImportBlockWithID,
			},
		},
	})
}

func testAccCheckCloudSSHKeyDestroy(s *terraform.State) error {
	return acctest.CheckResourceDestroyed(s, "gcore_cloud_ssh_key", func(client *gcore.Client, id string) error {
		_, err := client.Cloud.SSHKeys.Get(context.Background(), id, cloud.SSHKeyGetParams{})
		return err
	})
}

func testAccCloudSSHKeyConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_cloud_ssh_key" "test" {
  project_id = %[1]s
  name       = %[2]q
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFo9Yv8MkEsmt8mUm4gMRpOuIiLPkJfSdmR3lKA3GsOG tf-test@example.com"
}`, acctest.ProjectID(), name)
}

func testAccCloudSSHKeyConfigShared(name string, shared bool) string {
	return fmt.Sprintf(`
resource "gcore_cloud_ssh_key" "test" {
  project_id        = %[1]s
  name              = %[2]q
  public_key        = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFo9Yv8MkEsmt8mUm4gMRpOuIiLPkJfSdmR3lKA3GsOG tf-test@example.com"
  shared_in_project = %[3]t
}`, acctest.ProjectID(), name, shared)
}
