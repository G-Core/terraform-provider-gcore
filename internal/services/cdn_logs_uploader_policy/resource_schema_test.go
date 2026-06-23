// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_logs_uploader_policy_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cdn_logs_uploader_policy"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCDNLogsUploaderPolicyModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cdn_logs_uploader_policy.CDNLogsUploaderPolicyModel)(nil)
	schema := cdn_logs_uploader_policy.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
