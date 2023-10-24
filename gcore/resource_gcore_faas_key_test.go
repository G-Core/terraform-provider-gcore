package gcore

import (
	"fmt"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/faas/v1/faas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFaaSKey(t *testing.T) {
	type Params struct {
		Name        string
		Description string
	}

	create := Params{
		Name:        "test-key",
		Description: "description",
	}

	update := Params{
		Name:        "test-key",
		Description: "change in the description",
	}

	fullName := "gcore_faas_key.acctest"

	tpl := func(params *Params) string {
		return fmt.Sprintf(`
		resource "gcore_faas_key" "acctest" {
			name = "%s"
			description = "%s"
			%s
			%s
		}`, params.Name, params.Description, regionInfo(), projectInfo())
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      CheckDestroyFaaSKey,
		Steps: []resource.TestStep{
			{
				Config: tpl(&create),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", create.Name),
					resource.TestCheckResourceAttr(fullName, "description", create.Description),
				),
			},
			{
				Config: tpl(&update),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", update.Name),
					resource.TestCheckResourceAttr(fullName, "description", update.Description),
				),
			},
		},
	})
}

func CheckDestroyFaaSKey(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, faasKeysPoint, versionPointV1)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_faas_key" {
			continue
		}

		_, err := faas.GetKey(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("faas key still exists")
		}
	}

	return nil
}
