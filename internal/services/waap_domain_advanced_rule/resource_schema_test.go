// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_advanced_rule_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/waap_domain_advanced_rule"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestWaapDomainAdvancedRuleModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*waap_domain_advanced_rule.WaapDomainAdvancedRuleModel)(nil)
	schema := waap_domain_advanced_rule.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
