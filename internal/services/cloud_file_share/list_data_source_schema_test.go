// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_file_share_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_file_share"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudFileSharesDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_file_share.CloudFileSharesDataSourceModel)(nil)
	schema := cloud_file_share.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
