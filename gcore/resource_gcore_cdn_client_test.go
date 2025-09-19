package gcore

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccGcoreCDNClientConfig_basic(t *testing.T) {
	resourceName := "gcore_cdn_client_config.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVars(t, GCORE_PERMANENT_TOKEN_VAR, GCORE_CDN_URL_VAR)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
resource "gcore_cdn_client_config" "test" {
utilization_level = 70
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "utilization_level", "70"),
				),
			},
			{
				Config: `
resource "gcore_cdn_client_config" "test" {
utilization_level = 90
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "utilization_level", "90"),
				),
			},
		},
	})
}
