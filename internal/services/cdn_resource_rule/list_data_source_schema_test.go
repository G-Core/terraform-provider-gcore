// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_resource_rule_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cdn_resource_rule"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCDNResourceRulesDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cdn_resource_rule.CDNResourceRulesDataSourceModel)(nil)
	schema := cdn_resource_rule.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
