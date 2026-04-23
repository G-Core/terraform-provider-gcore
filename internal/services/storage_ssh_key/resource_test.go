package storage_ssh_key_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

const testAccStorageSSHKeyPublicKey = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFo9Yv8MkEsmt8mUm4gMRpOuIiLPkJfSdmR3lKA3GsOG tf-test@example.com"

func TestAccStorageSSHKey_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageSSHKeyConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_storage_ssh_key.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_storage_ssh_key.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_storage_ssh_key.test",
						tfjsonpath.New("created_at"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_storage_ssh_key.test",
						tfjsonpath.New("public_key"), knownvalue.StringExact(testAccStorageSSHKeyPublicKey)),
				},
			},
			{
				ResourceName:    "gcore_storage_ssh_key.test",
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithID,
			},
		},
	})
}

func testAccCheckStorageSSHKeyDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_storage_ssh_key" {
			continue
		}

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse id %q: %w", rs.Primary.ID, err)
		}

		_, err = client.Storage.SSHKeys.Get(context.Background(), id)
		if err == nil {
			return fmt.Errorf("storage ssh key %d still exists", id)
		}
		if !acctest.IsNotFoundError(err) {
			return fmt.Errorf("error checking ssh key deletion: %w", err)
		}
	}

	return nil
}

func testAccStorageSSHKeyConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_storage_ssh_key" "test" {
  name       = %[1]q
  public_key = %[2]q
}`, name, testAccStorageSSHKeyPublicKey)
}
