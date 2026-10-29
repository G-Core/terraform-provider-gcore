package cdn_origin_group_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/storage"
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

// TestAccCDNOriginGroup_gcoreStorage covers an S3 origin bound to a Gcore Object
// Storage (s3_type = "gcore"). The test provisions its own storage with two buckets,
// moves the origin from one bucket to the other in place, and imports the group. Such
// an origin has no write-only credentials, so it needs no s3_credentials_version and
// its sources round-trip on import.
func TestAccCDNOriginGroup_gcoreStorage(t *testing.T) {
	rName := acctest.RandomName()
	resourceName := "gcore_cdn_origin_group.test"

	compareIDSame := statecheck.CompareValue(compare.ValuesSame())
	compareStorageID := statecheck.CompareValuePairs(
		"gcore_storage_object_storage.test", tfjsonpath.New("id"),
		resourceName, tfjsonpath.New("sources").AtSliceIndex(0).AtMapKey("config").AtMapKey("storage_id"),
		compare.ValuesSame(),
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			testAccCheckCDNOriginGroupDestroy,
			testAccCheckObjectStorageDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCDNOriginGroupConfigGcoreStorage(rName, "a"),
				ConfigStateChecks: []statecheck.StateCheck{
					compareIDSame.AddStateValue(resourceName, tfjsonpath.New("id")),
					compareStorageID,
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("sources").AtSliceIndex(0).AtMapKey("config").AtMapKey("s3_type"),
						knownvalue.StringExact("gcore")),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("sources").AtSliceIndex(0).AtMapKey("config").AtMapKey("s3_bucket_name"),
						knownvalue.StringExact(rName+"-a")),
				},
			},
			// Update: move the origin to the other bucket, in place
			{
				Config: testAccCDNOriginGroupConfigGcoreStorage(rName, "b"),
				ConfigStateChecks: []statecheck.StateCheck{
					compareIDSame.AddStateValue(resourceName, tfjsonpath.New("id")),
					compareStorageID,
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("sources").AtSliceIndex(0).AtMapKey("config").AtMapKey("s3_bucket_name"),
						knownvalue.StringExact(rName+"-b")),
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: acctest.BuildImportID(resourceName, "id"),
				ImportStateVerify: true,
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

// testAccCheckObjectStorageDestroy verifies that the storage and buckets created for a
// Gcore Object Storage origin are gone.
func testAccCheckObjectStorageDestroy(s *terraform.State) error {
	client, err := acctest.NewGcoreClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "gcore_storage_object_storage_bucket":
			storageID, err := strconv.ParseInt(rs.Primary.Attributes["storage_id"], 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse storage_id %q: %w", rs.Primary.Attributes["storage_id"], err)
			}
			name := rs.Primary.Attributes["name"]
			_, err = client.Storage.ObjectStorages.Buckets.Get(context.Background(), name,
				storage.ObjectStorageBucketGetParams{StorageID: storageID})
			if err == nil {
				return fmt.Errorf("bucket %d/%s still exists", storageID, name)
			}
			if !acctest.IsNotFoundError(err) {
				return fmt.Errorf("error checking bucket deletion: %w", err)
			}
		case "gcore_storage_object_storage":
			id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse id %q: %w", rs.Primary.ID, err)
			}
			_, err = client.Storage.ObjectStorages.Get(context.Background(), id)
			if err == nil {
				return fmt.Errorf("object storage %d still exists", id)
			}
			if !acctest.IsNotFoundError(err) {
				return fmt.Errorf("error checking object storage deletion: %w", err)
			}
		}
	}

	return nil
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

// testAccCDNOriginGroupConfigGcoreStorage provisions a storage with the buckets
// <name>-a and <name>-b, and points the origin at the one named by bucketSuffix.
func testAccCDNOriginGroupConfigGcoreStorage(name, bucketSuffix string) string {
	return fmt.Sprintf(`
resource "gcore_storage_object_storage" "test" {
  name          = %[1]q
  location_name = "luxembourg-2"
}

resource "gcore_storage_object_storage_bucket" "a" {
  storage_id = gcore_storage_object_storage.test.id
  name       = "%[1]s-a"
}

resource "gcore_storage_object_storage_bucket" "b" {
  storage_id = gcore_storage_object_storage.test.id
  name       = "%[1]s-b"
}

resource "gcore_cdn_origin_group" "test" {
  name = %[1]q
  sources = [
    {
      origin_type = "s3"
      config = {
        s3_type        = "gcore"
        storage_id     = gcore_storage_object_storage.test.id
        s3_bucket_name = gcore_storage_object_storage_bucket.%[2]s.name
      }
    }
  ]
}`, name, bucketSuffix)
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
