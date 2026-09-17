// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_pool_member_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_load_balancer_pool_member"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudLoadBalancerPoolMemberDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_load_balancer_pool_member.CloudLoadBalancerPoolMemberDataSourceModel)(nil)
	schema := cloud_load_balancer_pool_member.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
