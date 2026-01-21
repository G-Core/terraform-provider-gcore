// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_advanced_rule_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/waap_domain_advanced_rule"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestWaapDomainAdvancedRulesDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*waap_domain_advanced_rule.WaapDomainAdvancedRulesDataSourceModel)(nil)
	schema := waap_domain_advanced_rule.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
