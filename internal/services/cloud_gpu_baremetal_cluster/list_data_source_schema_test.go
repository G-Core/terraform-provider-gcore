// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_baremetal_cluster_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_gpu_baremetal_cluster"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudGPUBaremetalClustersDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_gpu_baremetal_cluster.CloudGPUBaremetalClustersDataSourceModel)(nil)
	schema := cloud_gpu_baremetal_cluster.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
