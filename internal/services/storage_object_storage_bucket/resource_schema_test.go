// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_object_storage_bucket_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/storage_object_storage_bucket"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestStorageObjectStorageBucketModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*storage_object_storage_bucket.StorageObjectStorageBucketModel)(nil)
	schema := storage_object_storage_bucket.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
