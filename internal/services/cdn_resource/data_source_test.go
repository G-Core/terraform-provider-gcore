package cdn_resource_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func TestAccCDNResourceDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()
	cname := rName + ".example.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNResourceDataSourceConfig(rName, cname),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_cdn_resource.test",
						tfjsonpath.New("cname"), knownvalue.StringExact(cname)),
					statecheck.CompareValuePairs(
						"gcore_cdn_resource.test", tfjsonpath.New("id"),
						"data.gcore_cdn_resource.test", tfjsonpath.New("resource_id"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func testAccCDNResourceDataSourceConfig(name, cname string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_origin_group" "test" {
  name = %[1]q
  sources = [
    {
      source = "example.com"
    }
  ]
}

resource "gcore_cdn_resource" "test" {
  cname        = %[2]q
  origin_group = gcore_cdn_origin_group.test.id
}

data "gcore_cdn_resource" "test" {
  resource_id = gcore_cdn_resource.test.id
}`, name, cname)
}
