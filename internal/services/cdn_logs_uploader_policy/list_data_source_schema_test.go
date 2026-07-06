// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_logs_uploader_policy_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cdn_logs_uploader_policy"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCDNLogsUploaderPoliciesDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cdn_logs_uploader_policy.CDNLogsUploaderPoliciesDataSourceModel)(nil)
	schema := cdn_logs_uploader_policy.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
