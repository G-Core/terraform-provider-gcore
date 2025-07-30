// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_load_balancer_status_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_load_balancer_status"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudLoadBalancerStatusDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_load_balancer_status.CloudLoadBalancerStatusDataSourceModel)(nil)
	schema := cloud_load_balancer_status.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
