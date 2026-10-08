package storage_object_storage_bucket_test

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

func TestAccStorageObjectStorageBucketDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageObjectStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketDataSourceConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_storage_object_storage_bucket.test",
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("data.gcore_storage_object_storage_bucket.test",
						tfjsonpath.New("id"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("data.gcore_storage_object_storage_bucket.test",
						tfjsonpath.New("policy").AtMapKey("is_public"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("data.gcore_storage_object_storage_bucket.test",
						tfjsonpath.New("storage_object_storage_bucket_lifecycle").AtMapKey("expiration_days"), knownvalue.Int64Exact(30)),
					statecheck.CompareValuePairs(
						testAccBucketResource, tfjsonpath.New("storage_id"),
						"data.gcore_storage_object_storage_bucket.test", tfjsonpath.New("storage_id"),
						compare.ValuesSame(),
					),
				},
			},
		},
	})
}

func TestAccStorageObjectStorageBucketsDataSource_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageObjectStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketsDataSourceConfig(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.gcore_storage_object_storage_buckets.test",
						tfjsonpath.New("items"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue("data.gcore_storage_object_storage_buckets.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue("data.gcore_storage_object_storage_buckets.test",
						tfjsonpath.New("items").AtSliceIndex(0).AtMapKey("id"), knownvalue.StringExact(rName)),
				},
			},
		},
	})
}

func testAccBucketDataSourceConfig(name string) string {
	return testAccBucketConfigFull(name, []string{"https://a.example.com"}, true, 30) + `
data "gcore_storage_object_storage_bucket" "test" {
  storage_id = gcore_storage_object_storage_bucket.test.storage_id
  name       = gcore_storage_object_storage_bucket.test.name
}`
}

func testAccBucketsDataSourceConfig(name string) string {
	return testAccBucketConfigMinimal(name) + fmt.Sprintf(`
data "gcore_storage_object_storage_buckets" "test" {
  storage_id = gcore_storage_object_storage.test.id
  depends_on = [gcore_storage_object_storage_bucket.test]
}`)
}
