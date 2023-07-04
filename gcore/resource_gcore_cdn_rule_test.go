//go:build !cloud
// +build !cloud

package gcore

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCDNRule(t *testing.T) {
	fullName := "gcore_cdn_rule.acctest"

	type Params struct {
		Name    string
		Pattern string
		RawPart string
	}

	create := Params{
		Name:    "All images",
		Pattern: "/folder/images/*.png",
	}
	update := Params{
		Name:    "All scripts",
		Pattern: "/folder/scripts/*.js",
		RawPart: `
  options {
    host_header {
      enabled = true
      value = "rule-host.com"
    }
  }
		`,
	}

	template := func(params *Params) string {
		return fmt.Sprintf(`
resource "gcore_cdn_rule" "acctest" {
  resource_id = %s
  name = "%s"
  rule = "%s"
  rule_type = 0
  weight = 0
  %s
}
		`, GCORE_CDN_RESOURCE_ID, params.Name, params.Pattern, params.RawPart)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVars(t, GCORE_USERNAME_VAR, GCORE_PASSWORD_VAR, GCORE_CDN_URL_VAR, GCORE_CDN_RESOURCE_ID_VAR)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: template(&create),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", create.Name),
					resource.TestCheckResourceAttr(fullName, "rule", create.Pattern),
				),
			},
			{
				Config: template(&update),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", update.Name),
					resource.TestCheckResourceAttr(fullName, "rule", update.Pattern),
					resource.TestCheckResourceAttr(fullName, "options.0.host_header.0.value", "rule-host.com"),
				),
			},
			{
				ResourceName:      "gcore_cdn_rule.acctest",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
