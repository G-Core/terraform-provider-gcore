// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_baremetal_cluster_flavor_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_gpu_baremetal_cluster_flavor"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudGPUBaremetalClusterFlavorsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_gpu_baremetal_cluster_flavor.CloudGPUBaremetalClusterFlavorsDataSourceModel)(nil)
	schema := cloud_gpu_baremetal_cluster_flavor.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
