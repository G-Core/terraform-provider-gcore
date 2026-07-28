// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cloud_k8s_flavor_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_k8s_flavor"
	"github.com/G-Core/terraform-provider-gcore/internal/test_helpers"
)

func TestCloudK8SFlavorsDataSourceModelSchemaParity(t *testing.T) {
	t.Parallel()
	model := (*cloud_k8s_flavor.CloudK8SFlavorsDataSourceModel)(nil)
	schema := cloud_k8s_flavor.ListDataSourceSchema(context.TODO())
	errs := test_helpers.ValidateDataSourceModelSchemaIntegrity(model, schema)
	errs.Report(t)
}
