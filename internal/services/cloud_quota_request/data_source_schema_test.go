// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_quota_request_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_quota_request"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudQuotaRequestDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_quota_request.CloudQuotaRequestDataSourceModel)(nil)
	schema := cloud_quota_request.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
