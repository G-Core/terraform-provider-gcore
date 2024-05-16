//go:build !cloud
// +build !cloud

package gcore

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAppliedPreset(t *testing.T) {
	fullName := "gcore_cdn_applied_preset.acctest"

	type Params struct {
		PresetID string
		ObjectID string
	}

	create := Params{
		PresetID: GCORE_CDN_PRESET_ID,
		ObjectID: GCORE_CDN_PRESET_OBJECT_ID,
	}

	template := func(params *Params) string {
		return fmt.Sprintf(`
            resource "gcore_cdn_applied_preset" "acctest" {
				preset_id = %s
				object_id = %s
		`, params.PresetID, params.ObjectID)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVars(t, GCORE_USERNAME_VAR, GCORE_PASSWORD_VAR, GCORE_CDN_URL_VAR, GCORE_CDN_PRESET_ID_VAR, GCORE_CDN_PRESET_OBJECT_ID_VAR)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: template(&create),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "object_id", create.ObjectID),
					resource.TestCheckResourceAttr(fullName, "preset_id", create.PresetID),
				),
			},
		},
	})
}
