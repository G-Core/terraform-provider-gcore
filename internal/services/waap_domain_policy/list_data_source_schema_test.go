// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package waap_domain_policy_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/waap_domain_policy"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestWaapDomainPoliciesDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*waap_domain_policy.WaapDomainPoliciesDataSourceModel)(nil)
	schema := waap_domain_policy.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
