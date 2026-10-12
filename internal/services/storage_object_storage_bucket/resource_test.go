package storage_object_storage_bucket_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/gcore-go/storage"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/G-Core/terraform-provider-gcore/internal/acctest"
)

const (
	// S3-compatible location on production, where CI runs the acceptance
	// tests. Buckets need an object storage parent, so every test provisions
	// its own.
	testAccBucketLocation = "luxembourg-2"
	testAccBucketResource = "gcore_storage_object_storage_bucket.test"
)

// TestAccStorageObjectStorageBucket_basic covers the minimal resource: only the
// required attributes, no nested blocks. It pins the computed id to the bucket
// name, checks the plan converges, and imports the bucket with the composite
// <storage_id>/<name> identifier.
func TestAccStorageObjectStorageBucket_basic(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageObjectStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfigMinimal(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(testAccBucketResource,
						tfjsonpath.New("name"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue(testAccBucketResource,
						tfjsonpath.New("id"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue(testAccBucketResource,
						tfjsonpath.New("storage_id"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccBucketConfigMinimal(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				ResourceName:      testAccBucketResource,
				ImportState:       true,
				ImportStateIdFunc: acctest.BuildImportID(testAccBucketResource, "storage_id", "name"),
				ImportStateVerify: true,
			},
			{
				// Plans the configuration with an import block. Fails unless the
				// import is a no-op, so the config must match imported state.
				Config:            testAccBucketConfigMinimal(rName),
				ResourceName:      testAccBucketResource,
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				ImportStateIdFunc: acctest.BuildImportID(testAccBucketResource, "storage_id", "name"),
			},
		},
	})
}

// TestAccStorageObjectStorageBucket_fullUpdate creates a bucket with every
// nested block set, changes every value through PATCH, checks the plan
// converges, and imports the fully configured bucket.
func TestAccStorageObjectStorageBucket_fullUpdate(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageObjectStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfigFull(rName, []string{"https://a.example.com"}, true, 30),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(testAccBucketResource,
						tfjsonpath.New("id"), knownvalue.StringExact(rName)),
					statecheck.ExpectKnownValue(testAccBucketResource,
						tfjsonpath.New("cors").AtMapKey("allowed_origins"),
						knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact("https://a.example.com")})),
					statecheck.ExpectKnownValue(testAccBucketResource,
						tfjsonpath.New("policy").AtMapKey("is_public"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(testAccBucketResource,
						tfjsonpath.New("storage_object_storage_bucket_lifecycle").AtMapKey("expiration_days"), knownvalue.Int64Exact(30)),
				},
			},
			{
				Config: testAccBucketConfigFull(rName, []string{"https://b.example.com", "https://c.example.com"}, false, 7),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccBucketResource, plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(testAccBucketResource,
						tfjsonpath.New("cors").AtMapKey("allowed_origins"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("https://b.example.com"),
							knownvalue.StringExact("https://c.example.com"),
						})),
					statecheck.ExpectKnownValue(testAccBucketResource,
						tfjsonpath.New("policy").AtMapKey("is_public"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(testAccBucketResource,
						tfjsonpath.New("storage_object_storage_bucket_lifecycle").AtMapKey("expiration_days"), knownvalue.Int64Exact(7)),
				},
			},
			{
				Config: testAccBucketConfigFull(rName, []string{"https://b.example.com", "https://c.example.com"}, false, 7),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				// The API sorts origins and removes duplicates, so order and
				// repeats in the configuration must not produce a diff.
				Config: testAccBucketConfigFull(rName, []string{"https://c.example.com", "https://b.example.com", "https://c.example.com"}, false, 7),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				ResourceName:      testAccBucketResource,
				ImportState:       true,
				ImportStateIdFunc: acctest.BuildImportID(testAccBucketResource, "storage_id", "name"),
				ImportStateVerify: true,
			},
			{
				// Plans the configuration with an import block. Fails unless the
				// import is a no-op, so the config must match imported state.
				Config:            testAccBucketConfigFull(rName, []string{"https://b.example.com", "https://c.example.com"}, false, 7),
				ResourceName:      testAccBucketResource,
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				ImportStateIdFunc: acctest.BuildImportID(testAccBucketResource, "storage_id", "name"),
			},
		},
	})
}

