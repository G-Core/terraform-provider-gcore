// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_gpu_baremetal_cluster_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_gpu_baremetal_cluster"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudGPUBaremetalClusterModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_gpu_baremetal_cluster.CloudGPUBaremetalClusterModel)(nil)
	schema := cloud_gpu_baremetal_cluster.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
