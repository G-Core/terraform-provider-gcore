// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_access_key_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/storage_access_key"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestStorageAccessKeyModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*storage_access_key.StorageAccessKeyModel)(nil)
	schema := storage_access_key.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
