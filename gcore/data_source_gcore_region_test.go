//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/region/v1/regions"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRegionDataSource(t *testing.T) {
	cfg, err := createTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := CreateTestClient(cfg.Provider, regionPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	rs, err := regions.ListAll(client, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(rs) == 0 {
		t.Fatal("regions not found")
	}

	region := rs[0]

	fullName := "data.gcore_region.acctest"
	tpl := func(name string) string {
		return fmt.Sprintf(`
			data "gcore_region" "acctest" {
              name = "%s"
			}
		`, name)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tpl(region.DisplayName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", region.DisplayName),
					resource.TestCheckResourceAttr(fullName, "id", strconv.Itoa(region.ID)),
				),
			},
		},
	})
}
