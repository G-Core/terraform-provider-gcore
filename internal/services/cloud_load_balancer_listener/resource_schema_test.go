// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_listener_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_load_balancer_listener"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudLoadBalancerListenerModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_load_balancer_listener.CloudLoadBalancerListenerModel)(nil)
	schema := cloud_load_balancer_listener.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
