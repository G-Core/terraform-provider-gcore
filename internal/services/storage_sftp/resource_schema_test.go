// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package storage_sftp_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/storage_sftp"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestStorageSftpModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*storage_sftp.StorageSftpModel)(nil)
	schema := storage_sftp.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
