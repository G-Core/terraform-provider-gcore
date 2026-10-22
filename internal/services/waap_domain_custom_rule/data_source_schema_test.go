// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_custom_rule_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/waap_domain_custom_rule"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestWaapDomainCustomRuleDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*waap_domain_custom_rule.WaapDomainCustomRuleDataSourceModel)(nil)
	schema := waap_domain_custom_rule.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
