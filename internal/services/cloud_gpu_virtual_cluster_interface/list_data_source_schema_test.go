// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_virtual_cluster_interface_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_gpu_virtual_cluster_interface"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudGPUVirtualClusterInterfacesDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_gpu_virtual_cluster_interface.CloudGPUVirtualClusterInterfacesDataSourceModel)(nil)
	schema := cloud_gpu_virtual_cluster_interface.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
