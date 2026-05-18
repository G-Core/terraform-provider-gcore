// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_object_storage_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/storage_object_storage"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestStorageObjectStorageDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*storage_object_storage.StorageObjectStorageDataSourceModel)(nil)
	schema := storage_object_storage.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
