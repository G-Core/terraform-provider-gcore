// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_cdn_resource_rule_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cdn_cdn_resource_rule"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCdncdnResourceRuleDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cdn_cdn_resource_rule.CDNCDNResourceRuleDataSourceModel)(nil)
	schema := cdn_cdn_resource_rule.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
