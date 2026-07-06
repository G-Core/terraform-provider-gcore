// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_quota_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_quota"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudQuotaDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_quota.CloudQuotaDataSourceModel)(nil)
	schema := cloud_quota.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
