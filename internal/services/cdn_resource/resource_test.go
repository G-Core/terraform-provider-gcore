package cdn_resource_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/G-Core/gcore-go"
	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

func TestAccCDNResource_basic(t *testing.T) {
	rName := acctest.RandomName()
	cname := rName + ".example.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNResourceConfigBasic(rName, cname),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_resource.test",
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("gcore_cdn_resource.test",
						tfjsonpath.New("cname"), knownvalue.StringExact(cname)),
					statecheck.ExpectKnownValue("gcore_cdn_resource.test",
						tfjsonpath.New("active"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("gcore_cdn_resource.test",
						tfjsonpath.New("origin_protocol"), knownvalue.StringExact("HTTP")),
				},
			},
		},
	})
}

func TestAccCDNResource_update(t *testing.T) {
	rName := acctest.RandomName()
	cname := rName + ".example.com"

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNResourceConfigWithDescription(rName, cname, "initial description"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_resource.test",
						tfjsonpath.New("cname"), knownvalue.StringExact(cname)),
					statecheck.ExpectKnownValue("gcore_cdn_resource.test",
						tfjsonpath.New("description"), knownvalue.StringExact("initial description")),
					compareIDSame.AddStateValue(
						"gcore_cdn_resource.test",
						tfjsonpath.New("id"),
					),
				},
			},
			{
				Config: testAccCDNResourceConfigWithDescription(rName, cname, "updated description"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_resource.test",
						tfjsonpath.New("description"), knownvalue.StringExact("updated description")),
					// ID should not change — in-place update
					compareIDSame.AddStateValue(
						"gcore_cdn_resource.test",
						tfjsonpath.New("id"),
					),
				},
			},
		},
	})
}

// TestAccCDNResource_switchOriginGroup moves a CDN resource to another origin
// group in place. origin_group_name follows the group, so the plan must not
// carry the old name over.
func TestAccCDNResource_switchOriginGroup(t *testing.T) {
	rName := acctest.RandomName()
	cname := rName + ".example.com"

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNResourceConfigTwoGroups(rName, cname, "a"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_resource.test",
						tfjsonpath.New("origin_group_name"), knownvalue.StringExact(rName+"-a")),
					compareIDSame.AddStateValue("gcore_cdn_resource.test", tfjsonpath.New("id")),
				},
			},
			{
				Config: testAccCDNResourceConfigTwoGroups(rName, cname, "b"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cdn_resource.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectUnknownValue("gcore_cdn_resource.test", tfjsonpath.New("origin_group_name")),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_resource.test",
						tfjsonpath.New("origin_group_name"), knownvalue.StringExact(rName+"-b")),
					compareIDSame.AddStateValue("gcore_cdn_resource.test", tfjsonpath.New("id")),
				},
			},
		},
	})
}

// TestAccCDNResource_options sets a few options and leaves the rest to the API.
// The options the config does not set must keep their state values, so the
// plan after each apply is empty.
func TestAccCDNResource_options(t *testing.T) {
	rName := acctest.RandomName()
	cname := rName + ".example.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNResourceConfigOptions(rName, cname, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_resource.test",
						tfjsonpath.New("options").AtMapKey("gzip_on").AtMapKey("value"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("gcore_cdn_resource.test",
						tfjsonpath.New("options").AtMapKey("edge_cache_settings").AtMapKey("value"), knownvalue.StringExact("1d")),
				},
			},
			{
				Config: testAccCDNResourceConfigOptions(rName, cname, false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("gcore_cdn_resource.test", plancheck.ResourceActionUpdate),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_resource.test",
						tfjsonpath.New("options").AtMapKey("gzip_on").AtMapKey("value"), knownvalue.Bool(false)),
				},
			},
		},
	})
}

func TestAccCDNResource_import(t *testing.T) {
	rName := acctest.RandomName()
	cname := rName + ".example.com"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNResourceConfigBasic(rName, cname),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_resource.test",
						tfjsonpath.New("cname"), knownvalue.StringExact(cname)),
				},
			},
			{
				ResourceName:      "gcore_cdn_resource.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				ImportStateIdFunc: acctest.BuildImportID("gcore_cdn_resource.test", "id"),
			},
		},
	})
}

func testAccCheckCDNResourceDestroy(s *terraform.State) error {
	return acctest.CheckResourceDestroyed(s, "gcore_cdn_resource", func(client *gcore.Client, id string) error {
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing id: %w", err)
		}
		_, err = client.CDN.CDNResources.Get(context.Background(), idInt)
		return err
	})
}

func testAccCDNResourceConfigBasic(name, cname string) string {
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
}`, name, cname)
}

func testAccCDNResourceConfigWithDescription(name, cname, description string) string {
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
  description  = %[3]q
}`, name, cname, description)
}

func testAccCDNResourceConfigTwoGroups(name, cname, use string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_origin_group" "a" {
  name = "%[1]s-a"
  sources = [
    {
      source = "example.com"
    }
  ]
}

resource "gcore_cdn_origin_group" "b" {
  name = "%[1]s-b"
  sources = [
    {
      source = "example.org"
    }
  ]
}

resource "gcore_cdn_resource" "test" {
  cname        = %[2]q
  origin_group = gcore_cdn_origin_group.%[3]s.id
}`, name, cname, use)
}

func testAccCDNResourceConfigOptions(name, cname string, gzip bool) string {
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
  options = {
    gzip_on = {
      enabled = true
      value   = %[3]t
    }
    edge_cache_settings = {
      enabled = true
      value   = "1d"
    }
  }
}`, name, cname, gzip)
}
