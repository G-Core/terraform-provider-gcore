// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_virtual_cluster_image_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_gpu_virtual_cluster_image"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudGPUVirtualClusterImageDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_gpu_virtual_cluster_image.CloudGPUVirtualClusterImageDataSourceModel)(nil)
	schema := cloud_gpu_virtual_cluster_image.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
