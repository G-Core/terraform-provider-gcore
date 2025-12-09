// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_k8s_cluster_test

import (
	"context"
	"testing"

	"github.com/stainless-sdks/gcore-terraform/internal/services/cloud_k8s_cluster"
	"github.com/stainless-sdks/gcore-terraform/internal/test_helpers"
)

func TestCloudK8SClusterModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_k8s_cluster.CloudK8SClusterModel)(nil)
	schema := cloud_k8s_cluster.ResourceSchema(context.TODO())
	errs := test_helpers.ValidateResourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
