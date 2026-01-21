// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_firewall_rule_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/waap_domain_firewall_rule"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestWaapDomainFirewallRuleModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*waap_domain_firewall_rule.WaapDomainFirewallRuleModel)(nil)
	schema := waap_domain_firewall_rule.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
