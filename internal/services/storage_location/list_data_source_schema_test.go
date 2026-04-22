// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_location_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/storage_location"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestStorageLocationsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*storage_location.StorageLocationsDataSourceModel)(nil)
	schema := storage_location.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
