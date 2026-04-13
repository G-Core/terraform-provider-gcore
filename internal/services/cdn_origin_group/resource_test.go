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
						tfjsonpath.New("use_next"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("gcore_cdn_origin_group.test",
						tfjsonpath.New("proxy_next_upstream"), knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("error"),
							knownvalue.StringExact("timeout"),
						})),
					statecheck.ExpectKnownValue("gcore_cdn_origin_group.test",
						tfjsonpath.New("sources"), knownvalue.ListSizeExact(1)),
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
					compareIDSame.AddStateValue(
						"gcore_cdn_origin_group.test",
						tfjsonpath.New("id"),
					),
				},
			},
			{
				Config: testAccCDNOriginGroupConfigWithSource(rName, "example.org"),
				ConfigStateChecks: []statecheck.StateCheck{
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

func TestAccCDNOriginGroup_s3Inline(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNOriginGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNOriginGroupConfigS3Inline(rName, "test-bucket-a"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_origin_group.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cdn_origin_group.test",
						tfjsonpath.New("sources"), knownvalue.ListSizeExact(1)),
				},
			},
			// Update: change bucket name
			{
				Config: testAccCDNOriginGroupConfigS3Inline(rName, "test-bucket-b"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_origin_group.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
				},
			},
		},
	})
}

func TestAccCDNOriginGroup_s3InlineImport(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNOriginGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNOriginGroupConfigS3Inline(rName, "test-bucket"),
			},
			{
				ResourceName:      "gcore_cdn_origin_group.test",
				ImportState:       true,
				ImportStateId:     "",
				ImportStateIdFunc: acctest.BuildImportID("gcore_cdn_origin_group.test", "id"),
				ImportStateVerify: true,
				// s3_credentials_version is Terraform-only (not in API)
				// sources contains write-only fields that won't round-trip
				ImportStateVerifyIgnore: []string{"s3_credentials_version", "sources"},
			},
		},
	})
}

func TestAccCDNOriginGroup_mixed(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNOriginGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNOriginGroupConfigMixed(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("gcore_cdn_origin_group.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("gcore_cdn_origin_group.test",
						tfjsonpath.New("use_next"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("gcore_cdn_origin_group.test",
						tfjsonpath.New("sources"), knownvalue.ListSizeExact(2)),
				},
			},
		},
	})
}

func TestAccCDNOriginGroup_mixedImport(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCDNOriginGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCDNOriginGroupConfigMixed(rName),
			},
			{
				ResourceName:            "gcore_cdn_origin_group.test",
				ImportState:             true,
				ImportStateIdFunc:       acctest.BuildImportID("gcore_cdn_origin_group.test", "id"),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"s3_credentials_version", "sources"},
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
  name                = %[1]q
  use_next            = true
  proxy_next_upstream = ["error", "timeout"]
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

func testAccCDNOriginGroupConfigS3Inline(name string, bucket string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_origin_group" "test" {
  name                   = %[1]q
  s3_credentials_version = 1
  sources = [
    {
      origin_type = "s3"
      enabled     = true
      config = {
        s3_type              = "amazon"
        s3_bucket_name       = %[2]q
        s3_access_key_id     = "dummy-access-key-1234"
        s3_secret_access_key = "dummy-secret-key-12345678"
        s3_region            = "eu-west-1"
      }
    }
  ]
}`, name, bucket)
}

func testAccCDNOriginGroupConfigMixed(name string) string {
	return fmt.Sprintf(`
resource "gcore_cdn_origin_group" "test" {
  name                   = %[1]q
  use_next               = true
  s3_credentials_version = 1
  sources = [
    {
      source  = "example.com"
      enabled = true
    },
    {
      origin_type = "s3"
      enabled     = true
      backup      = true
      config = {
        s3_type              = "amazon"
        s3_bucket_name       = "dummy-bucket"
        s3_access_key_id     = "dummy-access-key-1234"
        s3_secret_access_key = "dummy-secret-key-12345678"
        s3_region            = "eu-west-1"
      }
    }
  ]
}`, name)
}