// TestAccStorageObjectStorageBucket_clearNestedBlocks is the acceptance guard
// for GCLOUD2-25203. Removing the cors and lifecycle blocks from config must
// send the v4 "clear" sentinels ({"allowed_origins": []} and
// {"expiration_days": 0}) rather than omit the fields, otherwise the API keeps
// the old values and state never converges.
func TestAccStorageObjectStorageBucket_clearNestedBlocks(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageObjectStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfigFull(rName, []string{"https://a.example.com"}, true, 30),
			},
			{
				Config: testAccBucketConfigPolicyOnly(rName, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccBucketResource, plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(testAccBucketResource,
						tfjsonpath.New("cors"), knownvalue.Null()),
					statecheck.ExpectKnownValue(testAccBucketResource,
						tfjsonpath.New("storage_object_storage_bucket_lifecycle"), knownvalue.Null()),
					statecheck.ExpectKnownValue(testAccBucketResource,
						tfjsonpath.New("policy").AtMapKey("is_public"), knownvalue.Bool(true)),
				},
			},
			{
				Config: testAccBucketConfigPolicyOnly(rName, true),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

// TestAccStorageObjectStorageBucket_corsWithoutAllowedOrigins runs the
// cors-block-without-allowed_origins shape end to end. On create the attribute
// is left unknown and resolved from the response; on update it must reuse the
// prior origins so the plan stays empty and the configured origins survive.
func TestAccStorageObjectStorageBucket_corsWithoutAllowedOrigins(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageObjectStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfigEmptyCors(rName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(testAccBucketResource,
						tfjsonpath.New("cors"), knownvalue.NotNull()),
				},
			},
			{
				Config: testAccBucketConfigEmptyCors(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				Config: testAccBucketConfigFull(rName, []string{"https://a.example.com"}, false, 30),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(testAccBucketResource,
						tfjsonpath.New("cors").AtMapKey("allowed_origins"),
						knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact("https://a.example.com")})),
				},
			},
			{
				// Dropping allowed_origins from a populated cors block must reuse
				// the prior value, not replan it as unknown.
				Config: testAccBucketConfigEmptyCors(rName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccBucketResource, plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(testAccBucketResource,
						tfjsonpath.New("cors").AtMapKey("allowed_origins"),
						knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact("https://a.example.com")})),
				},
			},
		},
	})
}

