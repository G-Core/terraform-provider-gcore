package gcore

import (
	"fmt"
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/faas/v1/faas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFaaSKeyDataSource(t *testing.T) {
	cfg, err := createTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := CreateTestClient(cfg.Provider, faasKeysPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	keyName := "test-key"
	if err := createTestKey(client, keyName); err != nil {
		t.Fatal(err)
	}
	defer faas.DeleteKey(client, keyName)

	key, err := faas.GetKey(client, keyName).Extract()
	if err != nil {
		t.Fatal(err)
	}

	fullName := "data.gcore_faas_key.acctest"
	tpl := func(n string) string {
		return fmt.Sprintf(`
			data "gcore_faas_key" "acctest" {
				%s
				%s
				name = "%s"
			}
		`, projectInfo(), regionInfo(), n)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tpl(keyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "id", key.Name),
					resource.TestCheckResourceAttr(fullName, "name", key.Name),
					resource.TestCheckResourceAttr(fullName, "description", key.Description),
				),
			},
		},
	})
}

func createTestKey(client *gcorecloud.ServiceClient, keyName string) error {
	opts := faas.CreateKeyOpts{
		Name:        keyName,
		Description: "test description",
	}

	_, err := faas.CreateKey(client, opts)
	if err != nil {
		return err
	}

	return nil
}
