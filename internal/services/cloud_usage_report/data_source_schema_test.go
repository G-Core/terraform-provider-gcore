// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_usage_report_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_usage_report"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudUsageReportDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_usage_report.CloudUsageReportDataSourceModel)(nil)
	schema := cloud_usage_report.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