// TestAccStorageObjectStorageBucket_addLifecycle adds a lifecycle rule to a
// bucket created without one, which goes through PATCH rather than create.
func TestAccStorageObjectStorageBucket_addLifecycle(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageObjectStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfigMinimal(rName),
			},
			{
				Config: testAccBucketConfigLifecycle(rName, "expiration_days = 30"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(testAccBucketResource, plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(testAccBucketResource,
						tfjsonpath.New("storage_object_storage_bucket_lifecycle").AtMapKey("expiration_days"), knownvalue.Int64Exact(30)),
				},
			},
			{
				Config: testAccBucketConfigLifecycle(rName, "expiration_days = 30"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

// TestAccStorageObjectStorageBucket_rejectsClearValues checks that the API's
// clear values are rejected at plan time. The API stores them as "not
// configured" and returns the block as null, so applying them would fail with
// an inconsistent result. Nothing is created.
func TestAccStorageObjectStorageBucket_rejectsClearValues(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccBucketConfigLifecycle(rName, ""),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Missing expiration_days`),
			},
			{
				Config:      testAccBucketConfigLifecycle(rName, "expiration_days = 0"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Invalid expiration_days`),
			},
			{
				Config:      testAccBucketConfigCors(rName, "allowed_origins = []"),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Empty allowed_origins`),
			},
		},
	})
}

// TestAccStorageObjectStorageBucket_drift changes each block outside
// Terraform and expects the next plan to show it.
func TestAccStorageObjectStorageBucket_drift(t *testing.T) {
	cases := map[string]string{
		"lifecycle": `{"lifecycle":{"expiration_days":60}}`,
		"cors":      `{"cors":{"allowed_origins":["https://drift.example.com"]}}`,
		"policy":    `{"policy":{"is_public":false}}`,
	}
	for name, patch := range cases {
		t.Run(name, func(t *testing.T) {
			testAccBucketDrift(t, patch)
		})
	}
}

func testAccBucketDrift(t *testing.T, patch string) {
	rName := acctest.RandomName()
	var storageID int64

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageObjectStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfigFull(rName, []string{"https://a.example.com"}, true, 30),
				Check: func(s *terraform.State) error {
					rs, ok := s.RootModule().Resources[testAccBucketResource]
					if !ok {
						return fmt.Errorf("%s not found in state", testAccBucketResource)
					}
					var err error
					storageID, err = strconv.ParseInt(rs.Primary.Attributes["storage_id"], 10, 64)
					return err
				},
			},
			{
				PreConfig:          func() { testAccPatchBucket(t, storageID, rName, patch) },
				Config:             testAccBucketConfigFull(rName, []string{"https://a.example.com"}, true, 30),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccStorageObjectStorageBucket_importErrors checks that a malformed
// import ID and a bucket that does not exist both fail the import.
func TestAccStorageObjectStorageBucket_importErrors(t *testing.T) {
	rName := acctest.RandomName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStorageObjectStorageBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfigMinimal(rName),
			},
			{
				// The bucket name alone, without the storage ID.
				ResourceName:  testAccBucketResource,
				ImportState:   true,
				ImportStateId: rName,
				ExpectError:   regexp.MustCompile(`invalid ID`),
			},
			{
				ResourceName: testAccBucketResource,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[testAccBucketResource]
					if !ok {
						return "", fmt.Errorf("%s not found in state", testAccBucketResource)
					}
					return rs.Primary.Attributes["storage_id"] + "/" + rName + "-missing", nil
				},
				ExpectError: regexp.MustCompile(`failed to make http request`),
			},
		},
	})
}

func testAccPatchBucket(t *testing.T, storageID int64, name, body string) {
	t.Helper()
	client, err := acctest.NewGcoreClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	_, err = client.Storage.ObjectStorages.Buckets.Update(context.Background(), name,
		storage.ObjectStorageBucketUpdateParams{StorageID: storageID},
		option.WithRequestBody("application/json", []byte(body)))
	if err != nil {
		t.Fatalf("failed to patch bucket %d/%s: %v", storageID, name, err)
	}
}

func testAccCheckStorageObjectStorageBucketDestroy(s *terraform.State) error {
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

func testAccBucketStorageConfig(name string) string {
	return fmt.Sprintf(`
resource "gcore_storage_object_storage" "test" {
  name          = %[1]q
  location_name = %[2]q
}`, name, testAccBucketLocation)
}

func testAccBucketConfigMinimal(name string) string {
	return testAccBucketStorageConfig(name) + fmt.Sprintf(`
resource "gcore_storage_object_storage_bucket" "test" {
  storage_id = gcore_storage_object_storage.test.id
  name       = %[1]q
}`, name)
}

func testAccBucketConfigFull(name string, origins []string, isPublic bool, expirationDays int) string {
	quoted := ""
	for i, o := range origins {
		if i > 0 {
			quoted += ", "
		}
		quoted += strconv.Quote(o)
	}
	return testAccBucketStorageConfig(name) + fmt.Sprintf(`
resource "gcore_storage_object_storage_bucket" "test" {
  storage_id = gcore_storage_object_storage.test.id
  name       = %[1]q

  cors = {
    allowed_origins = [%[2]s]
  }

  policy = {
    is_public = %[3]t
  }

  storage_object_storage_bucket_lifecycle = {
    expiration_days = %[4]d
  }
}`, name, quoted, isPublic, expirationDays)
}

func testAccBucketConfigPolicyOnly(name string, isPublic bool) string {
	return testAccBucketStorageConfig(name) + fmt.Sprintf(`
resource "gcore_storage_object_storage_bucket" "test" {
  storage_id = gcore_storage_object_storage.test.id
  name       = %[1]q

  policy = {
    is_public = %[2]t
  }
}`, name, isPublic)
}

func testAccBucketConfigEmptyCors(name string) string {
	return testAccBucketStorageConfig(name) + fmt.Sprintf(`
resource "gcore_storage_object_storage_bucket" "test" {
  storage_id = gcore_storage_object_storage.test.id
  name       = %[1]q

  cors = {}
}`, name)
}

func testAccBucketConfigLifecycle(name, body string) string {
	return testAccBucketStorageConfig(name) + fmt.Sprintf(`
resource "gcore_storage_object_storage_bucket" "test" {
  storage_id = gcore_storage_object_storage.test.id
  name       = %[1]q

  storage_object_storage_bucket_lifecycle = {
    %[2]s
  }
}`, name, body)
}

func testAccBucketConfigCors(name, body string) string {
	return testAccBucketStorageConfig(name) + fmt.Sprintf(`
resource "gcore_storage_object_storage_bucket" "test" {
  storage_id = gcore_storage_object_storage.test.id
  name       = %[1]q

  cors = {
    %[2]s
  }
}`, name, body)
}
