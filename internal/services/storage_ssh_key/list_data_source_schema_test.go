// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_ssh_key_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/storage_ssh_key"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestStorageSSHKeysDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*storage_ssh_key.StorageSSHKeysDataSourceModel)(nil)
	schema := storage_ssh_key.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
