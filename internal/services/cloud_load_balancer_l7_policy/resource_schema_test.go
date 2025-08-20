// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_l7_policy_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_load_balancer_l7_policy"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudLoadBalancerL7PolicyModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_load_balancer_l7_policy.CloudLoadBalancerL7PolicyModel)(nil)
	schema := cloud_load_balancer_l7_policy.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
