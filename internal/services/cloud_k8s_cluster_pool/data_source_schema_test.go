// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_k8s_cluster_pool_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_k8s_cluster_pool"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudK8SClusterPoolDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_k8s_cluster_pool.CloudK8SClusterPoolDataSourceModel)(nil)
	schema := cloud_k8s_cluster_pool.DataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
