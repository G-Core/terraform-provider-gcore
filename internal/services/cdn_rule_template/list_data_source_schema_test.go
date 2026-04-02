// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_rule_template_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cdn_rule_template"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCDNRuleTemplatesDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cdn_rule_template.CDNRuleTemplatesDataSourceModel)(nil)
	schema := cdn_rule_template.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
