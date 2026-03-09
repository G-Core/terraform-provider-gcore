// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_file_share_access_rule_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_file_share_access_rule"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudFileShareAccessRuleModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_file_share_access_rule.CloudFileShareAccessRuleModel)(nil)
	schema := cloud_file_share_access_rule.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
