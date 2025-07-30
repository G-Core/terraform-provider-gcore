// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_security_group_rule_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_security_group_rule"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudSecurityGroupRuleModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_security_group_rule.CloudSecurityGroupRuleModel)(nil)
	schema := cloud_security_group_rule.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
