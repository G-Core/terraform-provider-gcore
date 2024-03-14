//go:build !cloud
// +build !cloud

package gcore

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOriginShielding(t *testing.T) {
	fullName := "gcore_cdn_originshielding.acctest"

	type Params struct {
		ShieldingPop string
	}

	update := Params{"1"}

	template := func(params *Params) string {
		return fmt.Sprintf(`
            resource "gcore_cdn_originshielding" "acctest" {
				resource_id = %s
				shielding_pop = %s
		`, GCORE_CDN_RESOURCE_ID, params.ShieldingPop)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVars(t, GCORE_USERNAME_VAR, GCORE_PASSWORD_VAR, GCORE_CDN_URL_VAR)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: template(&update),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "shielding_pop", update.ShieldingPop),
				),
			},
		},
	})
}
