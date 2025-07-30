// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_l7_policy_rule_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_load_balancer_l7_policy_rule"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudLoadBalancerL7PolicyRuleModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_load_balancer_l7_policy_rule.CloudLoadBalancerL7PolicyRuleModel)(nil)
	schema := cloud_load_balancer_l7_policy_rule.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
