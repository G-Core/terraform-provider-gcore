// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_pool_health_monitor_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_load_balancer_pool_health_monitor"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudLoadBalancerPoolHealthMonitorModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_load_balancer_pool_health_monitor.CloudLoadBalancerPoolHealthMonitorModel)(nil)
	schema := cloud_load_balancer_pool_health_monitor.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
