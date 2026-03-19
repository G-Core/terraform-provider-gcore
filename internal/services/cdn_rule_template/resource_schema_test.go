// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cdn_rule_template_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cdn_rule_template"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCDNRuleTemplateModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cdn_rule_template.CDNRuleTemplateModel)(nil)
	schema := cdn_rule_template.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
