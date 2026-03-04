// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_cdn_resource_rule_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cdn_cdn_resource_rule"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCdncdnResourceRuleModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cdn_cdn_resource_rule.CDNCDNResourceRuleModel)(nil)
	schema := cdn_cdn_resource_rule.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
