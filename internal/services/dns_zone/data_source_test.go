package dns_zone_test

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

func TestAccDNSZoneDataSource_basic(t *testing.T) {
	rName := acctest.RandomName() + ".com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDNSZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSZoneDataSourceConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_dns_zone.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.CompareValuePairs(
						"gcore_dns_zone.test", tfjsonpath.New("name"),
						"data.gcore_dns_zone.test", tfjsonpath.New("name"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func testAccDNSZoneDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_dns_zone" "test" {
  name = %[1]q
}

data "gcore_dns_zone" "test" {
  name = gcore_dns_zone.test.name
}`, name)
}
