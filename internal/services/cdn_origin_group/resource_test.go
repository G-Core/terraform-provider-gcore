package cdn_origin_group_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/G-Core/gcore-go"
	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func TestAccCDNOriginGroup_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNOriginGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNOriginGroupConfigBasic(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_origin_group.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cdn_origin_group.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cdn_origin_group.test",
						tfjsonpath.New("has_related_resources"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue("gcore_cdn_origin_group.test",
						tfjsonpath.New("auth_type"), knownvalue.StringExact("none")),
					statecheck.ExpectKnownValue("gcore_cdn_origin_group.test",
						tfjsonpath.New("use_next"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("gcore_cdn_origin_group.test",
						tfjsonpath.New("proxy_next_upstream"), knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("error"),
							knownvalue.StringExact("timeout"),
						})),
					statecheck.ExpectKnownValue("gcore_cdn_origin_group.test",
						tfjsonpath.New("sources"), knownvalue.SetSizeExact(1)),
				},
			},
		},
	})
}

func TestAccCDNOriginGroup_update(t *testing.T) {
	rName := acctest.RandomName()

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNOriginGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNOriginGroupConfigWithSource(rName, "example.com"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_origin_group.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cdn_origin_group.test",
						tfjsonpath.New("sources"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"source":  knownvalue.StringExact("example.com"),
								"enabled": knownvalue.Bool(true),
								"backup":  knownvalue.Bool(false),
							}),
						})),
					compareIDSame.AddStateValue(
						"gcore_cdn_origin_group.test",
						tfjsonpath.New("id"),
					),
				},
			},
			{
				Config: testAccCDNOriginGroupConfigWithSource(rName, "example.org"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_origin_group.test",
						tfjsonpath.New("sources"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"source":  knownvalue.StringExact("example.org"),
								"enabled": knownvalue.Bool(true),
								"backup":  knownvalue.Bool(false),
							}),
						})),
					// ID should not change — in-place update
					compareIDSame.AddStateValue(
						"gcore_cdn_origin_group.test",
						tfjsonpath.New("id"),
					),
				},
			},
		},
	})
}

func TestAccCDNOriginGroup_import(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNOriginGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNOriginGroupConfigBasic(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_origin_group.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
				},
			},
			{
				ResourceName:      "gcore_cdn_origin_group.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cdn_origin_group.test", "id"),
			},
		},
	})
}

func testAccCheckCDNOriginGroupDestroy(s *terraform.State) error {
	return acctest.CheckResourceDestroyed(s, "gcore_cdn_origin_group", func(client *gcore.Client, id string) error {
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing id: %w", err)
		}
		_, err = client.CDN.OriginGroups.Get(context.Background(), idInt)
		return err
	})
}

func testAccCDNOriginGroupConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_origin_group" "test" {
  name                 = %[1]q
  auth_type            = "none"
  use_next             = true
  proxy_next_upstream  = ["error", "timeout"]
  sources = [
    {
      source = "example.com"
    }
  ]
}`, name)
}

func testAccCDNOriginGroupConfigWithSource(name string, source string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_origin_group" "test" {
  name = %[1]q
  sources = [
    {
      source = %[2]q
    }
  ]
}`, name, source)
}
