// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_object_storage_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/storage_object_storage"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestStorageObjectStoragesDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*storage_object_storage.StorageObjectStoragesDataSourceModel)(nil)
	schema := storage_object_storage.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
